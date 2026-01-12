package scope

import (
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// VarTree manages a stack of scopes (tree) and a global symbol table.
type VarTree struct {
	tree    []map[string]*tf.Var // Stack of local scopes
	globals map[string]*tf.Var   // Top-level global symbols
}

// NewVarTree initializes an empty variable tree.
func NewVarTree() *VarTree {
	return &VarTree{
		tree:    make([]map[string]*tf.Var, 0),
		globals: make(map[string]*tf.Var),
	}
}

// AddBlock pushes a new lexical block (e.g., if-statement, loop) onto the stack.
func (t *VarTree) AddBlock() {
	t.tree = append(t.tree, make(map[string]*tf.Var))
}

// AddFunc pushes a function boundary. It uses a nil marker or double block
// to differentiate function scope from nested blocks.
func (t *VarTree) AddFunc() {
	// Add a marker for function boundary
	t.tree = append(t.tree, nil)
	t.AddBlock()
}

// RemoveBlock pops the current innermost scope.
func (t *VarTree) RemoveBlock() {
	if len(t.tree) > 0 {
		t.tree = t.tree[:len(t.tree)-1]
	}
}

// RemoveFunc pops both the function block and the boundary marker.
func (t *VarTree) RemoveFunc() {
	t.RemoveBlock() // Remove local block
	t.RemoveBlock() // Remove function marker
}

// AddNewVar inserts a variable into the current (innermost) scope.
func (t *VarTree) AddNewVar(name string, v tf.Var) {
	if len(t.tree) > 0 && t.tree[len(t.tree)-1] != nil {
		t.tree[len(t.tree)-1][name] = &v
	} else {
		errorutils.Abort(errorutils.InternalError, "no active scope level found to add variable: "+name)
	}
}

// RegisterTypeHolders adds a symbol to the global scope.
func (t *VarTree) RegisterTypeHolders(block *bc.BlockHolder, name string, s tf.Var) {
	t.globals[name] = &s
}

func (t *VarTree) AddGlobal(name string, v *tf.Var) {
	t.globals[name] = v
}

// Search looks for a variable name from the innermost scope outwards to globals.
func (t *VarTree) Search(v string) (tf.Var, bool) {
	for i := len(t.tree) - 1; i >= 0; i-- {
		// If we hit a nil entry, it's a function boundary;
		// depending on language rules, you might stop or continue.
		if t.tree[i] == nil {
			continue
		}
		if x, ok := t.tree[i][v]; ok {
			return *x, true
		}
	}
	return t.searchGlobal(v)
}

// Replace finds an existing variable and updates its value.
func (t *VarTree) Replace(v string, by tf.Var) {
	for i := len(t.tree) - 1; i >= 0; i-- {
		if t.tree[i] == nil {
			continue
		}
		if _, ok := t.tree[i][v]; ok {
			t.tree[i][v] = &by
			return
		}
	}
}

// Exists checks if a variable is defined in the IMMEDIATE current scope.
func (t *VarTree) Exists(v string) bool {
	if len(t.tree) == 0 || t.tree[len(t.tree)-1] == nil {
		return false
	}
	_, ok := t.tree[len(t.tree)-1][v]
	return ok
}

// searchGlobal is an internal helper for global lookup.
func (t *VarTree) searchGlobal(v string) (tf.Var, bool) {
	x, ok := t.globals[v]
	if !ok {
		return nil, false
	}
	return *x, ok
}
