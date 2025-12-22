package interfaceh

import (
	"fmt"

	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
)

func (t *InterfaceHandler) DeclareInterface(ifs ast.InterfaceDeclarationStatement, sourcePkg state.PackageEntry) {

	ifName := identifier.NewIdentifierBuilder(sourcePkg.Name).Attach(ifs.Name)
	aliasName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(ifs.Name)

	if _, ok := t.st.Classes[aliasName]; ok {
		errorutils.Abort(errorutils.ClassRedeclaration, ifName)
	}

	udt := types.NewStruct() // opaque
	if _, ok := t.st.GlobalTypeList[ifName]; !ok {
		t.st.GlobalTypeList[ifName] = t.st.Module.NewTypeDef(ifName, udt)
	}
	mc := &tf.MetaClass{
		FieldIndexMap:     make(map[string]int),
		ArrayVarsEleTypes: make(map[int]types.Type),
		VarAST:            make(map[string]*ast.VariableDeclarationStatement),
		UDT:               types.NewPointer(udt),
		Methods:           make(map[string]*ir.Func),
		Returns:           map[string]ast.Type{},
	}
	t.st.Classes[aliasName] = mc

	mi := tf.NewMetaInterface()
	mi.UDT = types.NewPointer(udt)

	t.st.Interfaces[aliasName] = mi
	t.st.TypeHandler.RegisterInterface(aliasName, mi)
}

func (t *InterfaceHandler) DeclareFunctions(ifs ast.InterfaceDeclarationStatement, sourcePkg state.PackageEntry) {
	aliasName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(ifs.Name)
	for _, stI := range ifs.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			// store function signature for interface

			fqClsName := identifier.NewIdentifierBuilder(sourcePkg.Name).Attach(ifs.Name)
			aliasClsName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(ifs.Name)

			params := make([]*ir.Param, 0)
			for _, p := range st.Parameters {
				params = append(params, ir.NewParam(p.Name, t.st.TypeHandler.GetLLVMType(p.Type.Get())))
			}

			// at the end pass `this` parameter representing current object
			udt := t.st.Classes[aliasClsName].UDT
			params = append(params, ir.NewParam(constants.THIS, udt))

			fqFuncName := fmt.Sprintf("%s.%s", fqClsName, st.Name)
			aliasFuncName := fmt.Sprintf("%s.%s", aliasClsName, st.Name)

			var retType types.Type
			if st.ReturnType != nil {
				retType = t.st.TypeHandler.GetLLVMType(st.ReturnType.Get())
			} else {
				retType = t.st.TypeHandler.GetLLVMType("")
			}

			// store current functions so that later during class instantiation instance
			// can be made pointing to the functions.
			if _, ok := t.st.Classes[aliasClsName].Methods[aliasFuncName]; !ok {
				f, ok := t.st.GlobalFuncList[fqFuncName]
				if !ok {
					f = t.st.Module.NewFunc(fqFuncName, retType, params...)
					t.st.GlobalFuncList[fqFuncName] = f
					t.st.Interfaces[aliasName].Methods[st.Name] = tf.MethodSig{
						Hash:     utils.HashFuncSig(st.Parameters, st.ReturnType),
						Name:     st.Name,
						FuncType: f,
					}
				}
				t.st.Classes[aliasClsName].Methods[aliasFuncName] = f
				t.st.Classes[aliasClsName].Returns[aliasFuncName] = st.ReturnType
			}
		}
	}
}
