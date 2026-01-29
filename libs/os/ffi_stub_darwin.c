#include "platform/darwin/os.h"

/**
 * Force references to all public OS wrapper functions so Clang/LLVM 
 * emits the necessary declarations and symbols for FFI discovery.
 */
void* __ffi_force[] = {
    (void*)__public__os_errno,
    (void*)__public__os_getpid,
    (void*)__public__os_getppid,
    (void*)__public__os_gettid,
    (void*)__public__os_exit,
    (void*)__public__os_fork,
    (void*)__public__os_waitpid,
    (void*)__public__os_kill,
    (void*)__public__os_execve,
    (void*)__public__os_execvp,
    (void*)__public__os_environ,
    (void*)__public__os_getenv,
    (void*)__public__os_setenv,
    (void*)__public__os_unsetenv,
    (void*)__public__os_getcwd,
    (void*)__public__os_chdir,
    (void*)__public__os_getuid,
    (void*)__public__os_geteuid,
    (void*)__public__os_getgid,
    (void*)__public__os_getegid,
    (void*)__public__os_setuid,
    (void*)__public__os_setgid,
    (void*)__public__os_open,
    (void*)__public__os_close,
    (void*)__public__os_read,
    (void*)__public__os_write,
    (void*)__public__os_lseek,
    (void*)__public__os_dup,
    (void*)__public__os_dup2,
    (void*)__public__os_fcntl,
    (void*)__public__os_mmap,
    (void*)__public__os_munmap
};