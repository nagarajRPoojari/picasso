package parser

import (
	"fmt"
	"strconv"

	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/lexer"
)

func parse_expr(p *Parser, bp BindingPower) ast.Expression {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_table[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)
	for bp_table[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_table[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left
}

func parse_prefix_expr(p *Parser) ast.Expression {
	operatorToken := p.move()
	expr := parse_expr(p, unary)

	return ast.PrefixExpression{
		Operator: operatorToken,
		Operand:  expr,
	}
}

func parse_assignment_expr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.move()
	rhs := parse_expr(p, bp)

	return ast.AssignmentExpression{
		Assignee:      left,
		AssignedValue: rhs,
	}
}

func parse_range_expr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.move()
	return ast.RangeExpression{
		Lower: left,
		Upper: parse_expr(p, bp),
	}
}

func parse_binary_expr(p *Parser, left ast.Expression, _ BindingPower) ast.Expression {
	operatorToken := p.move()
	op_bp := bp_table[operatorToken.Kind]

	// For left-associative operators, parse RHS with lower precedence
	right := parse_expr(p, op_bp-1)

	return ast.BinaryExpression{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}
func parse_primary_expr(p *Parser) ast.Expression {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.move().Value, 64)
		return ast.NumberExpression{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringExpression{
			Value: p.move().Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpression{
			Value: p.move().Value,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary_expr from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parse_member_expr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	isComputed := p.move().Kind == lexer.OPEN_BRACKET

	if isComputed {
		rhs := parse_expr(p, bp)
		p.expect(lexer.CLOSE_BRACKET)
		return ast.ComputedExpression{
			Member:   left,
			Property: rhs,
		}
	}

	return ast.MemberExpression{
		Member:   left,
		Property: p.expect(lexer.IDENTIFIER).Value,
	}
}

func parse_array_literal_expr(p *Parser) ast.Expression {
	p.expect(lexer.OPEN_BRACKET)
	arrayContents := make([]ast.Expression, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		arrayContents = append(arrayContents, parse_expr(p, logical))

		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_BRACKET) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_BRACKET)

	return ast.ListExpression{
		Constants: arrayContents,
	}
}

func parse_grouping_expr(p *Parser) ast.Expression {
	p.expect(lexer.OPEN_PAREN)
	expr := parse_expr(p, default_bp)
	p.expect(lexer.OPEN_PAREN)
	return expr
}

func parse_call_expr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.move()
	arguments := make([]ast.Expression, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		arguments = append(arguments, parse_expr(p, assignment))

		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_PAREN) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PAREN)
	return ast.CallExpression{
		Method:    left,
		Arguments: arguments,
	}
}

func parse_fn_expr(p *Parser) ast.Expression {
	p.expect(lexer.FN)
	functionParams, returnType, functionBody := parse_fn_params_and_body(p)

	return ast.FunctionExpression{
		Parameters: functionParams,
		ReturnType: returnType,
		Body:       functionBody,
	}
}
