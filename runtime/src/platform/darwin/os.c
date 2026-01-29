#include "os.h"

#include <errno.h>
#include <limits.h>
#include <pthread.h>
#include <crt_externs.h>
#include <time.h>
#include <sys/time.h>
#include <sys/syscall.h>
#include <sys/attr.h>
#include <fcntl.h>
#include <stdio.h>
#include <dirent.h>
#include <string.h>
#include <unistd.h>

int __public__os_errno(void) {
    return errno;
}

int __public__os_getpid(void) {
    return getpid();
}

int __public__os_getppid(void) {
    return getppid();
}

uint64_t __public__os_gettid(void) {
    uint64_t tid = 0;
    pthread_threadid_np(NULL, &tid);
    return tid;
}

void __public__os_exit(int code) {
    _exit(code);
}

int __public__os_fork(void) {
    return fork();
}

int __public__os_waitpid(int pid, int *status, int options) {
    return waitpid(pid, status, options);
}

int __public__os_kill(int pid, int sig) {
    return kill(pid, sig);
}

int __public__os_execve(const char *path, char *const argv[], char *const envp[]) {
    return execve(path, argv, envp);
}

int __public__os_execvp(const char *file, char *const argv[]) {
    return execvp(file, argv);
}

char **__public__os_environ(void) {
    return *_NSGetEnviron();
}

const char *__public__os_getenv(const char *key) {
    return getenv(key);
}

int __public__os_setenv(const char *key, const char *value, int overwrite) {
    return setenv(key, value, overwrite);
}

int __public__os_unsetenv(const char *key) {
    return unsetenv(key);
}


int __public__os_open(const char *path, int flags, int mode) {
    int fd = open(path, flags, mode);

    /*
     * macOS has no O_DIRECT.
     * F_NOCACHE disables buffer cache for this FD.
     */
    if (fd != -1 && (flags & __public__os_O_DIRECT)) {
        fcntl(fd, F_NOCACHE, 1);
    }

    return fd;
}

int __public__os_close(int fd) {
    return close(fd);
}

ssize_t __public__os_read(int fd, void *buf, size_t n) {
    return read(fd, buf, n);
}

ssize_t __public__os_write(int fd, const void *buf, size_t n) {
    return write(fd, buf, n);
}

off_t __public__os_lseek(int fd, off_t offset, int whence) {
    return lseek(fd, offset, whence);
}

int __public__os_dup(int fd) {
    return dup(fd);
}

int __public__os_dup2(int oldfd, int newfd) {
    return dup2(oldfd, newfd);
}

int __public__os_fcntl(int fd, int cmd, long arg) {
    return fcntl(fd, cmd, arg);
}

int __public__os_fstat(int fd, struct stat *st) {
    return fstat(fd, st);
}

#ifndef RENAME_EXCL
#define RENAME_EXCL 0x00000001
#endif

#ifndef RENAME_SWAP
#define RENAME_SWAP 0x00000002
#endif

/*
 * Linux renameat2() compatibility shim.
 * macOS provides renameatx_np().
 */
int __public__os_renameat2(const char *oldpath,
                           const char *newpath,
                           int flags) {
    unsigned int native = 0;

    if (flags & __public__os_RENAME_NOREPLACE) native |= RENAME_EXCL;
    if (flags & __public__os_RENAME_EXCHANGE)  native |= RENAME_SWAP;

    return renameatx_np(AT_FDCWD, oldpath,
                        AT_FDCWD, newpath,
                        native);
}

/*
 * Darwin directory reading.
 * WARNING: layout is struct dirent (NOT linux_dirent64).
 */
int __public__os_getdents64(int fd, void *buf, size_t size) {
    DIR *dir;
    struct dirent *entry;
    char *ptr = buf;
    size_t remaining = size;

    // Convert fd to DIR* (fdopendir duplicates fd, so we need to close it)
    dir = fdopendir(fd);
    if (!dir) return -1;

    while ((entry = readdir(dir)) != NULL) {
        size_t reclen = entry->d_reclen;

        if (reclen > remaining)
            break;

        memcpy(ptr, entry, reclen);
        ptr += reclen;
        remaining -= reclen;
    }

    closedir(dir);

    return (int)(size - remaining);  // bytes written
}

void *__public__os_mmap(void *addr,
                        size_t len,
                        int prot,
                        int flags,
                        int fd,
                        off_t off) {
    return mmap(addr, len, prot, flags, fd, off);
}

int __public__os_munmap(void *addr, size_t len) {
    return munmap(addr, len);
}

size_t __public__os_page_size(void) {
    return (size_t)sysconf(_SC_PAGESIZE);
}

int __public__os_mprotect(void *addr, size_t len, int prot) {
    return mprotect(addr, len, prot);
}

int __public__os_madvise(void *addr, size_t len, int advice) {
    return madvise(addr, len, advice);
}

/*
 * Verified on Darwin 21–23 (macOS 12–14)
 */
#define SYS_ulock_wait 515
#define SYS_ulock_wake 516

#define UL_COMPARE_AND_WAIT 1
#define UL_WAKE_ALL         0x00000100

int __public__os_futex_wait(int *uaddr,
                            int val,
                            const struct timespec *timeout) {
    uint32_t timeout_us = 0;

    if (timeout) {
        timeout_us = (uint32_t)(
            timeout->tv_sec * 1000000 +
            timeout->tv_nsec / 1000
        );
    }

    return (int)syscall(SYS_ulock_wait,
                        UL_COMPARE_AND_WAIT,
                        uaddr,
                        (uint64_t)val,
                        timeout_us);
}

int __public__os_futex_wake_one(int *uaddr) {
    return (int)syscall(SYS_ulock_wake,
                        UL_COMPARE_AND_WAIT,
                        uaddr,
                        0);
}

int __public__os_futex_wake_all(int *uaddr) {
    return (int)syscall(SYS_ulock_wake,
                        UL_COMPARE_AND_WAIT | UL_WAKE_ALL,
                        uaddr,
                        0);
}

int __public__os_signal_install(int sig, void (*handler)(int)) {
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));

    sa.sa_handler = handler;
    sa.sa_flags   = SA_RESTART;

    sigemptyset(&sa.sa_mask);
    return sigaction(sig, &sa, NULL);
}

int __public__os_getcwd(char *buf, size_t size) {
    return getcwd(buf, size) ? 0 : -1;
}

int __public__os_chdir(const char *path) { return chdir(path); }

int __public__os_getuid(void)  { return getuid();  }
int __public__os_geteuid(void) { return geteuid(); }
int __public__os_getgid(void)  { return getgid();  }
int __public__os_getegid(void) { return getegid(); }

int __public__os_setuid(int uid) { return setuid(uid); }
int __public__os_setgid(int gid) { return setgid(gid); }

int __public__os_setpgid(int pid, int pgid) { return setpgid(pid, pgid); }
int __public__os_getpgid(int pid) { return getpgid(pid); }
int __public__os_getpgrp(void) { return getpgrp(); }
int __public__os_setsid(void) { return setsid(); }

int __public__os_getrlimit(int resource, void *rlim) {
    return getrlimit(resource, (struct rlimit *)rlim);
}

int __public__os_setrlimit(int resource, const void *rlim) {
    return setrlimit(resource, (const struct rlimit *)rlim);
}

int __public__os_mkdir(const char *path, int mode) { return mkdir(path, mode); }
int __public__os_rmdir(const char *path) { return rmdir(path); }
int __public__os_unlink(const char *path) { return unlink(path); }

int __public__os_rename(const char *oldpath, const char *newpath) {
    return rename(oldpath, newpath);
}

int __public__os_link(const char *oldpath, const char *newpath) {
    return link(oldpath, newpath);
}

int __public__os_symlink(const char *target, const char *linkpath) {
    return symlink(target, linkpath);
}

ssize_t __public__os_readlink(const char *path, char *buf, size_t size) {
    return readlink(path, buf, size);
}

int __public__os_stat(const char *path, struct stat *st) {
    return stat(path, st);
}

int __public__os_lstat(const char *path, struct stat *st) {
    return lstat(path, st);
}

int __public__os_access(const char *path, int mode) {
    return access(path, mode);
}
