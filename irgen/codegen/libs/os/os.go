package _os

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	function "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/libs/libutils"
	typedef "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

type OSHandler struct {
}

func NewOSHandler() *OSHandler {
	return &OSHandler{}
}

func (t *OSHandler) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)

	funcs[c.ALIAS_SYSCALL_ERRNO] = t.errno
	funcs[c.ALIAS_SYSCALL_GETPID] = t.getpid
	funcs[c.ALIAS_SYSCALL_GETPPID] = t.getppid
	funcs[c.ALIAS_SYSCALL_GETTID] = t.gettid
	funcs[c.ALIAS_SYSCALL_EXIT] = t.exit
	funcs[c.ALIAS_SYSCALL_FORK] = t.fork
	funcs[c.ALIAS_SYSCALL_WAITPID] = t.waitpid
	funcs[c.ALIAS_SYSCALL_KILL] = t.kill
	funcs[c.ALIAS_SYSCALL_EXECVE] = t.execve
	funcs[c.ALIAS_SYSCALL_EXECVP] = t.execvp
	funcs[c.ALIAS_SYSCALL_ENVIRON] = t.environ
	funcs[c.ALIAS_SYSCALL_GETENV] = t.getenv
	funcs[c.ALIAS_SYSCALL_SETENV] = t.setenv
	funcs[c.ALIAS_SYSCALL_UNSETENV] = t.unsetenv
	funcs[c.ALIAS_SYSCALL_GETCWD] = t.getcwd
	funcs[c.ALIAS_SYSCALL_CHDIR] = t.chdir
	funcs[c.ALIAS_SYSCALL_GETUID] = t.getuid
	funcs[c.ALIAS_SYSCALL_GETEUID] = t.geteuid
	funcs[c.ALIAS_SYSCALL_GETGID] = t.getgid
	funcs[c.ALIAS_SYSCALL_GETEGID] = t.getegid
	funcs[c.ALIAS_SYSCALL_SETUID] = t.setuid
	funcs[c.ALIAS_SYSCALL_SETGID] = t.setgid
	funcs[c.ALIAS_SYSCALL_SETPGID] = t.setpgid
	funcs[c.ALIAS_SYSCALL_GETPGID] = t.getpgid
	funcs[c.ALIAS_SYSCALL_GETPGRP] = t.getpgrp
	funcs[c.ALIAS_SYSCALL_SETSID] = t.setsid
	funcs[c.ALIAS_SYSCALL_GETRLIMIT] = t.getrlimit
	funcs[c.ALIAS_SYSCALL_SETRLIMIT] = t.setrlimit
	funcs[c.ALIAS_SYSCALL_SIGNAL_INSTALL] = t.signal_install
	funcs[c.ALIAS_SYSCALL_OPEN] = t.open
	funcs[c.ALIAS_SYSCALL_CLOSE] = t.close
	funcs[c.ALIAS_SYSCALL_READ] = t.read
	funcs[c.ALIAS_SYSCALL_WRITE] = t.write
	funcs[c.ALIAS_SYSCALL_LSEEK] = t.lseek
	funcs[c.ALIAS_SYSCALL_FSTAT] = t.fstat
	funcs[c.ALIAS_SYSCALL_DUP] = t.dup
	funcs[c.ALIAS_SYSCALL_DUP2] = t.dup2
	funcs[c.ALIAS_SYSCALL_FCNTL] = t.fcntl
	funcs[c.ALIAS_SYSCALL_MKDIR] = t.mkdir
	funcs[c.ALIAS_SYSCALL_RMDIR] = t.rmdir
	funcs[c.ALIAS_SYSCALL_UNLINK] = t.unlink
	funcs[c.ALIAS_SYSCALL_RENAME] = t.rename
	funcs[c.ALIAS_SYSCALL_RENAMEAT2] = t.renameat2
	funcs[c.ALIAS_SYSCALL_LINK] = t.link
	funcs[c.ALIAS_SYSCALL_SYMLINK] = t.symlink
	funcs[c.ALIAS_SYSCALL_READLINK] = t.readlink
	funcs[c.ALIAS_SYSCALL_STAT] = t.stat
	funcs[c.ALIAS_SYSCALL_LSTAT] = t.lstat
	funcs[c.ALIAS_SYSCALL_ACCESS] = t.access
	funcs[c.ALIAS_SYSCALL_GETDENTS64] = t.getdents64
	funcs[c.ALIAS_SYSCALL_MMAP] = t.mmap
	funcs[c.ALIAS_SYSCALL_MUNMAP] = t.munmap
	funcs[c.ALIAS_SYSCALL_MPROTECT] = t.mprotect
	funcs[c.ALIAS_SYSCALL_MADVISE] = t.madvise
	funcs[c.ALIAS_SYSCALL_MLOCK] = t.mlock
	funcs[c.ALIAS_SYSCALL_MUNLOCK] = t.munlock
	funcs[c.ALIAS_SYSCALL_MLOCKALL] = t.mlockall
	funcs[c.ALIAS_SYSCALL_MUNLOCKALL] = t.munlockall
	funcs[c.ALIAS_SYSCALL_PAGE_SIZE] = t.page_size

	return funcs
}

func (t *OSHandler) errno(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_ERRNO]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getpid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETPID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getppid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETPPID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) gettid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETTID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) exit(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_EXIT]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) fork(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_FORK]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) waitpid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_WAITPID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) kill(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_KILL]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) execve(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_EXECVE]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) execvp(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_EXECVP]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) environ(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_ENVIRON]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getenv(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETENV]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) setenv(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_SETENV]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) unsetenv(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_UNSETENV]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getcwd(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETCWD]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) chdir(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_CHDIR]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getuid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETUID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) geteuid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETEUID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getgid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETGID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getegid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETEGID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) setuid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_SETUID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) setgid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_SETGID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) setpgid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_SETPGID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getpgid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETPGID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getpgrp(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETPGRP]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) setsid(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_SETSID]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getrlimit(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETRLIMIT]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) setrlimit(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_SETRLIMIT]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) signal_install(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_SIGNAL_INSTALL]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) open(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_OPEN]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) close(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_CLOSE]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) read(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_READ]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) write(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_WRITE]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) lseek(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_LSEEK]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) fstat(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_FSTAT]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) dup(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_DUP]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) dup2(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_DUP2]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) fcntl(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_FCNTL]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) mkdir(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MKDIR]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) rmdir(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_RMDIR]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) unlink(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_UNLINK]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) rename(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_RENAME]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) renameat2(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_RENAMEAT2]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) link(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_LINK]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) symlink(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_SYMLINK]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) readlink(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_READLINK]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) stat(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_STAT]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) lstat(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_LSTAT]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) access(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_ACCESS]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) getdents64(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_GETDENTS64]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) mmap(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MMAP]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) munmap(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MUNMAP]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) mprotect(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MPROTECT]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) madvise(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MADVISE]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) mlock(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MLOCK]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) munlock(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MUNLOCK]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) mlockall(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MLOCKALL]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) munlockall(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_MUNLOCKALL]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}

func (t *OSHandler) page_size(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fn := c.Instance.Funcs[c.FUNC_SYSCALL_PAGE_SIZE]
	return libutils.CallCFunc(typeHandler, fn, bh, args)
}
