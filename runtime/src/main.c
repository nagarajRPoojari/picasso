
#define _GNU_SOURCE
#include "initutils.h"
#include "netio.h"

extern kernel_thread_t **kernel_thread_map;
extern struct io_uring **io_ring_map;

extern pthread_t sched_threads[SCHEDULER_THREAD_POOL_SIZE];
extern arena_t* __global__arena__;

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
    __global__arena__ = gc_create_global_arena();

    srand(time(NULL));

    init_io();
    init_scheduler();

    thread(start, 0);
    gc_init();

    wait_for_schedulers();

    clean_scheduler();
    return 0;
}
