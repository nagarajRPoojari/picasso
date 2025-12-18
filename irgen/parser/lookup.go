package parser

import (
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/lexer"
)

type BindingPower int

const (
	default_bp BindingPower = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type statement_handler func(p *Parser) ast.Statement
type nudHandler func(p *Parser) ast.Expression
type ledHandler func(p *Parser, left ast.Expression, bp BindingPower) ast.Expression

type stmtLookUp map[lexer.TokenKind]statement_handler
type nudLookUp map[lexer.TokenKind]nudHandler
type ledLookUp map[lexer.TokenKind]ledHandler

type bp_lookup map[lexer.TokenKind]BindingPower

var bp_table = bp_lookup{}
var nud_table = nudLookUp{}
var led_table = ledLookUp{}
var statement_table = stmtLookUp{}

func led(kind lexer.TokenKind, bp BindingPower, led_fn ledHandler) {
	bp_table[kind] = bp
	led_table[kind] = led_fn
}

func nud(kind lexer.TokenKind, nud_fn nudHandler) {
	bp_table[kind] = primary
	nud_table[kind] = nud_fn
}

func statement(kind lexer.TokenKind, stmt_fn statement_handler) {
	bp_table[kind] = default_bp
	statement_table[kind] = stmt_fn
}

func BuildTokensTable() {
	// Assignment
	led(lexer.ASSIGNMENT, assignment, parseAssignmentExpr)
	led(lexer.PLUS_EQUALS, assignment, parseAssignmentExpr)
	led(lexer.MINUS_EQUALS, assignment, parseAssignmentExpr)

	// Logical
	led(lexer.AND, logical, parseBinaryExpr)
	led(lexer.OR, logical, parseBinaryExpr)
	led(lexer.DOT_DOT, logical, parseRangeExpr)

	// Relational
	led(lexer.LESS, relational, parseBinaryExpr)
	led(lexer.LESS_EQUALS, relational, parseBinaryExpr)
	led(lexer.GREATER, relational, parseBinaryExpr)
	led(lexer.GREATER_EQUALS, relational, parseBinaryExpr)
	led(lexer.EQUALS, relational, parseBinaryExpr)
	led(lexer.NOT_EQUALS, relational, parseBinaryExpr)

	// Additive & Multiplicitave
	led(lexer.PLUS, additive, parseBinaryExpr)
	led(lexer.DASH, additive, parseBinaryExpr)
	led(lexer.SLASH, multiplicative, parseBinaryExpr)
	led(lexer.STAR, multiplicative, parseBinaryExpr)
	led(lexer.PERCENT, multiplicative, parseBinaryExpr)

	// Literals & Symbols
	nud(lexer.NUMBER, parsePrimaryExpr)
	nud(lexer.STRING, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, parsePrimaryExpr)

	// Unary/Prefix
	nud(lexer.TYPEOF, parsePrefixExpr)
	nud(lexer.DASH, parsePrefixExpr)
	nud(lexer.NOT, parsePrefixExpr)
	nud(lexer.OPEN_BRACKET, parseArrayLiteralExpr)

	// Member / Computed // Call
	led(lexer.DOT, member, parseMemberExpr)
	led(lexer.OPEN_BRACKET, member, parseMemberExpr)
	led(lexer.OPEN_PAREN, call, parseCallExpr)

	nud(lexer.NULL, parseNullExpr)
	// Grouping Expr
	nud(lexer.OPEN_PAREN, parseGroupingExpr)
	nud(lexer.FN, parseFuncExpr)
	nud(lexer.NEW, func(p *Parser) ast.Expression {
		p.move()
		classInstantiation := parseExpr(p, default_bp)

		return ast.NewExpression{
			Instantiation: ast.ExpectExpr[ast.CallExpression](classInstantiation),
		}
	})

	statement(lexer.OPEN_CURLY, parseBlockStmt)
	statement(lexer.SAY, parseVarDeclStmt)
	statement(lexer.CONST, parseVarDeclStmt)
	statement(lexer.FN, parseFuncDeclaration)
	statement(lexer.IF, parseIfStmt)
	statement(lexer.USING, parseImportStmt)
	statement(lexer.FOREACH, parseForeachStmt)
	statement(lexer.WHILE, parseWhileStmt)
	statement(lexer.CLASS, parseClassDeclStmt)
	statement(lexer.INTERFACE, parseInterfaceDeclStmt)
	statement(lexer.RETURN, parseFuncReturnStmt)
	statement(lexer.BREAK, parseBreakStmt)
}
