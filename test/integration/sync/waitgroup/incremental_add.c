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
#include "sync.h"

extern __thread arena_t* __arena__; 
extern arena_t* __global__arena__;

void setUp(void) {

}

void tearDown(void) {
    /* @todo: gracefull termination */
}

#define N 20
static atomic_int completed;
static atomic_int wait_returned;

__public__waitgroup_t* mocked_waitgroup_create() {
    __public__waitgroup_t* wg = (__public__waitgroup_t*)allocate(__global__arena__, sizeof(__public__waitgroup_t));
    safe_q_init(&wg->waiters, SCHEDULER_LOCAL_QUEUE_SIZE);
    pthread_mutex_init(&wg->lock, NULL);
    atomic_store_explicit(&wg->count, 0, memory_order_relaxed);
    return wg;
}

typedef struct {
    __public__waitgroup_t* wg;
    int task_id;
} task_arg_t;

void incremental_worker(task_arg_t* arg) {
    // Simulate work
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 500000; /* 0.5ms */
    nanosleep(&ts, NULL);
    
    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    __public__sync_waitgroup_done(arg->wg);
}

void waiter_task(__public__waitgroup_t* wg) {
    __public__sync_waitgroup_wait(wg);
    
    atomic_store_explicit(&wait_returned, 1, memory_order_release);
}

void test_waitgroup_incremental_add(void) {
    atomic_store(&completed, 0);
    atomic_store(&wait_returned, 0);

    __public__waitgroup_t* wg = mocked_waitgroup_create();
    
    // Add tasks incrementally
    for (int i = 0; i < N; i++) {
        __public__sync_waitgroup_add(wg, 1);
        
        task_arg_t* arg = (task_arg_t*)allocate(__global__arena__, sizeof(task_arg_t));
        arg->wg = wg;
        arg->task_id = i;
        thread(incremental_worker, 1, arg);
    }
    
    // Start waiter after all tasks are scheduled
    thread(waiter_task, 1, wg);

    /* bounded wait */
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000; /* 1ms */

    int spins = 0;
    while (atomic_load_explicit(&wait_returned, memory_order_acquire) == 0) {
        nanosleep(&ts, NULL);
        spins++;
        if (spins > 5000) {
            TEST_FAIL_MESSAGE("WaitGroup.Wait() did not return");
        }
    }

    TEST_ASSERT_EQUAL_INT(N, atomic_load(&completed));
    TEST_ASSERT_EQUAL_INT(1, atomic_load(&wait_returned));
}

int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();

    srand(time(NULL));
    gc_init();

    init_io();
    init_scheduler();

    gc_start();

    RUN_TEST(test_waitgroup_incremental_add);

    return UNITY_END();
}
