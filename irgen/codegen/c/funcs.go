package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// registerFuncs orchestrates the declaration of all external symbols within the
// current LLVM module. It partitions declarations into logical runtime domains.
func (t *Interface) registerFuncs(mod *ir.Module) {
	t.initRuntime(mod)     // Core lifecycle and memory
	t.initSyscalls(mod)    // init linux syscalls
	t.initStdio(mod)       // File and Console I/O
	t.initAtomicFuncs(mod) // Thread-safe primitives
	t.initStrs(mod)        // High-level string manipulation
}

// initRuntime declares core engine functions including memory management (malloc/alloc),
// threading, process synchronization, and essential libc utilities like memcpy.
func (t *Interface) initRuntime(mod *ir.Module) {

	t.Funcs[FUNC_HASH] = mod.NewFunc(FUNC_HASH, types.I64, ir.NewParam("data", types.I8Ptr), ir.NewParam("len", types.I64))

	t.Funcs[FUNC_STRLEN] = mod.NewFunc(FUNC_STRLEN, types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	t.Funcs[FUNC_STRCMP] = mod.NewFunc(FUNC_STRCMP, types.I32, ir.NewParam("", types.NewPointer(types.I8)), ir.NewParam("", types.NewPointer(types.I8)))

	t.Funcs[FUNC_MEMCPY] = mod.NewFunc("llvm.memcpy.p0i8.p0i8.i64",
		types.Void,
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("src", types.I8Ptr),
		ir.NewParam("len", types.I64),
		ir.NewParam("isvolatile", types.I1),
	)
	// @thread
	fnType := types.NewFunc(
		types.NewPointer(types.I8),
		types.NewPointer(types.I8),
	)
	t.Funcs[FUNC_THREAD] = mod.NewFunc(FUNC_THREAD, types.Void,
		ir.NewParam("", types.NewPointer(fnType)),
	)

	// @self_yield
	t.Funcs[FUNC_SELF_YIELD] = mod.NewFunc(FUNC_SELF_YIELD, types.Void)

	// @alloc
	t.Funcs[FUNC_ALLOC] = mod.NewFunc(FUNC_ALLOC, types.I8Ptr, ir.NewParam("", types.I64))

	t.Funcs[FUNC_MALLOC] = mod.NewFunc(FUNC_MALLOC, types.I8Ptr, ir.NewParam("", types.I64))

	// @runtime_init
	t.Funcs[FUNC_RUNTIME_INIT] = mod.NewFunc(FUNC_RUNTIME_INIT, types.Void)

	// @array_alloc
	t.Funcs[FUNC_ARRAY_ALLOC] = mod.NewFunc(FUNC_ARRAY_ALLOC, types.NewPointer(t.Types[TYPE_ARRAY]), ir.NewParam("", types.I32), ir.NewParam("", types.I32), ir.NewParam("", types.I32))

	t.Funcs[__UTILS__FUNC_DEBUG_ARRAY_INFO] = mod.NewFunc(
		__UTILS__FUNC_DEBUG_ARRAY_INFO,
		types.Void, ir.NewParam("", types.NewPointer(t.Types[TYPE_ARRAY])),
	)
}

// initAtomicFuncs defines the interface for thread-safe memory operations.
// It maps atomic aliases (like INT8/CHAR) to the same underlying runtime
// implementations and handles type-specific pointer signatures for atomic stores/loads.
func (t *Interface) initAtomicFuncs(mod *ir.Module) {
	// --- @bool ---
	t.Funcs[FUNC_ATOMIC_STORE_BOOL] = mod.NewFunc(FUNC_ATOMIC_STORE_BOOL,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_BOOL])),
		ir.NewParam("val", types.I1),
	)
	t.Funcs[FUNC_ATOMIC_LOAD_BOOL] = mod.NewFunc(FUNC_ATOMIC_LOAD_BOOL,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_BOOL])),
	)

	// --- @int8 ---
	t.Funcs[FUNC_ATOMIC_STORE_CHAR] = mod.NewFunc(FUNC_ATOMIC_STORE_CHAR,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[FUNC_ATOMIC_STORE_INT8] = t.Funcs[FUNC_ATOMIC_STORE_CHAR]
	t.Funcs[FUNC_ATOMIC_LOAD_CHAR] = mod.NewFunc(FUNC_ATOMIC_LOAD_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_CHAR])),
	)
	t.Funcs[FUNC_ATOMIC_LOAD_INT8] = t.Funcs[FUNC_ATOMIC_LOAD_CHAR]
	t.Funcs[FUNC_ATOMIC_ADD_CHAR] = mod.NewFunc(FUNC_ATOMIC_ADD_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[FUNC_ATOMIC_ADD_INT8] = t.Funcs[FUNC_ATOMIC_ADD_CHAR]
	t.Funcs[FUNC_ATOMIC_SUB_CHAR] = mod.NewFunc(FUNC_ATOMIC_SUB_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[FUNC_ATOMIC_SUB_INT8] = t.Funcs[FUNC_ATOMIC_SUB_CHAR]

	// --- @int16 ---
	t.Funcs[FUNC_ATOMIC_STORE_SHORT] = mod.NewFunc(FUNC_ATOMIC_STORE_SHORT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[FUNC_ATOMIC_STORE_INT16] = t.Funcs[FUNC_ATOMIC_STORE_SHORT]
	t.Funcs[FUNC_ATOMIC_LOAD_SHORT] = mod.NewFunc(FUNC_ATOMIC_LOAD_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_SHORT])),
	)
	t.Funcs[FUNC_ATOMIC_LOAD_INT16] = t.Funcs[FUNC_ATOMIC_LOAD_SHORT]
	t.Funcs[FUNC_ATOMIC_ADD_SHORT] = mod.NewFunc(FUNC_ATOMIC_ADD_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[FUNC_ATOMIC_ADD_INT16] = t.Funcs[FUNC_ATOMIC_ADD_SHORT]
	t.Funcs[FUNC_ATOMIC_SUB_SHORT] = mod.NewFunc(FUNC_ATOMIC_SUB_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[FUNC_ATOMIC_SUB_INT16] = t.Funcs[FUNC_ATOMIC_SUB_SHORT]

	// --- @int32 ---
	t.Funcs[FUNC_ATOMIC_STORE_INT] = mod.NewFunc(FUNC_ATOMIC_STORE_INT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[FUNC_ATOMIC_STORE_INT32] = t.Funcs[FUNC_ATOMIC_STORE_INT]
	t.Funcs[FUNC_ATOMIC_LOAD_INT] = mod.NewFunc(FUNC_ATOMIC_LOAD_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT])),
	)
	t.Funcs[FUNC_ATOMIC_LOAD_INT32] = t.Funcs[FUNC_ATOMIC_LOAD_INT]
	t.Funcs[FUNC_ATOMIC_ADD_INT] = mod.NewFunc(FUNC_ATOMIC_ADD_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[FUNC_ATOMIC_ADD_INT32] = t.Funcs[FUNC_ATOMIC_ADD_INT]
	t.Funcs[FUNC_ATOMIC_SUB_INT] = mod.NewFunc(FUNC_ATOMIC_SUB_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[FUNC_ATOMIC_SUB_INT32] = t.Funcs[FUNC_ATOMIC_SUB_INT]

	// --- @int64 ---
	t.Funcs[FUNC_ATOMIC_STORE_LONG] = mod.NewFunc(FUNC_ATOMIC_STORE_LONG,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[FUNC_ATOMIC_STORE_INT64] = t.Funcs[FUNC_ATOMIC_STORE_LONG]
	t.Funcs[FUNC_ATOMIC_LOAD_LONG] = mod.NewFunc(FUNC_ATOMIC_LOAD_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_LONG])),
	)
	t.Funcs[FUNC_ATOMIC_LOAD_INT64] = t.Funcs[FUNC_ATOMIC_LOAD_LONG]
	t.Funcs[FUNC_ATOMIC_ADD_LONG] = mod.NewFunc(FUNC_ATOMIC_ADD_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[FUNC_ATOMIC_ADD_INT64] = t.Funcs[FUNC_ATOMIC_ADD_LONG]
	t.Funcs[FUNC_ATOMIC_SUB_LONG] = mod.NewFunc(FUNC_ATOMIC_SUB_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[FUNC_ATOMIC_SUB_INT64] = t.Funcs[FUNC_ATOMIC_SUB_LONG]

	// --- floats and others unchanged ---
	t.Funcs[FUNC_ATOMIC_STORE_FLOAT] = mod.NewFunc(FUNC_ATOMIC_STORE_FLOAT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT])),
		ir.NewParam("val", types.Float),
	)
	t.Funcs[FUNC_ATOMIC_LOAD_FLOAT] = mod.NewFunc(FUNC_ATOMIC_LOAD_FLOAT,
		types.Float,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT])),
	)

	t.Funcs[FUNC_ATOMIC_STORE_DOUBLE] = mod.NewFunc(FUNC_ATOMIC_STORE_DOUBLE,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_DOUBLE])),
		ir.NewParam("val", types.Double),
	)
	t.Funcs[FUNC_ATOMIC_LOAD_DOUBLE] = mod.NewFunc(FUNC_ATOMIC_LOAD_DOUBLE,
		types.Double,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_DOUBLE])),
	)

	t.Funcs[FUNC_ATOMIC_STORE_PTR] = mod.NewFunc(FUNC_ATOMIC_STORE_PTR,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_PTR])),
		ir.NewParam("val", types.NewPointer(types.I8)),
	)
	t.Funcs[FUNC_ATOMIC_LOAD_PTR] = mod.NewFunc(FUNC_ATOMIC_LOAD_PTR,
		types.NewPointer(types.I8),
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_PTR])),
	)
}

func (t *Interface) initSyscalls(mod *ir.Module) {
	t.Funcs[FUNC_SYSCALL_ERRNO] = mod.NewFunc(FUNC_SYSCALL_ERRNO, types.I32)

	t.Funcs[FUNC_SYSCALL_GETPID] = mod.NewFunc(FUNC_SYSCALL_GETPID, types.I32)
	t.Funcs[FUNC_SYSCALL_GETPPID] = mod.NewFunc(FUNC_SYSCALL_GETPPID, types.I32)
	t.Funcs[FUNC_SYSCALL_GETTID] = mod.NewFunc(FUNC_SYSCALL_GETTID, types.I32)
	t.Funcs[FUNC_SYSCALL_EXIT] = mod.NewFunc(FUNC_SYSCALL_EXIT, types.Void, ir.NewParam("code", types.I32))
	t.Funcs[FUNC_SYSCALL_FORK] = mod.NewFunc(FUNC_SYSCALL_FORK, types.I32)
	t.Funcs[FUNC_SYSCALL_WAITPID] = mod.NewFunc(FUNC_SYSCALL_WAITPID, types.I32,
		ir.NewParam("pid", types.I32),
		ir.NewParam("status", types.I32Ptr),
		ir.NewParam("options", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_KILL] = mod.NewFunc(FUNC_SYSCALL_KILL, types.I32,
		ir.NewParam("pid", types.I32),
		ir.NewParam("sig", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_EXECVE] = mod.NewFunc(FUNC_SYSCALL_EXECVE, types.I32,
		ir.NewParam("path", types.I8Ptr),
		ir.NewParam("argv", types.NewPointer(types.I8Ptr)),
		ir.NewParam("envp", types.NewPointer(types.I8Ptr)),
	)
	t.Funcs[FUNC_SYSCALL_EXECVP] = mod.NewFunc(FUNC_SYSCALL_EXECVP, types.I32,
		ir.NewParam("file", types.I8Ptr),
		ir.NewParam("argv", types.NewPointer(types.I8Ptr)),
	)
	t.Funcs[FUNC_SYSCALL_ENVIRON] = mod.NewFunc(FUNC_SYSCALL_ENVIRON, types.NewPointer(types.I8Ptr))
	t.Funcs[FUNC_SYSCALL_GETENV] = mod.NewFunc(FUNC_SYSCALL_GETENV, types.I8Ptr, ir.NewParam("key", types.I8Ptr))
	t.Funcs[FUNC_SYSCALL_SETENV] = mod.NewFunc(FUNC_SYSCALL_SETENV, types.I32,
		ir.NewParam("key", types.I8Ptr),
		ir.NewParam("value", types.I8Ptr),
		ir.NewParam("overwrite", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_UNSETENV] = mod.NewFunc(FUNC_SYSCALL_UNSETENV, types.I32, ir.NewParam("key", types.I8Ptr))
	t.Funcs[FUNC_SYSCALL_GETCWD] = mod.NewFunc(FUNC_SYSCALL_GETCWD, types.I32,
		ir.NewParam("buf", types.I8Ptr),
		ir.NewParam("size", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_CHDIR] = mod.NewFunc(FUNC_SYSCALL_CHDIR, types.I32, ir.NewParam("path", types.I8Ptr))
	t.Funcs[FUNC_SYSCALL_GETUID] = mod.NewFunc(FUNC_SYSCALL_GETUID, types.I32)
	t.Funcs[FUNC_SYSCALL_GETEUID] = mod.NewFunc(FUNC_SYSCALL_GETEUID, types.I32)
	t.Funcs[FUNC_SYSCALL_GETGID] = mod.NewFunc(FUNC_SYSCALL_GETGID, types.I32)
	t.Funcs[FUNC_SYSCALL_GETEGID] = mod.NewFunc(FUNC_SYSCALL_GETEGID, types.I32)
	t.Funcs[FUNC_SYSCALL_SETUID] = mod.NewFunc(FUNC_SYSCALL_SETUID, types.I32, ir.NewParam("uid", types.I32))
	t.Funcs[FUNC_SYSCALL_SETGID] = mod.NewFunc(FUNC_SYSCALL_SETGID, types.I32, ir.NewParam("gid", types.I32))
	t.Funcs[FUNC_SYSCALL_SETPGID] = mod.NewFunc(FUNC_SYSCALL_SETPGID, types.I32,
		ir.NewParam("pid", types.I32),
		ir.NewParam("pgid", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_GETPGID] = mod.NewFunc(FUNC_SYSCALL_GETPGID, types.I32, ir.NewParam("pid", types.I32))
	t.Funcs[FUNC_SYSCALL_GETPGRP] = mod.NewFunc(FUNC_SYSCALL_GETPGRP, types.I32)
	t.Funcs[FUNC_SYSCALL_SETSID] = mod.NewFunc(FUNC_SYSCALL_SETSID, types.I32)
	t.Funcs[FUNC_SYSCALL_GETRLIMIT] = mod.NewFunc(FUNC_SYSCALL_GETRLIMIT, types.I32,
		ir.NewParam("resource", types.I32),
		ir.NewParam("rlim", types.I8Ptr),
	)
	t.Funcs[FUNC_SYSCALL_SETRLIMIT] = mod.NewFunc(FUNC_SYSCALL_SETRLIMIT, types.I32,
		ir.NewParam("resource", types.I32),
		ir.NewParam("rlim", types.I8Ptr),
	)
	t.Funcs[FUNC_SYSCALL_SIGNAL_INSTALL] = mod.NewFunc(FUNC_SYSCALL_SIGNAL_INSTALL, types.I32,
		ir.NewParam("sig", types.I32),
		ir.NewParam("handler", types.NewPointer(types.NewFunc(types.I32, types.I32))),
	)
	t.Funcs[FUNC_SYSCALL_OPEN] = mod.NewFunc(FUNC_SYSCALL_OPEN, types.I32,
		ir.NewParam("path", types.I8Ptr),
		ir.NewParam("flags", types.I32),
		ir.NewParam("mode", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_CLOSE] = mod.NewFunc(FUNC_SYSCALL_CLOSE, types.I32, ir.NewParam("fd", types.I32))
	t.Funcs[FUNC_SYSCALL_READ] = mod.NewFunc(FUNC_SYSCALL_READ, types.I64,
		ir.NewParam("fd", types.I32),
		ir.NewParam("buf", types.I8Ptr),
		ir.NewParam("n", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_WRITE] = mod.NewFunc(FUNC_SYSCALL_WRITE, types.I64,
		ir.NewParam("fd", types.I32),
		ir.NewParam("buf", types.I8Ptr),
		ir.NewParam("n", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_LSEEK] = mod.NewFunc(FUNC_SYSCALL_LSEEK, types.I64,
		ir.NewParam("fd", types.I32),
		ir.NewParam("offset", types.I64),
		ir.NewParam("whence", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_FSTAT] = mod.NewFunc(FUNC_SYSCALL_FSTAT, types.I32,
		ir.NewParam("fd", types.I32),
		ir.NewParam("st", types.I8Ptr),
	)
	t.Funcs[FUNC_SYSCALL_DUP] = mod.NewFunc(FUNC_SYSCALL_DUP, types.I32, ir.NewParam("fd", types.I32))
	t.Funcs[FUNC_SYSCALL_DUP2] = mod.NewFunc(FUNC_SYSCALL_DUP2, types.I32,
		ir.NewParam("oldfd", types.I32),
		ir.NewParam("newfd", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_FCNTL] = mod.NewFunc(FUNC_SYSCALL_FCNTL, types.I32,
		ir.NewParam("fd", types.I32),
		ir.NewParam("cmd", types.I32),
		ir.NewParam("arg", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_MKDIR] = mod.NewFunc(FUNC_SYSCALL_MKDIR, types.I32,
		ir.NewParam("path", types.I8Ptr),
		ir.NewParam("mode", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_RMDIR] = mod.NewFunc(FUNC_SYSCALL_RMDIR, types.I32, ir.NewParam("path", types.I8Ptr))
	t.Funcs[FUNC_SYSCALL_UNLINK] = mod.NewFunc(FUNC_SYSCALL_UNLINK, types.I32, ir.NewParam("path", types.I8Ptr))
	t.Funcs[FUNC_SYSCALL_RENAME] = mod.NewFunc(FUNC_SYSCALL_RENAME, types.I32,
		ir.NewParam("oldpath", types.I8Ptr),
		ir.NewParam("newpath", types.I8Ptr),
	)
	t.Funcs[FUNC_SYSCALL_RENAMEAT2] = mod.NewFunc(FUNC_SYSCALL_RENAMEAT2, types.I32,
		ir.NewParam("oldpath", types.I8Ptr),
		ir.NewParam("newpath", types.I8Ptr),
		ir.NewParam("flags", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_LINK] = mod.NewFunc(FUNC_SYSCALL_LINK, types.I32,
		ir.NewParam("oldpath", types.I8Ptr),
		ir.NewParam("newpath", types.I8Ptr),
	)
	t.Funcs[FUNC_SYSCALL_SYMLINK] = mod.NewFunc(FUNC_SYSCALL_SYMLINK, types.I32,
		ir.NewParam("target", types.I8Ptr),
		ir.NewParam("linkpath", types.I8Ptr),
	)
	t.Funcs[FUNC_SYSCALL_READLINK] = mod.NewFunc(FUNC_SYSCALL_READLINK, types.I64,
		ir.NewParam("path", types.I8Ptr),
		ir.NewParam("buf", types.I8Ptr),
		ir.NewParam("size", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_STAT] = mod.NewFunc(FUNC_SYSCALL_STAT, types.I32,
		ir.NewParam("path", types.I8Ptr),
		ir.NewParam("st", types.I8Ptr),
	)
	t.Funcs[FUNC_SYSCALL_LSTAT] = mod.NewFunc(FUNC_SYSCALL_LSTAT, types.I32,
		ir.NewParam("path", types.I8Ptr),
		ir.NewParam("st", types.I8Ptr),
	)
	t.Funcs[FUNC_SYSCALL_ACCESS] = mod.NewFunc(FUNC_SYSCALL_ACCESS, types.I32,
		ir.NewParam("path", types.I8Ptr),
		ir.NewParam("mode", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_GETDENTS64] = mod.NewFunc(FUNC_SYSCALL_GETDENTS64, types.I32,
		ir.NewParam("fd", types.I32),
		ir.NewParam("buf", types.I8Ptr),
		ir.NewParam("size", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_MMAP] = mod.NewFunc(FUNC_SYSCALL_MMAP, types.I8Ptr,
		ir.NewParam("addr", types.I8Ptr),
		ir.NewParam("len", types.I64),
		ir.NewParam("prot", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_MUNMAP] = mod.NewFunc(FUNC_SYSCALL_MUNMAP, types.I32,
		ir.NewParam("addr", types.I8Ptr),
		ir.NewParam("len", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_MPROTECT] = mod.NewFunc(FUNC_SYSCALL_MPROTECT, types.I32,
		ir.NewParam("addr", types.I8Ptr),
		ir.NewParam("len", types.I64),
		ir.NewParam("prot", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_MADVISE] = mod.NewFunc(FUNC_SYSCALL_MADVISE, types.I32,
		ir.NewParam("addr", types.I8Ptr),
		ir.NewParam("len", types.I64),
		ir.NewParam("advice", types.I32),
	)
	t.Funcs[FUNC_SYSCALL_MLOCK] = mod.NewFunc(FUNC_SYSCALL_MLOCK, types.I32,
		ir.NewParam("addr", types.I8Ptr),
		ir.NewParam("len", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_MUNLOCK] = mod.NewFunc(FUNC_SYSCALL_MUNLOCK, types.I32,
		ir.NewParam("addr", types.I8Ptr),
		ir.NewParam("len", types.I64),
	)
	t.Funcs[FUNC_SYSCALL_MLOCKALL] = mod.NewFunc(FUNC_SYSCALL_MLOCKALL, types.I32, ir.NewParam("flags", types.I32))
	t.Funcs[FUNC_SYSCALL_MUNLOCKALL] = mod.NewFunc(FUNC_SYSCALL_MUNLOCKALL, types.I32)
	t.Funcs[FUNC_SYSCALL_PAGE_SIZE] = mod.NewFunc(FUNC_SYSCALL_PAGE_SIZE, types.I64)
}

// initStdio declares standard input/output operations. It includes support
// for both blocking (Standard C) and non-blocking/asynchronous (prefixed with 'a')
// I/O calls, as well as variadic signatures for formatted printing.
func (t *Interface) initStdio(mod *ir.Module) {
	// @fopen
	t.Funcs[FUNC_FOPEN] = mod.NewFunc(FUNC_FOPEN, types.I8Ptr,
		ir.NewParam("filename", types.I8Ptr),
		ir.NewParam("mode", types.I8Ptr),
	)
	// @fclose
	t.Funcs[FUNC_FCLOSE] = mod.NewFunc(FUNC_FCLOSE, types.I32,
		ir.NewParam("stream", types.I8Ptr),
	)

	// @fflush
	t.Funcs[FUNC_FFLUSH] = mod.NewFunc(FUNC_FFLUSH, types.I32,
		ir.NewParam("stream", types.I8Ptr),
	)

	// @fseek
	t.Funcs[FUNC_FSEEK] = mod.NewFunc(FUNC_FSEEK, types.I32,
		ir.NewParam("stream", types.I8Ptr),
		ir.NewParam("offset", types.I64),
		ir.NewParam("whence", types.I32),
	)

	// @aprintf
	t.Funcs[FUNC_APRINTF] = mod.NewFunc(FUNC_APRINTF, types.I32, ir.NewParam("", types.I8Ptr))
	t.Funcs[FUNC_APRINTF].Sig.Variadic = true
	// @sprintf
	t.Funcs[FUNC_SPRINTF] = mod.NewFunc(FUNC_SPRINTF, types.I32, ir.NewParam("", types.I8Ptr))
	t.Funcs[FUNC_SPRINTF].Sig.Variadic = true

	// @ascanf
	t.Funcs[FUNC_ASCAN] = mod.NewFunc(FUNC_ASCAN, types.I32, ir.NewParam("format", types.I8Ptr))
	t.Funcs[FUNC_ASCAN].Sig.Variadic = true
	// @sscanf
	t.Funcs[FUNC_SSCAN] = mod.NewFunc(FUNC_SSCAN, types.I32, ir.NewParam("format", types.I8Ptr))
	t.Funcs[FUNC_SSCAN].Sig.Variadic = true

	// @afread
	t.Funcs[FUNC_AFREAD] = mod.NewFunc(FUNC_AFREAD, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)
	// @sfreed
	t.Funcs[FUNC_SFREAD] = mod.NewFunc(FUNC_SFREAD, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)

	// @afwrite
	t.Funcs[FUNC_AFWRITE] = mod.NewFunc(FUNC_AFWRITE, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)
	// @sfwrite
	t.Funcs[FUNC_SFWRITE] = mod.NewFunc(FUNC_SFWRITE, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)

	// netio
	t.Funcs[FUNC_NET_ACCEPT] = mod.NewFunc(FUNC_NET_ACCEPT, types.I64,
		ir.NewParam("epfd", types.I32),
	)

	t.Funcs[FUNC_NET_READ] = mod.NewFunc(FUNC_NET_READ, types.I64,
		ir.NewParam("fd", types.I32),
		ir.NewParam("buf", types.NewPointer(types.I8)),
		ir.NewParam("len", types.I64),
	)

	t.Funcs[FUNC_NET_WRITE] = mod.NewFunc(FUNC_NET_WRITE, types.I64,
		ir.NewParam("fd", types.I32),
		ir.NewParam("buf", types.NewPointer(types.I8)),
		ir.NewParam("len", types.I64),
	)

	t.Funcs[FUNC_NET_LISTEN] = mod.NewFunc(FUNC_NET_LISTEN, types.I64,
		ir.NewParam("addr", types.NewPointer(types.I8)),
		ir.NewParam("port", types.I64),
		ir.NewParam("backlog", types.I32),
	)
}

// initStrs defines high-level string utilities that operate on it's
// string representation, providing abstraction over raw C-string pointers.
func (t *Interface) initStrs(mod *ir.Module) {
	// @format
	t.Funcs[FUNC_FORMAT] = mod.NewFunc(FUNC_FORMAT, types.I8Ptr,
		ir.NewParam("fmt", types.I8Ptr),
	)

	// @len
	t.Funcs[FUNC_LEN] = mod.NewFunc(FUNC_LEN, types.I32,
		ir.NewParam("str", types.I8Ptr),
	)

	// @compare
	t.Funcs[FUNC_COMPARE] = mod.NewFunc(FUNC_COMPARE, types.I32,
		ir.NewParam("a", types.I8Ptr),
		ir.NewParam("b", types.I8Ptr),
	)
}
