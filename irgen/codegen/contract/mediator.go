package contract

// Mediator defines the interface that all handlers will use to talk to each other.
// We use 'any' here to avoid importing the handler packages into this interface file.
type Mediator interface {
	GetExpressionHandler() any
	GetStatementHandler() any
	GetFuncHandler() any
	GetBlockHandler() any
	GetClassHandler() any
	GetInterfaceHandler() any
}
