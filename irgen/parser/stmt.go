package parser

import (
	"hash/fnv"
	"strings"

	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorsx "github.com/nagarajRPoojari/niyama/irgen/error"
	"github.com/nagarajRPoojari/niyama/irgen/lexer"
)

func parseStmt(p *Parser) ast.Statement {
	stmt_fn, exists := statement_table[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	return parseExpressionStmt(p)
}

func parseExpressionStmt(p *Parser) ast.ExpressionStatement {
	expression := parseExpr(p, default_bp)
	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStatement{
		SourceLoc:  ast.SourceLoc(p.currentToken().Src),
		Expression: expression,
	}
}

func parseBlockStmt(p *Parser) ast.Statement {
	p.expect(lexer.OPEN_CURLY)
	body := []ast.Statement{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parseStmt(p))
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.BlockStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Body:      body,
	}
}

func parseVarDeclStmt(p *Parser) ast.Statement {
	var explicitType ast.Type
	p.expect(lexer.SAY)

	var isInternal bool
	if p.currentTokenKind() == lexer.INTERNAL {
		isInternal = true
		p.move()
	}

	var isStatic bool
	if p.currentTokenKind() == lexer.STATIC {
		isStatic = true
		p.move()
	}

	symbolName := p.currentToken()
	if p.currentTokenKind() != lexer.IDENTIFIER {
		errorsx.PanicParserError(
			"unexpected keyword in variable declaration",
			p.currentToken().Src.FilePath,
			p.currentToken().Src.Line,
			p.currentToken().Src.Col,
		)
	} else {
		p.move()
	}

	if p.currentTokenKind() == lexer.COLON {
		p.move()

		atomic := false
		if p.currentTokenKind() == lexer.ATOMIC {
			atomic = true
			p.move()
		}

		explicitType = parse_type(p, default_bp)
		if atomic {
			explicitType.SetAtomic()
		}
	}

	var assignmentValue ast.Expression
	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assignmentValue = parseExpr(p, assignment)
	} else if explicitType == nil {
		panic("Missing explicit type for variable declaration.")
	}

	p.expect(lexer.SEMI_COLON)

	return ast.VariableDeclarationStatement{
		SourceLoc:     ast.SourceLoc(p.currentToken().Src),
		Identifier:    symbolName.Value,
		AssignedValue: assignmentValue,
		ExplicitType:  explicitType,
		IsStatic:      isStatic,
		IsInternal:    isInternal,
	}
}

func parseFnParamsAndBody(p *Parser) ([]ast.Parameter, ast.Type, []ast.Statement) {
	functionParams := make([]ast.Parameter, 0)

	p.expect(lexer.OPEN_PAREN)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		paramName := p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)
		paramType := parse_type(p, default_bp)

		functionParams = append(functionParams, ast.Parameter{
			Name: paramName,
			Type: paramType,
		})

		if !p.currentToken().IsOneOfMany(lexer.CLOSE_PAREN, lexer.EOF) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PAREN)
	var returnType ast.Type

	if p.currentTokenKind() == lexer.COLON {
		p.move()
		returnType = parse_type(p, default_bp)
	}

	functionBody := ast.ExpectStmt[ast.BlockStatement](parseBlockStmt(p)).Body

	return functionParams, returnType, functionBody
}

func parseFuncDeclaration(p *Parser) ast.Statement {
	p.move()
	startToken := p.move()
	var isStatic bool
	var functionName string
	var isInternal bool

	if startToken.Kind == lexer.INTERNAL {
		isInternal = true
		startToken = p.move()
	}

	if startToken.Kind == lexer.STATIC {
		isStatic = true
		functionName = p.expect(lexer.IDENTIFIER).Value
	} else {
		if startToken.Kind == lexer.IDENTIFIER {
			functionName = startToken.Value
		} else {
			errorsx.PanicParserError(
				"unexpected keyword after fn",
				p.currentToken().Src.FilePath,
				p.currentToken().Src.Line,
				p.currentToken().Src.Col,
			)
		}
	}
	functionParams, returnType, functionBody := parseFnParamsAndBody(p)

	return ast.FunctionDefinitionStatement{
		SourceLoc:  ast.SourceLoc(p.currentToken().Src),
		Parameters: functionParams,
		ReturnType: returnType,
		Body:       functionBody,
		Name:       functionName,
		IsStatic:   isStatic,
		Hash:       funcHash(functionParams, returnType),
		IsInternal: isInternal,
	}
}

func funcHash(params []ast.Parameter, ret ast.Type) uint32 {
	var s = ""
	if ret != nil {
		s = ret.Get()
	}
	for _, i := range params {
		s += i.Type.Get()
	}

	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func parseIfStmt(p *Parser) ast.Statement {
	p.move()
	condition := parseExpr(p, assignment)
	consequent := parseBlockStmt(p)

	var alternate ast.Statement
	if p.currentTokenKind() == lexer.ELSE {
		p.move()

		if p.currentTokenKind() == lexer.IF {
			alternate = parseIfStmt(p)
		} else {
			alternate = parseBlockStmt(p)
		}
	}

	return ast.IfStatement{
		SourceLoc:  ast.SourceLoc(p.currentToken().Src),
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func parseImportStmt(p *Parser) ast.Statement {
	p.move()
	var importAlias string
	importName := strings.ReplaceAll(p.expect(lexer.STRING).Value, "/", ".")
	importName = importName[1 : len(importName)-1]

	if p.currentTokenKind() == lexer.AS {
		p.move()
		importAlias = p.expect(lexer.IDENTIFIER).Value
	} else {
		paths := strings.Split(importName, ".")
		importAlias = paths[len(paths)-1]
	}

	p.expect(lexer.SEMI_COLON)
	return ast.ImportStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Name:      importName,
		Alias:     importAlias,
	}
}

func parseForeachStmt(p *Parser) ast.Statement {
	p.move()
	valueName := p.expect(lexer.IDENTIFIER).Value

	var index bool
	if p.currentTokenKind() == lexer.COMMA {
		p.expect(lexer.COMMA)
		p.expect(lexer.IDENTIFIER)
		index = true
	}

	p.expect(lexer.IN)
	iterable := parseExpr(p, default_bp)
	body := ast.ExpectStmt[ast.BlockStatement](parseBlockStmt(p)).Body

	return ast.ForeachStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Value:     valueName,
		Index:     index,
		Iterable:  iterable,
		Body:      body,
	}
}

func parseWhileStmt(p *Parser) ast.Statement {
	p.move()
	condition := parseExpr(p, assignment)
	body := ast.ExpectStmt[ast.BlockStatement](parseBlockStmt(p)).Body

	return ast.WhileStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Condition: condition,
		Body:      body,
	}
}

func parseClassDeclStmt(p *Parser) ast.Statement {
	p.move()
	var isInternal bool
	if p.currentTokenKind() == lexer.INTERNAL {
		isInternal = true
		p.move()
	}
	className := p.expect(lexer.IDENTIFIER).Value
	var implements string
	if p.currentTokenKind() == lexer.COLON {
		p.move()
		implementsList := []string{p.expect(lexer.IDENTIFIER).Value}
		for p.currentTokenKind() == lexer.DOT {
			p.move()
			implementsList = append(implementsList, p.expect(lexer.IDENTIFIER).Value)
		}
		implements = strings.Join(implementsList, ".")
	}
	classBody := parseBlockStmt(p)

	return ast.ClassDeclarationStatement{
		SourceLoc:  ast.SourceLoc(p.currentToken().Src),
		Name:       className,
		Body:       ast.ExpectStmt[ast.BlockStatement](classBody).Body,
		Implements: implements,
		IsInternal: isInternal,
	}
}

func parseInterfaceDeclStmt(p *Parser) ast.Statement {
	p.move()
	interfaceName := p.expect(lexer.IDENTIFIER).Value
	interfaceBody := parseBlockStmt(p)

	return ast.InterfaceDeclarationStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Name:      interfaceName,
		Body:      ast.ExpectStmt[ast.BlockStatement](interfaceBody).Body,
	}
}

func parseFuncReturnStmt(p *Parser) ast.Statement {
	p.expect(lexer.RETURN)

	if p.currentTokenKind() == lexer.NULL {
		p.move()
		p.expect(lexer.SEMI_COLON)
		return ast.ReturnStatement{
			SourceLoc: ast.SourceLoc(p.currentToken().Src),
		}
	}
	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.move()
		return ast.ReturnStatement{
			SourceLoc: ast.SourceLoc(p.currentToken().Src),
			IsVoid:    true,
		}
	}
	exp := parseExpressionStmt(p)

	return ast.ReturnStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Value:     exp,
	}
}

func parseBreakStmt(p *Parser) ast.Statement {
	p.expect(lexer.BREAK)
	p.expect(lexer.SEMI_COLON)

	return ast.BreakStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
	}
}
