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
#include "queue.h"
#include "diskio.h"
#include "alloc.h"
#include "gc.h"
#include "initutils.h"
#include "sync.h"

extern arena_t* __global__arena__;

#define R 10
#define W 10

static atomic_int readers_in;
static atomic_int writers_in;
static atomic_int violations;
static atomic_int completed;

void setUp(void) {
    atomic_init(&readers_in, 0);
    atomic_init(&writers_in, 0);
    atomic_init(&violations, 0);
    atomic_init(&completed, 0);
}

void tearDown(void) {
    /* @todo: gracefull termination */
}

__public__rwmutex_t* mocked_rwmutex_create() {
    assert(__global__arena__ != NULL);
    __public__rwmutex_t* mux = (__public__rwmutex_t*)allocate(__global__arena__, sizeof(__public__rwmutex_t));
    safe_q_init(&mux->readers, SCHEDULER_LOCAL_QUEUE_SIZE);
    safe_q_init(&mux->writers, SCHEDULER_LOCAL_QUEUE_SIZE);

    pthread_mutex_init(&mux->lock, NULL);
    atomic_store_explicit(&mux->state, 0, memory_order_relaxed);

    return mux;
}

void* reader(__public__rwmutex_t* mu) {
    __public__rwmutex_rlock(mu);
    atomic_fetch_add(&readers_in, 1);
    
    if(atomic_load(&writers_in) > 0){
        atomic_store(&violations, 1);
    }

    atomic_fetch_sub(&readers_in, 1);
    
    atomic_fetch_add(&completed, 1);
    __public__rwmutex_runlock(mu);
    return NULL;
}

void* writer(__public__rwmutex_t* mu) {
    __public__rwmutex_rwlock(mu);
    atomic_fetch_add(&writers_in, 1);

    if(atomic_load(&writers_in) > 1){
        atomic_store(&violations, 1);
    }
    if(atomic_load(&readers_in) > 0){
        atomic_store(&violations, 1);
    }

    atomic_fetch_sub(&writers_in, 1);

    atomic_fetch_add(&completed, 1);
    __public__rwmutex_rwunlock(mu);
}

void test_rwmutex_concurrent_readers_writers(void) {
    atomic_store(&completed, 0);

    __public__rwmutex_t* mu = mocked_rwmutex_create();

    for(int i=0; i<R; i++) {
        thread(reader, 1, mu);
    }
    for(int i=0; i<W; i++) {
        thread(writer, 1, mu);
    }

    /* bounded wait — no deadlock */
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000; /* 1ms */

    int spins = 0;
    while (atomic_load_explicit(&completed, memory_order_acquire) < R+W) {
        nanosleep(&ts, NULL);
        spins++;
        if (spins > 5000) { /* ~5s */
            TEST_FAIL_MESSAGE("scheduler did not complete tasks");
        }
    }

    TEST_ASSERT_EQUAL_INT(R+W, atomic_load(&completed));
    TEST_ASSERT_EQUAL_INT(0, atomic_load(&violations));
}

int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();

    srand(time(NULL));
    gc_init();

    init_io();
    init_scheduler();

    gc_start();

    RUN_TEST(test_rwmutex_concurrent_readers_writers);
    return UNITY_END();
}