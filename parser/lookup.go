package parser

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/lexer"
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
type nud_handler func(p *Parser) ast.Expression
type led_handler func(p *Parser, left ast.Expression, bp BindingPower) ast.Expression

type statement_lookup map[lexer.TokenKind]statement_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler

type bp_lookup map[lexer.TokenKind]BindingPower

var bp_table = bp_lookup{}
var nud_table = nud_lookup{}
var led_table = led_lookup{}
var statement_table = statement_lookup{}

func led(kind lexer.TokenKind, bp BindingPower, led_fn led_handler) {
	bp_table[kind] = bp
	led_table[kind] = led_fn
}

func nud(kind lexer.TokenKind, nud_fn nud_handler) {
	bp_table[kind] = primary
	nud_table[kind] = nud_fn
}

func statement(kind lexer.TokenKind, stmt_fn statement_handler) {
	bp_table[kind] = default_bp
	statement_table[kind] = stmt_fn
}

func BuildTokensTable() {
	// Assignment
	led(lexer.ASSIGNMENT, assignment, parse_assignment_expr)
	led(lexer.PLUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.MINUS_EQUALS, assignment, parse_assignment_expr)

	// Logical
	led(lexer.AND, logical, parse_binary_expr)
	led(lexer.OR, logical, parse_binary_expr)
	led(lexer.DOT_DOT, logical, parse_range_expr)

	// Relational
	led(lexer.LESS, relational, parse_binary_expr)
	led(lexer.LESS_EQUALS, relational, parse_binary_expr)
	led(lexer.GREATER, relational, parse_binary_expr)
	led(lexer.GREATER_EQUALS, relational, parse_binary_expr)
	led(lexer.EQUALS, relational, parse_binary_expr)
	led(lexer.NOT_EQUALS, relational, parse_binary_expr)

	// Additive & Multiplicitave
	led(lexer.PLUS, additive, parse_binary_expr)
	led(lexer.DASH, additive, parse_binary_expr)
	led(lexer.SLASH, multiplicative, parse_binary_expr)
	led(lexer.STAR, multiplicative, parse_binary_expr)
	led(lexer.PERCENT, multiplicative, parse_binary_expr)

	// Literals & Symbols
	nud(lexer.NUMBER, parse_primary_expr)
	nud(lexer.STRING, parse_primary_expr)
	nud(lexer.IDENTIFIER, parse_primary_expr)

	// Unary/Prefix
	nud(lexer.TYPEOF, parse_prefix_expr)
	nud(lexer.DASH, parse_prefix_expr)
	nud(lexer.NOT, parse_prefix_expr)
	nud(lexer.OPEN_BRACKET, parse_array_literal_expr)

	// Member / Computed // Call
	led(lexer.DOT, member, parse_member_expr)
	led(lexer.OPEN_BRACKET, member, parse_member_expr)
	led(lexer.OPEN_PAREN, call, parse_call_expr)

	nud(lexer.NULL, parse_null_expr)
	// Grouping Expr
	nud(lexer.OPEN_PAREN, parse_grouping_expr)
	nud(lexer.FN, parse_fn_expr)
	nud(lexer.NEW, func(p *Parser) ast.Expression {
		p.move()
		classInstantiation := parse_expr(p, default_bp)

		return ast.NewExpression{
			Instantiation: ast.ExpectExpr[ast.CallExpression](classInstantiation),
		}
	})

	statement(lexer.OPEN_CURLY, parse_block_stmt)
	statement(lexer.SAY, parse_var_decl_stmt)
	statement(lexer.CONST, parse_var_decl_stmt)
	statement(lexer.FN, parse_fn_declaration)
	statement(lexer.IF, parse_if_stmt)
	statement(lexer.IMPORT, parse_import_stmt)
	statement(lexer.FOREACH, parse_foreach_stmt)
	statement(lexer.CLASS, parse_class_declaration_stmt)
	statement(lexer.RETURN, parse_function_return_stmt)
}
