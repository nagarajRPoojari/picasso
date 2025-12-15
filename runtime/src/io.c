#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>
#include <sys/epoll.h>
#include <liburing.h>
#include <stdarg.h>
#include <string.h>

#include "io.h"
#include "queue.h"
#include "task.h"
#include "scheduler.h"
#include "alloc.h"

extern __thread arena_t* __arena__;

/**
 * @brief Worker thread that waits for completed I/O events from io_uring.
 *
 * This thread continuously polls for completion events from the io_uring
 * submission queue. When an event completes, it retrieves the associated
 * task (stored in CQE user data), updates its `done_n` field with the result,
 * and pushes it back into the scheduler's ready queue.
 *
 * @param arg Unused.
 * @return void* Always returns NULL.
 */
void *io_worker(void *arg) {
    int id = (int)(intptr_t)arg;
    int ret = io_uring_queue_init(QUEUE_DEPTH, io_ring_map[id], 0);
    if (ret < 0) {
        perror("io_uring_queue_init");
        exit(1);
    }
    struct io_uring_cqe *cqe;

    while (1) {
        int ret = io_uring_wait_cqe(io_ring_map[id], &cqe);
        if (ret == 0) {
            task_t *t = io_uring_cqe_get_data(cqe);
            t->done_n = cqe->res;
            io_uring_cqe_seen(io_ring_map[id], cqe);
            safe_q_push(&(kernel_thread_map[t->sched_id]->ready_q), t);
        }
    }
    return NULL;
}

/**
 * @brief Internal helper to asynchronously read from STDIN using io_uring.
 *
 * Prepares and submits a read request for the current task’s buffer from STDIN.
 * STDIN is set to non-blocking mode for safety. Upon submission, the current
 * task yields control to the scheduler until the I/O operation completes.
 *
 * @return void
 */
void async_stdin_read() {
    struct io_uring_sqe *sqe;
    while ((sqe = io_uring_get_sqe(io_ring_map[current_task->sched_id])) == NULL) {
        task_yield(kernel_thread_map[current_task->sched_id]);
    }

    /** Make stdin non-blocking for safety */
    fcntl(STDIN_FILENO, F_SETFL, O_NONBLOCK);

    io_uring_prep_read(
        sqe,
        STDIN_FILENO,            /** fd */
        current_task->buf,       /** buffer */
        current_task->req_n,     /** bytes to read */
        0                        /** stdin is not seekable */
    );

    io_uring_sqe_set_data(sqe, current_task);

    int ret = io_uring_submit(io_ring_map[current_task->sched_id]);
    if (ret < 0) {
        fprintf(stderr, "io_uring_submit: %s\n", strerror(-ret));
        exit(1);
    }

    task_yield(kernel_thread_map[current_task->sched_id]);
}


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
void async_stdout_write() {
    struct io_uring_sqe *sqe;
    while ((sqe = io_uring_get_sqe(io_ring_map[current_task->sched_id])) == NULL) {
        task_yield(kernel_thread_map[current_task->sched_id]);
    }

    /** set stdout non blocking */
    fcntl(STDOUT_FILENO, F_SETFL, O_NONBLOCK);

    io_uring_prep_write(
        sqe,
        STDOUT_FILENO,
        current_task->buf,
        current_task->req_n,
        0  /** stdout is not seekable */
    );

    io_uring_sqe_set_data(sqe, current_task);

    int ret = io_uring_submit(io_ring_map[current_task->sched_id]);
    if (ret < 0) {
        fprintf(stderr, "io_uring_submit (stdout): %s\n", strerror(-ret));
        exit(1);
    }

    task_yield(kernel_thread_map[current_task->sched_id]);
}


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
void async_file_read() {
    struct io_uring_sqe *sqe;
    while ((sqe = io_uring_get_sqe(io_ring_map[current_task->sched_id])) == NULL) {
        task_yield(kernel_thread_map[current_task->sched_id]);
    }

    io_uring_prep_read(sqe,
                       current_task->fd,
                       current_task->buf,
                       current_task->req_n,
                       current_task->offset  
                    );

    io_uring_sqe_set_data(sqe, current_task);
    io_uring_submit(io_ring_map[current_task->sched_id]);

    task_yield(kernel_thread_map[current_task->sched_id]);
}

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
void async_file_write() {
    struct io_uring_sqe *sqe;
    while ((sqe = io_uring_get_sqe(io_ring_map[current_task->sched_id])) == NULL) {
        task_yield(kernel_thread_map[current_task->sched_id]);
    }

    io_uring_prep_write(
        sqe,
        current_task->fd,
        current_task->buf,
        current_task->req_n,
        current_task->offset
    );

    io_uring_sqe_set_data(sqe, current_task);

    int ret = io_uring_submit(io_ring_map[current_task->sched_id]);
    if (ret < 0) {
        fprintf(stderr, "io_uring_submit (write): %s\n", strerror(-ret));
        exit(1);
    }

    task_yield(kernel_thread_map[current_task->sched_id]);
}


/**
 * @brief Asynchronously read n bytes from STDIN.
 *
 * Allocates a buffer, configures the current task with the read parameters,
 * and submits an io_uring read request for STDIN. Yields until completion.
 *
 * @param n Number of bytes to read.
 *
 * @return Pointer to allocated buffer containing the read data.
 */
void* __public__ascan(int n) {
    char* buf = (char*)allocate(__arena__, n * sizeof(char));
    current_task->fd = STDIN_FILENO;
    current_task->buf = buf;
    current_task->req_n = n;
    async_stdin_read();
    return buf;
}

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
void* __public__aprintf(const char* fmt, ...) {
    va_list ap;
    va_start(ap, fmt);
    
    /** estimate needed size */
    char tmp[1];
    int len = vsnprintf(tmp, sizeof(tmp), fmt, ap);
    va_end(ap);
    
    if (len < 0) return NULL;
    
    char* buf = allocate(__arena__, len + 1);
    if (!buf) return NULL;
    
    va_start(ap, fmt);
    vsnprintf(buf, len + 1, fmt, ap);
    va_end(ap);
    
    current_task->fd = STDOUT_FILENO;
    current_task->buf = (char*)buf;
    current_task->req_n = len;
    
    async_stdout_write();

    return NULL;
}

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
void* __public__afread(char* f, char* buf, int n, int offset) {
    int fd = fileno((FILE*)f);
    current_task->fd = fd;
    current_task->buf = buf;
    current_task->req_n = n;
    current_task->offset = offset;
    async_file_read();
    return current_task->done_n; /* @todo: not tested return address */
}

/**
 * @brief Asynchronously write n bytes to a file at a given offset.
 *
 * Configures the current task with the file descriptor, buffer, byte count,
 * and offset, then submits an io_uring write request. Yields until completion.
 *
 * @param f      FILE pointer to write to.
 * @param buf    Buffer containing data to write.
 * @param n      Number of bytes to write.
 * @param offset File offset to start writing from.
 *
 * @return Pointer to bytes written count in current task context.
 */
void* __public__afwrite(char* f, char* buf, int n, int offset) {
    int fd = fileno((FILE*)f);
    current_task->fd = fd;
    current_task->buf = (char*)buf;
    current_task->req_n = n;
    current_task->offset = offset;

    async_file_write();
    return current_task->done_n; /* @todo: not tested return address */
}

/**
 * @brief Synchronously read n bytes from STDIN.
 *
 * uses blocking read() syscall to read n bytes. while reading from tty Input is line-buffered by 
 * the kernel read() typically returns as soon as a line is available. so it doesn't wait till n bytes 
 * are availanbe. __public__sscan can do multiple read() calls to read in case of EINTR errors.
 * @param n Number of bytes to read.
 *
 * @return Pointer to allocated buffer containing the read data.
 */
void* __public__sscan(int n) {
    if (n <= 0) return NULL;

    char *buf = allocate(__arena__, (size_t)n + 1);
    if (!buf) return NULL;

    ssize_t r;

    for (;;) {
        r = read(STDIN_FILENO, buf, (size_t)n);
        if (r < 0) {
            if (errno == EINTR)
                continue;
            return NULL;
        }
        break;
    }

    buf[r] = '\0';
    return buf;
}


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
int __public__sprintf(const char* fmt, ...) {
    if (!fmt) return -1;
    
    va_list ap;
    va_start(ap, fmt);
    
    /* estimate needed size */
    char tmp[1];
    int len = vsnprintf(tmp, sizeof(tmp), fmt, ap);
    va_end(ap);
    
    if (len < 0) return -1;
    
    char* buf = allocate(__arena__, len + 1);
    if (!buf) return -1;
    
    va_start(ap, fmt);
    vsnprintf(buf, len + 1, fmt, ap);
    va_end(ap);
    
    ssize_t bytes_written = write(STDOUT_FILENO, buf, len);
    
    if (bytes_written < 0) {
        return -1;
    }
    
    return bytes_written;
}

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
int __public__sfread(char* f, char* buf, int n, int offset) {
    if (!f || !buf || n <= 0 || offset < 0) return -1;
    
    int fd = fileno((FILE*)f);
    if (fd < 0) return -1;
    
    ssize_t bytes_read = pread(fd, buf, n, offset);
    if (bytes_read < 0) {
        return -1;
    }
    
    return bytes_read;
}

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
int __public__sfwrite(char* f, char* buf, int n, int offset) {
    if (!f || !buf || n <= 0 || offset < 0) return -1;
    
    int fd = fileno((FILE*)f);
    if (fd < 0) return -1;
    
    ssize_t bytes_written = pwrite(fd, buf, n, offset);
    if (bytes_written < 0) {
        return -1;
    }
    
    return bytes_written;
}