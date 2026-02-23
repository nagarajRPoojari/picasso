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
    #define CYCLES 100
    #define TASKS_PER_CYCLE 50
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define CYCLES 1000
    #define TASKS_PER_CYCLE 100
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define CYCLES 20000
    #define TASKS_PER_CYCLE 1000
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define CYCLES 1000
    #define TASKS_PER_CYCLE 1000
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define MAX_TIMEOUT ((CYCLES * TASKS_PER_CYCLE / 1000) * 50 + 100)

static atomic_int task_counter;

void short_task(void* arg) {
    (void)arg;
    atomic_fetch_add(&task_counter, 1);
}

void test_stress_rapid_cycles(void) {
    atomic_store(&task_counter, 0);
    int total_expected = CYCLES * TASKS_PER_CYCLE;
    
    printf("Spawning %d tasks rapidly (simulating %d cycles of %d tasks)...\n",
           total_expected, CYCLES, TASKS_PER_CYCLE);
    
    /* Spawn all tasks upfront to test rapid task creation */
    for (int i = 0; i < total_expected; i++) {
        thread(short_task, 1, NULL);
        
        /* Small delay every cycle to simulate rapid bursts */
        if ((i + 1) % TASKS_PER_CYCLE == 0) {
            usleep(100); /* 0.1ms between bursts */
        }
    }
    
    /* Wait for all tasks to complete */
    struct timespec ts = {0, 10000000}; /* 10ms */
    int timeout = 0;
    while (atomic_load(&task_counter) < total_expected && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int completed = atomic_load(&task_counter);
    printf("Completed: %d/%d tasks (%.1f%%)\n",
           completed, total_expected, (completed * 100.0) / total_expected);
    
    TEST_ASSERT_EQUAL_INT(total_expected, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_rapid_cycles);
    
    /* Scheduler shuts down automatically when all tasks complete */
    wait_for_schedulers();
    
    return UNITY_END();
}
