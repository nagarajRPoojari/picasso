package ast

type SymbolType struct {
	Value string
}

func (t SymbolType) Get() string {
	return t.Value
}

type ListType struct {
	Length     int
	Underlying Type
}

func (t ListType) GetEleType() string {
	switch underlying := t.Underlying.(type) {
	case ListType:
		return underlying.GetEleType()
	default:
		return underlying.Get()
	}
}
func (t ListType) Get() string { return "array" }
