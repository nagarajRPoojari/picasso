package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir/value"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

type InterfaceH struct {
	Class
	th *TypeHandler
}

func (s *InterfaceH) Update(bh *bc.BlockHolder, v value.Value) {
	block := bh.N

	if v == nil {
		errorutils.Abort(errorutils.InternalError, "cannot update object with nil value")
	}

	ptrType := s.UDT
	if s.Ptr == nil {
		// cases do exist where we need to create Class instance with memory allocated beforehand.
		// example in assignment, where i avoid NewClass to avoid heap allocation, in such cases s.Ptr
		// remains null, so allocate stack slot.
		s.Ptr = block.NewAlloca(ptrType)
	}

	// We simply overwrite the address stored in the slot.
	if v.Type().Equal(s.UDT) {
		block.NewStore(v, s.Ptr)
		return
	}

	// If the user passes a raw struct value, we can't store a struct into a pointer slot directly.
	// I must allocate new heap memory for this value, then point the slot to it.
	if v.Type().Equal(s.UDT.ElemType) {
		block.NewStore(v, s.Ptr)
		return
	}

	errorutils.Abort(errorutils.InternalError,
		fmt.Sprintf("Type mismatch in Update. Expected %s or %s, got %s",
			s.UDT, s.UDT.ElemType, v.Type()))
}
