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
    #define MIX_SHORT 100
    #define MIX_MEDIUM 20
    #define MIX_LONG 10
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define MIX_SHORT 1000
    #define MIX_MEDIUM 200
    #define MIX_LONG 100
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define MIX_SHORT 10000
    #define MIX_MEDIUM 2000
    #define MIX_LONG 1000
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define MIX_SHORT 4000
    #define MIX_MEDIUM 1000
    #define MIX_LONG 500
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif
    
#define MIX_TOTAL (MIX_SHORT + MIX_MEDIUM + MIX_LONG)
#define MAX_TIMEOUT ((MIX_TOTAL / 1000) * 50 + 100)


static atomic_int mixed_counter;

void mixed_task_short(void* arg) {
    (void)arg;
    atomic_fetch_add(&mixed_counter, 1);
}

void mixed_task_medium(void* arg) {
    (void)arg;
    volatile int sum = 0;
    for (int i = 0; i < 1000; i++) {
        sum += i;
    }
    atomic_fetch_add(&mixed_counter, 1);
}

void mixed_task_long(void* arg) {
    (void)arg;
    volatile int sum = 0;
    for (int i = 0; i < 10000; i++) {
        sum += i;
    }
    self_yield();
    atomic_fetch_add(&mixed_counter, 1);
}

void test_stress_mixed_workload(void) {
    atomic_store(&mixed_counter, 0);
    
    printf("Spawning mixed workload: %d short, %d medium, %d long tasks...\n",
           MIX_SHORT, MIX_MEDIUM, MIX_LONG);
    
    for (int i = 0; i < MIX_SHORT; i++) {
        thread(mixed_task_short, 1, NULL);
    }
    for (int i = 0; i < MIX_MEDIUM; i++) {
        thread(mixed_task_medium, 1, NULL);
    }
    for (int i = 0; i < MIX_LONG; i++) {
        thread(mixed_task_long, 1, NULL);
    }
    
    /* Wait for completion with longer timeout */
    struct timespec ts = {0, 50000000}; /* 50ms */
    int timeout = 0;
    while (atomic_load(&mixed_counter) < MIX_TOTAL && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int completed = atomic_load(&mixed_counter);
    printf("Completed: %d/%d mixed tasks\n", completed, MIX_TOTAL);
    
    TEST_ASSERT_EQUAL_INT(MIX_TOTAL, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_mixed_workload);
    
    /* Scheduler shuts down automatically when all tasks complete */
    wait_for_schedulers();
    
    return UNITY_END();
}
