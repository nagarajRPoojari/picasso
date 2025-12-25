/*
Package c provides the Foreign Function Interface (FFI) and Runtime Bridge for
compiler. Exposes set of functions & types from c code.
Notes:
  - it is expected to link corresponding c binaries during runtime
  - functions are accessed through libs defined in libs/ directory
  - __public__ identifier indicates c functions name while corresponding alias
    used as exposed name
*/
package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// Foreign Function and Runtime identifiers.
// These constants define the mapping between internal function
// aliases and the actual symbol names in the C runtime or standard library.
const (
	// Debug and Diagnostic Utilities
	__UTILS__FUNC_DEBUG_ARRAY_INFO = "__public__debug_array_info"

	// Concurrency and Threading
	FUNC_THREAD      = "thread"
	ALIAS_THREAD     = "thread"
	FUNC_SELF_YIELD  = "self_yield"
	ALIAS_SELF_YIELD = "self_yield"

	// Blocking I/O Calls
	// These map directly to the underlying OS file descriptors.
	ALIAS_PRINTF  = "printf"
	ALIAS_FPRINTF = "fprintf"
	ALIAS_SCANF   = "scanf"
	ALIAS_FSCANF  = "fscanf"
	ALIAS_FREAD   = "fread"
	ALIAS_FWRITE  = "fwrite"

	FUNC_FOPEN   = "fopen"
	ALIAS_FOPEN  = "fopen"
	FUNC_FCLOSE  = "fclose"
	ALIAS_FCLOSE = "fclose"
	FUNC_FFLUSH  = "fflush"
	ALIAS_FFLUSH = "fflush"
	FUNC_FSEEK   = "fseek"
	ALIAS_FSEEK  = "fseek"
	FUNC_FPUTS   = "fputs"
	ALIAS_FPUTS  = "fputs"
	FUNC_FGETS   = "fgets"
	ALIAS_FGETS  = "fgets"

	// Formatted I/O and Buffer operations
	FUNC_FSCANF  = "fscanf"
	FUNC_FPRINTF = "fprintf"
	FUNC_SPRINTF = "__public__sprintf"
	FUNC_SSCAN   = "__public__sscan"
	FUNC_SFREAD  = "__public__sfread"
	FUNC_SFWRITE = "__public__sfwrite"

	// Non-blocking I/O (Asynchronous Wrapper Calls)
	FUNC_APRINTF = "__public__aprintf"
	FUNC_ASCAN   = "__public__ascan"
	FUNC_AFREAD  = "__public__afread"
	FUNC_AFWRITE = "__public__afwrite"

	// Memory Management
	// Includes both standard C malloc and GC-tracked allocation.
	FUNC_MALLOC   = "malloc"
	ALIAS_MALLOC  = "malloc"
	FUNC_MEMCPY   = "memcpy"
	ALIAS_MEMCPY  = "memcpy"
	FUNC_MEMSET   = "memset"
	ALIAS_MEMSET  = "memset"
	FUNC_MEMMOVE  = "memmove"
	ALIAS_MEMMOVE = "memmove"

	FUNC_RUNTIME_INIT  = "runtime_init"
	ALIAS_RUNTIME_INIT = "runtime_init"

	FUNC_ALLOC  = "__public__alloc" // Garbage Collector tracked allocation
	ALIAS_ALLOC = "alloc"

	FUNC_ARRAY_ALLOC  = "__public__alloc_array"
	ALIAS_ARRAY_ALLOC = "alloc_array"

	TYPE_ARRAY = "array"

	// String and Container Operations
	FUNC_STRLEN   = "strlen"
	ALIAS_STRLEN  = "strlen"
	FUNC_FORMAT   = "format"
	ALIAS_FORMAT  = "format"
	FUNC_LEN      = "len"
	ALIAS_LEN     = "len"
	FUNC_COMPARE  = "compare"
	ALIAS_COMPARE = "compare"
	FUNC_STRCMP   = "strcmp"
	ALIAS_STRCMP  = "strcmp"

	// Process Control
	FUNC_HASH  = "hash"
	ALIAS_HASH = "hash"

	// Atomic Operations and Synchronization
	// These constants map to thread-safe primitives for various bit-widths.
	ALIAS_ATOMIC_STORE = "atomic_store"
	ALIAS_ATOMIC_LOAD  = "atomic_load"
	ALIAS_ATOMIC_ADD   = "atomic_add"
	ALIAS_ATOMIC_SUB   = "atomic_sub"

	FUNC_ATOMIC_STORE_BOOL = "__public__atomic_store_bool"
	FUNC_ATOMIC_LOAD_BOOL  = "__public__atomic_load_bool"

	// 8-bit Atomics
	FUNC_ATOMIC_STORE_CHAR = "__public__atomic_store_char"
	FUNC_ATOMIC_STORE_INT8 = "__public__atomic_store_int8"
	FUNC_ATOMIC_LOAD_CHAR  = "__public__atomic_load_char"
	FUNC_ATOMIC_LOAD_INT8  = "__public__atomic_load_int8"
	FUNC_ATOMIC_ADD_CHAR   = "__public__atomic_add_char"
	FUNC_ATOMIC_ADD_INT8   = "__public__atomic_add_int8"
	FUNC_ATOMIC_SUB_CHAR   = "__public__atomic_sub_char"
	FUNC_ATOMIC_SUB_INT8   = "__public__atomic_sub_int8"

	// 16-bit Atomics
	FUNC_ATOMIC_STORE_SHORT = "__public__atomic_store_short"
	FUNC_ATOMIC_STORE_INT16 = "__public__atomic_store_int16"
	FUNC_ATOMIC_LOAD_SHORT  = "__public__atomic_load_short"
	FUNC_ATOMIC_LOAD_INT16  = "__public__atomic_load_int16"
	FUNC_ATOMIC_ADD_SHORT   = "__public__atomic_add_short"
	FUNC_ATOMIC_ADD_INT16   = "__public__atomic_add_int16"
	FUNC_ATOMIC_SUB_SHORT   = "__public__atomic_sub_short"
	FUNC_ATOMIC_SUB_INT16   = "__public__atomic_sub_int16"

	// 32-bit Atomics
	FUNC_ATOMIC_STORE_INT   = "__public__atomic_store_int"
	FUNC_ATOMIC_STORE_INT32 = "__public__atomic_store_int32"
	FUNC_ATOMIC_LOAD_INT    = "__public__atomic_load_int"
	FUNC_ATOMIC_LOAD_INT32  = "__public__atomic_load_int32"
	FUNC_ATOMIC_ADD_INT     = "__public__atomic_add_int"
	FUNC_ATOMIC_ADD_INT32   = "__public__atomic_add_int32"
	FUNC_ATOMIC_SUB_INT     = "__public__atomic_sub_int"
	FUNC_ATOMIC_SUB_INT32   = "__public__atomic_sub_int32"

	// 64-bit Atomics
	FUNC_ATOMIC_STORE_LONG  = "__public__atomic_store_long"
	FUNC_ATOMIC_STORE_INT64 = "__public__atomic_store_int64"
	FUNC_ATOMIC_LOAD_LONG   = "__public__atomic_load_long"
	FUNC_ATOMIC_LOAD_INT64  = "__public__atomic_load_int64"
	FUNC_ATOMIC_ADD_LONG    = "__public__atomic_add_long"
	FUNC_ATOMIC_ADD_INT64   = "__public__atomic_add_int64"
	FUNC_ATOMIC_SUB_LONG    = "__public__atomic_sub_long"
	FUNC_ATOMIC_SUB_INT64   = "__public__atomic_sub_int64"

	// Pointer and Floating Point Atomics
	FUNC_ATOMIC_STORE_LLONG  = "__public__atomic_store_llong"
	FUNC_ATOMIC_LOAD_LLONG   = "__public__atomic_load_llong"
	FUNC_ATOMIC_ADD_LLONG    = "__public__atomic_add_llong"
	FUNC_ATOMIC_SUB_LLONG    = "__public__atomic_sub_llong"
	FUNC_ATOMIC_STORE_FLOAT  = "__public__atomic_store_float"
	FUNC_ATOMIC_LOAD_FLOAT   = "__public__atomic_load_float"
	FUNC_ATOMIC_STORE_DOUBLE = "__public__atomic_store_double"
	FUNC_ATOMIC_LOAD_DOUBLE  = "__public__atomic_load_double"
	FUNC_ATOMIC_STORE_PTR    = "__public__atomic_store_ptr"
	FUNC_ATOMIC_LOAD_PTR     = "__public__atomic_load_ptr"

	// Runtime Type Names
	// These identifiers are used when declaring opaque or alias types in LLVM IR.
	TYPE_ATOMIC_BOOL   = "atomic_bool_t"
	TYPE_ATOMIC_CHAR   = "atomic_char_t"
	TYPE_ATOMIC_INT8   = "atomic_int8_t"
	TYPE_ATOMIC_SHORT  = "atomic_short_t"
	TYPE_ATOMIC_INT16  = "atomic_int16_t"
	TYPE_ATOMIC_INT    = "atomic_int_t"
	TYPE_ATOMIC_INT32  = "atomic_int32_t"
	TYPE_ATOMIC_LONG   = "atomic_long_t"
	TYPE_ATOMIC_INT64  = "atomic_long_t"
	TYPE_ATOMIC_LLONG  = "atomic_llong_t"
	TYPE_ATOMIC_FLOAT  = "atomic_float_t"
	TYPE_ATOMIC_DOUBLE = "atomic_double_t"
	TYPE_ATOMIC_PTR    = "atomic_ptr_t"

	// syscalls
	FUNC_SYSCALL_ERRNO          = "__public__errno"
	FUNC_SYSCALL_GETPID         = "__public__getpid"
	FUNC_SYSCALL_GETPPID        = "__public__getppid"
	FUNC_SYSCALL_GETTID         = "__public__gettid"
	FUNC_SYSCALL_EXIT           = "__public__exit"
	FUNC_SYSCALL_FORK           = "__public__fork"
	FUNC_SYSCALL_WAITPID        = "__public__waitpid"
	FUNC_SYSCALL_KILL           = "__public__kill"
	FUNC_SYSCALL_EXECVE         = "__public__execve"
	FUNC_SYSCALL_EXECVP         = "__public__execvp"
	FUNC_SYSCALL_ENVIRON        = "__public__environ"
	FUNC_SYSCALL_GETENV         = "__public__getenv"
	FUNC_SYSCALL_SETENV         = "__public__setenv"
	FUNC_SYSCALL_UNSETENV       = "__public__unsetenv"
	FUNC_SYSCALL_GETCWD         = "__public__getcwd"
	FUNC_SYSCALL_CHDIR          = "__public__chdir"
	FUNC_SYSCALL_GETUID         = "__public__getuid"
	FUNC_SYSCALL_GETEUID        = "__public__geteuid"
	FUNC_SYSCALL_GETGID         = "__public__getgid"
	FUNC_SYSCALL_GETEGID        = "__public__getegid"
	FUNC_SYSCALL_SETUID         = "__public__setuid"
	FUNC_SYSCALL_SETGID         = "__public__setgid"
	FUNC_SYSCALL_SETPGID        = "__public__setpgid"
	FUNC_SYSCALL_GETPGID        = "__public__getpgid"
	FUNC_SYSCALL_GETPGRP        = "__public__getpgrp"
	FUNC_SYSCALL_SETSID         = "__public__setsid"
	FUNC_SYSCALL_GETRLIMIT      = "__public__getrlimit"
	FUNC_SYSCALL_SETRLIMIT      = "__public__setrlimit"
	FUNC_SYSCALL_SIGNAL_INSTALL = "__public__signal_install"
	FUNC_SYSCALL_OPEN           = "__public__open"
	FUNC_SYSCALL_CLOSE          = "__public__close"
	FUNC_SYSCALL_READ           = "__public__read"
	FUNC_SYSCALL_WRITE          = "__public__write"
	FUNC_SYSCALL_LSEEK          = "__public__lseek"
	FUNC_SYSCALL_FSTAT          = "__public__fstat"
	FUNC_SYSCALL_DUP            = "__public__dup"
	FUNC_SYSCALL_DUP2           = "__public__dup2"
	FUNC_SYSCALL_FCNTL          = "__public__fcntl"
	FUNC_SYSCALL_MKDIR          = "__public__mkdir"
	FUNC_SYSCALL_RMDIR          = "__public__rmdir"
	FUNC_SYSCALL_UNLINK         = "__public__unlink"
	FUNC_SYSCALL_RENAME         = "__public__rename"
	FUNC_SYSCALL_RENAMEAT2      = "__public__renameat2"
	FUNC_SYSCALL_LINK           = "__public__link"
	FUNC_SYSCALL_SYMLINK        = "__public__symlink"
	FUNC_SYSCALL_READLINK       = "__public__readlink"
	FUNC_SYSCALL_STAT           = "__public__stat"
	FUNC_SYSCALL_LSTAT          = "__public__lstat"
	FUNC_SYSCALL_ACCESS         = "__public__access"
	FUNC_SYSCALL_GETDENTS64     = "__public__getdents64"
	FUNC_SYSCALL_MMAP           = "__public__mmap"
	FUNC_SYSCALL_MUNMAP         = "__public__munmap"
	FUNC_SYSCALL_MPROTECT       = "__public__mprotect"
	FUNC_SYSCALL_MADVISE        = "__public__madvise"
	FUNC_SYSCALL_MLOCK          = "__public__mlock"
	FUNC_SYSCALL_MUNLOCK        = "__public__munlock"
	FUNC_SYSCALL_MLOCKALL       = "__public__mlockall"
	FUNC_SYSCALL_MUNLOCKALL     = "__public__munlockall"
	FUNC_SYSCALL_PAGE_SIZE      = "__public__page_size"

	ALIAS_SYSCALL_ERRNO          = "errno"
	ALIAS_SYSCALL_GETPID         = "getpid"
	ALIAS_SYSCALL_GETPPID        = "getppid"
	ALIAS_SYSCALL_GETTID         = "gettid"
	ALIAS_SYSCALL_EXIT           = "exit"
	ALIAS_SYSCALL_FORK           = "fork"
	ALIAS_SYSCALL_WAITPID        = "waitpid"
	ALIAS_SYSCALL_KILL           = "kill"
	ALIAS_SYSCALL_EXECVE         = "execve"
	ALIAS_SYSCALL_EXECVP         = "execvp"
	ALIAS_SYSCALL_ENVIRON        = "environ"
	ALIAS_SYSCALL_GETENV         = "getenv"
	ALIAS_SYSCALL_SETENV         = "setenv"
	ALIAS_SYSCALL_UNSETENV       = "unsetenv"
	ALIAS_SYSCALL_GETCWD         = "getcwd"
	ALIAS_SYSCALL_CHDIR          = "chdir"
	ALIAS_SYSCALL_GETUID         = "getuid"
	ALIAS_SYSCALL_GETEUID        = "geteuid"
	ALIAS_SYSCALL_GETGID         = "getgid"
	ALIAS_SYSCALL_GETEGID        = "getegid"
	ALIAS_SYSCALL_SETUID         = "setuid"
	ALIAS_SYSCALL_SETGID         = "setgid"
	ALIAS_SYSCALL_SETPGID        = "setpgid"
	ALIAS_SYSCALL_GETPGID        = "getpgid"
	ALIAS_SYSCALL_GETPGRP        = "getpgrp"
	ALIAS_SYSCALL_SETSID         = "setsid"
	ALIAS_SYSCALL_GETRLIMIT      = "getrlimit"
	ALIAS_SYSCALL_SETRLIMIT      = "setrlimit"
	ALIAS_SYSCALL_SIGNAL_INSTALL = "signal_install"
	ALIAS_SYSCALL_OPEN           = "open"
	ALIAS_SYSCALL_CLOSE          = "close"
	ALIAS_SYSCALL_READ           = "read"
	ALIAS_SYSCALL_WRITE          = "write"
	ALIAS_SYSCALL_LSEEK          = "lseek"
	ALIAS_SYSCALL_FSTAT          = "fstat"
	ALIAS_SYSCALL_DUP            = "dup"
	ALIAS_SYSCALL_DUP2           = "dup2"
	ALIAS_SYSCALL_FCNTL          = "fcntl"
	ALIAS_SYSCALL_MKDIR          = "mkdir"
	ALIAS_SYSCALL_RMDIR          = "rmdir"
	ALIAS_SYSCALL_UNLINK         = "unlink"
	ALIAS_SYSCALL_RENAME         = "rename"
	ALIAS_SYSCALL_RENAMEAT2      = "renameat2"
	ALIAS_SYSCALL_LINK           = "link"
	ALIAS_SYSCALL_SYMLINK        = "symlink"
	ALIAS_SYSCALL_READLINK       = "readlink"
	ALIAS_SYSCALL_STAT           = "stat"
	ALIAS_SYSCALL_LSTAT          = "lstat"
	ALIAS_SYSCALL_ACCESS         = "access"
	ALIAS_SYSCALL_GETDENTS64     = "getdents64"
	ALIAS_SYSCALL_MMAP           = "mmap"
	ALIAS_SYSCALL_MUNMAP         = "munmap"
	ALIAS_SYSCALL_MPROTECT       = "mprotect"
	ALIAS_SYSCALL_MADVISE        = "madvise"
	ALIAS_SYSCALL_MLOCK          = "mlock"
	ALIAS_SYSCALL_MUNLOCK        = "munlock"
	ALIAS_SYSCALL_MLOCKALL       = "mlockall"
	ALIAS_SYSCALL_MUNLOCKALL     = "munlockall"
	ALIAS_SYSCALL_PAGE_SIZE      = "page_size"

	ALIAS_OS_EAGAIN = "EAGAIN"
	ALIAS_OS_EINTR  = "EINTR"
	ALIAS_OS_EINVAL = "EINVAL"
	ALIAS_OS_EPERM  = "EPERM"
	ALIAS_OS_ENOENT = "ENOENT"
	ALIAS_OS_ENOMEM = "ENOMEM"

	ALIAS_OS_WNOHANG    = "WNOHANG"
	ALIAS_OS_WUNTRACED  = "WUNTRACED"
	ALIAS_OS_WCONTINUED = "WCONTINUED"

	ALIAS_OS_SIGINT  = "SIGINT"
	ALIAS_OS_SIGTERM = "SIGTERM"
	ALIAS_OS_SIGKILL = "SIGKILL"
	ALIAS_OS_SIGSEGV = "SIGSEGV"
	ALIAS_OS_SIGABRT = "SIGABRT"
	ALIAS_OS_SIGCHLD = "SIGCHLD"
	ALIAS_OS_SIGPIPE = "SIGPIPE"
	ALIAS_OS_SIGALRM = "SIGALRM"
	ALIAS_OS_SIGUSR1 = "SIGUSR1"
	ALIAS_OS_SIGUSR2 = "SIGUSR2"

	ALIAS_OS_RLIMIT_CPU    = "RLIMIT_CPU"
	ALIAS_OS_RLIMIT_FSIZE  = "RLIMIT_FSIZE"
	ALIAS_OS_RLIMIT_DATA   = "RLIMIT_DATA"
	ALIAS_OS_RLIMIT_STACK  = "RLIMIT_STACK"
	ALIAS_OS_RLIMIT_CORE   = "RLIMIT_CORE"
	ALIAS_OS_RLIMIT_NOFILE = "RLIMIT_NOFILE"
	ALIAS_OS_RLIMIT_AS     = "RLIMIT_AS"

	/* Standard file descriptor numbers */
	/** Standard input */
	ALIAS_OS_STDIN_FD = "0"
	/** Standard output */
	ALIAS_OS_STDOUT_FD = "1"
	/** Standard error */
	ALIAS_OS_STDERR_FD = "2"

	/*open() flags*/
	ALIAS_OS_O_RDONLY   = "O_RDONLY"
	ALIAS_OS_O_WRONLY   = "O_WRONLY"
	ALIAS_OS_O_RDWR     = "O_RDWR"
	ALIAS_OS_O_APPEND   = "O_APPEND"
	ALIAS_OS_O_CREAT    = "O_CREAT"
	ALIAS_OS_O_EXCL     = "O_EXCL"
	ALIAS_OS_O_TRUNC    = "O_TRUNC"
	ALIAS_OS_O_CLOEXEC  = "O_CLOEXEC"
	ALIAS_OS_O_NONBLOCK = "O_NONBLOCK"
	ALIAS_OS_O_SYNC     = "O_SYNC"
	ALIAS_OS_O_DSYNC    = "O_DSYNC"
	ALIAS_OS_O_DIRECT   = "O_DIRECT"

	/*seek constants*/
	ALIAS_OS_SEEK_SET = "SEEK_SET"
	ALIAS_OS_SEEK_CUR = "SEEK_CUR"
	ALIAS_OS_SEEK_END = "SEEK_END"

	/*fcntl commands*/
	ALIAS_OS_F_DUPFD         = "F_DUPFD"
	ALIAS_OS_F_DUPFD_CLOEXEC = "F_DUPFD_CLOEXEC"
	ALIAS_OS_F_GETFD         = "F_GETFD"
	ALIAS_OS_F_SETFD         = "F_SETFD"
	ALIAS_OS_F_GETFL         = "F_GETFL"
	ALIAS_OS_F_SETFL         = "F_SETFL"

	/*FD flags*/
	ALIAS_OS_FD_CLOEXEC = "FD_CLOEXEC"

	/*stat mode bits*/
	ALIAS_OS_S_IFREG  = "S_IFREG"
	ALIAS_OS_S_IFDIR  = "S_IFDIR"
	ALIAS_OS_S_IFCHR  = "S_IFCHR"
	ALIAS_OS_S_IFBLK  = "S_IFBLK"
	ALIAS_OS_S_IFIFO  = "S_IFIFO"
	ALIAS_OS_S_IFLNK  = "S_IFLNK"
	ALIAS_OS_S_IFSOCK = "S_IFSOCK"

	ALIAS_OS_S_IRUSR = "S_IRUSR"
	ALIAS_OS_S_IWUSR = "S_IWUSR"
	ALIAS_OS_S_IXUSR = "S_IXUSR"
	ALIAS_OS_S_IRGRP = "S_IRGRP"
	ALIAS_OS_S_IWGRP = "S_IWGRP"
	ALIAS_OS_S_IXGRP = "S_IXGRP"
	ALIAS_OS_S_IROTH = "S_IROTH"
	ALIAS_OS_S_IWOTH = "S_IWOTH"
	ALIAS_OS_S_IXOTH = "S_IXOTH"

	/* Errors (FD-relevant subset)*/
	ALIAS_OS_EBADF  = "EBADF"
	ALIAS_OS_EPIPE  = "EPIPE"
	ALIAS_OS_EIO    = "EIO"
	ALIAS_OS_ENOSPC = "ENOSPC"

	/*Special directory FDs*/
	/** Current working directory */
	ALIAS_OS_AT_FDCWD = "AT_FDCWD"

	/*unlinkat / renameat flags*/
	ALIAS_OS_AT_REMOVEDIR      = "AT_REMOVEDIR"
	ALIAS_OS_AT_SYMLINK_FOLLOW = "AT_SYMLINK_FOLLOW"

	/*link / rename flags*/
	OS_RENAME_NOREPLACE = "RENAME_NOREPLACE" // "RENAME_NOREPLACE"
	OS_RENAME_EXCHANGE  = "RENAME_EXCHANGE"  // "RENAME_EXCHANGE"
	OS_RENAME_WHITEOUT  = "RENAME_WHITEOUT"  // "RENAME_WHITEOUT"

	/*Access mode flags*/
	ALIAS_OS_F_OK = "F_OK"
	ALIAS_OS_R_OK = "R_OK"
	ALIAS_OS_W_OK = "W_OK"
	ALIAS_OS_X_OK = "X_OK"

	/*Directory entry types (d_type)*/
	ALIAS_OS_DT_UNKNOWN = "DT_UNKNOWN"
	ALIAS_OS_DT_FIFO    = "DT_FIFO"
	ALIAS_OS_DT_CHR     = "DT_CHR"
	ALIAS_OS_DT_DIR     = "DT_DIR"
	ALIAS_OS_DT_BLK     = "DT_BLK"
	ALIAS_OS_DT_REG     = "DT_REG"
	ALIAS_OS_DT_LNK     = "DT_LNK"
	ALIAS_OS_DT_SOCK    = "DT_SOCK"
	ALIAS_OS_DT_WHT     = "DT_WHT"

	/*Memory protection flags*/
	ALIAS_OS_PROT_NONE  = "PROT_NONE"
	ALIAS_OS_PROT_READ  = "PROT_READ"
	ALIAS_OS_PROT_WRITE = "PROT_WRITE"
	ALIAS_OS_PROT_EXEC  = "PROT_EXEC"

	/*mmap flags*/
	ALIAS_OS_MAP_SHARED    = "MAP_SHARED"
	ALIAS_OS_MAP_PRIVATE   = "MAP_PRIVATE"
	ALIAS_OS_MAP_FIXED     = "MAP_FIXED"
	ALIAS_OS_MAP_ANONYMOUS = "MAP_ANONYMOUS"
	ALIAS_OS_MAP_STACK     = "MAP_STACK"
	ALIAS_OS_MAP_NORESERVE = "MAP_NORESERVE"
	ALIAS_OS_MAP_POPULATE  = "MAP_POPULATE"
	ALIAS_OS_MAP_GROWSDOWN = "MAP_GROWSDOWN"

	/*madvise advice*/
	ALIAS_OS_MADV_NORMAL      = "MADV_NORMAL"
	ALIAS_OS_MADV_RANDOM      = "MADV_RANDOM"
	ALIAS_OS_MADV_SEQUENTIAL  = "MADV_SEQUENTIAL"
	ALIAS_OS_MADV_WILLNEED    = "MADV_WILLNEED"
	ALIAS_OS_MADV_DONTNEED    = "MADV_DONTNEED"
	ALIAS_OS_MADV_FREE        = "MADV_FREE"
	ALIAS_OS_MADV_DONTFORK    = "MADV_DONTFORK"
	ALIAS_OS_MADV_DOFORK      = "MADV_DOFORK"
	ALIAS_OS_MADV_MERGEABLE   = "MADV_MERGEABLE"
	ALIAS_OS_MADV_UNMERGEABLE = "MADV_UNMERGEABLE"
	ALIAS_OS_MADV_HUGEPAGE    = "MADV_HUGEPAGE"
	ALIAS_OS_MADV_NOHUGEPAGE  = "MADV_NOHUGEPAGE"

	ALIAS_OS_MCL_CURRENT = "MCL_CURRENT"
	ALIAS_OS_MCL_FUTURE  = "MCL_FUTURE"

	ALIAS_OS_EFAULT = "EFAULT"
	ALIAS_OS_EACCES = "EACCES"

	CONSTANT_OS_EAGAIN = "OS_EAGAIN"
	CONSTANT_OS_EINTR  = "OS_EINTR"
	CONSTANT_OS_EINVAL = "OS_EINVAL"
	CONSTANT_OS_EPERM  = "OS_EPERM"
	CONSTANT_OS_ENOENT = "OS_ENOENT"
	CONSTANT_OS_ENOMEM = "OS_ENOMEM"

	CONSTANT_OS_WNOHANG    = "OS_WNOHANG"
	CONSTANT_OS_WUNTRACED  = "OS_WUNTRACED"
	CONSTANT_OS_WCONTINUED = "OS_WCONTINUED"

	CONSTANT_OS_SIGINT  = "OS_SIGINT"
	CONSTANT_OS_SIGTERM = "OS_SIGTERM"
	CONSTANT_OS_SIGKILL = "OS_SIGKILL"
	CONSTANT_OS_SIGSEGV = "OS_SIGSEGV"
	CONSTANT_OS_SIGABRT = "OS_SIGABRT"
	CONSTANT_OS_SIGCHLD = "OS_SIGCHLD"
	CONSTANT_OS_SIGPIPE = "OS_SIGPIPE"
	CONSTANT_OS_SIGALRM = "OS_SIGALRM"
	CONSTANT_OS_SIGUSR1 = "OS_SIGUSR1"
	CONSTANT_OS_SIGUSR2 = "OS_SIGUSR2"

	CONSTANT_OS_RLIMIT_CPU    = "OS_RLIMIT_CPU"
	CONSTANT_OS_RLIMIT_FSIZE  = "OS_RLIMIT_FSIZE"
	CONSTANT_OS_RLIMIT_DATA   = "OS_RLIMIT_DATA"
	CONSTANT_OS_RLIMIT_STACK  = "OS_RLIMIT_STACK"
	CONSTANT_OS_RLIMIT_CORE   = "OS_RLIMIT_CORE"
	CONSTANT_OS_RLIMIT_NOFILE = "OS_RLIMIT_NOFILE"
	CONSTANT_OS_RLIMIT_AS     = "OS_RLIMIT_AS"

	/* Standard file descriptor numbers */
	/** Standard input */
	CONSTANT_OS_STDIN_FD = "OS_0"
	/** Standard output */
	CONSTANT_OS_STDOUT_FD = "OS_1"
	/** Standard error */
	CONSTANT_OS_STDERR_FD = "OS_2"

	/*open() flags*/
	CONSTANT_OS_O_RDONLY   = "OS_O_RDONLY"
	CONSTANT_OS_O_WRONLY   = "OS_O_WRONLY"
	CONSTANT_OS_O_RDWR     = "OS_O_RDWR"
	CONSTANT_OS_O_APPEND   = "OS_O_APPEND"
	CONSTANT_OS_O_CREAT    = "OS_O_CREAT"
	CONSTANT_OS_O_EXCL     = "OS_O_EXCL"
	CONSTANT_OS_O_TRUNC    = "OS_O_TRUNC"
	CONSTANT_OS_O_CLOEXEC  = "OS_O_CLOEXEC"
	CONSTANT_OS_O_NONBLOCK = "OS_O_NONBLOCK"
	CONSTANT_OS_O_SYNC     = "OS_O_SYNC"
	CONSTANT_OS_O_DSYNC    = "OS_O_DSYNC"
	CONSTANT_OS_O_DIRECT   = "OS_O_DIRECT"

	/*seek constants*/
	CONSTANT_OS_SEEK_SET = "OS_SEEK_SET"
	CONSTANT_OS_SEEK_CUR = "OS_SEEK_CUR"
	CONSTANT_OS_SEEK_END = "OS_SEEK_END"

	/*fcntl commands*/
	CONSTANT_OS_F_DUPFD         = "OS_F_DUPFD"
	CONSTANT_OS_F_DUPFD_CLOEXEC = "OS_F_DUPFD_CLOEXEC"
	CONSTANT_OS_F_GETFD         = "OS_F_GETFD"
	CONSTANT_OS_F_SETFD         = "OS_F_SETFD"
	CONSTANT_OS_F_GETFL         = "OS_F_GETFL"
	CONSTANT_OS_F_SETFL         = "OS_F_SETFL"

	/*FD flags*/
	CONSTANT_OS_FD_CLOEXEC = "OS_FD_CLOEXEC"

	/*stat mode bits*/
	CONSTANT_OS_S_IFREG  = "OS_S_IFREG"
	CONSTANT_OS_S_IFDIR  = "OS_S_IFDIR"
	CONSTANT_OS_S_IFCHR  = "OS_S_IFCHR"
	CONSTANT_OS_S_IFBLK  = "OS_S_IFBLK"
	CONSTANT_OS_S_IFIFO  = "OS_S_IFIFO"
	CONSTANT_OS_S_IFLNK  = "OS_S_IFLNK"
	CONSTANT_OS_S_IFSOCK = "OS_S_IFSOCK"

	CONSTANT_OS_S_IRUSR = "OS_S_IRUSR"
	CONSTANT_OS_S_IWUSR = "OS_S_IWUSR"
	CONSTANT_OS_S_IXUSR = "OS_S_IXUSR"
	CONSTANT_OS_S_IRGRP = "OS_S_IRGRP"
	CONSTANT_OS_S_IWGRP = "OS_S_IWGRP"
	CONSTANT_OS_S_IXGRP = "OS_S_IXGRP"
	CONSTANT_OS_S_IROTH = "OS_S_IROTH"
	CONSTANT_OS_S_IWOTH = "OS_S_IWOTH"
	CONSTANT_OS_S_IXOTH = "OS_S_IXOTH"

	/* Errors (FD-relevant subset)*/
	CONSTANT_OS_EBADF  = "OS_EBADF"
	CONSTANT_OS_EPIPE  = "OS_EPIPE"
	CONSTANT_OS_EIO    = "OS_EIO"
	CONSTANT_OS_ENOSPC = "OS_ENOSPC"

	/*Special directory FDs*/
	/** Current working directory */
	CONSTANT_OS_AT_FDCWD = "OS_AT_FDCWD"

	/*unlinkat / renameat flags*/
	CONSTANT_OS_AT_REMOVEDIR      = "OS_AT_REMOVEDIR"
	CONSTANT_OS_AT_SYMLINK_FOLLOW = "OS_AT_SYMLINK_FOLLOW"

	/*link / rename flags*/
	CONSTANT_OS_RENAME_NOREPLACE = "OS_RENAME_NOREPLACE" // "RENAME_NOREPLACE"
	CONSTANT_OS_RENAME_EXCHANGE  = "OS_RENAME_EXCHANGE"  // "RENAME_EXCHANGE"
	CONSTANT_OS_RENAME_WHITEOUT  = "OS_RENAME_WHITEOUT"  // "RENAME_WHITEOUT"

	/*Access mode flags*/
	CONSTANT_OS_F_OK = "OS_F_OK"
	CONSTANT_OS_R_OK = "OS_R_OK"
	CONSTANT_OS_W_OK = "OS_W_OK"
	CONSTANT_OS_X_OK = "OS_X_OK"

	/*Directory entry types (d_type)*/
	CONSTANT_OS_DT_UNKNOWN = "OS_DT_UNKNOWN"
	CONSTANT_OS_DT_FIFO    = "OS_DT_FIFO"
	CONSTANT_OS_DT_CHR     = "OS_DT_CHR"
	CONSTANT_OS_DT_DIR     = "OS_DT_DIR"
	CONSTANT_OS_DT_BLK     = "OS_DT_BLK"
	CONSTANT_OS_DT_REG     = "OS_DT_REG"
	CONSTANT_OS_DT_LNK     = "OS_DT_LNK"
	CONSTANT_OS_DT_SOCK    = "OS_DT_SOCK"
	CONSTANT_OS_DT_WHT     = "OS_DT_WHT"

	/*Memory protection flags*/
	CONSTANT_OS_PROT_NONE  = "OS_PROT_NONE"
	CONSTANT_OS_PROT_READ  = "OS_PROT_READ"
	CONSTANT_OS_PROT_WRITE = "OS_PROT_WRITE"
	CONSTANT_OS_PROT_EXEC  = "OS_PROT_EXEC"

	/*mmap flags*/
	CONSTANT_OS_MAP_SHARED    = "OS_MAP_SHARED"
	CONSTANT_OS_MAP_PRIVATE   = "OS_MAP_PRIVATE"
	CONSTANT_OS_MAP_FIXED     = "OS_MAP_FIXED"
	CONSTANT_OS_MAP_ANONYMOUS = "OS_MAP_ANONYMOUS"
	CONSTANT_OS_MAP_STACK     = "OS_MAP_STACK"
	CONSTANT_OS_MAP_NORESERVE = "OS_MAP_NORESERVE"
	CONSTANT_OS_MAP_POPULATE  = "OS_MAP_POPULATE"
	CONSTANT_OS_MAP_GROWSDOWN = "OS_MAP_GROWSDOWN"

	/*madvise advice*/
	CONSTANT_OS_MADV_NORMAL      = "OS_MADV_NORMAL"
	CONSTANT_OS_MADV_RANDOM      = "OS_MADV_RANDOM"
	CONSTANT_OS_MADV_SEQUENTIAL  = "OS_MADV_SEQUENTIAL"
	CONSTANT_OS_MADV_WILLNEED    = "OS_MADV_WILLNEED"
	CONSTANT_OS_MADV_DONTNEED    = "OS_MADV_DONTNEED"
	CONSTANT_OS_MADV_FREE        = "OS_MADV_FREE"
	CONSTANT_OS_MADV_DONTFORK    = "OS_MADV_DONTFORK"
	CONSTANT_OS_MADV_DOFORK      = "OS_MADV_DOFORK"
	CONSTANT_OS_MADV_MERGEABLE   = "OS_MADV_MERGEABLE"
	CONSTANT_OS_MADV_UNMERGEABLE = "OS_MADV_UNMERGEABLE"
	CONSTANT_OS_MADV_HUGEPAGE    = "OS_MADV_HUGEPAGE"
	CONSTANT_OS_MADV_NOHUGEPAGE  = "OS_MADV_NOHUGEPAGE"

	CONSTANT_OS_MCL_CURRENT = "OS_MCL_CURRENT"
	CONSTANT_OS_MCL_FUTURE  = "OS_MCL_FUTURE"

	CONSTANT_OS_EFAULT = "OS_EFAULT"
	CONSTANT_OS_EACCES = "OS_EACCES"

	FUNC_NET_LISTEN = "__public__net_listen"
	FUNC_NET_ACCEPT = "__public__net_accept"
	FUNC_NET_READ   = "__public__net_read"
	FUNC_NET_WRITE  = "__public__net_write"

	ALIAS_NET_LISTEN = "listen"
	ALIAS_NET_ACCEPT = "accept"
	ALIAS_NET_READ   = "read"
	ALIAS_NET_WRITE  = "write"
)

// Interface maintains a registry of available external functions and
// runtime types. It provides a centralized lookup table for the code
// generator to reference LLVM symbols.
type Interface struct {
	// Funcs maps symbol names to their corresponding LLVM IR function declarations.
	Funcs map[string]*ir.Func
	// Types maps type identifiers to their concrete LLVM IR type definitions.
	Types map[string]types.Type
	// constants
	Constants map[string]*ir.Global
}

// Instance is a global singleton providing access to the C runtime interface.
var Instance *Interface

// NewInterface initializes a new runtime registry for the given LLVM module.
// It populates the internal maps by registering all required external
// functions and built-in types.
func NewInterface(mod *ir.Module) *Interface {
	t := &Interface{}
	t.Funcs = make(map[string]*ir.Func)
	t.Types = make(map[string]types.Type)
	t.Constants = make(map[string]*ir.Global)

	t.registerTypes(mod)
	t.registerFuncs(mod)
	t.registerConstants(mod)

	Instance = t
	return Instance
}
