package expression

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/boolean"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/floats"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/ints"
)

// ProcessMemberExpression evaluates a member access expression (e.g., obj.field)
// and returns a typed runtime variable corresponding to the accessed field.
//
// Behavior:
//   - Recursively evaluates the base expression (`ex.Member`) to get the object instance.
//   - Validates that the base is a class instance (`*tf.Class`).
//   - Resolves the field index and type from class metadata.
//   - Obtains a pointer to the field and wraps it in the appropriate tf.Var type:
//   - Integer types (Int8, Int16, Int32, Int64, Boolean)
//   - Floating-point types (Float16, Float32, Float64)
//   - Struct pointers → wrapped as a new tf.Class instance
//   - Other pointer types → wrapped as a tf.String
//
// Parameters:
//
//	block - the current IR block for code generation
//	ex    - the AST member expression node
//
// Returns:
//
//	tf.Var - a runtime variable representing the field
func (t *ExpressionHandler) ProcessMemberExpression(bh *bc.BlockHolder, ex ast.MemberExpression) tf.Var {
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
	fieldID := t.st.IdentifierBuilder.Attach(cls.Name, ex.Property)
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
