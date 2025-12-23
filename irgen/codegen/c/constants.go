package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func (t *Interface) registerConstants(mod *ir.Module) {
	t.declareOSConstants(mod)
}

func declareExternInt(mod *ir.Module, name string) *ir.Global {
	g := ir.NewGlobal(name, types.I32)
	g.Linkage = enum.LinkageExternal
	mod.Globals = append(mod.Globals, g)
	return g
}

func (t *Interface) declareOSConstants(mod *ir.Module) {
	// errno / basic errors
	t.Constants[CONSTANT_OS_EAGAIN] = declareExternInt(mod, CONSTANT_OS_EAGAIN)
	t.Constants[CONSTANT_OS_EINTR] = declareExternInt(mod, CONSTANT_OS_EINTR)
	t.Constants[CONSTANT_OS_EINVAL] = declareExternInt(mod, CONSTANT_OS_EINVAL)
	t.Constants[CONSTANT_OS_EPERM] = declareExternInt(mod, CONSTANT_OS_EPERM)
	t.Constants[CONSTANT_OS_ENOENT] = declareExternInt(mod, CONSTANT_OS_ENOENT)
	t.Constants[CONSTANT_OS_ENOMEM] = declareExternInt(mod, CONSTANT_OS_ENOMEM)
	t.Constants[CONSTANT_OS_EBADF] = declareExternInt(mod, CONSTANT_OS_EBADF)
	t.Constants[CONSTANT_OS_EPIPE] = declareExternInt(mod, CONSTANT_OS_EPIPE)
	t.Constants[CONSTANT_OS_EIO] = declareExternInt(mod, CONSTANT_OS_EIO)
	t.Constants[CONSTANT_OS_ENOSPC] = declareExternInt(mod, CONSTANT_OS_ENOSPC)
	t.Constants[CONSTANT_OS_EFAULT] = declareExternInt(mod, CONSTANT_OS_EFAULT)
	t.Constants[CONSTANT_OS_EACCES] = declareExternInt(mod, CONSTANT_OS_EACCES)

	// wait() flags
	t.Constants[CONSTANT_OS_WNOHANG] = declareExternInt(mod, CONSTANT_OS_WNOHANG)
	t.Constants[CONSTANT_OS_WUNTRACED] = declareExternInt(mod, CONSTANT_OS_WUNTRACED)
	t.Constants[CONSTANT_OS_WCONTINUED] = declareExternInt(mod, CONSTANT_OS_WCONTINUED)

	// signals
	t.Constants[CONSTANT_OS_SIGINT] = declareExternInt(mod, CONSTANT_OS_SIGINT)
	t.Constants[CONSTANT_OS_SIGTERM] = declareExternInt(mod, CONSTANT_OS_SIGTERM)
	t.Constants[CONSTANT_OS_SIGKILL] = declareExternInt(mod, CONSTANT_OS_SIGKILL)
	t.Constants[CONSTANT_OS_SIGSEGV] = declareExternInt(mod, CONSTANT_OS_SIGSEGV)
	t.Constants[CONSTANT_OS_SIGABRT] = declareExternInt(mod, CONSTANT_OS_SIGABRT)
	t.Constants[CONSTANT_OS_SIGCHLD] = declareExternInt(mod, CONSTANT_OS_SIGCHLD)
	t.Constants[CONSTANT_OS_SIGPIPE] = declareExternInt(mod, CONSTANT_OS_SIGPIPE)
	t.Constants[CONSTANT_OS_SIGALRM] = declareExternInt(mod, CONSTANT_OS_SIGALRM)
	t.Constants[CONSTANT_OS_SIGUSR1] = declareExternInt(mod, CONSTANT_OS_SIGUSR1)
	t.Constants[CONSTANT_OS_SIGUSR2] = declareExternInt(mod, CONSTANT_OS_SIGUSR2)

	// rlimit
	t.Constants[CONSTANT_OS_RLIMIT_CPU] = declareExternInt(mod, CONSTANT_OS_RLIMIT_CPU)
	t.Constants[CONSTANT_OS_RLIMIT_FSIZE] = declareExternInt(mod, CONSTANT_OS_RLIMIT_FSIZE)
	t.Constants[CONSTANT_OS_RLIMIT_DATA] = declareExternInt(mod, CONSTANT_OS_RLIMIT_DATA)
	t.Constants[CONSTANT_OS_RLIMIT_STACK] = declareExternInt(mod, CONSTANT_OS_RLIMIT_STACK)
	t.Constants[CONSTANT_OS_RLIMIT_CORE] = declareExternInt(mod, CONSTANT_OS_RLIMIT_CORE)
	t.Constants[CONSTANT_OS_RLIMIT_NOFILE] = declareExternInt(mod, CONSTANT_OS_RLIMIT_NOFILE)
	t.Constants[CONSTANT_OS_RLIMIT_AS] = declareExternInt(mod, CONSTANT_OS_RLIMIT_AS)

	// stdio fds
	t.Constants[CONSTANT_OS_STDIN_FD] = declareExternInt(mod, CONSTANT_OS_STDIN_FD)
	t.Constants[CONSTANT_OS_STDOUT_FD] = declareExternInt(mod, CONSTANT_OS_STDOUT_FD)
	t.Constants[CONSTANT_OS_STDERR_FD] = declareExternInt(mod, CONSTANT_OS_STDERR_FD)

	// open flags
	t.Constants[CONSTANT_OS_O_RDONLY] = declareExternInt(mod, CONSTANT_OS_O_RDONLY)
	t.Constants[CONSTANT_OS_O_WRONLY] = declareExternInt(mod, CONSTANT_OS_O_WRONLY)
	t.Constants[CONSTANT_OS_O_RDWR] = declareExternInt(mod, CONSTANT_OS_O_RDWR)
	t.Constants[CONSTANT_OS_O_APPEND] = declareExternInt(mod, CONSTANT_OS_O_APPEND)
	t.Constants[CONSTANT_OS_O_CREAT] = declareExternInt(mod, CONSTANT_OS_O_CREAT)
	t.Constants[CONSTANT_OS_O_EXCL] = declareExternInt(mod, CONSTANT_OS_O_EXCL)
	t.Constants[CONSTANT_OS_O_TRUNC] = declareExternInt(mod, CONSTANT_OS_O_TRUNC)
	t.Constants[CONSTANT_OS_O_CLOEXEC] = declareExternInt(mod, CONSTANT_OS_O_CLOEXEC)
	t.Constants[CONSTANT_OS_O_NONBLOCK] = declareExternInt(mod, CONSTANT_OS_O_NONBLOCK)
	t.Constants[CONSTANT_OS_O_SYNC] = declareExternInt(mod, CONSTANT_OS_O_SYNC)
	t.Constants[CONSTANT_OS_O_DSYNC] = declareExternInt(mod, CONSTANT_OS_O_DSYNC)
	t.Constants[CONSTANT_OS_O_DIRECT] = declareExternInt(mod, CONSTANT_OS_O_DIRECT)

	// seek
	t.Constants[CONSTANT_OS_SEEK_SET] = declareExternInt(mod, CONSTANT_OS_SEEK_SET)
	t.Constants[CONSTANT_OS_SEEK_CUR] = declareExternInt(mod, CONSTANT_OS_SEEK_CUR)
	t.Constants[CONSTANT_OS_SEEK_END] = declareExternInt(mod, CONSTANT_OS_SEEK_END)

	// fcntl
	t.Constants[CONSTANT_OS_F_DUPFD] = declareExternInt(mod, CONSTANT_OS_F_DUPFD)
	t.Constants[CONSTANT_OS_F_DUPFD_CLOEXEC] = declareExternInt(mod, CONSTANT_OS_F_DUPFD_CLOEXEC)
	t.Constants[CONSTANT_OS_F_GETFD] = declareExternInt(mod, CONSTANT_OS_F_GETFD)
	t.Constants[CONSTANT_OS_F_SETFD] = declareExternInt(mod, CONSTANT_OS_F_SETFD)
	t.Constants[CONSTANT_OS_F_GETFL] = declareExternInt(mod, CONSTANT_OS_F_GETFL)
	t.Constants[CONSTANT_OS_F_SETFL] = declareExternInt(mod, CONSTANT_OS_F_SETFL)
	t.Constants[CONSTANT_OS_FD_CLOEXEC] = declareExternInt(mod, CONSTANT_OS_FD_CLOEXEC)

	// stat types
	t.Constants[CONSTANT_OS_S_IFREG] = declareExternInt(mod, CONSTANT_OS_S_IFREG)
	t.Constants[CONSTANT_OS_S_IFDIR] = declareExternInt(mod, CONSTANT_OS_S_IFDIR)
	t.Constants[CONSTANT_OS_S_IFCHR] = declareExternInt(mod, CONSTANT_OS_S_IFCHR)
	t.Constants[CONSTANT_OS_S_IFBLK] = declareExternInt(mod, CONSTANT_OS_S_IFBLK)
	t.Constants[CONSTANT_OS_S_IFIFO] = declareExternInt(mod, CONSTANT_OS_S_IFIFO)
	t.Constants[CONSTANT_OS_S_IFLNK] = declareExternInt(mod, CONSTANT_OS_S_IFLNK)
	t.Constants[CONSTANT_OS_S_IFSOCK] = declareExternInt(mod, CONSTANT_OS_S_IFSOCK)

	// permissions
	t.Constants[CONSTANT_OS_S_IRUSR] = declareExternInt(mod, CONSTANT_OS_S_IRUSR)
	t.Constants[CONSTANT_OS_S_IWUSR] = declareExternInt(mod, CONSTANT_OS_S_IWUSR)
	t.Constants[CONSTANT_OS_S_IXUSR] = declareExternInt(mod, CONSTANT_OS_S_IXUSR)
	t.Constants[CONSTANT_OS_S_IRGRP] = declareExternInt(mod, CONSTANT_OS_S_IRGRP)
	t.Constants[CONSTANT_OS_S_IWGRP] = declareExternInt(mod, CONSTANT_OS_S_IWGRP)
	t.Constants[CONSTANT_OS_S_IXGRP] = declareExternInt(mod, CONSTANT_OS_S_IXGRP)
	t.Constants[CONSTANT_OS_S_IROTH] = declareExternInt(mod, CONSTANT_OS_S_IROTH)
	t.Constants[CONSTANT_OS_S_IWOTH] = declareExternInt(mod, CONSTANT_OS_S_IWOTH)
	t.Constants[CONSTANT_OS_S_IXOTH] = declareExternInt(mod, CONSTANT_OS_S_IXOTH)

	// *at() and rename
	t.Constants[CONSTANT_OS_AT_FDCWD] = declareExternInt(mod, CONSTANT_OS_AT_FDCWD)
	t.Constants[CONSTANT_OS_AT_REMOVEDIR] = declareExternInt(mod, CONSTANT_OS_AT_REMOVEDIR)
	t.Constants[CONSTANT_OS_AT_SYMLINK_FOLLOW] = declareExternInt(mod, CONSTANT_OS_AT_SYMLINK_FOLLOW)
	t.Constants[CONSTANT_OS_RENAME_NOREPLACE] = declareExternInt(mod, CONSTANT_OS_RENAME_NOREPLACE)
	t.Constants[CONSTANT_OS_RENAME_EXCHANGE] = declareExternInt(mod, CONSTANT_OS_RENAME_EXCHANGE)
	t.Constants[CONSTANT_OS_RENAME_WHITEOUT] = declareExternInt(mod, CONSTANT_OS_RENAME_WHITEOUT)

	// access()
	t.Constants[CONSTANT_OS_F_OK] = declareExternInt(mod, CONSTANT_OS_F_OK)
	t.Constants[CONSTANT_OS_R_OK] = declareExternInt(mod, CONSTANT_OS_R_OK)
	t.Constants[CONSTANT_OS_W_OK] = declareExternInt(mod, CONSTANT_OS_W_OK)
	t.Constants[CONSTANT_OS_X_OK] = declareExternInt(mod, CONSTANT_OS_X_OK)

	// dirent d_type
	t.Constants[CONSTANT_OS_DT_UNKNOWN] = declareExternInt(mod, CONSTANT_OS_DT_UNKNOWN)
	t.Constants[CONSTANT_OS_DT_FIFO] = declareExternInt(mod, CONSTANT_OS_DT_FIFO)
	t.Constants[CONSTANT_OS_DT_CHR] = declareExternInt(mod, CONSTANT_OS_DT_CHR)
	t.Constants[CONSTANT_OS_DT_DIR] = declareExternInt(mod, CONSTANT_OS_DT_DIR)
	t.Constants[CONSTANT_OS_DT_BLK] = declareExternInt(mod, CONSTANT_OS_DT_BLK)
	t.Constants[CONSTANT_OS_DT_REG] = declareExternInt(mod, CONSTANT_OS_DT_REG)
	t.Constants[CONSTANT_OS_DT_LNK] = declareExternInt(mod, CONSTANT_OS_DT_LNK)
	t.Constants[CONSTANT_OS_DT_SOCK] = declareExternInt(mod, CONSTANT_OS_DT_SOCK)
	t.Constants[CONSTANT_OS_DT_WHT] = declareExternInt(mod, CONSTANT_OS_DT_WHT)

	// mmap / mprotect
	t.Constants[CONSTANT_OS_PROT_NONE] = declareExternInt(mod, CONSTANT_OS_PROT_NONE)
	t.Constants[CONSTANT_OS_PROT_READ] = declareExternInt(mod, CONSTANT_OS_PROT_READ)
	t.Constants[CONSTANT_OS_PROT_WRITE] = declareExternInt(mod, CONSTANT_OS_PROT_WRITE)
	t.Constants[CONSTANT_OS_PROT_EXEC] = declareExternInt(mod, CONSTANT_OS_PROT_EXEC)

	t.Constants[CONSTANT_OS_MAP_SHARED] = declareExternInt(mod, CONSTANT_OS_MAP_SHARED)
	t.Constants[CONSTANT_OS_MAP_PRIVATE] = declareExternInt(mod, CONSTANT_OS_MAP_PRIVATE)
	t.Constants[CONSTANT_OS_MAP_FIXED] = declareExternInt(mod, CONSTANT_OS_MAP_FIXED)
	t.Constants[CONSTANT_OS_MAP_ANONYMOUS] = declareExternInt(mod, CONSTANT_OS_MAP_ANONYMOUS)
	t.Constants[CONSTANT_OS_MAP_STACK] = declareExternInt(mod, CONSTANT_OS_MAP_STACK)
	t.Constants[CONSTANT_OS_MAP_NORESERVE] = declareExternInt(mod, CONSTANT_OS_MAP_NORESERVE)
	t.Constants[CONSTANT_OS_MAP_POPULATE] = declareExternInt(mod, CONSTANT_OS_MAP_POPULATE)
	t.Constants[CONSTANT_OS_MAP_GROWSDOWN] = declareExternInt(mod, CONSTANT_OS_MAP_GROWSDOWN)

	// madvise
	t.Constants[CONSTANT_OS_MADV_NORMAL] = declareExternInt(mod, CONSTANT_OS_MADV_NORMAL)
	t.Constants[CONSTANT_OS_MADV_RANDOM] = declareExternInt(mod, CONSTANT_OS_MADV_RANDOM)
	t.Constants[CONSTANT_OS_MADV_SEQUENTIAL] = declareExternInt(mod, CONSTANT_OS_MADV_SEQUENTIAL)
	t.Constants[CONSTANT_OS_MADV_WILLNEED] = declareExternInt(mod, CONSTANT_OS_MADV_WILLNEED)
	t.Constants[CONSTANT_OS_MADV_DONTNEED] = declareExternInt(mod, CONSTANT_OS_MADV_DONTNEED)
	t.Constants[CONSTANT_OS_MADV_FREE] = declareExternInt(mod, CONSTANT_OS_MADV_FREE)
	t.Constants[CONSTANT_OS_MADV_DONTFORK] = declareExternInt(mod, CONSTANT_OS_MADV_DONTFORK)
	t.Constants[CONSTANT_OS_MADV_DOFORK] = declareExternInt(mod, CONSTANT_OS_MADV_DOFORK)
	t.Constants[CONSTANT_OS_MADV_MERGEABLE] = declareExternInt(mod, CONSTANT_OS_MADV_MERGEABLE)
	t.Constants[CONSTANT_OS_MADV_UNMERGEABLE] = declareExternInt(mod, CONSTANT_OS_MADV_UNMERGEABLE)
	t.Constants[CONSTANT_OS_MADV_HUGEPAGE] = declareExternInt(mod, CONSTANT_OS_MADV_HUGEPAGE)
	t.Constants[CONSTANT_OS_MADV_NOHUGEPAGE] = declareExternInt(mod, CONSTANT_OS_MADV_NOHUGEPAGE)

	// mlock
	t.Constants[CONSTANT_OS_MCL_CURRENT] = declareExternInt(mod, CONSTANT_OS_MCL_CURRENT)
	t.Constants[CONSTANT_OS_MCL_FUTURE] = declareExternInt(mod, CONSTANT_OS_MCL_FUTURE)
}
