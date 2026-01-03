package class

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	typedef "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
)

// DefineClassFuncs triggers the emission of concrete LLVM IR function bodies for
// all methods belonging to a class. It handles the traversal of local
// definitions using a tracking map to ensure that overridden methods are
// defined only once using the most specific implementation.
func (t *ClassHandler) DefineClassFuncs(cls ast.ClassDeclarationStatement) {
	avoid := make(map[string]struct{}, 0)
	fqClsName := t.st.IdentifierBuilder.Attach(cls.Name)
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			t.m.GetFuncHandler().(*funcs.FuncHandler).DefineFunc(fqClsName, &st, avoid)
			avoid[st.Name] = struct{}{}
		}
	}
}

// DefineClass resolves a previously declared opaque struct into a concrete
// LLVM struct layout. It maps class members (fields and methods) to numerical
// indices for GEP (GetElementPtr) instructions and populates the struct's
// memory footprint.
//
// Key Logic:
//   - Interface Validation: If the class implements an interface, it verifies
//     that all required methods exist and their signatures (hashed) match
//     the interface definition.
//   - Member Mapping: Iterates through the class body to assign indices to
//     variables and function pointers, storing clsMetadata in the ClassclsMeta (clsMeta).
//   - Type Resolution: Converts AST types into concrete LLVM types for each
//     struct field.
//   - Opaque Completion: Updates the underlying LLVM struct type (stored in
//     clsMeta.UDT) with the finalized field list, completing the type definition.
func (t *ClassHandler) DefineClass(cls ast.ClassDeclarationStatement, sourcePkg state.PackageEntry) {
	aliasClsName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(cls.Name)
	clsMeta := t.st.Classes[aliasClsName]

	fieldTypes := make([]types.Type, 0)
	vars := make(map[string]struct{}, 0)
	funcs := make(map[string]struct{}, 0)

	// map each fields with corresponding udt struct index
	i := 0

	// Opaque resolution: define concrete types of all fields
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			t.defineField(i, aliasClsName, clsMeta, &fieldTypes, vars, st)
			i++

		case ast.FunctionDefinitionStatement:
			t.defineMethod(i, aliasClsName, clsMeta, &fieldTypes, funcs, st)
			i++
		}
	}

	ptr, ok := clsMeta.UDT.(*types.PointerType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalUDTDefinitionError, "udt must be a pointer")
	}
	st, ok := ptr.ElemType.(*types.StructType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalUDTDefinitionError, "udt must be pointer to a struct")
	}
	st.Fields = fieldTypes
}

// defineField registers all needed info about class field in class metadata
func (t *ClassHandler) defineField(i int, aliasClsName string, clsMeta *typedef.MetaClass, fieldTypes *[]types.Type, vars map[string]struct{}, st ast.VariableDeclarationStatement) {
	aliasVarName := fmt.Sprintf("%s.%s", aliasClsName, st.Identifier)
	if _, ok := vars[aliasVarName]; ok {
		errorutils.Abort(errorutils.VariableRedeclaration, st.Identifier)
	}

	clsMeta.FieldIndexMap[aliasVarName] = i
	clsMeta.VarAST[aliasVarName] = &st

	*fieldTypes = append(*fieldTypes, t.st.TypeHandler.GetLLVMType(st.ExplicitType.Get()))
	vars[aliasVarName] = struct{}{}

	if st.ExplicitType.GetUnderlyingType() != "" {
		clsMeta.ArrayVarsEleTypes[i] = t.st.TypeHandler.GetLLVMType(st.ExplicitType.GetUnderlyingType())
	}

	// mark access mode
	if st.IsInternal {
		clsMeta.InternalFields[aliasVarName] = struct{}{}
	}
}

// defineField registers all needed info about class method in class metadata
func (t *ClassHandler) defineMethod(i int, aliasClsName string, clsMeta *typedef.MetaClass, fieldTypes *[]types.Type, funcs map[string]struct{}, st ast.FunctionDefinitionStatement) {
	aliasFuncName := fmt.Sprintf("%s.%s", aliasClsName, st.Name)
	if _, ok := funcs[aliasFuncName]; ok {
		errorutils.Abort(errorutils.MethodRedeclaration, st.Name)
	}

	var retType types.Type
	if st.ReturnType != nil {
		retType = t.st.TypeHandler.GetLLVMType(st.ReturnType.Get())
	} else {
		// empty string is expected to give a void type
		retType = t.st.TypeHandler.GetLLVMType("")
	}

	args := make([]types.Type, 0)
	for _, p := range st.Parameters {
		args = append(args, t.st.TypeHandler.GetLLVMType(p.Type.Get()))
	}

	// args = append(args, clsMeta.UDT)
	funcType := types.NewFunc(retType, args...)
	*fieldTypes = append(*fieldTypes, types.NewPointer(funcType))

	clsMeta.FieldIndexMap[aliasFuncName] = i
	funcs[aliasFuncName] = struct{}{}

	// mark access mode
	if st.IsInternal {
		clsMeta.InternalFields[aliasFuncName] = struct{}{}
	}
}
