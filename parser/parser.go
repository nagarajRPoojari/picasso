package parser

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokens []lexer.Token) *Parser {
	BuildTokensTable()
	BuildTypeTokensTable()

	p := &Parser{
		tokens: tokens,
		pos:    0,
	}

	return p
}

func Parse(source string) ast.BlockStatement {
	tokens := lexer.Tokenize(source)
	p := createParser(tokens)
	body := make([]ast.Statement, 0)

	for p.hasTokens() {
		body = append(body, parse_stmt(p))
	}

	return ast.BlockStatement{
		Body: body,
	}
}
