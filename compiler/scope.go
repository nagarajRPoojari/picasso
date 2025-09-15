package compiler

import (
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

type VarTree struct {
	tree []map[string]*tf.Var
}

func NewVarTree() *VarTree {
	return &VarTree{
		tree: make([]map[string]*tf.Var, 0),
	}
}

func (t *VarTree) AddLevel() {
	t.tree = append(t.tree, make(map[string]*tf.Var))
}

func (t *VarTree) RemoveLevel() {
	t.tree = t.tree[:len(t.tree)-1]
}

func (t *VarTree) AddNewVar(name string, v tf.Var) {
	if len(t.tree) > 0 {
		t.tree[len(t.tree)-1][name] = &v
	} else {
		errorsx.PanicCompilationError("no level found in var tree")
	}
}

func (t *VarTree) Search(v string) (tf.Var, bool) {
	for i := len(t.tree) - 1; i >= 0; i-- {
		if x, ok := t.tree[i][v]; ok {
			return *x, true
		}
	}
	return nil, false
}
