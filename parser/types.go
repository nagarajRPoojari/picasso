package parser

import (
	"fmt"
	"strconv"

	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/lexer"
)

type type_nud_handler func(p *Parser) ast.Type
type type_led_handler func(p *Parser, left ast.Type, bp BindingPower) ast.Type

type type_nud_lookup map[lexer.TokenKind]type_nud_handler
type type_led_lookup map[lexer.TokenKind]type_led_handler
type type_bp_lookup map[lexer.TokenKind]BindingPower

var type_by_table = type_bp_lookup{}
var type_nud_table = type_nud_lookup{}
var type_led_table = type_led_lookup{}

func type_led(kind lexer.TokenKind, bp BindingPower, led_fn type_led_handler) {
	type_by_table[kind] = bp
	type_led_table[kind] = led_fn
}

func type_nud(kind lexer.TokenKind, bp BindingPower, nud_fn type_nud_handler) {
	type_by_table[kind] = primary
	type_nud_table[kind] = nud_fn
}

func BuildTypeTokensTable() {

	type_nud(lexer.IDENTIFIER, primary, func(p *Parser) ast.Type {
		return ast.SymbolType{
			Value: p.move().Value,
		}
	})

	type_nud(lexer.OPEN_BRACKET, member, func(p *Parser) ast.Type {
		p.move()
		token := p.currentToken()
		var size int
		if token.Kind != lexer.NUMBER {
			panic("expected size of array")
		} else {
			num, err := strconv.Atoi(token.Value)
			if err != nil {
				panic("unable to parse size of array")
			}
			size = num
			p.move()
		}
		p.expect(lexer.CLOSE_BRACKET)
		insideType := parse_type(p, default_bp)

		return ast.ListType{
			Underlying: insideType,
			Length:     size,
		}
	})
}

func parse_type(p *Parser, bp BindingPower) ast.Type {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_table[tokenKind]

	if !exists {
		panic(fmt.Sprintf("type: NUD Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)

	for type_by_table[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := type_led_table[tokenKind]

		if !exists {
			panic(fmt.Sprintf("type: LED Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left
}
