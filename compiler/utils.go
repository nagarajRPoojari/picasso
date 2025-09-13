package compiler

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (t *LLVM) ensurePrintf() *ir.Func {
	for _, f := range t.module.Funcs {
		if f.Name() == "printf" {
			return f
		}
	}
	printf := t.module.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
	printf.Sig.Variadic = true
	return printf
}

func (t *LLVM) print(block *ir.Block, format string, args ...interface{}) {
	printf := t.ensurePrintf()

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if !strings.HasSuffix(format, "\x00") {
		format += "\x00"
	}

	strConst := constant.NewCharArrayFromString(format)
	t.strCounter++
	global := t.module.NewGlobalDef(fmt.Sprintf(".str%d", t.strCounter), strConst)
	global.Linkage = enum.LinkagePrivate

	zero := constant.NewInt(types.I64, 0)
	ptr := block.NewGetElementPtr(strConst.Typ, global, zero, zero)

	var callArgs []value.Value
	callArgs = append(callArgs, ptr)

	for _, a := range args {
		switch v := a.(type) {
		case value.Value:
			// Already an LLVM value (float, int, etc.)
			callArgs = append(callArgs, v)

		case string:
			// Turn Go string into global char array + pointer
			str := v + "\x00"
			strConst := constant.NewCharArrayFromString(str)

			t.strCounter++
			g := t.module.NewGlobalDef(fmt.Sprintf(".str%d", t.strCounter), strConst)
			g.Linkage = enum.LinkagePrivate

			zero := constant.NewInt(types.I64, 0)
			p := block.NewGetElementPtr(strConst.Typ, g, zero, zero)
			callArgs = append(callArgs, p)

		case int:
			callArgs = append(callArgs, constant.NewInt(types.I32, int64(v)))

		case float64:
			callArgs = append(callArgs, constant.NewFloat(types.Double, v))

		case float32:
			callArgs = append(callArgs, constant.NewFloat(types.Float, float64(v)))

		default:
			panic(fmt.Sprintf("unsupported print arg type: %T", v))
		}
	}

	block.NewCall(printf, callArgs...)
}

type BridgeIdentifier struct {
	val string
}

func NewBridgeIdentifier() *BridgeIdentifier {
	return &BridgeIdentifier{val: ""}
}

func (t *BridgeIdentifier) Attach(name string) *BridgeIdentifier {
	t.val += "." + name
	return t
}

func (t *BridgeIdentifier) Ret() string {
	return t.val
}
