package parser

import (
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/lexer"
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

func ParseAll(path string) ast.BlockStatement {
	tokens := lexer.Tokenize(path)
	p := createParser(tokens)
	body := make([]ast.Statement, 0)

	for p.hasTokens() {
		body = append(body, parseStmt(p))
	}

	return ast.BlockStatement{
		Body: body,
	}
}

func ParseImports(filePath string) ast.BlockStatement {
	tokens := lexer.Tokenize(filePath)
	p := createParser(tokens)
	body := make([]ast.Statement, 0)

	for p.hasTokens() {
		stmt := parseStmt(p)
		// as soon as I see non-import statement then stop parsing
		// @todo: ideally don't accept source string, it could be too big
		// I should read incrementely
		// @todo: don't even need full tokenizing
		if _, ok := stmt.(ast.ImportStatement); !ok {
			break
		}
		body = append(body, stmt)
	}

	return ast.BlockStatement{
		Body: body,
	}
}
