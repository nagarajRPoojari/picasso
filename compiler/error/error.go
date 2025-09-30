package errorutils

import "fmt"

const (
	InternalError                   = "[internal]: %s: %s"
	InvalidStatement                = "invalid statement"
	InvalidExpression               = "invalid expression"
	InvalidBinaryExpressionOperand  = "invalid operand for binary operation"
	InvalidBinaryExpressionOperator = "invalid binary operator %s"
	BinaryOperationError            = "binary operation failed: %s"
	ImplicitTypeCastError           = "failed to implicitly type cast: %s to %s"
	MemberExpressionError           = "member expression error: %s"
	UnknownMethod                   = "unknown method %s"
	ClassRedeclaration              = "class %s already defined"
	VariableRedeclaration           = "variable %s already defined"
	FunctionSignatureMisMatch       = "function signatue should match it's parent type: %s"
	UnknownClassField               = "unknown class field %s"
	UnknownClass                    = "unknown class %s"
	UnknownVariable                 = "unknown variable %s"
	MainFuncError                   = "main function error: %s"
	TypeError                       = "type error: %s: %s"
)

const (
	InvalidNativeType = "invalid native type"
	InvalidLLVMType   = "invalid llvm type"
	InvalidTargetType = "invalid target type"
)

const (
	InternalFuncCallError      = "function call error"
	InternalUDTDefinitionError = "udt definition error"
	InternalMemberExprError    = "member expression errror"
	InternalInstantiationError = "intantiation error"
	InternalTypeError          = "type error"
)

func Abort(msg string, args ...any) {
	format := fmt.Sprintf(msg, args...)
	panic(format)
}
