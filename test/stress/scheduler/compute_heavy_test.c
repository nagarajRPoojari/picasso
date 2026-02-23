#include "unity/unity.h"
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdatomic.h>
#include <time.h>
#include "alloc.h"
#include "gc.h"
#include "initutils.h"
#include "scheduler.h"

extern __thread arena_t* __arena__;
extern arena_t* __global__arena__;
extern atomic_int task_count;
extern kernel_thread_t **kernel_thread_map;

void setUp(void) {}
void tearDown(void) {}

/* Stress levels - configurable via compile-time macros
 * Define one of: STRESS_LEVEL_1, STRESS_LEVEL_2, STRESS_LEVEL_3, STRESS_LEVEL_4
 * Or use default (STRESS_LEVEL_3)
 */
#if defined(STRESS_LEVEL_1)
    #define COMPUTE_TASKS 1000
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define COMPUTE_TASKS 5000
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define COMPUTE_TASKS 20000
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define COMPUTE_TASKS 10000
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define MAX_TIMEOUT ((COMPUTE_TASKS / 1000) * 50 + 100)

static atomic_int compute_counter;

void compute_heavy_task(void* arg) {
    int iterations = (int)(uintptr_t)arg;
    volatile uint64_t sum = 0;
    
    for (int i = 0; i < iterations; i++) {
        sum += i * i;
    }
    
    atomic_fetch_add(&compute_counter, 1);
}

void test_stress_compute_heavy_tasks(void) {
    atomic_store(&compute_counter, 0);

    printf("Spawning %d compute-heavy tasks...\n", COMPUTE_TASKS);
    
    for (int i = 0; i < COMPUTE_TASKS; i++) {
        thread(compute_heavy_task, 1, (void*)(uintptr_t)10000);
    }
    
    /* Wait for completion with longer timeout */
    struct timespec ts = {0, 100000000}; /* 100ms */
    int timeout = 0;
    while (atomic_load(&compute_counter) < COMPUTE_TASKS && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int completed = atomic_load(&compute_counter);
    printf("Completed: %d/%d compute tasks\n", completed, COMPUTE_TASKS);
    
    TEST_ASSERT_EQUAL_INT(COMPUTE_TASKS, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_compute_heavy_tasks);
    
    /* Scheduler shuts down automatically when all tasks complete */
    wait_for_schedulers();
    
    return UNITY_END();
}
