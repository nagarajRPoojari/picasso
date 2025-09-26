package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/boolean"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/floats"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/ints"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *ExpressionHandler) ProcessMemberExpression(block *ir.Block, ex ast.MemberExpression) tf.Var {
	// Evaluate the base expression
	baseVar, safe := t.ProcessExpression(block, ex.Member)
	block = safe

	if baseVar == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("nil base in member expression: %v", ex.Member))
	}

	// Base must be a class instance
	cls, ok := baseVar.(*tf.Class)
	if !ok {
		errorsx.PanicCompilationError(fmt.Sprintf("member access base is not a class instance, got %T, while, %v", baseVar, ex))
	}

	// Get metadata for base class
	classMeta, ok := t.st.Classes[cls.Name]
	if !ok {
		errorsx.PanicCompilationError(fmt.Sprintf("unknown class metadata: %s", cls.Name))
	}

	// Compute field name in identifier map
	fieldID := t.st.IdentifierBuilder.Attach(cls.Name, ex.Property)
	idx, ok := classMeta.VarIndexMap[fieldID]
	if !ok {
		errorsx.PanicCompilationError(fmt.Sprintf("unknown field %s on class %s", ex.Property, cls.Name))
	}

	// Get field type from struct UDT
	st := classMeta.StructType()
	fieldType := st.Fields[idx]

	// Get pointer to the field
	fieldPtr := cls.FieldPtr(block, idx)
	// return t.typeHandler.BuildVar(block, "", fieldPtr)

	// Determine the class name if the field is a struct
	getClassName := func(tt types.Type) string {
		for cname, meta := range t.st.Classes {
			if meta.UDT == tt {
				return cname
			}
		}
		for cname, meta := range t.st.TypeHandler.Udts {
			if meta.UDT == tt {
				return cname
			}
		}
		return ""
	}

	// Wrap into appropriate Var
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
			panic(fmt.Sprintf("unsupported int size %d", ft.BitSize))
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
			panic(fmt.Sprintf("unsupported float kind %v", ft.Kind))
		}

	case *types.PointerType:
		if _, ok := ft.ElemType.(*types.StructType); ok {
			c := tf.NewClass(
				block, getClassName(fieldType), ft,
			)
			c.Update(block, block.NewLoad(fieldType, fieldPtr))
			return c
		} else {
			return tf.NewString(block, block.NewLoad(types.I8Ptr, fieldPtr))
		}

	default:
		errorsx.PanicCompilationError(fmt.Sprintf("unsupported field type %T in member expression", fieldType))
	}
	return nil
}
