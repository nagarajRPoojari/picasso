#include "platform.h"
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
#include <assert.h>

#include "array.h"
#include "diskio.h"
#include "queue.h"
#include "task.h"
#include "scheduler.h"
#include "alloc.h"

extern __thread arena_t* __arena__;
extern arena_t* __global__arena__;
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
void *diskio_worker(void *arg) {
    int id = (int)(intptr_t)arg;

    struct io_uring_cqe *cqe;

    for (;;) {
        int ret = io_uring_wait_cqe(diskio_ring_map[id], &cqe);
        if (ret < 0)
            continue;

        task_t *t = io_uring_cqe_get_data(cqe);

        if (cqe->res < 0) {
            t->io.io_err = -cqe->res;
            t->io.done_n = -1;
        } else {
            t->io.done_n = cqe->res;
            t->io.io_err = 0;
        }

        t->io.io_done = 1;

        io_uring_cqe_seen(diskio_ring_map[id], cqe);

        safe_q_push(&kernel_thread_map[t->sched_id]->ready_q, t);
    }
}


/**
 * @brief Submit an asynchronous STDIN read for the current task.
 *
 * Prepares and submits an io_uring read request on STDIN for the
 * current task’s buffer. The read follows normal read(2) semantics:
 * it may complete with fewer bytes than requested and may block
 * internally if STDIN is a TTY.
 *
 * This function does NOT modify STDIN file descriptor flags
 * (e.g. O_NONBLOCK) and does not provide true asynchronous behavior
 * for terminal input. The calling task explicitly yields and will be
 * resumed by the scheduler once the io_uring completion is processed
 * by the I/O worker thread.
 *
 * Completion status (bytes read or error) is stored in the task
 * structure and must be examined after the task resumes.
 *
 * @return void
 */
void async_stdin_read(task_t *t) {
    struct io_uring *ring = diskio_ring_map[t->sched_id];
    struct io_uring_sqe *sqe;

    while ((sqe = io_uring_get_sqe(ring)) == NULL) {
        task_yield(kernel_thread_map[t->sched_id]);
    }

    io_uring_prep_read(
        sqe,
        STDIN_FILENO,
        t->io.buf,
        t->io.req_n,
        0   /* ignored for stdin */
    );

    io_uring_sqe_set_data(sqe, t);

    unsafe_q_push(&kernel_thread_map[t->sched_id]->wait_q, t);

    int ret = io_uring_submit(ring);
    if (ret < 0) {
        t->io.io_err = -ret;
        t->io.done_n = -1;
        t->io.io_done = 1;
        return;
    }



    task_yield(kernel_thread_map[t->sched_id]);
}



/**
 * @brief Submit an io_uring write request to STDOUT for the current task.
 *
 * Prepares and submits a write request for the current task’s buffer
 * to STDOUT at the current file offset. The task yields execution
 * to allow other tasks to run while the I/O completes.
 *
 * Notes:
 *  - Partial writes are possible; actual bytes written are stored in
 *    current_task->done_n after completion.
 *  - I/O errors are reported via current_task->io_err.
 *  - The task will be resumed by the scheduler once the I/O worker
 *    processes the completion.
 *  - True asynchronous non-blocking behavior depends on the underlying file descriptor
 *    (TTY writes may still block internally).
 *  - The function does **not** exit the process on SQE allocation or submission failure;
 *    errors are propagated to the task structure instead.
 *
 * @return void
 */
void async_stdout_write() {
    task_t *t = current_task;
    struct io_uring_sqe *sqe;
    struct io_uring *ring = diskio_ring_map[t->sched_id];

    /* acquire SQE */
    while ((sqe = io_uring_get_sqe(ring)) == NULL) {
        task_yield(kernel_thread_map[t->sched_id]);
    }

    /* prepare write (stdout is not seekable) */
    io_uring_prep_write(sqe, STDOUT_FILENO, t->io.buf, t->io.req_n, 0);
    io_uring_sqe_set_data(sqe, t);

    unsafe_q_push(&kernel_thread_map[t->sched_id]->wait_q, t);

    int ret = io_uring_submit(ring);
    if (ret < 0) {
        t->io.io_err = -ret;
        t->io.done_n = -1;
        t->io.io_done = 1;
        return;
    }

    task_yield(kernel_thread_map[t->sched_id]);
}



/**
 * @brief Submit an io_uring read request on a file descriptor for the current task.
 *
 * Prepares and submits a read request for the current task’s buffer
 * from the specified file descriptor at the given offset. The task
 * yields execution to allow other tasks to run while the I/O completes.
 *
 * Notes:
 *  - Partial reads are possible; the actual number of bytes read is
 *    stored in current_task->done_n after completion.
 *  - I/O errors are reported via current_task->io_err.
 *  - The task will be resumed by the scheduler once the I/O worker
 *    processes the completion.
 *
 * @return void
 */
void async_file_read() {
    task_t *t = current_task;
    struct io_uring *ring = diskio_ring_map[t->sched_id];
    struct io_uring_sqe *sqe;

    /* wait for a submission queue entry */
    while ((sqe = io_uring_get_sqe(ring)) == NULL) {
        task_yield(kernel_thread_map[t->sched_id]);
    }

    /* prepare read at specified offset */
    io_uring_prep_read(sqe,
                       t->io.fd,
                       t->io.buf,
                       t->io.req_n,
                       t->io.offset);

    io_uring_sqe_set_data(sqe, t);

    unsafe_q_push(&kernel_thread_map[t->sched_id]->wait_q, t);

    int ret = io_uring_submit(ring);
    if (ret < 0) {
        t->io.io_err  = -ret;
        t->io.done_n  = -1;
        t->io.io_done = 1;
    }

    task_yield(kernel_thread_map[t->sched_id]);
}


/**
 * @brief Submit an io_uring write request on a file descriptor for the current task.
 *
 * Prepares and submits a write request for the current task’s buffer
 * to the specified file descriptor at the given offset. The task
 * yields execution to allow other tasks to run while the I/O completes.
 *
 * Notes:
 *  - Partial writes may occur; actual bytes written are stored in
 *    current_task->done_n.
 *  - I/O errors are reported via current_task->io_err.
 *  - The task will be resumed by the scheduler once the I/O worker
 *    processes the completion.
 *  - True asynchronous non-blocking behavior depends on the file type.
 *
 * @return void
 */
void async_file_write() {
    task_t *t = current_task;
    struct io_uring *ring = diskio_ring_map[t->sched_id];
    struct io_uring_sqe *sqe;

    /* wait for a submission queue entry */
    while ((sqe = io_uring_get_sqe(ring)) == NULL) {
        task_yield(kernel_thread_map[t->sched_id]);
    }

    /* prepare write at specified offset */
    io_uring_prep_write(
        sqe,
        t->io.fd,
        t->io.buf,
        t->io.req_n,
        t->io.offset
    );

    io_uring_sqe_set_data(sqe, t);

    unsafe_q_push(&kernel_thread_map[t->sched_id]->wait_q, t);

    int ret = io_uring_submit(ring);
    if (ret < 0) {
        t->io.io_err  = -ret;
        t->io.done_n  = -1;
        t->io.io_done = 1;
        return;
    }

    task_yield(kernel_thread_map[t->sched_id]);
}



/**
 * @brief Read up to n bytes from STDIN and suspend the current task.
 *
 * Allocates a buffer and submits an io_uring read request on STDIN for the
 * current task. The task yields execution and is resumed by the scheduler
 * once the I/O completion is processed by the I/O worker thread.
 *
 * The read follows read(2) semantics: it may complete with fewer than n
 * bytes and does not guarantee true asynchronous behavior for terminal
 * input. The returned buffer is NUL-terminated based on the number of
 * bytes actually read.
 *
 * @param n Maximum number of bytes to read.
 *
 * @return Pointer to the allocated buffer on success, or NULL on error.
 */
__public__array_t* __public__ascan(int n) {
    if (n <= 0)
        return NULL;

    task_t *t = current_task;

    /* +1 for NUL */
    __public__array_t* buf = __public__alloc_array((size_t)n + 1, sizeof(size_t), 1);
    if (!buf)
        return NULL;

    /* setup task I/O state */
    t->io = (io_metadata_t){
        .fd = STDIN_FILENO,
        .buf = buf->data,
        .req_n = n,
        .done_n  = 0,
        .io_err  = 0,
        .io_done = 0,
    };

    /* submit async read */
    async_stdin_read(t);

    /* resumed after io_worker enqueues us */
    if (t->io.done_n < 0)
        return NULL;

    /* NUL terminate */
    if ((size_t)t->io.done_n < (size_t)n)
        buf->data[t->io.done_n] = '\0';
    else
        buf->data[n] = '\0';

    return buf;
}


/**
 * @brief Format and write output to STDOUT, suspending the current task until done.
 *
 * Formats the input string and arguments into a dynamically allocated buffer,
 * sets up the current task for an io_uring write to STDOUT, and submits it.
 * The task yields execution and will be resumed by the scheduler once the
 * write completes.
 *
 * Note:
 *  - The write may be partial; actual bytes written are stored in
 *    current_task->done_n.
 *  - The returned buffer is owned by the task/runtime; do not free manually.
 *  - True asynchronous non-blocking behavior is only possible for
 *    non-TTY fds.
 *
 * @param fmt Format string (printf-style).
 * @param ... Arguments matching the format string.
 *
 * @return Number of bytes successfully written on success, or -1 on error.
 */
ssize_t __public__aprintf(const char* fmt, ...) {
    if (!fmt) return NULL;

    va_list ap;
    va_start(ap, fmt);

    /* estimate needed buffer size */
    char tmp[1];
    int len = vsnprintf(tmp, sizeof(tmp), fmt, ap);
    va_end(ap);

    if (len < 0) return NULL;

    char* buf = allocate(__arena__, len + 1);
    if (!buf) return NULL;

    va_start(ap, fmt);
    vsnprintf(buf, len + 1, fmt, ap);
    va_end(ap);

    task_t *t = current_task;

    /* setup task for stdout write */
    t->io = (io_metadata_t){
        .fd = STDOUT_FILENO,
        .buf = buf,
        .req_n = len,
        .done_n = 0,
        .io_err = 0,
        .io_done = 0,
    };

    async_stdout_write();

    /* task resumes here after write completion */
    if (t->io.done_n < 0)
        return NULL;

    return t->io.done_n;
}


/**
 * @brief Read up to n bytes from a file at a given offset, suspending the current task.
 *
 * Configures the current task with the file descriptor, buffer, byte count,
 * and offset, then submits an io_uring read request. The task yields execution
 * and will be resumed by the scheduler once the read completes.
 *
 * Notes:
 *  - Partial reads are possible; actual bytes read are stored in
 *    current_task->done_n.
 *  - I/O errors are reported via current_task->io_err.
 *  - True asynchronous non-blocking behavior depends on the file type.
 *
 * @param f      FILE pointer to read from.
 * @param buf    Buffer to store read data.
 * @param n      Maximum number of bytes to read.
 * @param offset File offset to start reading from.
 *
 * @return Number of bytes read on success (ssize_t), or -1 on error.
 */
ssize_t __public__afread(char* f, __public__array_t* buf, int n, int offset) {
    if (!f || !buf || n <= 0 || offset < 0)
        return -1;

    int fd = fileno((FILE*)f);
    if (fd < 0)
        return -1;

    task_t *t = current_task;

    /* setup task I/O state */
    t->io = (io_metadata_t){
        .fd      = fd,
        .buf     = buf->data,
        .req_n   = n,
        .offset  = offset,
        .done_n  = 0,
        .io_err  = 0,
        .io_done = 0,
    };

    /* submit async read */
    async_file_read();

    /* task resumes here after I/O completion */
    if (t->io.done_n < 0)
        return -1;

    return t->io.done_n;
}


/**
 * @brief Write up to n bytes to a file at a given offset, suspending the current task.
 *
 * Configures the current task with the file descriptor, buffer, byte count,
 * and offset, then submits an io_uring write request. The task yields execution
 * and will be resumed by the scheduler once the write completes.
 *
 * Notes:
 *  - Partial writes are possible; actual bytes written are stored in
 *    current_task->done_n.
 *  - I/O errors are reported via current_task->io_err.
 *  - True asynchronous non-blocking behavior depends on the file type.
 *
 * @param f      FILE pointer to write to.
 * @param buf    Buffer containing data to write.
 * @param n      Maximum number of bytes to write.
 * @param offset File offset to start writing from.
 *
 * @return Number of bytes written on success (ssize_t), or -1 on error.
 */
ssize_t __public__afwrite(char* f, __public__array_t* buf, int n, int offset) {
    if (!f || !buf || n <= 0 || offset < 0)
        return -1;

    int fd = fileno((FILE*)f);
    if (fd < 0)
        return -1;

    task_t *t = current_task;

    /* setup task I/O state */
    t->io = (io_metadata_t){
        .fd      = fd,
        .buf     = (char*)buf->data,
        .req_n   = n,
        .offset  = offset,
        .done_n  = 0,
        .io_err  = 0,
        .io_done = 0,
    };

    /* submit async write */
    async_file_write();

    /* task resumes here after I/O completion */
    if (t->io.done_n < 0)
        return -1;

    return t->io.done_n;
}


/**
 * @brief Synchronously read up to n bytes from STDIN.
 *
 * Allocates a buffer of size n+1 and performs a blocking read() from STDIN.
 * Notes on behavior:
 *  - When reading from a TTY, input is typically line-buffered. The read() syscall
 *    may return fewer than n bytes as soon as a line is available.
 *  - If read() is interrupted by a signal (errno == EINTR), the function retries
 *    the read automatically.
 *  - The returned buffer is NUL-terminated at the position corresponding to
 *    the number of bytes actually read.
 *
 * @param n Maximum number of bytes to read.
 *
 * @return Pointer to the allocated buffer containing the read data on success,
 *         or NULL on allocation failure or read error.
 */
__public__array_t* __public__sscan(int n) {
    if (n <= 0) return NULL;

    __public__array_t* buf = __public__alloc_array((size_t)n + 1, sizeof(size_t), 1);
    if (!buf) return NULL;

    ssize_t r;

    for (;;) {
        r = read(STDIN_FILENO, buf->data, (size_t)n);
        if (r < 0) {
            if (errno == EINTR)
                continue;
            return NULL;
        }
        break;
    }

    buf->data[r] = '\0';
    return buf;
}


/**
 * @brief Synchronously write formatted output to STDOUT.
 *
 * Formats the input string and arguments into a dynamically allocated buffer,
 * then writes the entire result to STDOUT using the write() syscall. The function
 * ensures that all bytes are written, retrying as needed in case of partial writes
 * or EINTR interruptions.
 *
 * Notes:
 *  - The buffer is NUL-terminated internally for formatting purposes,
 *    but the NUL byte is not written to STDOUT.
 *  - The function may perform multiple write() syscalls to complete the output.
 *
 * @param fmt Format string (printf-style).
 * @param ... Arguments matching the format string.
 *
 * @return Number of bytes successfully written on success, or -1 on error.
 */
ssize_t __public__sprintf(const char *fmt, ...) {
    if (!fmt)
        return -1;

    va_list ap;
    va_start(ap, fmt);

    /* get required length (excluding NUL) */
    int len = vsnprintf(NULL, 0, fmt, ap);
    va_end(ap);

    if (len < 0)
        return -1;

    char *buf = allocate(__arena__, (size_t)len + 1);
    if (!buf)
        return -1;

    va_start(ap, fmt);
    vsnprintf(buf, (size_t)len + 1, fmt, ap);
    va_end(ap);

    /* write all bytes */
    size_t total = 0;
    while (total < (size_t)len) {
        ssize_t w = write(STDOUT_FILENO, buf + total, (size_t)len - total);
        if (w < 0) {
            if (errno == EINTR)
                continue;
            return -1;
        }
        total += (size_t)w;
    }

    return (ssize_t)total;
}

/**
 * @brief Synchronously read up to n bytes from a file at a given offset.
 *
 * Reads from the specified FILE pointer into the provided buffer using the
 * pread() syscall. The function handles partial reads and EINTR interruptions,
 * and continues reading until either n bytes are read or end-of-file is reached.
 *
 * Notes:
 *  - The file descriptor offset is updated internally after each read.
 *  - Partial reads are handled transparently; the function loops until
 *    the requested byte count is satisfied or EOF occurs.
 *  - The buffer is not NUL-terminated; it contains raw file data.
 *
 * @param f      FILE pointer to read from.
 * @param buf    Buffer to store read data.
 * @param n      Maximum number of bytes to read.
 * @param offset File offset to start reading from.
 *
 * @return Number of bytes actually read on success (0 indicates EOF),
 *         or -1 on error.
 */
ssize_t __public__sfread(char* f, __public__array_t* buf, int n, int offset) {
    if (!f || !buf || n <= 0 || offset < 0) return -1;
    
    int fd = fileno((FILE*)f);
    if (fd < 0) return -1;
    
    size_t total = 0;
    char *p = buf->data;

    while (total < n) {
        ssize_t r = pread(fd, p + total, n - total, offset);
        if (r < 0) {
            if (errno == EINTR) continue;
            return -1;
        }
        if (r == 0) break; // EOF
        total += r;
        offset += r;
    }
    return total;
}

/**
 * @brief Synchronously write up to n bytes to a file at a given offset.
 *
 * Writes data from the provided buffer to the specified FILE pointer using
 * the pwrite() syscall. The function handles partial writes and EINTR
 * interruptions, looping until either all n bytes are written or an error occurs.
 *
 * Notes:
 *  - The file offset is updated internally for each pwrite() call.
 *  - Partial writes are handled transparently; the function ensures that
 *    the full requested byte count is written unless an error occurs.
 *  - The buffer is not modified by the function.
 *
 * @param f      FILE pointer to write to.
 * @param buf    Buffer containing data to write.
 * @param n      Number of bytes to write.
 * @param offset File offset to start writing from.
 *
 * @return Number of bytes actually written on success,
 *         or -1 on error.
 */
ssize_t __public__sfwrite(char* f, __public__array_t* buf, int n, int offset) {
    if (!f || !buf || n <= 0 || offset < 0) return -1;
    
    int fd = fileno((FILE*)f);
    if (fd < 0) return -1;
    
    size_t total = 0;
    const char *p = buf->data;

    while (total < n) {
        ssize_t w = pwrite(fd, p + total, n - total, offset + total);
        if (w < 0) {
            if (errno == EINTR)
                continue;
            return -1;
        }
        total += (size_t)w;
    }

    return (ssize_t)total;
}