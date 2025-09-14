package typedef

import (
	"testing"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/stretchr/testify/assert"
)

func TestTypeHandler_getPrimitiveVar(t *testing.T) {
	block := ir.NewBlock("")
	tests := []struct {
		name  string
		block *ir.Block
		_type Type
		init  value.Value
		want  Var
	}{
		{
			name:  "expected Boolean set to false",
			block: block,
			_type: Type("boolean"),
			init:  constant.NewInt(types.I1, 0),
			want:  NewBooleanVar(block, false),
		},
		{
			name:  "expected int8",
			block: block,
			_type: Type("int8"),
			init:  constant.NewInt(types.I8, 19),
			want:  NewInt8Var(block, 19),
		},
		{
			name:  "expected int16",
			block: block,
			_type: Type("int16"),
			init:  constant.NewInt(types.I16, 19),
			want:  NewInt16Var(block, 19),
		},
		{
			name:  "expected int32",
			block: block,
			_type: Type("int32"),
			init:  constant.NewInt(types.I32, 19),
			want:  NewInt32Var(block, 19),
		},
		{
			name:  "expected int64",
			block: block,
			_type: Type("int64"),
			init:  constant.NewInt(types.I64, 19),
			want:  NewInt64Var(block, 19),
		},
		{
			name:  "expected float16",
			block: block,
			_type: Type("float16"),
			init:  constant.NewFloat(types.Half, 80.0),
			want:  NewFloat16Var(block, 80.0),
		},
		{
			name:  "expected float32",
			block: block,
			_type: Type("float32"),
			init:  constant.NewFloat(types.Float, 80.0),
			want:  NewFloat32Var(block, 80.0),
		},
		{
			name:  "expected float64",
			block: block,
			_type: Type("float64"),
			init:  constant.NewFloat(types.Double, 80.0),
			want:  NewFloat64Var(block, 80.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ty := NewTypeHandler()
			got := ty.getPrimitiveVar(tt.block, tt._type, tt.init)
			assert.Equal(t, tt.want, got)
		})

	}
}

func TestBuildVar(t *testing.T) {
	f := ir.NewFunc("test", types.Void)
	block := f.NewBlock("entry")

	tests := []struct {
		name      string
		paramType Type
		param     value.Value
		setupUDT  func(*TypeHandler)
		want      Var
	}{
		{
			name:      "int32 primitive",
			paramType: Type("int32"),
			param:     constant.NewInt(types.I32, 42),
			want:      NewInt32Var(block, int32(42)),
		},
		{
			name:      "float64 primitive",
			paramType: Type("float64"),
			param:     constant.NewFloat(types.Double, 3.14),
			want:      NewFloat64Var(block, float64(3.14)),
		},
		{
			name:      "custom pointer type",
			paramType: Type("MyClass"),
			param:     ir.NewParam("p", types.NewStruct(types.I1)),
			setupUDT: func(th *TypeHandler) {
				th.Udts["MyClass"] = &MetaClass{
					UDT: types.NewStruct(types.I1),
				}
			},
			want: NewClass(block, "MyClass", types.NewStruct(types.I1)),
		},
		{
			name:      "custom struct type",
			paramType: Type("Point"),
			param: constant.NewStruct(
				types.NewStruct(types.I32, types.I32),
				constant.NewInt(types.I32, 1), constant.NewInt(types.I32, 2),
			),
			setupUDT: func(th *TypeHandler) {
				th.Udts["Point"] = &MetaClass{
					UDT: types.NewStruct(types.I32, types.I32),
				}
			},
			want: NewClass(block, "Point", types.NewStruct(types.I32, types.I32)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := NewTypeHandler()
			if tt.setupUDT != nil {
				tt.setupUDT(th)
			}

			got := th.BuildVar(block, tt.paramType, tt.param)

			assert.NotNil(t, got)
			assert.IsType(t, got, got, "should return correct Var type")
			assert.Equal(t, tt.want, got)
		})
	}
}
