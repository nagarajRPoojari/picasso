#ifndef OS_H
#define OS_H

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
#include <dirent.h>
#include <linux/futex.h>


/** Operation would block */
const int OS_EAGAIN = EAGAIN;
/** Interrupted system call */
const int OS_EINTR  = EINTR;
/** Invalid argument */
const int OS_EINVAL = EINVAL;
/** Permission denied */
const int OS_EPERM  = EPERM;
/** No such file or process */
const int OS_ENOENT = ENOENT;
/** Out of memory */
const int OS_ENOMEM = ENOMEM;


const int OS_WNOHANG   = WNOHANG;
const int OS_WUNTRACED = WUNTRACED;
const int OS_WCONTINUED = WCONTINUED;

const int OS_SIGINT  = SIGINT;
const int OS_SIGTERM = SIGTERM;
const int OS_SIGKILL = SIGKILL;
const int OS_SIGSEGV = SIGSEGV;
const int OS_SIGABRT = SIGABRT;
const int OS_SIGCHLD = SIGCHLD;
const int OS_SIGPIPE = SIGPIPE;
const int OS_SIGALRM = SIGALRM;
const int OS_SIGUSR1 = SIGUSR1;
const int OS_SIGUSR2 = SIGUSR2;

const int OS_RLIMIT_CPU    = RLIMIT_CPU;
const int OS_RLIMIT_FSIZE = RLIMIT_FSIZE;
const int OS_RLIMIT_DATA  = RLIMIT_DATA;
const int OS_RLIMIT_STACK = RLIMIT_STACK;
const int OS_RLIMIT_CORE  = RLIMIT_CORE;
const int OS_RLIMIT_NOFILE = RLIMIT_NOFILE;
const int OS_RLIMIT_AS    = RLIMIT_AS;

/* Standard file descriptor numbers */
/** Standard input */
const int OS_STDIN_FD  = 0;
/** Standard output */
const int OS_STDOUT_FD = 1;
/** Standard error */
const int OS_STDERR_FD = 2;

/*open() flags*/
const int OS_O_RDONLY   = O_RDONLY;
const int OS_O_WRONLY   = O_WRONLY;
const int OS_O_RDWR     = O_RDWR;
const int OS_O_APPEND   = O_APPEND;
const int OS_O_CREAT    = O_CREAT;
const int OS_O_EXCL     = O_EXCL;
const int OS_O_TRUNC    = O_TRUNC;
const int OS_O_CLOEXEC  = O_CLOEXEC;
const int OS_O_NONBLOCK = O_NONBLOCK;
const int OS_O_SYNC     = O_SYNC;
const int OS_O_DSYNC    = O_DSYNC;
const int OS_O_DIRECT   = O_DIRECT;

/*seek constants*/
const int OS_SEEK_SET = SEEK_SET;
const int OS_SEEK_CUR = SEEK_CUR;
const int OS_SEEK_END = SEEK_END;

/*fcntl commands*/
const int OS_F_DUPFD        = F_DUPFD;
const int OS_F_DUPFD_CLOEXEC = F_DUPFD_CLOEXEC;
const int OS_F_GETFD        = F_GETFD;
const int OS_F_SETFD        = F_SETFD;
const int OS_F_GETFL        = F_GETFL;
const int OS_F_SETFL        = F_SETFL;

/*FD flags*/
const int OS_FD_CLOEXEC = FD_CLOEXEC;

/*stat mode bits*/
const int OS_S_IFREG = S_IFREG;
const int OS_S_IFDIR = S_IFDIR;
const int OS_S_IFCHR = S_IFCHR;
const int OS_S_IFBLK = S_IFBLK;
const int OS_S_IFIFO = S_IFIFO;
const int OS_S_IFLNK = S_IFLNK;
const int OS_S_IFSOCK = S_IFSOCK;

const int OS_S_IRUSR = S_IRUSR;
const int OS_S_IWUSR = S_IWUSR;
const int OS_S_IXUSR = S_IXUSR;
const int OS_S_IRGRP = S_IRGRP;
const int OS_S_IWGRP = S_IWGRP;
const int OS_S_IXGRP = S_IXGRP;
const int OS_S_IROTH = S_IROTH;
const int OS_S_IWOTH = S_IWOTH;
const int OS_S_IXOTH = S_IXOTH;

/* Errors (FD-relevant subset)*/
const int OS_EBADF  = EBADF;
const int OS_EPIPE  = EPIPE;
const int OS_EIO    = EIO;
const int OS_ENOSPC = ENOSPC;

/*Special directory FDs*/
/** Current working directory */
const int OS_AT_FDCWD = AT_FDCWD;

/*unlinkat / renameat flags*/
const int OS_AT_REMOVEDIR = AT_REMOVEDIR;
const int OS_AT_SYMLINK_FOLLOW = AT_SYMLINK_FOLLOW;

/*link / rename flags*/
const int OS_RENAME_NOREPLACE = 1; // RENAME_NOREPLACE
const int OS_RENAME_EXCHANGE  = 2; // RENAME_EXCHANGE
const int OS_RENAME_WHITEOUT  = 4; // RENAME_WHITEOUT

/*Access mode flags*/
const int OS_F_OK = F_OK;
const int OS_R_OK = R_OK;
const int OS_W_OK = W_OK;
const int OS_X_OK = X_OK;

/*Directory entry types (d_type)*/
const int OS_DT_UNKNOWN = DT_UNKNOWN;
const int OS_DT_FIFO    = DT_FIFO;
const int OS_DT_CHR     = DT_CHR;
const int OS_DT_DIR     = DT_DIR;
const int OS_DT_BLK     = DT_BLK;
const int OS_DT_REG     = DT_REG;
const int OS_DT_LNK     = DT_LNK;
const int OS_DT_SOCK    = DT_SOCK;
const int OS_DT_WHT     = DT_WHT;


/*Memory protection flags*/
const int OS_PROT_NONE  = PROT_NONE;
const int OS_PROT_READ  = PROT_READ;
const int OS_PROT_WRITE = PROT_WRITE;
const int OS_PROT_EXEC  = PROT_EXEC;

/*mmap flags*/
const int OS_MAP_SHARED    = MAP_SHARED;
const int OS_MAP_PRIVATE   = MAP_PRIVATE;
const int OS_MAP_FIXED     = MAP_FIXED;
const int OS_MAP_ANONYMOUS = MAP_ANONYMOUS;
const int OS_MAP_STACK     = MAP_STACK;
const int OS_MAP_NORESERVE = MAP_NORESERVE;
const int OS_MAP_POPULATE  = MAP_POPULATE;
const int OS_MAP_GROWSDOWN = MAP_GROWSDOWN;

/*madvise advice*/
const int OS_MADV_NORMAL     = MADV_NORMAL;
const int OS_MADV_RANDOM     = MADV_RANDOM;
const int OS_MADV_SEQUENTIAL = MADV_SEQUENTIAL;
const int OS_MADV_WILLNEED   = MADV_WILLNEED;
const int OS_MADV_DONTNEED   = MADV_DONTNEED;
const int OS_MADV_FREE       = MADV_FREE;
const int OS_MADV_DONTFORK   = MADV_DONTFORK;
const int OS_MADV_DOFORK     = MADV_DOFORK;
const int OS_MADV_MERGEABLE  = MADV_MERGEABLE;
const int OS_MADV_UNMERGEABLE = MADV_UNMERGEABLE;
const int OS_MADV_HUGEPAGE   = MADV_HUGEPAGE;
const int OS_MADV_NOHUGEPAGE = MADV_NOHUGEPAGE;

const int OS_MCL_CURRENT = MCL_CURRENT;
const int OS_MCL_FUTURE  = MCL_FUTURE;

const int OS_EFAULT = EFAULT;
const int OS_EACCES = EACCES;


/**
 * @brief Get the current thread-local errno value.
 *
 * @return Current errno.
 */
int __public__errno(void);
/**
 * @brief Get current process ID.
 *
 * @return Process ID.
 */
int __public__getpid(void);

/**
 * @brief Get parent process ID.
 *
 * @return Parent process ID.
 */
int __public__getppid(void);

/**
 * @brief Get calling thread ID.
 *
 * @return Thread ID.
 */
int __public__gettid(void);

/**
 * @brief Terminate the current process immediately.
 *
 * @param code Exit status.
 */
void __public__exit(int code);

/**
 * @brief Create a new process.
 *
 * @return 0 in child, child PID in parent, or -1 on error.
 */
int __public__fork(void);

/**
 * @brief Wait for process state change.
 *
 * @param pid     Process ID to wait for.
 * @param status  Pointer to store exit status.
 * @param options wait options.
 *
 * @return PID of child or -1 on error.
 */
int __public__waitpid(int pid, int *status, int options);

/**
 * @brief Send a signal to a process.
 *
 * @param pid Target process ID.
 * @param sig Signal number.
 *
 * @return 0 on success or -1 on error.
 */
int __public__kill(int pid, int sig);

/**
 * @brief Execute a program with explicit environment.
 *
 * @param path Path to executable.
 * @param argv Argument vector.
 * @param envp Environment vector.
 *
 * @return -1 on error (does not return on success).
 */
int __public__execve(const char *path, char *const argv[], char *const envp[]);

/**
 * @brief Execute a program using current environment.
 *
 * @param file Executable file.
 * @param argv Argument vector.
 *
 * @return -1 on error (does not return on success).
 */
int __public__execvp(const char *file, char *const argv[]);

extern char **environ;

/**
 * @brief Get pointer to environment array.
 *
 * @return Pointer to environment vector.
 */
char **__public__environ(void);
/**
 * @brief Get environment variable value.
 *
 * @param key Environment variable name.
 *
 * @return Value string or NULL if not found.
 */
const char *__public__getenv(const char *key);

/**
 * @brief Set an environment variable.
 *
 * @param key        Variable name.
 * @param value      Variable value.
 * @param overwrite Whether to overwrite existing value.
 *
 * @return 0 on success or -1 on error.
 */
int __public__setenv(const char *key, const char *value, int overwrite);
/**
 * @brief Remove an environment variable.
 *
 * @param key Variable name.
 *
 * @return 0 on success or -1 on error.
 */
int __public__unsetenv(const char *key);

/**
 * @brief Get current working directory.
 *
 * @param buf  Buffer to store path.
 * @param size Buffer size.
 *
 * @return Length of path or -1 on error.
 */
int __public__getcwd(char *buf, size_t size);

/**
 * @brief Change current working directory.
 *
 * @param path New working directory.
 *
 * @return 0 on success or -1 on error.
 */
int __public__chdir(const char *path);

/**
 * @brief Get real user ID.
 */
int __public__getuid(void);

/**
 * @brief Get effective user ID.
 */
int __public__geteuid(void);

/**
 * @brief Get real group ID.
 */
int __public__getgid(void);

/**
 * @brief Get effective group ID.
 */
int __public__getegid(void);

/**
 * @brief Set user ID.
 */
int __public__setuid(int uid);

/**
 * @brief Set group ID.
 */
int __public__setgid(int gid);

/**
 * @brief Set process group ID.
 */
int __public__setpgid(int pid, int pgid);

/**
 * @brief Get process group ID.
 */
int __public__getpgid(int pid);

/**
 * @brief Get current process group ID.
 */
int __public__getpgrp(void);

/**
 * @brief Create a new session.
 */
int __public__setsid(void);

/**
 * @brief Get resource limit.
 *
 * @param resource Resource type.
 * @param rlim     Pointer to rlimit struct.
 */
int __public__getrlimit(int resource, void *rlim);

/**
 * @brief Set resource limit.
 *
 * @param resource Resource type.
 * @param rlim     Pointer to rlimit struct.
 */
int __public__setrlimit(int resource, const void *rlim);

/**
 * @brief Install a signal handler.
 *
 * @param sig     Signal number.
 * @param handler Signal handler function.
 *
 * @return 0 on success or -1 on error.
 */
int __public__signal_install(int sig, void (*handler)(int));

/**
 * @brief Open a file.
 *
 * @param path  File path.
 * @param flags Open flags (OS_O_*).
 * @param mode  File mode (for create).
 *
 * @return File descriptor or -1 on error.
 */
int __public__open(const char *path, int flags, int mode);

/**
 * @brief Close a file descriptor.
 *
 * @param fd File descriptor.
 *
 * @return 0 on success or -1 on error.
 */
int __public__close(int fd);

/**
 * @brief Read from a file descriptor.
 *
 * @param fd  File descriptor.
 * @param buf Buffer to fill.
 * @param n   Number of bytes.
 *
 * @return Bytes read, 0 on EOF, or -1 on error.
 */
ssize_t __public__read(int fd, void *buf, size_t n);

/**
 * @brief Write to a file descriptor.
 *
 * @param fd  File descriptor.
 * @param buf Data buffer.
 * @param n   Number of bytes.
 *
 * @return Bytes written or -1 on error.
 */
ssize_t __public__write(int fd, const void *buf, size_t n);

/**
 * @brief Reposition file offset.
 *
 * @param fd     File descriptor.
 * @param offset Offset.
 * @param whence OS_SEEK_*.
 *
 * @return New offset or -1 on error.
 */
off_t __public__lseek(int fd, off_t offset, int whence);

/**
 * @brief Get file status.
 *
 * @param fd File descriptor.
 * @param st Pointer to struct stat.
 *
 * @return 0 on success or -1 on error.
 */
int __public__fstat(int fd, struct stat *st);
/**
 * @brief Duplicate a file descriptor.
 *
 * @param fd File descriptor.
 *
 * @return New file descriptor or -1 on error.
 */
int __public__dup(int fd);

/**
 * @brief Duplicate a file descriptor to a specific value.
 *
 * @param oldfd Existing FD.
 * @param newfd Target FD.
 *
 * @return New FD or -1 on error.
 */
int __public__dup2(int oldfd, int newfd);

/**
 * @brief Control file descriptor behavior.
 *
 * @param fd  File descriptor.
 * @param cmd Command (OS_F_*).
 * @param arg Command-specific argument.
 *
 * @return Result or -1 on error.
 */
int __public__fcntl(int fd, int cmd, long arg);

/**
 * @brief Create a directory.
 *
 * @param path Directory path.
 * @param mode Permission bits.
 *
 * @return 0 on success or -1 on error.
 */
int __public__mkdir(const char *path, int mode);

/**
 * @brief Remove an empty directory.
 *
 * @param path Directory path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__rmdir(const char *path);

/**
 * @brief Remove a file.
 *
 * @param path File path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__unlink(const char *path);

/**
 * @brief Rename a filesystem object.
 *
 * @param oldpath Source path.
 * @param newpath Destination path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__rename(const char *oldpath, const char *newpath);

/**
 * @brief Rename with flags.
 *
 * @param oldpath Source path.
 * @param newpath Destination path.
 * @param flags   OS_RENAME_* flags.
 *
 * @return 0 on success or -1 on error.
 */
int __public__renameat2(const char *oldpath, const char *newpath, int flags);

/**
 * @brief Create a hard link.
 *
 * @param oldpath Existing file.
 * @param newpath New link path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__link(const char *oldpath, const char *newpath);

/**
 * @brief Create a symbolic link.
 *
 * @param target Target path.
 * @param linkpath Symlink path.
 *
 * @return 0 on success or -1 on error.
 */
int __public__symlink(const char *target, const char *linkpath);

/**
 * @brief Read a symbolic link.
 *
 * @param path Symlink path.
 * @param buf  Buffer to receive target.
 * @param size Buffer size.
 *
 * @return Number of bytes written or -1 on error.
 */
ssize_t __public__readlink(const char *path, char *buf, size_t size);

/**
 * @brief Get file metadata (follow symlinks).
 *
 * @param path File path.
 * @param st   Stat buffer.
 *
 * @return 0 on success or -1 on error.
 */
int __public__stat(const char *path, struct stat *st);

/**
 * @brief Get file metadata (do not follow symlinks).
 *
 * @param path File path.
 * @param st   Stat buffer.
 *
 * @return 0 on success or -1 on error.
 */
int __public__lstat(const char *path, struct stat *st);

/**
 * @brief Check access permissions.
 *
 * @param path File path.
 * @param mode OS_*_OK flags.
 *
 * @return 0 on success or -1 on error.
 */
int __public__access(const char *path, int mode);

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
int __public__getdents64(int fd, void *buf, size_t size);

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
void *__public__mmap(void *addr, size_t len, int prot,
              int flags, int fd, size_t off);

/**
 * @brief Unmap virtual memory.
 *
 * @param addr Mapped address.
 * @param len  Length in bytes.
 *
 * @return 0 on success or -1 on error.
 */
int __public__munmap(void *addr, size_t len);

/**
 * @brief Change memory protection.
 *
 * @param addr Mapped address.
 * @param len  Length in bytes.
 * @param prot New protection flags.
 *
 * @return 0 on success or -1 on error.
 */
int __public__mprotect(void *addr, size_t len, int prot);

/**
 * @brief Advise kernel about memory usage.
 *
 * @param addr   Address range.
 * @param len    Length in bytes.
 * @param advice OS_MADV_*.
 *
 * @return 0 on success or -1 on error.
 */
int __public__madvise(void *addr, size_t len, int advice);
/**
 * @brief Lock memory into RAM.
 *
 * @param addr Address range.
 * @param len  Length in bytes.
 *
 * @return 0 on success or -1 on error.
 */
int __public__mlock(void *addr, size_t len);

/**
 * @brief Unlock memory.
 *
 * @param addr Address range.
 * @param len  Length in bytes.
 *
 * @return 0 on success or -1 on error.
 */
int __public__munlock(void *addr, size_t len);

/**
 * @brief Control future/current memory locking.
 *
 * @param flags OS_MCL_* flags.
 *
 * @return 0 on success or -1 on error.
 */
int __public__mlockall(int flags);

/**
 * @brief Unlock all locked memory.
 *
 * @return 0 on success or -1 on error.
 */
int __public__munlockall(void);

/**
 * @brief Get system page size.
 *
 * @return Page size in bytes.
 */
size_t __public__page_size(void);

/** @futex: */

/**
 * @brief Wait on a futex word.
 *
 * The calling thread sleeps if *uaddr == val.
 *
 * @param uaddr Pointer to futex word.
 * @param val Expected value.
 * @param timeout Optional timeout (NULL for infinite wait).
 * @return 0 on success or -1 on error.
 */
int __public__futex_wait(int *uaddr, int val, const struct timespec *timeout);

/**
 * @brief Wake up threads waiting on a futex word.
 *
 * @param uaddr Pointer to futex word.
 * @param count Maximum number of waiters to wake.
 * @return Number of woken threads or -1 on error.
 */
int __public__futex_wake(int *uaddr, int count);

/**
 * @brief Wait on a futex word with bitmask.
 *
 * The calling thread sleeps if (*uaddr & mask) == val.
 *
 * @param uaddr Pointer to futex word.
 * @param val Expected value.
 * @param timeout Optional timeout.
 * @param mask Bitmask.
 * @return 0 on success or -1 on error.
 */
int __public__futex_wait_bitset(int *uaddr,int val,const struct timespec *timeout,int mask);

/**
 * @brief Wake threads waiting on a futex word using a bitmask.
 *
 * @param uaddr Pointer to futex word.
 * @param count Maximum number of waiters to wake.
 * @param mask Bitmask.
 * @return Number of woken threads or -1 on error.
 */
int __public__futex_wake_bitset(int *uaddr, int count, int mask);

/**
 * @brief Requeue waiters from one futex to another.
 *
 * Wakes up to wake_count waiters and requeues the rest to uaddr2.
 *
 * @param uaddr Source futex.
 * @param wake_count Number of waiters to wake.
 * @param requeue_count Number of waiters to requeue.
 * @param uaddr2 Target futex.
 * @return Number of affected waiters or -1 on error.
 */
int __public__futex_requeue( int *uaddr, int wake_count, int requeue_count, int *uaddr2);

/**
 * @brief Wake one waiter and requeue remaining waiters.
 *
 * @param uaddr Source futex.
 * @param uaddr2 Target futex.
 * @param wake_count Number of waiters to wake.
 * @param requeue_count Number of waiters to requeue.
 * @return Number of affected waiters or -1 on error.
 */
int __public__futex_cmp_requeue( int *uaddr, int *uaddr2, int wake_count, int requeue_count, int val);


/**
 * @brief Wake a single waiter (optimized common case).
 *
 * @param uaddr Pointer to futex word.
 * @return Number of woken threads or -1 on error.
 */
int __public__futex_wake_one(int *uaddr);

/**
 * @brief Wake all waiters.
 *
 * @param uaddr Pointer to futex word.
 * @return Number of woken threads or -1 on error.
 */
int __public__futex_wake_all(int *uaddr);
#endif