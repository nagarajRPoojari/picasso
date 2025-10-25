package errorutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

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
	UnknownClassField               = "unknown class field %s in class %s"
	UnknownClass                    = "unknown class %s"
	UnknownVariable                 = "unknown variable %s"
	UnknownModule                   = "unknown module %s"
	InvalidModulerSource            = "invalid source of module %s: %s"
	MainFuncError                   = "main function error: %s"
	TypeError                       = "type error: %s: %s"
	ParamsError                     = "params mismatch in %s: expected %v"
	GlobalVarsNotAllowedError       = "global vars not allowed"
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
	redBold := color.New(color.FgRed, color.Bold).SprintFunc()
	mainColor := color.New(color.FgYellow).SprintFunc()
	argColor := color.New(color.FgRed).SprintFunc()

	formattedMsg := msg
	argPositions := []struct {
		start int
		end   int
	}{}

	for _, arg := range args {
		placeholder := "%s"
		idx := strings.Index(formattedMsg, placeholder)
		if idx == -1 {
			break
		}
		argStr := fmt.Sprint(arg)
		formattedMsg = formattedMsg[:idx] + argStr + formattedMsg[idx+len(placeholder):]
		argPositions = append(argPositions, struct{ start, end int }{idx, idx + len(argStr)})
	}

	coloredMsg := ""
	lastIndex := 0
	for _, pos := range argPositions {
		coloredMsg += mainColor(formattedMsg[lastIndex:pos.start]) // normal text
		coloredMsg += argColor(formattedMsg[pos.start:pos.end])    // red arg
		lastIndex = pos.end
	}
	coloredMsg += mainColor(formattedMsg[lastIndex:])

	fmt.Print(redBold("ERROR: "))
	fmt.Println(coloredMsg)

	lineLen := len("ERROR: ") + len([]rune(formattedMsg))
	squiggleRunes := make([]rune, lineLen)
	for i := 0; i < lineLen; i++ {
		squiggleRunes[i] = ' '
	}
	for _, pos := range argPositions {
		for i := pos.start + len("ERROR: "); i < pos.end+len("ERROR: "); i++ {
			squiggleRunes[i] = '~'
		}
	}
	fmt.Println(argColor(string(squiggleRunes)))

	os.Exit(1)
}
