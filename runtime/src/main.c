
#include "platform.h"
#include "initutils.h"

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

    gc_init();
    thread(start, 0);

    gc_start();
    wait_for_schedulers();

    clean_scheduler();
    return 0;
}
