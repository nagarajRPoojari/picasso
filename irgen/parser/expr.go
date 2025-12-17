package parser

import (
	"fmt"
	"strconv"

	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorsx "github.com/nagarajRPoojari/niyama/irgen/error"
	"github.com/nagarajRPoojari/niyama/irgen/lexer"
)

func parseExpr(p *Parser, bp BindingPower) ast.Expression {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_table[tokenKind]

	if !exists {
		errorsx.PanicParserError(fmt.Sprintf("NUD Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)
	for bp_table[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := led_table[tokenKind]

		if !exists {
			errorsx.PanicParserError(fmt.Sprintf("LED Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left
}

func parsePrefixExpr(p *Parser) ast.Expression {
	operatorToken := p.move()
	expr := parseExpr(p, unary)

	return ast.PrefixExpression{
		Operator: operatorToken,
		Operand:  expr,
	}
}

func parseAssignmentExpr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.move()
	rhs := parseExpr(p, bp)

	return ast.AssignmentExpression{
		Assignee:      left,
		AssignedValue: rhs,
	}
}

func parseRangeExpr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.move()
	return ast.RangeExpression{
		Lower: left,
		Upper: parseExpr(p, bp),
	}
}

func parseBinaryExpr(p *Parser, left ast.Expression, _ BindingPower) ast.Expression {
	operatorToken := p.move()
	op_bp := bp_table[operatorToken.Kind]

	// For left-associative operators, parse RHS with lower precedence
	right := parseExpr(p, op_bp-1)

	return ast.BinaryExpression{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}
func parsePrimaryExpr(p *Parser) ast.Expression {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.move().Value, 64)
		return ast.NumberExpression{
			Value: number,
		}
	case lexer.STRING:
		str := p.move().Value
		unescaped, err := strconv.Unquote(str)
		if err != nil {
			panic(fmt.Sprintf("unexpected str format %s", str))
		}
		return ast.StringExpression{
			Value: unescaped,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpression{
			Value: p.move().Value,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary_expr from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parseMemberExpr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	isComputed := p.move().Kind == lexer.OPEN_BRACKET

	if isComputed {
		rhsList := make([]ast.Expression, 0)
		rhsList = append(rhsList, parseExpr(p, bp))
		for {
			if p.currentTokenKind() == lexer.CLOSE_BRACKET {
				break
			}
			p.expect(lexer.COMMA)
			rhsList = append(rhsList, parseExpr(p, bp))
		}
		p.expect(lexer.CLOSE_BRACKET)
		return ast.ComputedExpression{
			Member:  left,
			Indices: rhsList,
		}
	}

	return ast.MemberExpression{
		Member:   left,
		Property: p.expect(lexer.IDENTIFIER).Value,
	}
}

func parseArrayLiteralExpr(p *Parser) ast.Expression {
	p.expect(lexer.OPEN_BRACKET)
	arrayContents := make([]ast.Expression, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		arrayContents = append(arrayContents, parseExpr(p, logical))

		if !p.currentToken().IsOneOfMany(lexer.EOF, lexer.CLOSE_BRACKET) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_BRACKET)

	return ast.ListExpression{
		Constants: arrayContents,
	}
}

func parseGroupingExpr(p *Parser) ast.Expression {
	p.expect(lexer.OPEN_PAREN)
	expr := parseExpr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	return expr
}

func parseNullExpr(p *Parser) ast.Expression {
	p.expect(lexer.NULL)
	return ast.NullExpression{}
}

func parse_call_expr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.move()
	arguments := make([]ast.Expression, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		arguments = append(arguments, parseExpr(p, assignment))

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

func parseFuncExpr(p *Parser) ast.Expression {
	p.expect(lexer.FN)
	functionParams, returnType, functionBody := parseFnParamsAndBody(p)

	return ast.FunctionExpression{
		Parameters: functionParams,
		ReturnType: returnType,
		Body:       functionBody,
	}
}
