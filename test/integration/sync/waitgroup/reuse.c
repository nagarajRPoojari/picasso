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

static atomic_int completed;
static atomic_int batch_done;

__public__waitgroup_t* mocked_waitgroup_create() {
    __public__waitgroup_t* wg = (__public__waitgroup_t*)allocate(__global__arena__, sizeof(__public__waitgroup_t));
    safe_q_init(&wg->waiters, SCHEDULER_LOCAL_QUEUE_SIZE);
    pthread_mutex_init(&wg->lock, NULL);
    atomic_store_explicit(&wg->count, 0, memory_order_relaxed);
    return wg;
}

void batch_worker(__public__waitgroup_t* wg) {
    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    __public__sync_waitgroup_done(wg);
}

void reuse_coordinator(__public__waitgroup_t* wg) {
    // First batch
    __public__sync_waitgroup_add(wg, 5);
    
    for (int i = 0; i < 5; i++) {
        thread(batch_worker, 1, wg);
    }
    
    __public__sync_waitgroup_wait(wg);
    atomic_store_explicit(&batch_done, 1, memory_order_release);
    

    // Reuse for second batch
    atomic_store(&completed, 0);
    __public__sync_waitgroup_add(wg, 5);
    
    for (int i = 0; i < 5; i++) {
        thread(batch_worker, 1, wg);
    }
    
    __public__sync_waitgroup_wait(wg);
    atomic_store_explicit(&batch_done, 2, memory_order_release);
}

void test_waitgroup_reuse(void) {
    atomic_store(&completed, 0);
    atomic_store(&batch_done, 0);

    __public__waitgroup_t* wg = mocked_waitgroup_create();
    
    thread(reuse_coordinator, 1, wg);
    
    /* bounded wait */
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000; /* 1ms */

    int spins = 0;
    while (atomic_load_explicit(&batch_done, memory_order_acquire) < 2) {
        nanosleep(&ts, NULL);
        spins++;
        if (spins > 5000) {
            int done = atomic_load(&batch_done);
            if (done == 0) {
                TEST_FAIL_MESSAGE("First batch did not complete");
            } else {
                TEST_FAIL_MESSAGE("Second batch did not complete");
            }
        }
    }
    
    TEST_ASSERT_EQUAL_INT(5, atomic_load(&completed));
    TEST_ASSERT_EQUAL_INT(2, atomic_load(&batch_done));
}

int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();

    srand(time(NULL));
    gc_init();

    init_io();
    init_scheduler();

    gc_start();

    RUN_TEST(test_waitgroup_reuse);

    return UNITY_END();
}
