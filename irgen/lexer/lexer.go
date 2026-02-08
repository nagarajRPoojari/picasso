package lexer

import (
	"bufio"
	"os"
	"regexp"

	errorsx "github.com/nagarajRPoojari/picasso/irgen/error"
)

const readChunkSize = 4096

type RegexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []RegexPattern
	Tokens   []Token

	reader   *bufio.Reader
	buffer   []byte
	eof      bool
	filePath string

	line int
	col  int
}

type regexHandler func(lex *lexer, regex *regexp.Regexp)

func Tokenize(path string) []Token {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lex := newLexer(path, bufio.NewReader(f))

	for {
		if lex.eof && len(lex.buffer) == 0 {
			break
		}

		lex.fillBuffer()

		matched := false
		for _, pattern := range lex.patterns {
			loc := pattern.regex.FindIndex(lex.buffer)
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			errorsx.PanicLexerError(
				"lexer error: unrecognized token",
				lex.filePath,
				lex.line,
				lex.col,
			)
		}
	}

	return lex.Tokens
}

func newLexer(path string, reader *bufio.Reader) *lexer {
	return &lexer{
		reader:   reader,
		filePath: path,
		line:     1,
		col:      1,
		buffer:   make([]byte, 0),
		Tokens:   make([]Token, 0),
		patterns: []RegexPattern{
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\/\/.*`), commentHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},

			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},

			{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},

			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},

			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},

			// bitwise op should be checked after logical op
			{regexp.MustCompile(`\|`), defaultHandler(BITWISE_OR, "|")},
			{regexp.MustCompile(`&`), defaultHandler(BITWISE_AND, "&")},
			{regexp.MustCompile(`\^`), defaultHandler(BITWISE_XOR, "^")},
			{regexp.MustCompile(`\~`), defaultHandler(BITWISE_NOT, "~")},

			{regexp.MustCompile(`\.\.`), defaultHandler(DOT_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},

			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},

			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(DASH, "-")},
			{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
			{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
		},
	}
}

/* =========================
   Buffer management
   ========================= */

func (l *lexer) fillBuffer() {
	if l.eof || len(l.buffer) >= readChunkSize {
		return
	}

	tmp := make([]byte, readChunkSize)
	n, err := l.reader.Read(tmp)
	if err != nil {
		l.eof = true
	}

	l.buffer = append(l.buffer, tmp[:n]...)
}

func (l *lexer) advance(n int) {
	for i := 0; i < n; i++ {
		if l.buffer[i] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
	}
	l.buffer = l.buffer[n:]
}

func (l *lexer) srcLoc() SourceLoc {
	return SourceLoc{
		FilePath: l.filePath,
		Line:     l.line,
		Col:      l.col,
	}
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, _ *regexp.Regexp) {
		lex.Tokens = append(
			lex.Tokens,
			newUniqueToken(kind, value, lex.srcLoc()),
		)
		lex.advance(len(value))
	}
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	loc := regex.FindIndex(lex.buffer)
	lit := string(lex.buffer[loc[0]:loc[1]])

	lex.Tokens = append(
		lex.Tokens,
		newUniqueToken(STRING, lit, lex.srcLoc()),
	)
	lex.advance(len(lit))
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.Find(lex.buffer)
	lex.Tokens = append(
		lex.Tokens,
		newUniqueToken(NUMBER, string(match), lex.srcLoc()),
	)
	lex.advance(len(match))
}

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	match := string(regex.Find(lex.buffer))

	if kind, ok := reserved_keywords[match]; ok {
		lex.Tokens = append(
			lex.Tokens,
			newUniqueToken(kind, match, lex.srcLoc()),
		)
	} else {
		lex.Tokens = append(
			lex.Tokens,
			newUniqueToken(IDENTIFIER, match, lex.srcLoc()),
		)
	}
	lex.advance(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	loc := regex.FindIndex(lex.buffer)
	lex.advance(loc[1])
}

func commentHandler(lex *lexer, regex *regexp.Regexp) {
	loc := regex.FindIndex(lex.buffer)
	lex.advance(loc[1])
}

func preview(buf []byte) string {
	if len(buf) > 20 {
		return string(buf[:20]) + "..."
	}
	return string(buf)
}
