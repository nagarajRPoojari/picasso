/**
 * This file exposes minimal set of linux specific syscalls.
 *
 * All functions are thin syscall wrappers.
 * No retries, no buffering, no allocation policies.
 *
 * Linux-specific implementation.
 */

#define _GNU_SOURCE
#include <stdint.h>
#include <stddef.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/syscall.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/mman.h>
#include <sys/resource.h>
#include <sched.h>
#include <signal.h>
#include <string.h>
#include <errno.h>
#include <stdlib.h>
#include "os.h"

extern char **environ;

/**
 * @brief Get the current thread-local errno value.
 *
 * @return Current errno.
 */
int __public__errno(void) {
    return errno;
}

/**
 * @brief Get current process ID.
 *
 * @return Process ID.
 */
int __public__getpid(void) {
    return syscall(SYS_getpid);
}

/**
 * @brief Get parent process ID.
 *
 * @return Parent process ID.
 */
int __public__getppid(void) {
    return syscall(SYS_getppid);
}

/**
 * @brief Get calling thread ID.
 *
 * @return Thread ID.
 */
int __public__gettid(void) {
    return syscall(SYS_gettid);
}

/**
 * @brief Terminate the current process immediately.
 *
 * @param code Exit status.
 */
void __public__exit(int code) {
    syscall(SYS_exit, code);
    __builtin_unreachable();
}

/**
 * @brief Create a new process.
 *
 * @return 0 in child, child PID in parent, or -1 on error.
 */
int __public__fork(void) {
    return syscall(SYS_clone, SIGCHLD, 0, 0, 0, 0);
}

/**
 * @brief Wait for process state change.
 *
 * @param pid     Process ID to wait for.
 * @param status  Pointer to store exit status.
 * @param options wait options.
 *
 * @return PID of child or -1 on error.
 */
int __public__waitpid(int pid, int *status, int options) {
    return syscall(SYS_wait4, pid, status, options, NULL);
}

/**
 * @brief Send a signal to a process.
 *
 * @param pid Target process ID.
 * @param sig Signal number.
 *
 * @return 0 on success or -1 on error.
 */
int __public__kill(int pid, int sig) {
    return syscall(SYS_kill, pid, sig);
}

/**
 * @brief Execute a program with explicit environment.
 *
 * @param path Path to executable.
 * @param argv Argument vector.
 * @param envp Environment vector.
 *
 * @return -1 on error (does not return on success).
 */
int __public__execve(const char *path, char *const argv[], char *const envp[]) {
    return syscall(SYS_execve, path, argv, envp);
}

/**
 * @brief Execute a program using current environment.
 *
 * @param file Executable file.
 * @param argv Argument vector.
 *
 * @return -1 on error (does not return on success).
 */
int __public__execvp(const char *file, char *const argv[]) {
    return syscall(SYS_execve, file, argv, environ);
}

/**
 * @brief Get pointer to environment array.
 *
 * @return Pointer to environment vector.
 */
char **__public__environ(void) {
    return environ;
}

/**
 * @brief Get environment variable value.
 *
 * @param key Environment variable name.
 *
 * @return Value string or NULL if not found.
 */
const char *__public__getenv(const char *key) {
    size_t klen = strlen(key);
    for (char **e = environ; *e; e++) {
        if (!strncmp(*e, key, klen) && (*e)[klen] == '=') {
            return *e + klen + 1;
        }
    }
    return NULL;
}

/**
 * @brief Set an environment variable.
 *
 * @param key        Variable name.
 * @param value      Variable value.
 * @param overwrite Whether to overwrite existing value.
 *
 * @return 0 on success or -1 on error.
 */
int __public__setenv(const char *key, const char *value, int overwrite) {
    if (!overwrite && __public__getenv(key))
        return 0;

    size_t klen = strlen(key), vlen = strlen(value);
    char *buf = malloc(klen + vlen + 2);
    if (!buf) return -1;
    memcpy(buf, key, klen);
    buf[klen] = '=';
    memcpy(buf+klen+1, value, vlen+1);

    for (char **e = environ; *e; e++) {
        if (!strncmp(*e, key, klen) && (*e)[klen] == '=') {
            *e = buf;
            return 0;
        }
    }

    int count = 0;
    while (environ[count]) count++;
    char **newenv = realloc(environ, sizeof(char*)*(count+2));
    if (!newenv) return -1;
    newenv[count] = buf;
    newenv[count+1] = NULL;
    environ = newenv;
    return 0;
}

/**
 * @brief Remove an environment variable.
 *
 * @param key Variable name.
 *
 * @return 0 on success or -1 on error.
 */
int __public__unsetenv(const char *key) {
    size_t klen = strlen(key);
    for (int i = 0; environ[i]; i++) {
        if (!strncmp(environ[i], key, klen) && environ[i][klen] == '=') {
            free(environ[i]);
            for (int j = i; environ[j]; j++)
                environ[j] = environ[j+1];
            return 0;
        }
    }
    return 0;
}

/**
 * @brief Get current working directory.
 *
 * @param buf  Buffer to store path.
 * @param size Buffer size.
 *
 * @return Length of path or -1 on error.
 */
int __public__getcwd(char *buf, size_t size) {
    return syscall(SYS_getcwd, buf, size);
}

/**
 * @brief Change current working directory.
 *
 * @param path New working directory.
 *
 * @return 0 on success or -1 on error.
 */
int __public__chdir(const char *path) {
    return syscall(SYS_chdir, path);
}

/**
 * @brief Get real user ID.
 */
int __public__getuid(void) {
    return syscall(SYS_getuid);
}

/**
 * @brief Get effective user ID.
 */
int __public__geteuid(void) {
    return syscall(SYS_geteuid);
}

/**
 * @brief Get real group ID.
 */
int __public__getgid(void) {
    return syscall(SYS_getgid);
}

/**
 * @brief Get effective group ID.
 */
int __public__getegid(void) {
    return syscall(SYS_getegid);
}

/**
 * @brief Set user ID.
 */
int __public__setuid(int uid) {
    return syscall(SYS_setuid, uid);
}

/**
 * @brief Set group ID.
 */
int __public__setgid(int gid) {
    return syscall(SYS_setgid, gid);
}

/**
 * @brief Set process group ID.
 */
int __public__setpgid(int pid, int pgid) {
    return syscall(SYS_setpgid, pid, pgid);
}

/**
 * @brief Get process group ID.
 */
int __public__getpgid(int pid) {
    return syscall(SYS_getpgid, pid);
}

/**
 * @brief Get current process group ID.
 */
int __public__getpgrp(void) {
    return syscall(SYS_getpgid, 0);
}

/**
 * @brief Create a new session.
 */
int __public__setsid(void) {
    return syscall(SYS_setsid);
}

/**
 * @brief Get resource limit.
 *
 * @param resource Resource type.
 * @param rlim     Pointer to rlimit struct.
 */
int __public__getrlimit(int resource, void *rlim) {
    return syscall(SYS_getrlimit, resource, rlim);
}

/**
 * @brief Set resource limit.
 *
 * @param resource Resource type.
 * @param rlim     Pointer to rlimit struct.
 */
int __public__setrlimit(int resource, const void *rlim) {
    return syscall(SYS_setrlimit, resource, rlim);
}

/**
 * @brief Install a signal handler.
 *
 * @param sig     Signal number.
 * @param handler Signal handler function.
 *
 * @return 0 on success or -1 on error.
 */
int __public__signal_install(int sig, void (*handler)(int)) {
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = handler;
    return syscall(SYS_rt_sigaction, sig, &sa, NULL, sizeof(sigset_t));
}

/**
 * @brief Open a file.
 *
 * @param path  File path.
 * @param flags Open flags (OS_O_*).
 * @param mode  File mode (for create).
 *
 * @return File descriptor or -1 on error.
 */
int __public__open(const char *path, int flags, int mode) {
    return syscall(SYS_openat, AT_FDCWD, path, flags, mode);
}

/**
 * @brief Close a file descriptor.
 *
 * @param fd File descriptor.
 *
 * @return 0 on success or -1 on error.
 */
int __public__close(int fd) {
    return syscall(SYS_close, fd);
}

/**
 * @brief Read from a file descriptor.
 *
 * @param fd  File descriptor.
 * @param buf Buffer to fill.
 * @param n   Number of bytes.
 *
 * @return Bytes read, 0 on EOF, or -1 on error.
 */
ssize_t __public__read(int fd, void *buf, size_t n) {
    return syscall(SYS_read, fd, buf, n);
}

/**
 * @brief Write to a file descriptor.
 *
 * @param fd  File descriptor.
 * @param buf Data buffer.
 * @param n   Number of bytes.
 *
 * @return Bytes written or -1 on error.
 */
ssize_t __public__write(int fd, const void *buf, size_t n) {
    return syscall(SYS_write, fd, buf, n);
}

/**
 * @brief Reposition file offset.
 *
 * @param fd     File descriptor.
 * @param offset Offset.
 * @param whence OS_SEEK_*.
 *
 * @return New offset or -1 on error.
 */
off_t __public__lseek(int fd, off_t offset, int whence) {
    return syscall(SYS_lseek, fd, offset, whence);
}

/**
 * @brief Get file status.
 *
 * @param fd File descriptor.
 * @param st Pointer to struct stat.
 *
 * @return 0 on success or -1 on error.
 */
int __public__fstat(int fd, struct stat *st) {
    return syscall(SYS_fstat, fd, st);
}

/**
 * @brief Duplicate a file descriptor.
 *
 * @param fd File descriptor.
 *
 * @return New file descriptor or -1 on error.
 */
int __public__dup(int fd) {
    return syscall(SYS_dup, fd);
}

/**
 * @brief Duplicate a file descriptor to a specific value.
 *
 * @param oldfd Existing FD.
 * @param newfd Target FD.
 *
 * @return New FD or -1 on error.
 */
int __public__dup2(int oldfd, int newfd) {
    return syscall(SYS_dup3, oldfd, newfd, 0);
}

/**
 * @brief Control file descriptor behavior.
 *
 * @param fd  File descriptor.
 * @param cmd Command (OS_F_*).
 * @param arg Command-specific argument.
 *
 * @return Result or -1 on error.
 */
int __public__fcntl(int fd, int cmd, long arg) {
    return syscall(SYS_fcntl, fd, cmd, arg);
}

/**
 * @brief Create a directory.
 *
 * @param path Directory path.
 * @param mode Permission bits.
 *
 * @return 0 on success or -1 on error.
 */
int __public__mkdir(const char *path, int mode) {
    return syscall(SYS_mkdirat, AT_FDCWD, path, mode);
}

/**
 * @brief Remove an empty directory.
 *
 * @param path Directory path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__rmdir(const char *path) {
    return syscall(SYS_unlinkat, AT_FDCWD, path, AT_REMOVEDIR);
}

/**
 * @brief Remove a file.
 *
 * @param path File path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__unlink(const char *path) {
    return syscall(SYS_unlinkat, AT_FDCWD, path, 0);
}

/**
 * @brief Rename a filesystem object.
 *
 * @param oldpath Source path.
 * @param newpath Destination path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__rename(const char *oldpath, const char *newpath) {
    return syscall(SYS_renameat, AT_FDCWD, oldpath, AT_FDCWD, newpath);
}

/**
 * @brief Rename with flags.
 *
 * @param oldpath Source path.
 * @param newpath Destination path.
 * @param flags   OS_RENAME_* flags.
 *
 * @return 0 on success or -1 on error.
 */
int __public__renameat2(const char *oldpath, const char *newpath, int flags) {
    return syscall(SYS_renameat2, AT_FDCWD, oldpath, AT_FDCWD, newpath, flags);
}

/**
 * @brief Create a hard link.
 *
 * @param oldpath Existing file.
 * @param newpath New link path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__link(const char *oldpath, const char *newpath) {
    return syscall(SYS_linkat, AT_FDCWD, oldpath, AT_FDCWD, newpath, 0);
}

/**
 * @brief Create a symbolic link.
 *
 * @param target Target path.
 * @param linkpath Symlink path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__symlink(const char *target, const char *linkpath) {
    return syscall(SYS_symlinkat, target, AT_FDCWD, linkpath);
}

/**
 * @brief Read a symbolic link.
 *
 * @param path Symlink path.
 * @param buf  Buffer to receive target.
 * @param size Buffer size.
 *
 * @return Number of bytes written or -1 on error.
 */
ssize_t __public__readlink(const char *path, char *buf, size_t size) {
    return syscall(SYS_readlinkat, AT_FDCWD, path, buf, size);
}

/**
 * @brief Check access permissions.
 *
 * @param path File path.
 * @param mode OS_*_OK flags.
 *
 * @return 0 on success or -1 on error.
 */
int __public__access(const char *path, int mode) {
    return syscall(SYS_faccessat, AT_FDCWD, path, mode, 0);
}

/**
 * @brief Get file metadata (follow symlinks).
 *
 * @param path File path.
 * @param st   Stat buffer.
 *
 * @return 0 on success or -1 on error.
 */
int __public__stat(const char *path, struct stat *st) {
    return syscall(SYS_newfstatat, AT_FDCWD, path, st, 0);
}

/**
 * @brief Get file metadata (do not follow symlinks).
 *
 * @param path File path.
 * @param st   Stat buffer.
 *
 * @return 0 on success or -1 on error.
 */
int __public__lstat(const char *path, struct stat *st) {
    return syscall(SYS_newfstatat, AT_FDCWD, path, st, AT_SYMLINK_NOFOLLOW);
}

/**
 * @brief Read directory entries.
 *
 * This is a low-level primitive. The runtime must parse
 * linux_dirent64 structures manually.
 *
 * @param fd   Directory file descriptor.
 * @param buf  Buffer for entries.
 * @param size Buffer size.
 *
 * @return Number of bytes read or -1 on error.
 */
int __public__getdents64(int fd, void *buf, size_t size) {
    return syscall(SYS_getdents64, fd, buf, size);
}

/**
 * @brief Map virtual memory.
 *
 * @param addr  Requested address (or NULL).
 * @param len   Length in bytes.
 * @param prot  Protection flags (OS_PROT_*).
 * @param flags Mapping flags (OS_MAP_*).
 * @param fd    File descriptor or -1.
 * @param off   File offset.
 *
 * @return Pointer to mapped memory or MAP_FAILED.
 */
void *__public__mmap(void *addr, size_t len, int prot, int flags, int fd, size_t off) {
    return (void *)syscall(SYS_mmap, addr, len, prot, flags, fd, off);
}

/**
 * @brief Unmap virtual memory.
 *
 * @param addr Mapped address.
 * @param len  Length in bytes.
 *
 * @return 0 on success or -1 on error.
 */
int __public__munmap(void *addr, size_t len) {
    return syscall(SYS_munmap, addr, len);
}

/**
 * @brief Change memory protection.
 *
 * @param addr Mapped address.
 * @param len  Length in bytes.
 * @param prot New protection flags.
 *
 * @return 0 on success or -1 on error.
 */
int __public__mprotect(void *addr, size_t len, int prot) {
    return syscall(SYS_mprotect, addr, len, prot);
}

/**
 * @brief Advise kernel about memory usage.
 *
 * @param addr   Address range.
 * @param len    Length in bytes.
 * @param advice OS_MADV_*.
 *
 * @return 0 on success or -1 on error.
 */
int __public__madvise(void *addr, size_t len, int advice) {
    return syscall(SYS_madvise, addr, len, advice);
}

/**
 * @brief Lock memory into RAM.
 *
 * @param addr Address range.
 * @param len  Length in bytes.
 *
 * @return 0 on success or -1 on error.
 */
int __public__mlock(void *addr, size_t len) {
    return syscall(SYS_mlock, addr, len);
}

/**
 * @brief Unlock memory.
 *
 * @param addr Address range.
 * @param len  Length in bytes.
 *
 * @return 0 on success or -1 on error.
 */
int __public__munlock(void *addr, size_t len) {
    return syscall(SYS_munlock, addr, len);
}

/**
 * @brief Control future/current memory locking.
 *
 * @param flags OS_MCL_* flags.
 *
 * @return 0 on success or -1 on error.
 */
int __public__mlockall(int flags) {
    return syscall(SYS_mlockall, flags);
}

/**
 * @brief Unlock all locked memory.
 *
 * @return 0 on success or -1 on error.
 */
int __public__munlockall(void) {
    return syscall(SYS_munlockall);
}

/**
 * @brief Get system page size.
 *
 * @return Page size in bytes.
 */
size_t __public__page_size(void) {
    return sysconf(_SC_PAGESIZE);
}
