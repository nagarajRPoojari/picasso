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

void worker_task(task_arg_t* arg) {
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000 * (rand() % 10); /* 0-10ms random delay */
    nanosleep(&ts, NULL);
    
    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    
    // signal completion
    __public__sync_waitgroup_done(arg->wg);
}

void waiter_task(__public__waitgroup_t* wg) {
    __public__sync_waitgroup_wait(wg);
}

void test_waitgroup_multiple_waiters(void) {
    atomic_store(&completed, 0);

    __public__waitgroup_t* wg = mocked_waitgroup_create();
    
    // Add N to the wait group
    __public__sync_waitgroup_add(wg, N);
    
    // Start multiple waiter tasks
    thread(waiter_task, 1, wg);
    thread(waiter_task, 1, wg);
    thread(waiter_task, 1, wg);
    
    // Start N worker tasks
    task_arg_t* args[N];
    for (int i = 0; i < N; i++) {
        args[i] = (task_arg_t*)allocate(__global__arena__, sizeof(task_arg_t));
        args[i]->wg = wg;
        args[i]->task_id = i;
        thread(worker_task, 1, args[i]);
    }

    /* bounded wait — no deadlock */
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000; /* 1ms */

    int spins = 0;
    while (atomic_load_explicit(&completed, memory_order_acquire) < N) {
        nanosleep(&ts, NULL);
        spins++;
        if (spins > 5000) { /* ~5s */
            TEST_FAIL_MESSAGE("Workers did not complete");
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

    RUN_TEST(test_waitgroup_multiple_waiters);

    return UNITY_END();
}
