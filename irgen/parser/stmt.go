package parser

import (
	"hash/fnv"

	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/lexer"
)

func parse_stmt(p *Parser) ast.Statement {
	stmt_fn, exists := statement_table[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	return parse_expression_stmt(p)
}

func parse_expression_stmt(p *Parser) ast.ExpressionStatement {
	expression := parse_expr(p, default_bp)
	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStatement{
		Expression: expression,
	}
}

func parse_block_stmt(p *Parser) ast.Statement {
	p.expect(lexer.OPEN_CURLY)
	body := []ast.Statement{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parse_stmt(p))
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.BlockStatement{
		Body: body,
	}
}

func parse_var_decl_stmt(p *Parser) ast.Statement {
	var explicitType ast.Type

	startToken := p.move().Kind

	isConstant := startToken == lexer.CONST

	nextToken := p.move()
	var isStatic bool
	if nextToken.Kind == lexer.STATIC {
		isStatic = true
		nextToken = p.move()
	}
	if nextToken.Kind != lexer.IDENTIFIER {
		panic("unexpected keyword in variable declaration")
	}

	symbolName := nextToken
	if p.currentTokenKind() == lexer.COLON {
		p.expect(lexer.COLON)

		atomic := false
		if p.currentTokenKind() == lexer.ATOMIC {
			atomic = true
			nextToken = p.move()
		}

		explicitType = parse_type(p, default_bp)
		if atomic {
			explicitType.SetAtomic()
		}
	}

	var assignmentValue ast.Expression
	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assignmentValue = parse_expr(p, assignment)
	} else if explicitType == nil {
		panic("Missing explicit type for variable declaration.")
	}

	p.expect(lexer.SEMI_COLON)

	if isConstant && assignmentValue == nil {
		panic("Cannot define constant variable without providing default value.")
	}

	return ast.VariableDeclarationStatement{
		Constant:      isConstant,
		Identifier:    symbolName.Value,
		AssignedValue: assignmentValue,
		ExplicitType:  explicitType,
		IsStatic:      isStatic,
	}
}

func parse_fn_params_and_body(p *Parser) ([]ast.Parameter, ast.Type, []ast.Statement) {
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

	functionBody := ast.ExpectStmt[ast.BlockStatement](parse_block_stmt(p)).Body

	return functionParams, returnType, functionBody
}

func parse_fn_declaration(p *Parser) ast.Statement {
	p.move()
	startToken := p.move()
	var isStatic bool
	var functionName string
	if startToken.Kind == lexer.STATIC {
		isStatic = true
		functionName = p.expect(lexer.IDENTIFIER).Value
	} else {
		if startToken.Kind == lexer.IDENTIFIER {
			functionName = startToken.Value
		} else {
			panic("unexpected keyword after fn")
		}
	}
	functionParams, returnType, functionBody := parse_fn_params_and_body(p)

	return ast.FunctionDefinitionStatement{
		Parameters: functionParams,
		ReturnType: returnType,
		Body:       functionBody,
		Name:       functionName,
		IsStatic:   isStatic,
		Hash:       func_hash(functionParams, returnType),
	}
}

func func_hash(params []ast.Parameter, ret ast.Type) uint32 {
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

func parse_if_stmt(p *Parser) ast.Statement {
	p.move()
	condition := parse_expr(p, assignment)
	consequent := parse_block_stmt(p)

	var alternate ast.Statement
	if p.currentTokenKind() == lexer.ELSE {
		p.move()

		if p.currentTokenKind() == lexer.IF {
			alternate = parse_if_stmt(p)
		} else {
			alternate = parse_block_stmt(p)
		}
	}

	return ast.IfStatement{
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func parse_import_stmt(p *Parser) ast.Statement {
	p.move()
	var importFrom string
	importName := p.expect(lexer.IDENTIFIER).Value

	if p.currentTokenKind() == lexer.FROM {
		p.move()
		importFrom = p.expect(lexer.IDENTIFIER).Value
	} else {
		importFrom = importName
	}

	p.expect(lexer.SEMI_COLON)
	return ast.ImportStatement{
		Name: importName,
		From: importFrom,
	}
}

func parse_foreach_stmt(p *Parser) ast.Statement {
	p.move()
	valueName := p.expect(lexer.IDENTIFIER).Value

	var index bool
	if p.currentTokenKind() == lexer.COMMA {
		p.expect(lexer.COMMA)
		p.expect(lexer.IDENTIFIER)
		index = true
	}

	p.expect(lexer.IN)
	iterable := parse_expr(p, default_bp)
	body := ast.ExpectStmt[ast.BlockStatement](parse_block_stmt(p)).Body

	return ast.ForeachStatement{
		Value:    valueName,
		Index:    index,
		Iterable: iterable,
		Body:     body,
	}
}

func parse_while_stmt(p *Parser) ast.Statement {
	p.move()
	condition := parse_expr(p, assignment)
	body := ast.ExpectStmt[ast.BlockStatement](parse_block_stmt(p)).Body

	return ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func parse_class_declaration_stmt(p *Parser) ast.Statement {
	p.move()
	className := p.expect(lexer.IDENTIFIER).Value
	var implements string
	if p.currentTokenKind() == lexer.COLON {
		p.move()
		implements = p.expect(lexer.IDENTIFIER).Value
	}
	classBody := parse_block_stmt(p)

	return ast.ClassDeclarationStatement{
		Name:       className,
		Body:       ast.ExpectStmt[ast.BlockStatement](classBody).Body,
		Implements: implements,
	}
}

func parse_interface_declaration_stmt(p *Parser) ast.Statement {
	p.move()
	interfaceName := p.expect(lexer.IDENTIFIER).Value
	interfaceBody := parse_block_stmt(p)

	return ast.InterfaceDeclarationStatement{
		Name: interfaceName,
		Body: ast.ExpectStmt[ast.BlockStatement](interfaceBody).Body,
	}
}

func parse_function_return_stmt(p *Parser) ast.Statement {
	p.expect(lexer.RETURN)

	if p.currentTokenKind() == lexer.NULL {
		p.move()
		p.expect(lexer.SEMI_COLON)
		return ast.ReturnStatement{}
	}
	if p.currentTokenKind() == lexer.SEMI_COLON {
		p.move()
		return ast.ReturnStatement{
			IsVoid: true,
		}
	}
	exp := parse_expression_stmt(p)

	return ast.ReturnStatement{
		Value: exp,
	}
}

func parse_break_stmt(p *Parser) ast.Statement {
	p.expect(lexer.BREAK)
	p.expect(lexer.SEMI_COLON)

	return ast.BreakStatement{}
}
