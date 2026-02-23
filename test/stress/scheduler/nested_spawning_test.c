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
    #define SPAWN_DEPTH 5
    #define SPAWN_WIDTH 3
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define SPAWN_DEPTH 20
    #define SPAWN_WIDTH 10
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define SPAWN_DEPTH 500
    #define SPAWN_WIDTH 30
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define SPAWN_DEPTH 50
    #define SPAWN_WIDTH 30
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define MAX_TIMEOUT ((SPAWN_DEPTH * SPAWN_WIDTH / 1000) * 50 + 100)

static atomic_int nested_counter;

void nested_spawner(void* arg) {
    int depth = (int)(uintptr_t)arg;
    atomic_fetch_add(&nested_counter, 1);
    
    if (depth < SPAWN_DEPTH) {
        for (int i = 0; i < SPAWN_WIDTH; i++) {
            thread(nested_spawner, 1, (void*)(uintptr_t)(depth + 1));
        }
    }
}

void test_stress_nested_task_spawning(void) {
    atomic_store(&nested_counter, 0);
    
    int expected = 1 + SPAWN_WIDTH;
    
    printf("Spawning nested tasks (2 levels)...\n");
    thread(nested_spawner, 1, (void*)(uintptr_t)SPAWN_DEPTH-1);  /* Start at depth - 1 */
    
    /* Wait for completion with longer timeout */
    struct timespec ts = {0, 50000000}; /* 50ms */
    int timeout = 0;
    while (atomic_load(&nested_counter) < expected && timeout < 200) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int completed = atomic_load(&nested_counter);
    printf("Completed: %d/%d nested tasks\n", completed, expected);
    
    TEST_ASSERT_EQUAL_INT(expected, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_nested_task_spawning);
    
    /* Scheduler shuts down automatically when all tasks complete */
    wait_for_schedulers();
    
    return UNITY_END();
}
