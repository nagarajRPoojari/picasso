#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>
#include <sys/epoll.h>
#include <liburing.h>
#include <stdarg.h>
#include <gc.h>

#include "io.h"
#include "queue.h"
#include "task.h"
#include "scheduler.h"


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
void _async_stdin_read() {
    struct io_uring_sqe *sqe;
    while ((sqe = io_uring_get_sqe(io_ring_map[0])) == NULL) {
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

    int ret = io_uring_submit(io_ring_map[0]);
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
void _async_stdout_write() {
    struct io_uring_sqe *sqe;
    while ((sqe = io_uring_get_sqe(io_ring_map[0])) == NULL) {
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

    int ret = io_uring_submit(io_ring_map[0]);
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
void _async_file_read() {
    struct io_uring_sqe *sqe;
    while ((sqe = io_uring_get_sqe(io_ring_map[0])) == NULL) {
        task_yield(kernel_thread_map[current_task->sched_id]);
    }

    io_uring_prep_read(sqe,
                       current_task->fd,
                       current_task->buf,
                       current_task->req_n,
                       current_task->offset  
                    );

    io_uring_sqe_set_data(sqe, current_task);
    io_uring_submit(io_ring_map[0]);

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
void _async_file_write() {
    struct io_uring_sqe *sqe;
    while ((sqe = io_uring_get_sqe(io_ring_map[0])) == NULL) {
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

    int ret = io_uring_submit(io_ring_map[0]);
    if (ret < 0) {
        fprintf(stderr, "io_uring_submit (write): %s\n", strerror(-ret));
        exit(1);
    }

    task_yield(kernel_thread_map[current_task->sched_id]);
}


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
 * @return void* Pointer to `current_task->done_n` indicating number of bytes read.
 */
void* async_stdin_read(char* buf, int n) {
    current_task->fd = STDIN_FILENO;
    current_task->buf = buf;
    current_task->req_n = n;
    _async_stdin_read();
    return &(current_task->done_n);
}

void* ascan(int n) {
    char* buf = (char*)GC_MALLOC(n * sizeof(char));
    async_stdin_read(buf, n);
    return buf;
}

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
 * @return void* Pointer to `current_task->done_n`, which holds the number of bytes written.
 *
 * @note STDOUT is set to non-blocking mode before submission.
 */
void* async_stdout_write(const char* buf, int n) {
    current_task->fd = STDOUT_FILENO;
    current_task->buf = (char*)buf;
    current_task->req_n = n;
    _async_stdout_write();
    return &(current_task->done_n);
}

void* aprintf(const char* fmt, ...) {
    va_list ap;
    va_start(ap, fmt);

    /** estimate needed size */
    char tmp[1];
    int len = vsnprintf(tmp, sizeof(tmp), fmt, ap);
    va_end(ap);

    if (len < 0) return NULL;

    char* buf = malloc(len + 1);
    if (!buf) return NULL;

    va_start(ap, fmt);
    vsnprintf(buf, len + 1, fmt, ap);
    va_end(ap);

    void* handle = async_stdout_write(buf, len);

    free(buf);
    return handle;
}

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
 * @return void*  Pointer to `current_task->done_n` indicating number of bytes read.
 */
void* async_file_read(int fd, char* buf, int n, int offset) {
    current_task->fd = fd;
    current_task->buf = buf;
    current_task->req_n = n;
    current_task->offset = offset;
    _async_file_read();
    return &(current_task->done_n);
}

void* afread(char* f, char* buf, int n, int offset) {
    int fd = fileno((FILE*)f);
    async_file_read(fd, buf, n, offset);
}

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
 * @return void*  Pointer to `current_task->done_n`, which holds the number of bytes written.
 *
 * @note Assumes `current_task` and its scheduler context are properly initialized.
 */
void* async_file_write(int fd, const char* buf, int n, int offset) {
    current_task->fd = fd;
    current_task->buf = (char*)buf;
    current_task->req_n = n;
    current_task->offset = offset;
    _async_file_write();
    return &(current_task->done_n);
}
