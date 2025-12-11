#ifndef IO_H
#include "queue.h"
#include "task.h"


/** Number of I/O worker threads in the pool, must be kept equal to number of scheduler threads */
#define IO_THREAD_POOL_SIZE 4

/**  @depricated: Maximum number of events returned by epoll_wait at once */
#define MAX_EVENTS 16

/**  @depricated: Maximum number of tasks in the I/O queue */
#define IO_QUEUE_SIZE 256

/** Queue depth */
#define QUEUE_DEPTH  256

/** @deprecated: io_uring submission done when req hits this threshold */
#define SUBMIT_THRESHOLD (256/2)

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
void async_stdin_read(void);

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
void async_stdout_write();

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
void async_file_read(void);

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
void async_file_write();

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
void* __public__ascan(int n);

/**
 * @brief Asynchronously write formatted output to STDOUT using io_uring.
 *
 * Formats the input string and arguments, allocates a buffer for the result,
 * configures the current task with write parameters, and submits an io_uring
 * write request. Yields until the operation completes.
 *
 * @param fmt Format string (printf-style).
 * @param ... Variable arguments matching the format string.
 *
 * @return NULL on success, NULL on allocation or formatting failure.
 */
void* __public__aprintf(const char* fmt, ...);

/**
 * @brief Asynchronously read n bytes from a file at a given offset.
 *
 * Configures the current task with the file descriptor, buffer, byte count,
 * and offset, then submits an io_uring read request. Yields until completion.
 *
 * @param f      FILE pointer to read from.
 * @param buf    Buffer to store read data.
 * @param n      Number of bytes to read.
 * @param offset File offset to start reading from.
 *
 * @return Pointer to bytes read count in current task context.
 */
void* __public__afread(char* f, char* buf, int n, int offset);

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
void* __public__afwrite(char* f, char* buf, int n, int offset);

/**
 * @brief Synchronously read n bytes from STDIN.
 *
 * Allocates a buffer, configures the current task with the read parameters,
 * and submits an io_uring read request for STDIN. Yields until completion.
 *
 * @param n Number of bytes to read.
 *
 * @return Pointer to allocated buffer containing the read data.
 */
void* __public__sscan(int n);

/**
 * @brief Write formatted output to STDOUT synchronously.
 *
 * Formats the input string and arguments, then writes the result to STDOUT
 * using the write syscall. Handles partial writes and errors gracefully.
 *
 * @param fmt Format string (printf-style).
 * @param ... Variable arguments matching the format string.
 *
 * @return Number of bytes written on success, -1 on error.
 */
int __public__sprintf(const char* fmt, ...);

/**
 * @brief Synchronously read n bytes from a file at a given offset.
 *
 * Reads up to n bytes from the specified file at the given offset using
 * the pread syscall. Handles errors and partial reads gracefully.
 *
 * @param f      FILE pointer to read from.
 * @param buf    Buffer to store read data.
 * @param n      Number of bytes to read.
 * @param offset File offset to start reading from.
 *
 * @return Number of bytes read on success, -1 on error.
 */
int __public__sfread(char* f, char* buf, int n, int offset);

/**
 * @brief Synchronously write n bytes to a file at a given offset.
 *
 * Writes up to n bytes to the specified file at the given offset using
 * the pwrite syscall. Handles errors and partial writes gracefully.
 *
 * @param f      FILE pointer to write to.
 * @param buf    Buffer containing data to write.
 * @param n      Number of bytes to write.
 * @param offset File offset to start writing from.
 *
 * @return Number of bytes written on success, -1 on error.
 */
int __public__sfwrite(char* f, char* buf, int n, int offset);

#endif