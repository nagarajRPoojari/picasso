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

#define YIELD_TASKS 20
#define YIELDS_PER_TASK 5

/* Stress levels - configurable via compile-time macros
 * Define one of: STRESS_LEVEL_1, STRESS_LEVEL_2, STRESS_LEVEL_3, STRESS_LEVEL_4
 * Or use default (STRESS_LEVEL_3)
 */
#if defined(STRESS_LEVEL_1)
    #define YIELD_TASKS 20
    #define YIELDS_PER_TASK 5
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define YIELD_TASKS 200
    #define YIELDS_PER_TASK 10
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define YIELD_TASKS 5000
    #define YIELDS_PER_TASK 500
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define YIELD_TASKS 1000
    #define YIELDS_PER_TASK 50
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define MAX_TIMEOUT ((YIELD_TASKS * YIELDS_PER_TASK / 1000) * 50 + 100)

static atomic_int yield_counter;

void yielding_task(void* arg) {
    int yields = (int)(uintptr_t)arg;
    
    for (int i = 0; i < yields; i++) {
        self_yield();
    }
    
    atomic_fetch_add(&yield_counter, 1);
}

void test_stress_task_yielding(void) {
    atomic_store(&yield_counter, 0);
    
    printf("Spawning %d tasks with %d yields each...\n", YIELD_TASKS, YIELDS_PER_TASK);
    
    for (int i = 0; i < YIELD_TASKS; i++) {
        thread(yielding_task, 1, (void*)(uintptr_t)YIELDS_PER_TASK);
    }
    
    /* Wait for completion with longer timeout */
    struct timespec ts = {0, 50000000}; /* 50ms */
    int timeout = 0;
    while (atomic_load(&yield_counter) < YIELD_TASKS && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int completed = atomic_load(&yield_counter);
    printf("Completed: %d/%d yielding tasks\n", completed, YIELD_TASKS);
    
    TEST_ASSERT_EQUAL_INT(YIELD_TASKS, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_task_yielding);
    
    /* Scheduler shuts down automatically when all tasks complete */
    wait_for_schedulers();
    
    return UNITY_END();
}
