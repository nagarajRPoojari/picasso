package error

import "fmt"

type CompileTimeError string
type LexerError string
type ParserError string

func RaiseCompileError(msg string, args ...any) {
	m := "CompilerError: " + fmt.Sprintf(msg, args...)
	panic(m)
}

func RaiseLexerError(msg string, args ...any) {
	m := "LexerError: " + fmt.Sprintf(msg, args...)
	panic(m)
}

func RaiseParserError(msg string, args ...any) {
	m := "ParserError: " + fmt.Sprintf(msg, args...)
	panic(m)
}
