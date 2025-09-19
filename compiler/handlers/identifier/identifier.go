package identifier

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
