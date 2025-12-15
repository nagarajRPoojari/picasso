package parser

import (
	"fmt"

	errorsx "github.com/nagarajRPoojari/niyama/irgen/error"
	"github.com/nagarajRPoojari/niyama/irgen/lexer"
)

func (t *Parser) currentToken() lexer.Token {
	return t.tokens[t.pos]
}

func (t *Parser) move() lexer.Token {
	tk := t.currentToken()
	t.pos++
	return tk
}

func (t *Parser) hasTokens() bool {
	return t.pos < len(t.tokens) && t.currentTokenKind() != lexer.EOF
}

func (t *Parser) nextToken() lexer.Token {
	return t.tokens[t.pos+1]
}

func (t *Parser) previousToken() lexer.Token {
	return t.tokens[t.pos-1]
}

func (t *Parser) currentTokenKind() lexer.TokenKind {
	return t.tokens[t.pos].Kind
}

func (t *Parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	token := t.currentToken()
	kind := token.Kind

	if kind != expectedKind {
		if err == nil {
			errString := fmt.Sprintf("Expected %s but recieved %s instead\n", lexer.TokenKindString(expectedKind), lexer.TokenKindString(kind))
			errorsx.PanicParserError(errString)
		}
	}

	return t.move()
}

// expect will consume current token
func (t *Parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return t.expectError(expectedKind, nil)
}
