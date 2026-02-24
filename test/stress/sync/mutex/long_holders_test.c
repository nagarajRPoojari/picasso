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


/* Stress levels - configurable via compile-time macros
 * Define one of: STRESS_LEVEL_1, STRESS_LEVEL_2, STRESS_LEVEL_3, STRESS_LEVEL_4
 * Or use default (STRESS_LEVEL_3)
 */
#if defined(STRESS_LEVEL_1)
    #define R 2
    #define W 1
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define R 20
    #define W 10
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define R 100
    #define W 50
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define R 1000
    #define W 100
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define MAX_TIMEOUT ((R + W / 1000) * 50 + 100)


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

__public__mutex_t* mocked_mutex_create() {
    __public__mutex_t* mux = (__public__mutex_t*)allocate(__global__arena__, sizeof(__public__mutex_t));
    pthread_mutex_init(&mux->lock, NULL);
    safe_q_init(&mux->waiters, SCHEDULER_LOCAL_QUEUE_SIZE);
    atomic_store_explicit(&mux->state, 0, memory_order_relaxed);
    return mux;
}

void short_reader(__public__mutex_t* mu) {
    __public__sync_mutex_lock(mu);
    atomic_fetch_add(&readers_in, 1);
    
    if(atomic_load(&writers_in) > 0){
        atomic_store(&violations, 1);
    }

    atomic_fetch_sub(&readers_in, 1);
    
    atomic_fetch_add(&completed, 1);
    __public__sync_mutex_unlock(mu);
}

void short_writer(__public__mutex_t* mu) {
    __public__sync_mutex_lock(mu);
    atomic_fetch_add(&writers_in, 1);

    if(atomic_load(&writers_in) > 1){
        atomic_store(&violations, 1);
    }
    if(atomic_load(&readers_in) > 0){
        atomic_store(&violations, 1);
    }

    atomic_fetch_sub(&writers_in, 1);

    atomic_fetch_add(&completed, 1);
    __public__sync_mutex_unlock(mu);

}

void long_writer(__public__mutex_t* mu) {
    __public__sync_mutex_lock(mu);
    atomic_fetch_add(&writers_in, 1);

    if(atomic_load(&writers_in) > 1){
        atomic_store(&violations, 1);
    }
    if(atomic_load(&readers_in) > 0){
        atomic_store(&violations, 1);
    }

    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000; /* 1ms */
    nanosleep(&ts, NULL);


    atomic_fetch_sub(&writers_in, 1);

    atomic_fetch_add(&completed, 1);
    __public__sync_mutex_unlock(mu);

}

void test_mutex_concurrent_readers_writers(void) {
    atomic_store(&completed, 0);

    __public__mutex_t* mu = mocked_mutex_create();

    int i=0;
    int j=0;
    for(; i < R || j < W ;) {
        if(i < R) {
            thread(short_reader, 1, mu);
            i++;
        }
        if(j < W) {
            if(random()%2)
                thread(short_writer, 1, mu);
            else thread(long_writer, 1, mu);
            j++;
        }
    }

    /* bounded wait — no deadlock */
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000; /* 1ms */

    int spins = 0;
    while (atomic_load_explicit(&completed, memory_order_acquire) < R+W ) {
        nanosleep(&ts, NULL);
        spins++;
        if (spins > MAX_TIMEOUT) { /* ~5s */
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

    RUN_TEST(test_mutex_concurrent_readers_writers);
    return UNITY_END();
}