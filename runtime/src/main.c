
#define _GNU_SOURCE
#include "initutils.h"

extern kernel_thread_t **kernel_thread_map;
extern struct io_uring **io_ring_map;

extern pthread_t sched_threads[SCHEDULER_THREAD_POOL_SIZE];


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

    thread(start, NULL);
    gc_init();

    for (int i = 0; i < SCHEDULER_THREAD_POOL_SIZE; i++) {
        pthread_join(sched_threads[i], NULL);
    }

    clean_scheduler();
    return 0;
}
