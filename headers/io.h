#ifndef IO_H
#define IO_H

#include "queue.h"
#include "task.h"


/** Number of I/O worker threads in the pool */
#define IO_THREAD_POOL_SIZE 1

/**  @depricated: Maximum number of events returned by epoll_wait at once */
#define MAX_EVENTS 16

/**  @depricated: Maximum number of tasks in the I/O queue */
#define IO_QUEUE_SIZE 256

/** Queue depth */
#define QUEUE_DEPTH  32

extern struct io_uring **io_ring_map;

/**
 * @brief Worker thread that waits for completed I/O events from io_uring.
 *
 * This thread continuously polls for completion events from the io_uring
 * submission queue. When an event completes, it retrieves the associated
 * task (stored in CQE user data), updates its `nread` field with the result,
 * and pushes it back into the scheduler's ready queue.
 *
 * @param arg Unused.
 * @return void* Always returns NULL.
 */
void *io_worker(void *arg);

/**
 * @brief Internal helper to asynchronously read from STDIN using io_uring.
 *
 * Prepares and submits a read request for the current task’s buffer from STDIN.
 * STDIN is set to non-blocking mode for safety. Upon submission, the current
 * task yields control to the scheduler until the I/O operation completes.
 *
 * @return void
 */
void _async_stdin_read(void);

/**
 * @brief Internal helper to asynchronously write to STDOUT using io_uring.
 *
 * This function prepares and submits a non-blocking write request 
 * for the current task’s buffer to STDOUT. STDOUT is set to non-blocking 
 * mode before the operation. Once submitted, the current task yields 
 * control back to the scheduler until the write completes.
 *
 * @return void
 *
 * @note Exits the process if SQE allocation or submission fails.
 */
void _async_stdout_write();

/**
 * @brief Internal helper to asynchronously read from a file descriptor using io_uring.
 *
 * Prepares and submits a read request for the current task’s buffer.
 * The request is configured to read from the given file descriptor at
 * the specified offset. After submission, the task yields execution
 * to allow other tasks to run.
 *
 * @return void
 */
void _async_file_read(void);

/**
 * @brief Internal helper to asynchronously write to a file descriptor using io_uring.
 *
 * Prepares and submits an asynchronous write request for the current task’s buffer
 * to the specified file descriptor at the given offset. Once the request is submitted,
 * the current task yields execution to the scheduler until the write operation completes.
 *
 * @return void
 *
 * @note Exits the process if SQE allocation or submission fails.
 */
void _async_file_write();

/**
 * @brief Public API for asynchronous STDIN read.
 *
 * Configures the current task for reading `n` bytes from STDIN into
 * the provided buffer `buf`. Submits the read request through io_uring
 * and yields until the operation completes.
 *
 * @param buf Pointer to buffer where data will be stored.
 * @param n   Number of bytes to read.
 *
 * @return void* Pointer to `current_task->nread` indicating number of bytes read.
 */
void* async_stdin_read(char* buf, int n);

void* ascan(int n);

/**
 * @brief Public API for asynchronous STDOUT write using io_uring.
 *
 * Configures the current task to write `n` bytes from buffer `buf`
 * to STDOUT. Submits the write request via io_uring and yields execution
 * until the operation completes.
 *
 * @param buf Pointer to buffer containing data to write.
 * @param n   Number of bytes to write.
 *
 * @return void* Pointer to `current_task->nwrite`, which holds the number of bytes written.
 *
 * @note STDOUT is set to non-blocking mode before submission.
 */
void* async_stdout_write(const char* buf, int n);

void* aprintf(const char* fmt, ...);

/**
 * @brief Public API for asynchronous file read.
 *
 * Configures the current task context for reading `n` bytes from the
 * given file descriptor `fd` starting at `offset` into `buf`. Submits
 * the operation through io_uring and yields the current task until
 * the operation completes.
 *
 * @param fd      File descriptor to read from.
 * @param buf     Pointer to buffer where data will be stored.
 * @param n       Number of bytes to read.
 * @param offset  Offset in file to begin reading.
 * 
 * @return void*  Pointer to `current_task->nread` indicating number of bytes read.
 */
void* async_file_read(int fd, char* buf, int n, int offset);

void* afread(char* f, char* buf, int n, int offset);

/**
 * @brief Public API for asynchronous file write using io_uring.
 *
 * Configures the current task to write `n` bytes from buffer `buf`
 * to the specified file descriptor `fd` starting at the given `offset`.
 * Submits the write request via io_uring and yields execution until
 * the operation completes.
 *
 * @param fd      File descriptor to write to.
 * @param buf     Pointer to buffer containing data to write.
 * @param n       Number of bytes to write.
 * @param offset  File offset to begin writing from.
 *
 * @return void*  Pointer to `current_task->nwrite`, which holds the number of bytes written.
 *
 * @note Assumes `current_task` and its scheduler context are properly initialized.
 */
void* async_file_write(int fd, const char* buf, int n, int offset);

void* afwrite(char* f, char* buf, int n, int offset);

#endif