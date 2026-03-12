package parser

import (
	"hash/fnv"
	"strings"

	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorsx "github.com/nagarajRPoojari/picasso/irgen/error"
	"github.com/nagarajRPoojari/picasso/irgen/lexer"
)

func parseStmt(p *Parser) ast.Statement {
	stmt_fn, exists := statement_table[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	return parseExpressionStmt(p)
}

func parseExpressionStmt(p *Parser) ast.ExpressionStatement {
	// Parse first expression - for single assignments this will consume the whole thing
	// For multi-assignments, we need to detect the pattern
	firstExpr := parseExpr(p, default_bp)

	// Check if firstExpr is an assignment expression (single assignment case)
	if assignExpr, ok := firstExpr.(ast.AssignmentExpression); ok {
		// This is a single assignment that was already parsed
		p.expect(lexer.SEMI_COLON)
		return ast.ExpressionStatement{
			SourceLoc:  ast.SourceLoc(p.currentToken().Src),
			Expression: assignExpr,
		}
	}

	// Check if we have comma after first expression (multiple LHS for assignment)
	if p.currentTokenKind() == lexer.COMMA {
		assignees := []ast.Expression{firstExpr}

		// Parse remaining assignees - use assignment bp to stop before '='
		for p.currentTokenKind() == lexer.COMMA {
			p.move() // consume comma
			assignees = append(assignees, parseExpr(p, assignment))
		}

		// Now we must have '=' for multi-assignment
		if p.currentTokenKind() == lexer.ASSIGNMENT {
			p.move() // consume '='

			// Parse RHS values
			assignedValues := []ast.Expression{}
			assignedValues = append(assignedValues, parseExpr(p, assignment))

			for p.currentTokenKind() == lexer.COMMA {
				p.move()
				assignedValues = append(assignedValues, parseExpr(p, assignment))
			}

			p.expect(lexer.SEMI_COLON)

			// Create multi-assignment expression
			return ast.ExpressionStatement{
				SourceLoc: ast.SourceLoc(p.currentToken().Src),
				Expression: ast.AssignmentExpression{
					SourceLoc:      ast.SourceLoc(p.currentToken().Src),
					Assignees:      assignees,
					AssignedValues: assignedValues,
				},
			}
		}
	}

	// Regular expression statement (function call, etc.)
	p.expect(lexer.SEMI_COLON)
	return ast.ExpressionStatement{
		SourceLoc:  ast.SourceLoc(p.currentToken().Src),
		Expression: firstExpr,
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

	// Parse first identifier
	identifiers := []string{}
	explicitTypes := []ast.Type{}

	symbolName := p.currentToken()
	if p.currentTokenKind() != lexer.IDENTIFIER {
		errorsx.PanicParserError(
			"unexpected keyword in variable declaration",
			p.currentToken().Src.FilePath,
			p.currentToken().Src.Line,
			p.currentToken().Src.Col,
		)
	}
	identifiers = append(identifiers, symbolName.Value)
	p.move()

	// Check for multiple declarations: say a: int, b: int
	if p.currentTokenKind() == lexer.COLON {
		p.move()

		atomic := false
		if p.currentTokenKind() == lexer.ATOMIC {
			atomic = true
			p.move()
		}

		explicitType := parse_type(p, default_bp)
		if atomic {
			explicitType.SetAtomic()
		}
		explicitTypes = append(explicitTypes, explicitType)

		// Parse additional variable declarations
		for p.currentTokenKind() == lexer.COMMA {
			p.move() // consume comma

			nextName := p.expect(lexer.IDENTIFIER).Value
			identifiers = append(identifiers, nextName)

			if p.currentTokenKind() == lexer.COLON {
				p.move()

				atomic := false
				if p.currentTokenKind() == lexer.ATOMIC {
					atomic = true
					p.move()
				}

				nextType := parse_type(p, default_bp)
				if atomic {
					nextType.SetAtomic()
				}
				explicitTypes = append(explicitTypes, nextType)
			} else {
				// No type specified for this variable
				explicitTypes = append(explicitTypes, nil)
			}
		}
	}

	// Parse assignment values
	var assignmentValues []ast.Expression
	if p.currentTokenKind() == lexer.ASSIGNMENT {
		p.move()

		assignmentValues = append(assignmentValues, parseExpr(p, assignment))

		for p.currentTokenKind() == lexer.COMMA {
			p.move()
			assignmentValues = append(assignmentValues, parseExpr(p, assignment))
		}
	} else if len(explicitTypes) == 0 || explicitTypes[0] == nil {
		panic("Missing explicit type for variable declaration.")
	}

	p.expect(lexer.SEMI_COLON)

	// Check for reserved keywords
	for _, id := range identifiers {
		if _, ok := reserved_keywords[id]; ok {
			errorsx.PanicParserError(
				"use of reserved keyword",
				p.currentToken().Src.FilePath,
				p.currentToken().Src.Line,
				p.currentToken().Src.Col,
			)
		}
	}

	// Single variable declaration (backward compatibility)
	if len(identifiers) == 1 {
		var explicitType ast.Type
		if len(explicitTypes) > 0 {
			explicitType = explicitTypes[0]
		}
		var assignmentValue ast.Expression
		if len(assignmentValues) > 0 {
			assignmentValue = assignmentValues[0]
		}

		return ast.VariableDeclarationStatement{
			SourceLoc:     ast.SourceLoc(p.currentToken().Src),
			Identifier:    identifiers[0],
			AssignedValue: assignmentValue,
			ExplicitType:  explicitType,
			IsStatic:      isStatic,
			IsInternal:    isInternal,
		}
	}

	// Multiple variable declarations
	return ast.VariableDeclarationStatement{
		SourceLoc:      ast.SourceLoc(p.currentToken().Src),
		Identifiers:    identifiers,
		ExplicitTypes:  explicitTypes,
		AssignedValues: assignmentValues,
		IsStatic:       isStatic,
		IsInternal:     isInternal,
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

	if _, ok := reserved_keywords[functionName]; ok {
		errorsx.PanicParserError(
			"use of reserved keyword",
			p.currentToken().Src.FilePath,
			p.currentToken().Src.Line,
			p.currentToken().Src.Col,
		)
	}

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

	// Parse multiple return values separated by commas
	values := []ast.Expression{}
	values = append(values, parseExpr(p, assignment))

	for p.currentTokenKind() == lexer.COMMA {
		p.move() // consume comma
		values = append(values, parseExpr(p, assignment))
	}

	p.expect(lexer.SEMI_COLON)

	// If single value, use old format for backward compatibility
	if len(values) == 1 {
		return ast.ReturnStatement{
			SourceLoc: ast.SourceLoc(p.currentToken().Src),
			Value: ast.ExpressionStatement{
				SourceLoc:  ast.SourceLoc(p.currentToken().Src),
				Expression: values[0],
			},
		}
	}

	// Multiple values
	return ast.ReturnStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Values:    values,
	}
}

func parseBreakStmt(p *Parser) ast.Statement {
	p.expect(lexer.BREAK)
	p.expect(lexer.SEMI_COLON)

	return ast.BreakStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
	}
}

func parseAtomicBlockStmt(p *Parser) ast.Statement {
	p.expect(lexer.OPEN_ATOMIC)
	body := []ast.Statement{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_ATOMIC {
		body = append(body, parseStmt(p))
	}

	p.expect(lexer.CLOSE_ATOMIC)
	return ast.AtomicBlockStatement{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Body:      body,
	}
}
