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

/* Function Prototypes */
/**
 * @brief Get current errno value.
 * @return errno.
 */
int __public__os_errno(void) {
    return errno;
}

/**
 * @brief Get current process ID.
 * @return PID.
 */
int __public__os_getpid(void) {
    return getpid();
}

/**
 * @brief Get parent process ID.
 * @return Parent PID.
 */
int __public__os_getppid(void) {
    return getppid();
}

/**
 * @brief Get current thread ID.
 * @note macOS implementation uses pthread_threadid_np().
 * @return Thread ID.
 */
uint64_t __public__os_gettid(void) {
    uint64_t tid = 0;
    pthread_threadid_np(NULL, &tid);
    return tid;
}

/**
 * @brief Terminate the current process.
 * @param code Exit status.
 */
void __public__os_exit(int code) {
    _exit(code);
}

/**
 * @brief Create a child process.
 * @return 0 in child, child PID in parent, -1 on error.
 */
int __public__os_fork(void) {
    return fork();
}

/**
 * @brief Wait for a child process.
 * @param pid     Process ID.
 * @param status  Exit status.
 * @param options Wait options.
 * @return PID or -1 on error.
 */
int __public__os_waitpid(int pid, int *status, int options) {
    return waitpid(pid, status, options);
}

/**
 * @brief Send a signal to a process.
 * @param pid Process ID.
 * @param sig Signal number.
 * @return 0 on success, -1 on error.
 */
int __public__os_kill(int pid, int sig) {
    return kill(pid, sig);
}

/**
 * @brief Execute a program.
 * @param path Executable path.
 * @param argv Argument vector.
 * @param envp Environment.
 * @return -1 on error.
 */
int __public__os_execve(const char *path, char *const argv[], char *const envp[]) {
    return execve(path, argv, envp);
}

/**
 * @brief Execute a program using PATH lookup.
 * @param file Executable name.
 * @param argv Argument vector.
 * @return -1 on error.
 */
int __public__os_execvp(const char *file, char *const argv[]) {
    return execvp(file, argv);
}

/**
 * @brief Get environment variable array.
 * @return Environment pointer.
 */
char **__public__os_environ(void) {
    return *_NSGetEnviron();
}

/**
 * @brief Get environment variable value.
 * @param key Variable name.
 * @return Value or NULL.
 */
const char *__public__os_getenv(const char *key) {
    return getenv(key);
}

/**
 * @brief Set environment variable.
 * @param key       Variable name.
 * @param value     Value.
 * @param overwrite Overwrite if exists.
 * @return 0 on success, -1 on error.
 */
int __public__os_setenv(const char *key, const char *value, int overwrite) {
    return setenv(key, value, overwrite);
}

/**
 * @brief Remove environment variable.
 * @param key Variable name.
 * @return 0 on success, -1 on error.
 */
int __public__os_unsetenv(const char *key) {
    return unsetenv(key);
}


/**
 * @brief Get current working directory.
 * @param buf  Output buffer.
 * @param size Buffer size.
 * @return 0 on success, -1 on error.
 */
int __public__os_getcwd(char *buf, size_t size) {
    return getcwd(buf, size) ? 0 : -1;
}

/**
 * @brief Change mode of given path.
 * @param path Directory/File path.
 * @param mode access mode.
 * @return 0 on success, -1 on error.
 */
int __public__os_chmod(const char *path, int64_t mode) {
    return chmod(path, mode);
}

/**
 * @brief Chown changes the numeric uid and gid of the named file. 
 * If the file is a symbolic link, it changes the uid and gid of the link's target.
 * A uid or gid of -1 means to not change that value.
 * @param name Directory/File path.
 * @param uid uid.
 * @param gid uid.
 * @return 0 on success, -1 on error.
 */
int __public__os_chown(const char *name, int64_t uid, int64_t gid) {
    return chown(name, uid, gid);
}

/**
 * @brief Change working directory.
 * @param path Directory path.
 * @return 0 on success, -1 on error.
 */
int __public__os_chdir(const char *path) { return chdir(path); }

/**
 * @brief Get user ID.
 * @return UID.
 */
int __public__os_getuid(void)  { return getuid();  }
/**
 * @brief Get effective user ID.
 * @return Effective UID.
 */
int __public__os_geteuid(void) { return geteuid(); }
/**
 * @brief Get group ID.
 * @return GID.
 */
int __public__os_getgid(void)  { return getgid();  }
/**
 * @brief Get effective group ID.
 * @return Effective GID.
 */
int __public__os_getegid(void) { return getegid(); }

/**
 * @brief Set user ID.
 * @param uid User ID.
 * @return 0 on success, -1 on error.
 */
int __public__os_setuid(int uid) { return setuid(uid); }
/**
 * @brief Set group ID.
 * @param gid Group ID.
 * @return 0 on success, -1 on error.
 */
int __public__os_setgid(int gid) { return setgid(gid); }

/**
 * @brief Set process group ID.
 * @param pid  Process ID.
 * @param pgid Process group ID.
 * @return 0 on success, -1 on error.
 */
int __public__os_setpgid(int pid, int pgid) { return setpgid(pid, pgid); }

/**
 * @brief Get process group ID.
 * @param pid Process ID.
 * @return Process group ID or -1 on error.
 */
int __public__os_getpgid(int pid) { return getpgid(pid); }

/**
 * @brief Get process group ID of calling process.
 * @return Process group ID.
 */
int __public__os_getpgrp(void) { return getpgrp(); }

/**
 * @brief Create a new session.
 * @return Session ID or -1 on error.
 */
int __public__os_setsid(void) { return setsid(); }

/**
 * @brief Get resource limits.
 * @param resource Resource type.
 * @param rlim     Output rlimit structure.
 * @return 0 on success, -1 on error.
 */
int __public__os_getrlimit(int resource, void *rlim) {
    return getrlimit(resource, (struct rlimit *)rlim);
}

/**
 * @brief Set resource limits.
 * @param resource Resource type.
 * @param rlim     Input rlimit structure.
 * @return 0 on success, -1 on error.
 */
int __public__os_setrlimit(int resource, const void *rlim) {
    return setrlimit(resource, (const struct rlimit *)rlim);
}

/**
 * @brief Open a file.
 * @param path  File path.
 * @param flags Open flags.
 * @param mode  File mode.
 * @return File descriptor or -1.
 */
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

/**
 * @brief Close a file descriptor.
 * @param fd File descriptor.
 * @return 0 on success, -1 on error.
 */
int __public__os_close(int fd) {
    return close(fd);
}

/**
 * @brief Read from a file descriptor.
 * @param fd  File descriptor.
 * @param buf Buffer.
 * @param n   Bytes to read.
 * @return Bytes read or -1.
 */
ssize_t __public__os_read(int fd, void *buf, size_t n) {
    return read(fd, buf, n);
}

/**
 * @brief Write to a file descriptor.
 * @param fd  File descriptor.
 * @param buf Buffer.
 * @param n   Bytes to write.
 * @return Bytes written or -1.
 */
ssize_t __public__os_write(int fd, const void *buf, size_t n) {
    return write(fd, buf, n);
}

/**
 * @brief Reposition file offset.
 * @param fd     File descriptor.
 * @param offset Offset.
 * @param whence Seek mode.
 * @return New offset or -1.
 */
off_t __public__os_lseek(int fd, off_t offset, int whence) {
    return lseek(fd, offset, whence);
}

/**
 * @brief Get file status.
 * @param fd File descriptor.
 * @param st Pointer to struct stat.
 * @return 0 on success or -1 on error.
 */
int __public__os_fstat(int fd, struct stat *st) {
    return fstat(fd, st);
}

/**
 * @brief Duplicate a file descriptor.
 * @param fd File descriptor.
 * @return New FD or -1.
 */
int __public__os_dup(int fd) {
    return dup(fd);
}

/**
 * @brief Duplicate a file descriptor to a specific value.
 *
 * @param oldfd Existing FD.
 * @param newfd Target FD.
 *
 * @return New FD or -1 on error.
 */
int __public__os_dup2(int oldfd, int newfd) {
    return dup2(oldfd, newfd);
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
int __public__os_fcntl(int fd, int cmd, long arg) {
    return fcntl(fd, cmd, arg);
}

#ifndef RENAME_EXCL
#define RENAME_EXCL 0x00000001
#endif

#ifndef RENAME_SWAP
#define RENAME_SWAP 0x00000002
#endif


/**
 * @brief Create a directory.
 * @param path Directory path.
 * @param mode Permissions.
 * @return 0 on success, -1 on error.
 */
int __public__os_mkdir(const char *path, int mode) { return mkdir(path, mode); }

/**
 * @brief Remove a directory.
 * @param path Directory path.
 * @return 0 on success, -1 on error.
 */
int __public__os_rmdir(const char *path) { return rmdir(path); }

/**
 * @brief Delete a file.
 * @param path File path.
 * @return 0 on success, -1 on error.
 */
int __public__os_unlink(const char *path) { return unlink(path); }

/**
 * @brief Rename a file or directory.
 * @param oldpath Old path.
 * @param newpath New path.
 * @return 0 on success, -1 on error.
 */
int __public__os_rename(const char *oldpath, const char *newpath) {
    return rename(oldpath, newpath);
}

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

/**
 * @brief Create a hard link.
 * @param oldpath Existing path.
 * @param newpath Link path.
 * @return 0 on success, -1 on error.
 */
int __public__os_link(const char *oldpath, const char *newpath) {
    return link(oldpath, newpath);
}

/**
 * @brief Create a symbolic link.
 * @param target   Target path.
 * @param linkpath Link path.
 * @return 0 on success, -1 on error.
 */
int __public__os_symlink(const char *target, const char *linkpath) {
    return symlink(target, linkpath);
}

/**
 * @brief Read symbolic link contents.
 * @param path Symbolic link path.
 * @param buf  Output buffer.
 * @param size Buffer size.
 * @return Bytes read or -1 on error.
 */
ssize_t __public__os_readlink(const char *path, char *buf, size_t size) {
    return readlink(path, buf, size);
}

/**
 * @brief Get file status.
 * @param path File path.
 * @param st   Output stat structure.
 * @return 0 on success, -1 on error.
 */
int __public__os_stat(const char *path, struct stat *st) {
    return stat(path, st);
}

/**
 * @brief Get file status (don't follow symlinks).
 * @param path File path.
 * @param st   Output stat structure.
 * @return 0 on success, -1 on error.
 */
int __public__os_lstat(const char *path, struct stat *st) {
    return lstat(path, st);
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

/**
 * @brief Map memory.
 *
 * @param addr  Address hint.
 * @param len   Length.
 * @param prot  Protection flags.
 * @param flags Mapping flags.
 * @param fd    File descriptor.
 * @param off   Offset.
 *
 * @return Mapped address or MAP_FAILED.
 */
void *__public__os_mmap(void *addr,
                        size_t len,
                        int prot,
                        int flags,
                        int fd,
                        off_t off) {
    return mmap(addr, len, prot, flags, fd, off);
}

/**
 * @brief Unmap memory.
 * @param addr Address.
 * @param len  Length.
 * @return 0 on success, -1 on error.
 */
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

/**
 * @brief Check file accessibility.
 * @param path File path.
 * @param mode Access mode.
 * @return 0 on success, -1 on error.
 */
int __public__os_access(const char *path, int mode) {
    return access(path, mode);
}