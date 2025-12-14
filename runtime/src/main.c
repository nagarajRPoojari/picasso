
#define _GNU_SOURCE
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

kernel_thread_t **kernel_thread_map;
struct io_uring **io_ring_map = NULL;

/**
 * @brief Create and schedule a new task on a random scheduler thread.
 * 
 * Allocates a task with its own stack and context, assigns it a random
 * ID, and pushes it onto a scheduler thread's ready queue.
 * 
 * @param fn   Function pointer for the task to execute.
 * @param this Argument to pass to the task function.
 */
void thread(void*(*fn)(void*), void *this) {
    int kernel_thread_id = rand() % SCHEDULER_THREAD_POOL_SIZE;
    task_t *t1 = task_create(fn, this, kernel_thread_map[kernel_thread_id]);
    t1->id = rand();
    safe_q_push(&(kernel_thread_map[kernel_thread_id]->ready_q), t1);
}


/**
 * @brief Initialize the I/O subsystem.
 * 
 * - Initializes the global I/O queue.
 * - Creates an epoll instance for monitoring file descriptors.
 * - Launches a pool of I/O worker threads.
 * 
 * @return 0 on success, 1 on failure.
 */
int init_io() {
    io_ring_map = calloc(IO_THREAD_POOL_SIZE, sizeof(struct io_uring*));
    if (!io_ring_map) {
        perror("calloc io_ring_map");
        exit(1);
    }

    pthread_t io_threads[IO_THREAD_POOL_SIZE];
    for (int i = 0; i < IO_THREAD_POOL_SIZE; i++) {
        io_ring_map[i] = calloc(1, sizeof(struct io_uring));
        if (!io_ring_map[i]) {
            perror("calloc ring");
            exit(1);
        }
    }

    for (int i = 0; i < IO_THREAD_POOL_SIZE; i++) {
        int rc = pthread_create(&io_threads[i], NULL, io_worker, (void*)(intptr_t)i);
        if (rc != 0) {
            fprintf(stderr, "pthread_create(%d) failed: %s\n", i, strerror(rc));
            exit(1);
        }
    }
    return 0;
}

pthread_t sched_threads[SCHEDULER_THREAD_POOL_SIZE];

/**
 * @brief Initialize scheduler threads.
 * 
 * - Allocates and initializes kernel_thread_t structures.
 * - Initializes each scheduler's local ready queue.
 * - Creates threads running the scheduler_run() loop.
 * 
 * @return 0 on success.
 */
int init_scheduler() {
    kernel_thread_map = calloc(SCHEDULER_THREAD_POOL_SIZE, sizeof(kernel_thread_t*));
    for (int i=0;i<SCHEDULER_THREAD_POOL_SIZE;i++) {
        kernel_thread_map[i] = calloc(1, sizeof(kernel_thread_t));
        kernel_thread_map[i]->id = i;
        kernel_thread_map[i]->current = NULL;
        safe_q_init(&kernel_thread_map[i]->ready_q, SCHEDULER_LOCAL_QUEUE_SIZE);

        pthread_create(&sched_threads[i], NULL, scheduler_run, kernel_thread_map[i]);
    }

    return 0;
}

/**
 * @brief Cleanup resources used by the scheduler.
 * 
 * Currently frees only the first kernel thread. In production, all
 * threads and queues should be properly deallocated.
 */
void clean_scheduler() {
    // free(kernel_thread_map[0]);
}


/**
 * @brief Program entry point.
 * 
 * - Initializes garbage collector (Boehm GC).
 * - Initializes I/O subsystem and scheduler threads.
 * - Creates the first task to run the 'start' function.
 * - Waits for all scheduler threads to complete.
 * 
 * @todo identify all task finish & return
 */
int main(void) {
    srand(time(NULL));

    init_io();
    init_scheduler();

    gc_init();
    thread(start, NULL);

    for (int i = 0; i < SCHEDULER_THREAD_POOL_SIZE; i++) {
        pthread_join(sched_threads[i], NULL);
    }

    clean_scheduler();
    return 0;
}
