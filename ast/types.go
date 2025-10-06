package ast

type SymbolType struct {
	Value string
}

func (t SymbolType) GetUnderlyingType() string {
	return ""
}

func (t SymbolType) Get() string {
	return t.Value
}

type ListType struct {
	Length     int
	Underlying Type
}

func (t ListType) GetUnderlyingType() string {
	switch underlying := t.Underlying.(type) {
	case ListType:
		return underlying.GetUnderlyingType()
	default:
		return underlying.Get()
	}
}
func (t ListType) Get() string { return "array" }
