package generator

import (
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	typedef "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
)

func GetTypeString(t types.Type) string {
	var target string
	switch et := t.(type) {
	case *types.PointerType:
		if st, ok := et.ElemType.(*types.StructType); ok {
			target = st.Name()
		} else {
			target = t.String()
		}
	case *types.StructType:
		target = et.Name()
	default:
		target = t.String()
	}
	return target
}

type IdentifierBuilder struct {
	module string
}

func NewIdentifierBuilder(module string) *IdentifierBuilder {
	return &IdentifierBuilder{module}
}

func (t *IdentifierBuilder) Attach(name ...string) string {
	res := t.module
	for _, n := range name {
		res += "." + n
	}
	return res
}

func RegisterDeclarations(st *state.State, pkg state.PackageEntry, dst, src *ir.Module) {
	funcs := indexFuncs(dst)
	globals := indexGlobals(dst)
	types_ := indexTypes(dst)

	st.FFIModules[pkg.Name] = state.FFIDeclarations{
		Methods: make(map[string]*ir.Func),
		Types:   make(map[string]*types.Type),
		Globals: make(map[string]*ir.Global),
	}

	copyFuncDecls(st.FFIModules[pkg.Name].Methods, dst, src, funcs)
	copyGlobalDecls(st, st.FFIModules[pkg.Name].Globals, dst, src, globals)
	copyTypeDecls(st, st.FFIModules[pkg.Name].Types, dst, src, types_)
}

func indexFuncs(m *ir.Module) map[string]struct{} {
	idx := make(map[string]struct{}, len(m.Funcs))
	for _, fn := range m.Funcs {
		idx[fn.Name()] = struct{}{}
	}
	return idx
}

func indexGlobals(m *ir.Module) map[string]struct{} {
	idx := make(map[string]struct{}, len(m.Globals))
	for _, g := range m.Globals {
		idx[g.Name()] = struct{}{}
	}
	return idx
}

func indexTypes(m *ir.Module) map[string]struct{} {
	idx := make(map[string]struct{}, len(m.TypeDefs))
	for _, td := range m.TypeDefs {
		idx[td.Name()] = struct{}{}
	}
	return idx
}

func copyFuncDecls(funcs map[string]*ir.Func, dst, src *ir.Module, existing map[string]struct{}) {
	for _, fn := range src.Funcs {
		// declarations only
		if len(fn.Blocks) != 0 {
			continue
		}

		if _, ok := existing[fn.Name()]; ok {
			continue
		}

		decl := &ir.Func{
			GlobalIdent: fn.GlobalIdent,
			Sig:         fn.Sig,
			Params:      fn.Params,
			FuncAttrs:   fn.FuncAttrs,
			CallingConv: fn.CallingConv,
			Linkage:     fn.Linkage,
			Visibility:  fn.Visibility,
			UnnamedAddr: fn.UnnamedAddr,
			Align:       fn.Align,
			Section:     fn.Section,
			Comdat:      fn.Comdat,
			GC:          fn.GC,
			Metadata:    fn.Metadata,
		}

		dst.Funcs = append(dst.Funcs, decl)
		existing[fn.Name()] = struct{}{}

		funcs[fn.Name()] = fn
	}
}

func copyGlobalDecls(_ *state.State, globals map[string]*ir.Global, dst, src *ir.Module, existing map[string]struct{}) {
	for _, g := range src.Globals {
		// declarations only
		if g.Init == nil {
			continue
		}

		if !strings.HasPrefix(g.Name(), "__public__") {
			continue
		}

		if _, ok := existing[g.Name()]; ok {
			continue
		}

		decl := &ir.Global{
			GlobalIdent: g.GlobalIdent,
			ContentType: g.ContentType,
			AddrSpace:   g.AddrSpace,
			Linkage:     enum.LinkageExternal,
			Visibility:  g.Visibility,
			UnnamedAddr: g.UnnamedAddr,
			Align:       g.Align,
			Section:     g.Section,
			Comdat:      g.Comdat,
			Metadata:    g.Metadata,
		}

		dst.Globals = append(dst.Globals, decl)
		existing[g.Name()] = struct{}{}

		globals[g.Name()] = decl
	}
}

func copyTypeDecls(st *state.State, typesOut map[string]*types.Type, dst, src *ir.Module, existing map[string]struct{}) {
	for _, td := range src.TypeDefs {
		if _, ok := existing[td.Name()]; ok {
			continue
		}

		dst.TypeDefs = append(dst.TypeDefs, td)
		existing[td.Name()] = struct{}{}

		typesOut[td.Name()] = &td

		// wrap every struct type as pointer to struct
		st.TypeHandler.RegisterClass(td.Name(), typedef.NewMetaClass(types.NewPointer(td), ""))
	}
}
