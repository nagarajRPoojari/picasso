#ifndef TASK_H
#define TASK_H

#include <ucontext.h>
#include <signal.h>

/* Size of per-task I/O buffer (bytes) */
#define TASK_IO_BUFFER 256

typedef enum {
    TASK_RUNNING, 
    TASK_YIELDED,
    TASK_FINISHED
} task_state_t;

/**
 * @struct task_t
 * @brief Represents a single task/coroutine managed by the scheduler.
 * 
 * Tasks have their own stack, CPU context, and optionally perform I/O operations.
 * They are scheduled cooperatively or preemptively on a kernel thread.
 */
typedef struct task {

    /* Unique identifier for the task */
    int id;

    /* CPU context used for saving/restoring execution state */
    ucontext_t ctx;

    /* Size of the private stack (usable bytes, excluding guard page) */
    size_t stack_size;

    /* Scheduler/kernel thread ID that owns this task */
    int sched_id;

    /* Function to execute when task is scheduled */
    void* (*fn)(void *);

    /* Pointer to the task's private stack (after guard page) */
    char *stack;

    /* File descriptor if the task performs I/O (otherwise -1) */
    int fd;

    /* Buffer for I/O operations */
    char *buf;

    /* Number of bytes requested to read/write */
    ssize_t req_n;

    /* Number of bytes actually read or written */
    ssize_t done_n;

    /* seek offset */
    ssize_t offset;

    /* Flag indicating whether epoll is used for this task */
    int use_epoll; 

    /* task state: RUNNING, YIELDED, TERMINATED */
    task_state_t state;
} task_t;

#endif