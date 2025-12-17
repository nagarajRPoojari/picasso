package parser

import (
	"fmt"

	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/lexer"
)

type typeNudHandler func(p *Parser) ast.Type
type typeLedHandler func(p *Parser, left ast.Type, bp BindingPower) ast.Type

type typeNudLookUp map[lexer.TokenKind]typeNudHandler
type typeLedLookUp map[lexer.TokenKind]typeLedHandler
type typeBpLookUp map[lexer.TokenKind]BindingPower

var typeByTable = typeBpLookUp{}
var typeNudTable = typeNudLookUp{}
var typeLedTable = typeLedLookUp{}

func typeLed(kind lexer.TokenKind, bp BindingPower, led_fn typeLedHandler) {
	typeByTable[kind] = bp
	typeLedTable[kind] = led_fn
}

func typeNud(kind lexer.TokenKind, bp BindingPower, nud_fn typeNudHandler) {
	typeByTable[kind] = primary
	typeNudTable[kind] = nud_fn
}

func BuildTypeTokensTable() {

	typeNud(lexer.IDENTIFIER, primary, func(p *Parser) ast.Type {
		return &ast.SymbolType{
			Value: p.move().Value,
		}
	})

	typeNud(lexer.OPEN_BRACKET, member, func(p *Parser) ast.Type {
		p.move()
		// token := p.currentToken()
		// var size int
		// if token.Kind != lexer.NUMBER {
		// 	panic("expected size of array")
		// } else {
		// 	num, err := strconv.Atoi(token.Value)
		// 	if err != nil {
		// 		panic("unable to parse size of array")
		// 	}
		// 	size = num
		// 	p.move()
		// }
		p.expect(lexer.CLOSE_BRACKET)
		insideType := parse_type(p, default_bp)

		return &ast.ListType{
			Underlying: insideType,
		}
	})
}

func parse_type(p *Parser, bp BindingPower) ast.Type {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := typeNudTable[tokenKind]

	if !exists {
		panic(fmt.Sprintf("type: NUD Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
	}

	left := nud_fn(p)

	for typeByTable[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		led_fn, exists := typeLedTable[tokenKind]

		if !exists {
			panic(fmt.Sprintf("type: LED Handler expected for token %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = led_fn(p, left, bp)
	}

	return left
}
