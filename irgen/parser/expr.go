package parser

import (
	"fmt"
	"strconv"

	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorsx "github.com/nagarajRPoojari/picasso/irgen/error"
	"github.com/nagarajRPoojari/picasso/irgen/lexer"
)

func parseExpr(p *Parser, bp BindingPower) ast.Expression {
	tokenKind := p.currentTokenKind()
	nudFn, exists := nud_table[tokenKind]

	if !exists {
		errorsx.PanicParserError(
			fmt.Sprintf("Unrecognized %s\n", lexer.TokenKindString(tokenKind)),
			p.currentToken().Src.FilePath,
			p.currentToken().Src.Line,
			p.currentToken().Src.Col,
		)
	}

	left := nudFn(p)
	for bp_table[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		ledFn, exists := led_table[tokenKind]

		if !exists {
			errorsx.PanicParserError(fmt.Sprintf("Unrecognized %s\n", lexer.TokenKindString(tokenKind)),
				p.currentToken().Src.FilePath,
				p.currentToken().Src.Line,
				p.currentToken().Src.Col)
		}

		left = ledFn(p, left, bp)
	}

	return left
}

func parsePrefixExpr(p *Parser) ast.Expression {
	operatorToken := p.move()
	expr := parseExpr(p, unary)

	return ast.PrefixExpression{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Operator:  operatorToken,
		Operand:   expr,
	}
}

func parseAssignmentExpr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.move()
	rhs := parseExpr(p, bp)

	return ast.AssignmentExpression{
		SourceLoc:     ast.SourceLoc(p.currentToken().Src),
		Assignee:      left,
		AssignedValue: rhs,
	}
}

func parseRangeExpr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
	p.move()
	return ast.RangeExpression{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Lower:     left,
		Upper:     parseExpr(p, bp),
	}
}

func parseBinaryExpr(p *Parser, left ast.Expression, _ BindingPower) ast.Expression {
	operatorToken := p.move()
	op_bp := bp_table[operatorToken.Kind]

	// For left-associative operators, parse RHS with lower precedence
	right := parseExpr(p, op_bp-1)

	return ast.BinaryExpression{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Left:      left,
		Operator:  operatorToken,
		Right:     right,
	}
}
func parsePrimaryExpr(p *Parser) ast.Expression {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.move().Value, 64)
		return ast.NumberExpression{
			SourceLoc: ast.SourceLoc(p.currentToken().Src),
			Value:     number,
		}
	case lexer.STRING:
		str := p.move().Value
		unescaped, err := strconv.Unquote(str)
		if err != nil {
			panic(fmt.Sprintf("unexpected str format %s", str))
		}
		return ast.StringExpression{
			SourceLoc: ast.SourceLoc(p.currentToken().Src),
			Value:     unescaped,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpression{
			SourceLoc: ast.SourceLoc(p.currentToken().Src),
			Value:     p.move().Value,
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
			SourceLoc: ast.SourceLoc(p.currentToken().Src),
			Member:    left,
			Indices:   rhsList,
		}
	}

	return ast.MemberExpression{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Member:    left,
		Property:  p.expect(lexer.IDENTIFIER).Value,
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
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
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
	return ast.NullExpression{
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
	}
}

func parseCallExpr(p *Parser, left ast.Expression, bp BindingPower) ast.Expression {
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
		SourceLoc: ast.SourceLoc(p.currentToken().Src),
		Method:    left,
		Arguments: arguments,
	}
}

func parseFuncExpr(p *Parser) ast.Expression {
	p.expect(lexer.FN)
	functionParams, returnType, functionBody := parseFnParamsAndBody(p)

	return ast.FunctionExpression{
		SourceLoc:  ast.SourceLoc(p.currentToken().Src),
		Parameters: functionParams,
		ReturnType: returnType,
		Body:       functionBody,
	}
}
