package scope

import (
	errorutils "github.com/nagarajRPoojari/niyama/frontend/codegen/error"
	tf "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
)

type VarTree struct {
	tree    []map[string]*tf.Var
	globals map[string]*tf.Var
}

func NewVarTree() *VarTree {
	return &VarTree{
		tree:    make([]map[string]*tf.Var, 0),
		globals: make(map[string]*tf.Var, 0),
	}
}

func (t *VarTree) AddBlock() {
	t.tree = append(t.tree, make(map[string]*tf.Var))
}

func (t *VarTree) RegisterTypeHolders(block *bc.BlockHolder, name string, s tf.Var) {
	t.globals[name] = &s
}

func (t *VarTree) AddFunc() {
	t.tree = append(t.tree, nil)
	t.AddBlock()
}

func (t *VarTree) RemoveBlock() {
	t.tree = t.tree[:len(t.tree)-1]
}
func (t *VarTree) RemoveFunc() {
	t.RemoveBlock()
	t.RemoveBlock()
}

func (t *VarTree) AddNewVar(name string, v tf.Var) {
	if len(t.tree) > 0 {
		t.tree[len(t.tree)-1][name] = &v
	} else {
		errorutils.Abort(errorutils.InternalError, "no level found in var tree")
	}
}

func (t *VarTree) Search(v string) (tf.Var, bool) {
	for i := len(t.tree) - 1; i >= 0; i-- {
		if t.tree == nil {
			break
		}
		if x, ok := t.tree[i][v]; ok {
			return *x, true
		}
	}
	return t.searchGlobal(v)
}

func (t *VarTree) Replace(v string, by tf.Var) {
	for i := len(t.tree) - 1; i >= 0; i-- {
		if t.tree == nil {
			break
		}
		if _, ok := t.tree[i][v]; ok {
			t.tree[i][v] = &by
		}
	}
}

func (t *VarTree) searchGlobal(v string) (tf.Var, bool) {
	x, ok := t.globals[v]
	if !ok {
		return nil, false
	}
	return *x, ok
}

func (t *VarTree) Exists(v string) bool {
	if len(t.tree) == 0 {
		return false
	}
	_, ok := t.tree[len(t.tree)-1][v]
	return ok
}
