package expression

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/type/primitives/boolean"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/type/primitives/floats"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/type/primitives/ints"
)

// ProcessMemberExpression generates LLVM IR to access a specific field of a class instance.
// It performs a metadata lookup to find the structural offset of the property,
// calculates the field's memory address, and wraps the resulting pointer into
// a Niyama-specific type container (Var) based on the field's underlying LLVM type.
//
// Technical Logic:
//   - Property Mapping: Uses the IdentifierBuilder to resolve the mangled field name
//     and retrieves its index from the class's FieldIndexMap.
//   - Pointer Arithmetic: Invokes FieldPtr to emit a GetElementPtr (GEP) instruction,
//     targeting the specific index within the class struct.
//   - Type Wrapping: Inspects the LLVM type of the field (Int, Float, Pointer, etc.)
//     and initializes the corresponding high-level wrapper (e.g., ints.Int32 or tf.Class).
//   - Recursive Resolution: For complex fields like nested classes or arrays, it
//     performs the necessary 'load' instructions to return an addressable instance.
func (t *ExpressionHandler) ProcessMemberExpression(bh *bc.BlockHolder, ex ast.MemberExpression) tf.Var {
	// check imported base modules for method resolution
	x, ok := ex.Member.(ast.SymbolExpression)
	if ok {
		if _, ok := t.st.Imports[x.Value]; ok {
			// fName := fmt.Sprintf("%s.%s", t.st.Imports[x.Value].Name, ex.Property)
			v, ok := t.st.CI.Constants[ex.Property]
			if !ok {
				errorutils.Abort(errorutils.UnknownVariable, ex.Property)
			}
			val := bh.N.NewLoad(types.I32, v)
			return t.st.TypeHandler.BuildVar(bh, tf.NewType("int32"), val)
		}
	}

	// Evaluate the base expression
	baseVar := t.ProcessExpression(bh, ex.Member)

	if baseVar == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalMemberExprError, "nil base for member expression")
	}

	// Base must be a class instance
	cls, ok := baseVar.(*tf.Class)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalMemberExprError, "member access base is not a class instance")
	}

	// Get metadata for base class
	classMeta, ok := t.st.Classes[cls.Name]
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalMemberExprError, "unknown class metadata: "+cls.Name)
	}

	// Compute field name in identifier map
	fieldID := fmt.Sprintf("%s.%s", cls.Name, ex.Property)
	idx, ok := classMeta.FieldIndexMap[fieldID]

	if !ok {
		errorutils.Abort(errorutils.UnknownClassField, ex.Property, cls.Name)
	}

	// Get field type from struct UDT
	st := classMeta.StructType()
	fieldType := st.Fields[idx]

	// Get pointer to the field
	fieldPtr := cls.FieldPtr(bh, idx)
	// return t.typeHandler.BuildVar(block, "", fieldPtr)

	// Determine the class name if the field is a struct
	getClassName := func(tt types.Type) string {
		for cname, meta := range t.st.TypeHandler.ClassUDTS {
			if meta.UDT == tt {
				return cname
			}
		}
		for cname, meta := range t.st.TypeHandler.InterfaceUDTS {
			if meta.UDT == tt {
				return cname
			}
		}
		return ""
	}

	// Wrap into appropriate Var
	// @fix: i already have reusable piece of code for this job, need to replace
	// this below code & test.
	// @fix: test atomic var return.
	switch ft := fieldType.(type) {
	case *types.IntType:
		switch ft.BitSize {
		case 1:
			return &boolean.Boolean{NativeType: types.I1, Value: fieldPtr}
		case 8:
			return &ints.Int8{NativeType: types.I8, Value: fieldPtr}
		case 16:
			return &ints.Int16{NativeType: types.I16, Value: fieldPtr}
		case 32:
			return &ints.Int32{NativeType: types.I32, Value: fieldPtr}
		case 64:
			return &ints.Int64{NativeType: types.I64, Value: fieldPtr}
		default:
			errorutils.Abort(errorutils.InternalError, errorutils.InternalTypeError, fmt.Sprintf("unsupported int size %d", ft.BitSize))
		}

	case *types.FloatType:
		switch ft.Kind {
		case types.FloatKindHalf:
			return &floats.Float16{NativeType: types.Half, Value: fieldPtr}
		case types.FloatKindFloat:
			return &floats.Float32{NativeType: types.Float, Value: fieldPtr}
		case types.FloatKindDouble:
			return &floats.Float64{NativeType: types.Double, Value: fieldPtr}
		default:
			errorutils.Abort(errorutils.InternalError, errorutils.InternalTypeError, fmt.Sprintf("unsupported float kind %v", ft.Kind))
		}

	case *types.PointerType:
		if ele, ok := ft.ElemType.(*types.StructType); ok {
			if ele.Name() == constants.ARRAY {
				f := bh.N.NewLoad(types.NewPointer(tf.ARRAYSTRUCT), fieldPtr)
				return &tf.Array{
					Ptr:       f,
					ArrayType: tf.ARRAYSTRUCT,
					ElemType:  classMeta.ArrayVarsEleTypes[idx],
				}
			}

			c := &tf.Class{
				Name: getClassName(fieldType),
				UDT:  ft,
			}
			c.Update(bh, bh.N.NewLoad(fieldType, fieldPtr))
			return c
		} else {
			return tf.NewString(bh, bh.N.NewLoad(types.I8Ptr, fieldPtr))
		}

	default:
		errorutils.Abort(errorutils.InternalError, errorutils.InternalTypeError, fmt.Sprintf("unsupported field type %T in member expression", fieldType))
	}
	return nil
}
