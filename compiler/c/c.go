package c

import (
	"sync"

	"github.com/llir/llvm/ir"
)

const (
	ALLOC        = "lang_alloc"
	RUNTIME_INIT = "runtime_init"
	ARRAY_ALLOC  = "lang_alloc_array"

	PRINTF  = "printf"
	SCANF   = "scanf"
	FPRINTF = "fprintf"
	FSCANF  = "fscanf"
	FPUTS   = "fputs"
	FGETS   = "fgets"
	FOPEN   = "fopen"
	FCLOSE  = "fclose"
	FFLUSH  = "fflush"
	FSEEK   = "fseek"
	MALLOC  = "malloc"
	FREE    = "free"
	STRLEN  = "strlen"
	MEMCPY  = "memcpy"
	MEMSET  = "memset"
	MEMMOVE = "memmove"
	EXIT    = "exit"
)

type Interface struct {
	Funcs map[string]*ir.Func
}

var Instance *Interface
var once sync.Once

func NewInterface(mod *ir.Module) *Interface {
	t := &Interface{}
	once.Do(func() {
		t.Funcs = make(map[string]*ir.Func)
		t.registerFuncs(mod)
		Instance = t
	})
	return Instance
}
