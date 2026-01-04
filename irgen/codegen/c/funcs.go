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

	// @runtime error
	// t.Funcs[FUNC_RUNTIME_ERROR] =

	// @array_alloc
	t.Funcs[FUNC_ARRAY_ALLOC] = mod.NewFunc(FUNC_ARRAY_ALLOC, types.NewPointer(t.Types[TYPE_ARRAY]), ir.NewParam("", types.I32), ir.NewParam("", types.I32), ir.NewParam("", types.I32))

	t.Funcs[__UTILS__FUNC_DEBUG_ARRAY_INFO] = mod.NewFunc(
		__UTILS__FUNC_DEBUG_ARRAY_INFO,
		types.Void, ir.NewParam("", types.NewPointer(t.Types[TYPE_ARRAY])),
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

	// --- Futex syscalls ---

	t.Funcs[FUNC_SYSCALL_FUTEX_WAIT] = mod.NewFunc(
		FUNC_SYSCALL_FUTEX_WAIT,
		types.I32,
		ir.NewParam("uaddr", types.I32Ptr),
		ir.NewParam("val", types.I32),
		ir.NewParam("timeout", types.I8Ptr),
	)

	t.Funcs[FUNC_SYSCALL_FUTEX_WAKE] = mod.NewFunc(
		FUNC_SYSCALL_FUTEX_WAKE,
		types.I32,
		ir.NewParam("uaddr", types.I32Ptr),
		ir.NewParam("count", types.I32),
	)

	t.Funcs[FUNC_SYSCALL_FUTEX_WAIT_BITSET] = mod.NewFunc(
		FUNC_SYSCALL_FUTEX_WAIT_BITSET,
		types.I32,
		ir.NewParam("uaddr", types.I32Ptr),
		ir.NewParam("val", types.I32),
		ir.NewParam("timeout", types.I8Ptr),
		ir.NewParam("mask", types.I32),
	)

	t.Funcs[FUNC_SYSCALL_FUTEX_WAKE_BITSET] = mod.NewFunc(
		FUNC_SYSCALL_FUTEX_WAKE_BITSET,
		types.I32,
		ir.NewParam("uaddr", types.I32Ptr),
		ir.NewParam("count", types.I32),
		ir.NewParam("mask", types.I32),
	)

	t.Funcs[FUNC_SYSCALL_FUTEX_REQUEUE] = mod.NewFunc(
		FUNC_SYSCALL_FUTEX_REQUEUE,
		types.I32,
		ir.NewParam("uaddr", types.I32Ptr),
		ir.NewParam("wake_count", types.I32),
		ir.NewParam("requeue_count", types.I32),
		ir.NewParam("uaddr2", types.I32Ptr),
	)

	t.Funcs[FUNC_SYSCALL_FUTEX_CMP_REQUEUE] = mod.NewFunc(
		FUNC_SYSCALL_FUTEX_CMP_REQUEUE,
		types.I32,
		ir.NewParam("uaddr", types.I32Ptr),
		ir.NewParam("uaddr2", types.I32Ptr),
		ir.NewParam("wake_count", types.I32),
		ir.NewParam("requeue_count", types.I32),
		ir.NewParam("val", types.I32),
	)

	t.Funcs[FUNC_SYSCALL_FUTEX_WAKE_ONE] = mod.NewFunc(
		FUNC_SYSCALL_FUTEX_WAKE_ONE,
		types.I32,
		ir.NewParam("uaddr", types.I32Ptr),
	)

	t.Funcs[FUNC_SYSCALL_FUTEX_WAKE_ALL] = mod.NewFunc(
		FUNC_SYSCALL_FUTEX_WAKE_ALL,
		types.I32,
		ir.NewParam("uaddr", types.I32Ptr),
	)

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
	t.Funcs[FUNC_ASCAN] = mod.NewFunc(FUNC_ASCAN, t.Types[TYPE_ARRAY], ir.NewParam("format", types.I8Ptr))
	t.Funcs[FUNC_ASCAN].Sig.Variadic = true
	// @sscanf
	t.Funcs[FUNC_SSCAN] = mod.NewFunc(FUNC_SSCAN, t.Types[TYPE_ARRAY], ir.NewParam("format", types.I8Ptr))
	t.Funcs[FUNC_SSCAN].Sig.Variadic = true

	// @afread
	t.Funcs[FUNC_AFREAD] = mod.NewFunc(FUNC_AFREAD, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", t.Types[TYPE_ARRAY]),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)
	// @sfreed
	t.Funcs[FUNC_SFREAD] = mod.NewFunc(FUNC_SFREAD, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", t.Types[TYPE_ARRAY]),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)

	// @afwrite
	t.Funcs[FUNC_AFWRITE] = mod.NewFunc(FUNC_AFWRITE, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("src", t.Types[TYPE_ARRAY]),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)
	// @sfwrite
	t.Funcs[FUNC_SFWRITE] = mod.NewFunc(FUNC_SFWRITE, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("src", t.Types[TYPE_ARRAY]),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)

	// netio
	t.Funcs[FUNC_NET_ACCEPT] = mod.NewFunc(FUNC_NET_ACCEPT, types.I64,
		ir.NewParam("epfd", types.I32),
	)

	t.Funcs[FUNC_NET_READ] = mod.NewFunc(FUNC_NET_READ, types.I64,
		ir.NewParam("fd", types.I64),
		ir.NewParam("buf", types.NewPointer(t.Types[TYPE_ARRAY])),
		ir.NewParam("len", types.I64),
	)

	t.Funcs[FUNC_NET_WRITE] = mod.NewFunc(FUNC_NET_WRITE, types.I64,
		ir.NewParam("fd", types.I64),
		ir.NewParam("buf", types.NewPointer(t.Types[TYPE_ARRAY])),
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

// initAtomicFuncs defines the interface for thread-safe memory operations.
// It maps atomic aliases (like INT8/CHAR) to the same underlying runtime
// implementations and handles type-specific pointer signatures for atomic stores/loads.
func (t *Interface) initAtomicFuncs(mod *ir.Module) {
	// atomic_store_bool
	t.Funcs[FUNC_ATOMIC_STORE_BOOLEAN] = mod.NewFunc(
		FUNC_ATOMIC_STORE_BOOLEAN,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_BOOL])),
		ir.NewParam("val", types.I1),
	)

	// atomic_load_bool
	t.Funcs[FUNC_ATOMIC_LOAD_BOOLEAN] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_BOOLEAN,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_BOOL])),
	)

	// atomic_exchange_bool
	t.Funcs[FUNC_ATOMIC_EXCHANGE_BOOLEAN] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_BOOLEAN,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_BOOL])),
		ir.NewParam("val", types.I1),
	)

	// atomic_compare_exchange_bool
	t.Funcs[FUNC_ATOMIC_CAS_BOOLEAN] = mod.NewFunc(
		FUNC_ATOMIC_CAS_BOOLEAN,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_BOOL])),
		ir.NewParam("expected", types.I1),
		ir.NewParam("desired", types.I1),
	)

	// atomic_store_uint8
	t.Funcs[FUNC_ATOMIC_STORE_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_STORE_UINT8,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_load_uint8
	t.Funcs[FUNC_ATOMIC_LOAD_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_UINT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
	)

	// atomic_add_uint8
	t.Funcs[FUNC_ATOMIC_ADD_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_ADD_UINT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_sub_uint8
	t.Funcs[FUNC_ATOMIC_SUB_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_SUB_UINT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_and_uint8
	t.Funcs[FUNC_ATOMIC_AND_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_AND_UINT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_or_uint8
	t.Funcs[FUNC_ATOMIC_OR_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_OR_UINT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_xor_uint8
	t.Funcs[FUNC_ATOMIC_XOR_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_XOR_UINT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_exchange_uint8
	t.Funcs[FUNC_ATOMIC_EXCHANGE_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_UINT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_compare_exchange_uint8
	t.Funcs[FUNC_ATOMIC_CAS_UINT8] = mod.NewFunc(
		FUNC_ATOMIC_CAS_UINT8,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT8])),
		ir.NewParam("expected", types.I8),
		ir.NewParam("desired", types.I8),
	)

	// atomic_store_uint16
	t.Funcs[FUNC_ATOMIC_STORE_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_STORE_UINT16,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
		ir.NewParam("val", types.I16),
	)

	// atomic_load_uint16
	t.Funcs[FUNC_ATOMIC_LOAD_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_UINT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
	)

	// atomic_add_uint16
	t.Funcs[FUNC_ATOMIC_ADD_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_ADD_UINT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
		ir.NewParam("val", types.I16),
	)

	// atomic_sub_uint16
	t.Funcs[FUNC_ATOMIC_SUB_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_SUB_UINT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
		ir.NewParam("val", types.I16),
	)

	// atomic_and_uint16
	t.Funcs[FUNC_ATOMIC_AND_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_AND_UINT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
		ir.NewParam("val", types.I16),
	)

	// atomic_or_uint16
	t.Funcs[FUNC_ATOMIC_OR_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_OR_UINT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
		ir.NewParam("val", types.I16),
	)

	// atomic_xor_uint16
	t.Funcs[FUNC_ATOMIC_XOR_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_XOR_UINT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
		ir.NewParam("val", types.I16),
	)

	// atomic_exchange_uint16
	t.Funcs[FUNC_ATOMIC_EXCHANGE_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_UINT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
		ir.NewParam("val", types.I16),
	)

	// atomic_compare_exchange_uint16
	t.Funcs[FUNC_ATOMIC_CAS_UINT16] = mod.NewFunc(
		FUNC_ATOMIC_CAS_UINT16,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT16])),
		ir.NewParam("expected", types.I16),
		ir.NewParam("desired", types.I16),
	)

	// atomic_store_uint32
	t.Funcs[FUNC_ATOMIC_STORE_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_STORE_UINT32,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
		ir.NewParam("val", types.I32),
	)

	// atomic_load_uint32
	t.Funcs[FUNC_ATOMIC_LOAD_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_UINT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
	)

	// atomic_add_uint32
	t.Funcs[FUNC_ATOMIC_ADD_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_ADD_UINT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
		ir.NewParam("val", types.I32),
	)

	// atomic_sub_uint32
	t.Funcs[FUNC_ATOMIC_SUB_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_SUB_UINT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
		ir.NewParam("val", types.I32),
	)

	// atomic_and_uint32
	t.Funcs[FUNC_ATOMIC_AND_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_AND_UINT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
		ir.NewParam("val", types.I32),
	)

	// atomic_or_uint32
	t.Funcs[FUNC_ATOMIC_OR_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_OR_UINT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
		ir.NewParam("val", types.I32),
	)

	// atomic_xor_uint32
	t.Funcs[FUNC_ATOMIC_XOR_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_XOR_UINT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
		ir.NewParam("val", types.I32),
	)

	// atomic_exchange_uint32
	t.Funcs[FUNC_ATOMIC_EXCHANGE_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_UINT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
		ir.NewParam("val", types.I32),
	)

	// atomic_compare_exchange_uint32
	t.Funcs[FUNC_ATOMIC_CAS_UINT32] = mod.NewFunc(
		FUNC_ATOMIC_CAS_UINT32,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT32])),
		ir.NewParam("expected", types.I32),
		ir.NewParam("desired", types.I32),
	)

	// atomic_store_uint64
	t.Funcs[FUNC_ATOMIC_STORE_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_STORE_UINT64,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
		ir.NewParam("val", types.I64),
	)

	// atomic_load_uint64
	t.Funcs[FUNC_ATOMIC_LOAD_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_UINT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
	)

	// atomic_add_uint64
	t.Funcs[FUNC_ATOMIC_ADD_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_ADD_UINT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
		ir.NewParam("val", types.I64),
	)

	// atomic_sub_uint64
	t.Funcs[FUNC_ATOMIC_SUB_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_SUB_UINT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
		ir.NewParam("val", types.I64),
	)

	// atomic_and_uint64
	t.Funcs[FUNC_ATOMIC_AND_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_AND_UINT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
		ir.NewParam("val", types.I64),
	)

	// atomic_or_uint64
	t.Funcs[FUNC_ATOMIC_OR_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_OR_UINT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
		ir.NewParam("val", types.I64),
	)

	// atomic_xor_uint64
	t.Funcs[FUNC_ATOMIC_XOR_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_XOR_UINT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
		ir.NewParam("val", types.I64),
	)

	// atomic_exchange_uint64
	t.Funcs[FUNC_ATOMIC_EXCHANGE_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_UINT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
		ir.NewParam("val", types.I64),
	)

	// atomic_compare_exchange_uint64
	t.Funcs[FUNC_ATOMIC_CAS_UINT64] = mod.NewFunc(
		FUNC_ATOMIC_CAS_UINT64,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_UINT64])),
		ir.NewParam("expected", types.I64),
		ir.NewParam("desired", types.I64),
	)

	// ---------------- int8 ----------------

	// atomic_store_int8
	t.Funcs[FUNC_ATOMIC_STORE_INT8] = mod.NewFunc(
		FUNC_ATOMIC_STORE_INT8,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_load_int8
	t.Funcs[FUNC_ATOMIC_LOAD_INT8] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_INT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
	)

	// atomic_add_int8
	t.Funcs[FUNC_ATOMIC_ADD_INT8] = mod.NewFunc(
		FUNC_ATOMIC_ADD_INT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_sub_int8
	t.Funcs[FUNC_ATOMIC_SUB_INT8] = mod.NewFunc(
		FUNC_ATOMIC_SUB_INT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_and_int8
	t.Funcs[FUNC_ATOMIC_AND_INT8] = mod.NewFunc(
		FUNC_ATOMIC_AND_INT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_or_int8
	t.Funcs[FUNC_ATOMIC_OR_INT8] = mod.NewFunc(
		FUNC_ATOMIC_OR_INT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_xor_int8
	t.Funcs[FUNC_ATOMIC_XOR_INT8] = mod.NewFunc(
		FUNC_ATOMIC_XOR_INT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_exchange_int8
	t.Funcs[FUNC_ATOMIC_EXCHANGE_INT8] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_INT8,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
		ir.NewParam("val", types.I8),
	)

	// atomic_compare_exchange_int8
	t.Funcs[FUNC_ATOMIC_CAS_INT8] = mod.NewFunc(
		FUNC_ATOMIC_CAS_INT8,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT8])),
		ir.NewParam("expected", types.I8),
		ir.NewParam("desired", types.I8),
	)

	// ---------------- int16 ----------------

	t.Funcs[FUNC_ATOMIC_STORE_INT16] = mod.NewFunc(
		FUNC_ATOMIC_STORE_INT16,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
		ir.NewParam("val", types.I16),
	)

	t.Funcs[FUNC_ATOMIC_LOAD_INT16] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_INT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
	)

	t.Funcs[FUNC_ATOMIC_ADD_INT16] = mod.NewFunc(
		FUNC_ATOMIC_ADD_INT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
		ir.NewParam("val", types.I16),
	)

	t.Funcs[FUNC_ATOMIC_SUB_INT16] = mod.NewFunc(
		FUNC_ATOMIC_SUB_INT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
		ir.NewParam("val", types.I16),
	)

	t.Funcs[FUNC_ATOMIC_AND_INT16] = mod.NewFunc(
		FUNC_ATOMIC_AND_INT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
		ir.NewParam("val", types.I16),
	)

	t.Funcs[FUNC_ATOMIC_OR_INT16] = mod.NewFunc(
		FUNC_ATOMIC_OR_INT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
		ir.NewParam("val", types.I16),
	)

	t.Funcs[FUNC_ATOMIC_XOR_INT16] = mod.NewFunc(
		FUNC_ATOMIC_XOR_INT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
		ir.NewParam("val", types.I16),
	)

	t.Funcs[FUNC_ATOMIC_EXCHANGE_INT16] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_INT16,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
		ir.NewParam("val", types.I16),
	)

	t.Funcs[FUNC_ATOMIC_CAS_INT16] = mod.NewFunc(
		FUNC_ATOMIC_CAS_INT16,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT16])),
		ir.NewParam("expected", types.I16),
		ir.NewParam("desired", types.I16),
	)

	// ---------------- int32 ----------------

	t.Funcs[FUNC_ATOMIC_STORE_INT32] = mod.NewFunc(
		FUNC_ATOMIC_STORE_INT32,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
		ir.NewParam("val", types.I32),
	)

	t.Funcs[FUNC_ATOMIC_LOAD_INT32] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_INT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
	)

	t.Funcs[FUNC_ATOMIC_ADD_INT32] = mod.NewFunc(
		FUNC_ATOMIC_ADD_INT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
		ir.NewParam("val", types.I32),
	)

	t.Funcs[FUNC_ATOMIC_SUB_INT32] = mod.NewFunc(
		FUNC_ATOMIC_SUB_INT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
		ir.NewParam("val", types.I32),
	)

	t.Funcs[FUNC_ATOMIC_AND_INT32] = mod.NewFunc(
		FUNC_ATOMIC_AND_INT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
		ir.NewParam("val", types.I32),
	)

	t.Funcs[FUNC_ATOMIC_OR_INT32] = mod.NewFunc(
		FUNC_ATOMIC_OR_INT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
		ir.NewParam("val", types.I32),
	)

	t.Funcs[FUNC_ATOMIC_XOR_INT32] = mod.NewFunc(
		FUNC_ATOMIC_XOR_INT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
		ir.NewParam("val", types.I32),
	)

	t.Funcs[FUNC_ATOMIC_EXCHANGE_INT32] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_INT32,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
		ir.NewParam("val", types.I32),
	)

	t.Funcs[FUNC_ATOMIC_CAS_INT32] = mod.NewFunc(
		FUNC_ATOMIC_CAS_INT32,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT32])),
		ir.NewParam("expected", types.I32),
		ir.NewParam("desired", types.I32),
	)

	// ---------------- int64 ----------------

	t.Funcs[FUNC_ATOMIC_STORE_INT64] = mod.NewFunc(
		FUNC_ATOMIC_STORE_INT64,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
		ir.NewParam("val", types.I64),
	)

	t.Funcs[FUNC_ATOMIC_LOAD_INT64] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_INT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
	)

	t.Funcs[FUNC_ATOMIC_ADD_INT64] = mod.NewFunc(
		FUNC_ATOMIC_ADD_INT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
		ir.NewParam("val", types.I64),
	)

	t.Funcs[FUNC_ATOMIC_SUB_INT64] = mod.NewFunc(
		FUNC_ATOMIC_SUB_INT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
		ir.NewParam("val", types.I64),
	)

	t.Funcs[FUNC_ATOMIC_AND_INT64] = mod.NewFunc(
		FUNC_ATOMIC_AND_INT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
		ir.NewParam("val", types.I64),
	)

	t.Funcs[FUNC_ATOMIC_OR_INT64] = mod.NewFunc(
		FUNC_ATOMIC_OR_INT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
		ir.NewParam("val", types.I64),
	)

	t.Funcs[FUNC_ATOMIC_XOR_INT64] = mod.NewFunc(
		FUNC_ATOMIC_XOR_INT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
		ir.NewParam("val", types.I64),
	)

	t.Funcs[FUNC_ATOMIC_EXCHANGE_INT64] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_INT64,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
		ir.NewParam("val", types.I64),
	)

	t.Funcs[FUNC_ATOMIC_CAS_INT64] = mod.NewFunc(
		FUNC_ATOMIC_CAS_INT64,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT64])),
		ir.NewParam("expected", types.I64),
		ir.NewParam("desired", types.I64),
	)

	// atomic_store_float
	t.Funcs[FUNC_ATOMIC_STORE_FLOAT32] = mod.NewFunc(
		FUNC_ATOMIC_STORE_FLOAT32,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT32])),
		ir.NewParam("val", types.Float),
	)

	// atomic_load_float
	t.Funcs[FUNC_ATOMIC_LOAD_FLOAT32] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_FLOAT32,
		types.Float,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT32])),
	)

	// atomic_store_float16
	t.Funcs[FUNC_ATOMIC_STORE_FLOAT16] = mod.NewFunc(
		FUNC_ATOMIC_STORE_FLOAT16,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT16])),
		ir.NewParam("val", types.Half),
	)

	// atomic_load_float16
	t.Funcs[FUNC_ATOMIC_LOAD_FLOAT16] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_FLOAT16,
		types.Half,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT16])),
	)

	// atomic_exchange_float16
	t.Funcs[FUNC_ATOMIC_EXCHANGE_FLOAT16] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_FLOAT16,
		types.Half,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT16])),
		ir.NewParam("val", types.Half),
	)

	// atomic_compare_exchange_float16
	t.Funcs[FUNC_ATOMIC_CAS_FLOAT16] = mod.NewFunc(
		FUNC_ATOMIC_CAS_FLOAT16,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT16])),
		ir.NewParam("expected", types.Half),
		ir.NewParam("desired", types.Half),
	)

	// atomic_exchange_float
	t.Funcs[FUNC_ATOMIC_EXCHANGE_FLOAT32] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_FLOAT32,
		types.Float,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT32])),
		ir.NewParam("val", types.Float),
	)

	// atomic_compare_exchange_float
	t.Funcs[FUNC_ATOMIC_CAS_FLOAT32] = mod.NewFunc(
		FUNC_ATOMIC_CAS_FLOAT32,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT32])),
		ir.NewParam("expected", types.Float),
		ir.NewParam("desired", types.Float),
	)

	// atomic_store_double
	t.Funcs[FUNC_ATOMIC_STORE_FLOAT64] = mod.NewFunc(
		FUNC_ATOMIC_STORE_FLOAT64,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT64])),
		ir.NewParam("val", types.Double),
	)

	// atomic_load_double
	t.Funcs[FUNC_ATOMIC_LOAD_FLOAT64] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_FLOAT64,
		types.Double,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT64])),
	)

	// atomic_exchange_double
	t.Funcs[FUNC_ATOMIC_EXCHANGE_FLOAT64] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_FLOAT64,
		types.Double,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT64])),
		ir.NewParam("val", types.Double),
	)

	// atomic_compare_exchange_double
	t.Funcs[FUNC_ATOMIC_CAS_FLOAT64] = mod.NewFunc(
		FUNC_ATOMIC_CAS_FLOAT64,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT64])),
		ir.NewParam("expected", types.Double),
		ir.NewParam("desired", types.Double),
	)

	// atomic_store_ptr
	t.Funcs[FUNC_ATOMIC_STORE_PTR] = mod.NewFunc(
		FUNC_ATOMIC_STORE_PTR,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_PTR])),
		ir.NewParam("val", types.NewPointer(types.I8)),
	)

	// atomic_load_ptr
	t.Funcs[FUNC_ATOMIC_LOAD_PTR] = mod.NewFunc(
		FUNC_ATOMIC_LOAD_PTR,
		types.NewPointer(types.I8),
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_PTR])),
	)

	// atomic_exchange_ptr
	t.Funcs[FUNC_ATOMIC_EXCHANGE_PTR] = mod.NewFunc(
		FUNC_ATOMIC_EXCHANGE_PTR,
		types.NewPointer(types.I8),
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_PTR])),
		ir.NewParam("val", types.NewPointer(types.I8)),
	)

	// atomic_compare_exchange_ptr
	t.Funcs[FUNC_ATOMIC_CAS_PTR] = mod.NewFunc(
		FUNC_ATOMIC_CAS_PTR,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_PTR])),
		ir.NewParam("expected", types.NewPointer(types.I8)),
		ir.NewParam("desired", types.NewPointer(types.I8)),
	)

}
