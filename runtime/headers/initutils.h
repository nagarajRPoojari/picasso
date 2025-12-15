#ifndef INITUTILS_H
#define INITUTILS_H


#include <ucontext.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <pthread.h>
#include <sys/epoll.h>
#include <signal.h>
#include <liburing.h>
#include <string.h> 


#include "start.h"
#include "array.h"
#include "ggc.h"
#include "io.h"
#include "queue.h"
#include "scheduler.h"
#include "task.h"
#include "crypto.h"
#include "str.h"
#include "alloc.h"
#include "gc.h"

extern kernel_thread_t **kernel_thread_map;
extern struct io_uring **io_ring_map;

extern pthread_t sched_threads[SCHEDULER_THREAD_POOL_SIZE];

/**
 * @brief Create and schedule a new task on a random scheduler thread.
 * 
 * Allocates a task with its own stack and context, assigns it a random
 * ID, and pushes it onto a scheduler thread's ready queue.
 * 
 * @param fn   Function pointer for the task to execute.
 * @param this Argument to pass to the task function.
 */
void thread(void*(*fn)(void*), void *this);

/**
 * @brief Initialize the I/O subsystem.
 * 
 * - Initializes the global I/O queue.
 * - Creates an epoll instance for monitoring file descriptors.
 * - Launches a pool of I/O worker threads.
 * 
 * @return 0 on success, 1 on failure.
 */
int init_io();

/**
 * @brief Initialize scheduler threads.
 * 
 * - Allocates and initializes kernel_thread_t structures.
 * - Initializes each scheduler's local ready queue.
 * - Creates threads running the scheduler_run() loop.
 * 
 * @return 0 on success.
 */
int init_scheduler();


/**
 * @brief Cleanup resources used by the scheduler.
 * 
 * Currently frees only the first kernel thread. In production, all
 * threads and queues should be properly deallocated.
 */
void clean_scheduler();


#endif