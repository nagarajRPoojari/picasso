package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	errorsx "github.com/nagarajRPoojari/x-lang/error"

	"github.com/llir/llvm/ir/value"
)

// TypeHolder is a special type that just hold native as string
// e.g, int, Math, string etc.
type TypeHolder struct {
	NativeType *types.PointerType // i8*
	Value      value.Value        // slot (alloca of i8*)
	GoVal      string             // optional Go-side constant
}

// IMP: init must be i8*, pointer to a global string constant
func NewTypeHolder(block *ir.Block, init value.Value) *TypeHolder {
	slot := block.NewAlloca(types.I8Ptr)
	block.NewStore(init, slot)

	var s string
	if c, ok := init.(*constant.CharArray); ok {
		s = constantToGoString(c)
	}
	return &TypeHolder{
		NativeType: types.I8Ptr,
		Value:      slot,
		GoVal:      s,
	}
}
func (s *TypeHolder) Update(block *ir.Block, v value.Value) {
	block.NewStore(v, s.Value)
}

// return i8*
func (s *TypeHolder) Load(block *ir.Block) value.Value {
	return block.NewLoad(types.I8Ptr, s.Value)
}

func (s *TypeHolder) Constant() constant.Constant {
	return constant.NewNull(types.I8Ptr) // constants handled via global TypeHolders
}

func (s *TypeHolder) Slot() value.Value {
	return s.Value
}

func (s *TypeHolder) Type() types.Type {
	return s.NativeType
}

func (s *TypeHolder) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.PointerType:
		if v.Type().Equal(types.I8Ptr) {
			return v, nil
		}
		return block.NewBitCast(v, types.I8Ptr), nil
	default:
		return nil, errorsx.NewCompilationError(
			fmt.Sprintf("cannot cast %v to string", v.Type()))
	}
}
func (f *TypeHolder) NativeTypeString() string { return "string" }
