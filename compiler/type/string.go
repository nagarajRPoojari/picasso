package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// String hold i8**, pointer to a global string pointer
type String struct {
	NativeType *types.PointerType // i8*
	Value      value.Value        // slot (alloca of i8*)
	GoVal      string             // optional Go-side constant
}

// IMP: init must be i8*, pointer to a global string constant
func NewString(block *bc.BlockHolder, init value.Value) *String {
	slot := block.V.NewAlloca(types.I8Ptr)
	block.N.NewStore(init, slot)

	var s string
	if c, ok := init.(*constant.CharArray); ok {
		s = constantToGoString(c)
	}
	return &String{
		NativeType: types.I8Ptr,
		Value:      slot,
		GoVal:      s,
	}
}

func constantToGoString(c *constant.CharArray) string {
	bytes := make([]byte, len(c.X))
	for i, ch := range c.X {
		bytes[i] = byte(ch)
	}
	return string(bytes)
}

func (s *String) Update(block *bc.BlockHolder, v value.Value) {
	block.N.NewStore(v, s.Value)
}

// return i8*
func (s *String) Load(block *bc.BlockHolder) value.Value {
	return block.N.NewLoad(types.I8Ptr, s.Value)
}

func (s *String) Constant() constant.Constant {
	return constant.NewNull(types.I8Ptr) // constants handled via global strings
}

func (s *String) Slot() value.Value {
	return s.Value
}

func (s *String) Type() types.Type {
	return s.NativeType
}

func (s *String) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.PointerType:
		if v.Type().Equal(types.I8Ptr) {
			return v, nil
		}
		return block.N.NewBitCast(v, types.I8Ptr), nil
	default:
		return nil, errorsx.NewCompilationError(
			fmt.Sprintf("cannot cast %v to string", v.Type()))
	}
}
func (f *String) NativeTypeString() string { return "string" }
