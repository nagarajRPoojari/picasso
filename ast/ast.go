package ast

import "github.com/nagarajRPoojari/x-lang/utils"

type Statement interface {
	stmt()
}

type Expression interface {
	expr()
}

type Type interface {
	Get() string
}

func ExpectExpr[T Expression](expr Expression) T {
	return utils.ExpectType[T](expr)
}

func ExpectStmt[T Statement](expr Statement) T {
	return utils.ExpectType[T](expr)
}
