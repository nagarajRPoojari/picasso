package ast

import "fmt"

// SymbolType represents a named type in the Niyama type system.
// It can represent either a primitive/built-in type (Atomic) or
// a user-defined type (like a class or interface).
type SymbolType struct {
	// Atomic indicates if this is a atomic language type (e.g.,atomic int, atomic float).
	Atomic bool
	// Value is the raw name of the type (e.g., "int" or "MyClass").
	Value string
}

// GetUnderlyingType for a SymbolType returns an empty string as it
// is the base of the type hierarchy and has no nested types.
func (t *SymbolType) GetUnderlyingType() string {
	return ""
}

// Get returns the string representation of the type.
// If Atomic is true, it formats the name as a system-level type (e.g., "atomic_int_t").
func (t *SymbolType) Get() string {
	if t.Atomic {
		return fmt.Sprintf("atomic_%s_t", t.Value)
	}
	return t.Value
}

// IsAtomic reports whether the type is a atomic unit.
func (t *SymbolType) IsAtomic() bool {
	return t.Atomic
}

// SetAtomic marks the type as a atomic unit.
func (t *SymbolType) SetAtomic() {
	t.Atomic = true
}

// ListType represents a fixed-length or dynamic collection of a specific type.
// It implements a recursive structure to support multi-dimensional arrays.
type ListType struct {
	Atomic bool
	// Length specifies the number of elements; used for static array allocation.
	Length int
	// Underlying points to the Type of the elements contained in the list.
	Underlying Type
}

// IsAtomic reports whether the list is treated as a atomic unit.
func (t *ListType) IsAtomic() bool {
	return t.Atomic
}

// SetAtomic marks the list as a atomic unit.
func (t *ListType) SetAtomic() {
	t.Atomic = true
}

// GetUnderlyingType traverses nested ListTypes to find the base element type.
// Example: For a [][][]int, it will return the string representation of "int".
func (t *ListType) GetUnderlyingType() string {
	switch underlying := t.Underlying.(type) {
	case *ListType:
		return underlying.GetUnderlyingType()
	default:
		return underlying.Get()
	}
}

// Get returns the generic string identifier for list structures.
func (t *ListType) Get() string { return "array" }
