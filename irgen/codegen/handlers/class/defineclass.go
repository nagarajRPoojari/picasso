package class

import (
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
)

// DefineClass generates IR definitions for all methods of a class.
// It defines both the class's own functions and inherited ones from its parent,
// ensuring that duplicates are avoided by tracking already-defined names.
//
// Params:
//
//	cls – the AST ClassDeclarationStatement representing the class
func (t *ClassHandler) DefineClass(cls ast.ClassDeclarationStatement) {
	avoid := make(map[string]struct{}, 0)
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DefineFunc(cls.Name, &st, avoid)
			avoid[st.Name] = struct{}{}
		}
	}

	for _, stI := range t.st.TypeHeirarchy.ClassDefs[t.st.TypeHeirarchy.Parent[cls.Name]].Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DefineFunc(cls.Name, &st, avoid)
		}
	}
}

// DefineClassUDT finalizes the LLVM struct definition for a user-defined class.
// It assigns struct field indices for variables and methods, validates overrides,
// incorporates fields and methods from the parent class, and updates the
// opaque UDT created earlier in DeclareClassUDT with concrete field types.
//
// Behavior:
//   - Inherits fields and methods from the parent class, extending the struct.
//   - Registers new fields from the class body, ensuring no duplicates.
//   - Adds function pointers for methods, validating method overrides
//     (signature must match).
//   - Updates the underlying LLVM struct with all collected field types.
//
// Params:
//
//	cls – the AST ClassDeclarationStatement representing the class
func (t *ClassHandler) DefineClassUDT(cls ast.ClassDeclarationStatement) {
	mc := t.st.Classes[cls.Name]
	fieldTypes := make([]types.Type, 0)
	vars := make(map[string]struct{}, 0)

	funcs := make(map[string]uint32, 0)
	// map each fields with corresponding udt struct index
	i := 0

	parentClass := t.st.Classes[cls.Implements]
	if parentClass != nil {
		for _, stI := range t.st.TypeHeirarchy.ClassDefs[t.st.TypeHeirarchy.Parent[cls.Name]].Body {
			switch st := stI.(type) {
			case ast.FunctionDefinitionStatement:
				fqName := t.st.IdentifierBuilder.Attach(cls.Name, st.Name)
				mc.FieldIndexMap[fqName] = i
				funcs[st.Name] = st.Hash
				i++
			case ast.VariableDeclarationStatement:
				fqName := t.st.IdentifierBuilder.Attach(cls.Name, st.Identifier)
				if _, ok := vars[fqName]; ok {
					errorutils.Abort(errorutils.VariableRedeclaration, st.Identifier)
				}

				mc.FieldIndexMap[fqName] = i
				mc.VarAST[fqName] = &st
				vars[fqName] = struct{}{}
				i++
			}
		}
		ptr, ok := parentClass.UDT.(*types.PointerType)
		if !ok {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalUDTDefinitionError, "udt must be a pointer")
		}
		st, ok := ptr.ElemType.(*types.StructType)
		if !ok {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalUDTDefinitionError, "udt must be pointer to a struct")
		}
		fieldTypes = append(fieldTypes, st.Fields...)
	}

	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			fqName := t.st.IdentifierBuilder.Attach(cls.Name, st.Identifier)
			if _, ok := vars[fqName]; ok {
				errorutils.Abort(errorutils.VariableRedeclaration, st.Identifier)
			}

			mc.FieldIndexMap[fqName] = i
			mc.VarAST[fqName] = &st

			fieldType := t.st.TypeHandler.GetLLVMType(st.ExplicitType.Get())
			fieldTypes = append(fieldTypes, fieldType)
			vars[fqName] = struct{}{}

			if st.ExplicitType.GetUnderlyingType() != "" {
				mc.ArrayVarsEleTypes[i] = t.st.TypeHandler.GetLLVMType(st.ExplicitType.GetUnderlyingType())
			}
			i++

		case ast.FunctionDefinitionStatement:
			fqName := t.st.IdentifierBuilder.Attach(cls.Name, st.Name)
			var retType types.Type
			if st.ReturnType != nil {
				retType = t.st.TypeHandler.GetLLVMType(st.ReturnType.Get())
			} else {
				retType = t.st.TypeHandler.GetLLVMType("")
			}

			args := make([]types.Type, 0)
			for _, p := range st.Parameters {
				args = append(args, t.st.TypeHandler.GetLLVMType(p.Type.Get()))
			}

			funcType := types.NewFunc(retType, args...)
			funcPtrType := types.NewPointer(funcType)

			if sh, ok := funcs[st.Name]; ok {
				if sh != st.Hash {
					errorutils.Abort(errorutils.FunctionSignatureMisMatch, st.Name)
				}
				fieldTypes[mc.FieldIndexMap[fqName]] = funcPtrType
			} else {
				fieldTypes = append(fieldTypes, funcPtrType)
				mc.FieldIndexMap[fqName] = i
				i++
			}
		}
	}

	ptr, ok := mc.UDT.(*types.PointerType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalUDTDefinitionError, "udt must be a pointer")
	}
	st, ok := ptr.ElemType.(*types.StructType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalUDTDefinitionError, "udt must be pointer to a struct")
	}
	st.Fields = fieldTypes

}
