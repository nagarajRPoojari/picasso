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
    #define TASK_COUNT 1000
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define TASK_COUNT 5000
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define TASK_COUNT 20000
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define TASK_COUNT 10000
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define MAX_TIMEOUT ((TASK_COUNT / 1000) * 50 + 100)

static atomic_int task_counter;

void short_task(void* arg) {
    (void)arg;
    atomic_fetch_add(&task_counter, 1);
}

void test_stress_many_short_tasks(void) {
    atomic_store(&task_counter, 0);
    
    printf("Stress Level: %s\n", STRESS_LEVEL);
    printf("Spawning %d short tasks...\n", TASK_COUNT);
    
    for (int i = 0; i < TASK_COUNT; i++) {
        thread(short_task, 1, NULL);
    }
    
    /* Wait for completion */
    struct timespec ts = {0, 10000000}; /* 10ms */
    int timeout = 0;
    while (atomic_load(&task_counter) < TASK_COUNT && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int completed = atomic_load(&task_counter);
    printf("Completed: %d/%d tasks\n", completed, TASK_COUNT);
    
    TEST_ASSERT_EQUAL_INT(TASK_COUNT, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_many_short_tasks);
    
    /* Scheduler shuts down automatically when all tasks complete */
    wait_for_schedulers();
    
    return UNITY_END();
}
