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

static atomic_int wait_returned;

__public__waitgroup_t* mocked_waitgroup_create() {
    __public__waitgroup_t* wg = (__public__waitgroup_t*)allocate(__global__arena__, sizeof(__public__waitgroup_t));
    safe_q_init(&wg->waiters, SCHEDULER_LOCAL_QUEUE_SIZE);
    pthread_mutex_init(&wg->lock, NULL);
    atomic_store_explicit(&wg->count, 0, memory_order_relaxed);
    return wg;
}

void test_waitgroup_immediate_return(void) {
    atomic_store(&wait_returned, 0);

    __public__waitgroup_t* wg = mocked_waitgroup_create();
    
    // Don't add anything - counter is already 0
    // Wait should return immediately
    __public__sync_waitgroup_wait(wg);
    
    atomic_store(&wait_returned, 1);
    
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

    RUN_TEST(test_waitgroup_immediate_return);

    return UNITY_END();
}
