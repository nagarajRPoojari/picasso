
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

#include <gc.h>
#include <gc/gc.h> 

#include "start.h"
#include "array.h"
#include "ggc.h"
#include "io.h"
#include "queue.h"
#include "scheduler.h"
#include "task.h"



kernel_thread_t **kernel_thread_map;

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
    safe_q_init(&io_queue, IO_QUEUE_SIZE);
    
    epfd = epoll_create1(0);
    if (epfd == -1) { perror("epoll_create1"); return 1; }
    
    pthread_t io_threads[IO_THREAD_POOL_SIZE];
    // Thread pool
    for (int i=0;i<IO_THREAD_POOL_SIZE;i++){
        pthread_create(&io_threads[i], NULL, io_worker, NULL);
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
    kernel_thread_map = calloc(4, sizeof(kernel_thread_t*));
    for (int i=0;i<SCHEDULER_THREAD_POOL_SIZE;i++) {
        kernel_thread_map[i] = calloc(1, sizeof(kernel_thread_t));
        kernel_thread_map[i]->id = i;
        kernel_thread_map[i]->current = NULL;
        safe_queue_t ready_q;
        safe_q_init(&ready_q, SCHEDULER_LOCAL_QUEUE_SIZE);

        kernel_thread_map[i]->ready_q = ready_q;
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
    free(kernel_thread_map[0]);
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

    GC_INIT();
    GC_allow_register_threads(); 

    init_io();
    init_scheduler();

    thread(start, NULL);

    for (int i = 0; i < SCHEDULER_THREAD_POOL_SIZE; i++) {
        pthread_join(sched_threads[i], NULL);
    }

    clean_scheduler();
    return 0;
}
