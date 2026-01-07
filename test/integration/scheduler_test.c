#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include <pthread.h>
#include <stdatomic.h>
#include <time.h>
#include "diskio.h"
#include "alloc.h"
#include "gc.h"
#include "initutils.h"

extern __thread arena_t* __arena__; 
extern arena_t* __global__arena__;

void setUp(void) {

}

void tearDown(void) {
    /* @todo: gracefull termination */
}

#define N 20
static atomic_int completed;
static atomic_int tasks_n;

void* test_func(void* arg);

void* test_func(void* arg) {
    (void)arg;

    if(atomic_load(&tasks_n) > 0) {
        atomic_fetch_sub(&tasks_n, 1);
        thread(test_func, 1, NULL);
    }

    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    return NULL;
}

void test_scheduler_executes_tasks(void) {
    atomic_store(&completed, 0);
    atomic_store(&tasks_n, N-1);

    thread(test_func, 1, NULL);

    /* bounded wait — no deadlock */
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000; /* 1ms */

    int spins = 0;
    while (atomic_load_explicit(&completed, memory_order_acquire) < N) {
        nanosleep(&ts, NULL);
        spins++;
        if (spins > 5000) { /* ~5s */
            TEST_FAIL_MESSAGE("scheduler did not complete tasks");
        }
    }

    TEST_ASSERT_EQUAL_INT(N, atomic_load(&completed));
}

int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();

    srand(time(NULL));
    gc_init();

    init_io();
    init_scheduler();

    gc_start();

    RUN_TEST(test_scheduler_executes_tasks);

    return UNITY_END();
}