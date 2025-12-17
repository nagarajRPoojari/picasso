package class

import (
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
)

// DefineClass triggers the emission of concrete LLVM IR function bodies for
// all methods belonging to a class. It handles the traversal of both local
// definitions and inherited parent methods, using a tracking map to ensure
// that overridden methods are defined only once using the most specific
// implementation.
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

// DefineClassUDT performs the "lowering" of high-level class definitions into
// concrete LLVM struct layouts. It handles the structural aspects of inheritance
// by flattening parent fields into the child struct, mapping field identifiers
// to numerical indices (for GEP instructions), and calculating the memory
// footprint for function pointers used in method dispatch.
//
// Key Logic:
//   - Structural Inheritance: Deeply copies field types from the parent UDT to
//     ensure binary compatibility for polymorphism.
//   - Member Indexing: Assigns monotonically increasing indices to fields and
//     methods to facilitate GetElementPtr (GEP) offset calculations.
//   - Signature Validation: Uses method hashes to verify that overrides match
//     the parent signature, aborting on interface mismatches.
//   - Opaque Resolution: Updates the previously declared opaque struct with
//     the finalized field set, completing the type definition in the LLVM module.
func (t *ClassHandler) DefineClassUDT(cls ast.ClassDeclarationStatement) {
	mc := t.st.Classes[cls.Name]
	fieldTypes := make([]types.Type, 0)
	vars := make(map[string]struct{}, 0)

	funcs := make(map[string]uint32, 0)
	// map each fields with corresponding udt struct index
	i := 0

	parentClass := t.st.Classes[cls.Implements]

	// If current class inherits pull its fields to current class
	// DefineClassUDT is expected to be called in the order of inheritance, so
	// no need to recursively pull higher parent classes.
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

	// Opaque resolution: define concrete types of all fields
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
