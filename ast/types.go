package ast

import "fmt"

type SymbolType struct {
	Atomic bool
	Value  string
}

func (t *SymbolType) GetUnderlyingType() string {
	return ""
}

func (t *SymbolType) Get() string {
	if t.Atomic {
		return fmt.Sprintf("atomic_%s_t", t.Value)
	}
	return t.Value
}

func (t *SymbolType) IsAtomic() bool {
	return t.Atomic
}

func (t *SymbolType) SetAtomic() {
	t.Atomic = true
}

type ListType struct {
	Atomic     bool
	Length     int
	Underlying Type
}

func (t *ListType) IsAtomic() bool {
	return t.Atomic
}

func (t *ListType) SetAtomic() {
	t.Atomic = true
}

func (t *ListType) GetUnderlyingType() string {
	switch underlying := t.Underlying.(type) {
	case *ListType:
		return underlying.GetUnderlyingType()
	default:
		return underlying.Get()
	}
}
func (t *ListType) Get() string { return "array" }
