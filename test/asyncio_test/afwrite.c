#include "test/unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include <pthread.h>
#include <stdatomic.h>
#include <time.h>
#include "diskio.h"
#include "alloc.h"
#include "gc.h"
#include "ggc.h"
#include "initutils.h"

extern arena_t* __global__arena__;

void setUp(void) {}
void tearDown(void) {
    /* @todo: graceful termination */
}
static atomic_int completed;

static void submit_task(void*(*fn)(void*), int count, int timeout_sec) {
    atomic_store(&completed, 0);

    for (int i = 0; i < count; i++) {
        thread(fn, 1, NULL);
    }

    struct timespec ts = {0, 1000000}; // 1ms
    int spins = 0;
    int max_spins = timeout_sec * 1000; // approximate timeout in ms

    while (atomic_load_explicit(&completed, memory_order_acquire) < count) {
        nanosleep(&ts, NULL);
        spins++;
        if (spins > max_spins) {
            TEST_FAIL_MESSAGE("scheduler did not complete tasks");
        }
    }

    TEST_ASSERT_EQUAL_INT(count, atomic_load(&completed));
}

static void* __public__afwrite_thread_func(void* arg) {
    (void)arg;

    FILE* file = fopen("test/data/test__public__swrite.txt", "w");
    TEST_ASSERT_NOT_NULL(file);

    char buf[10];
    for(int i=0; i<10; i++) buf[i] = 'a' + i%26;
    ssize_t r = __public__sfwrite(file, buf, 10, 0);
    fclose(file);

    TEST_ASSERT_EQUAL(10, r);
    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    return NULL;
}

void test__public__afwrite(void) {
    submit_task(__public__afwrite_thread_func, 1, 5);
}

int main(void) {
    srand(time(NULL));

    __global__arena__ = gc_create_global_arena();

    init_io();
    init_scheduler();
    gc_init();

    UNITY_BEGIN();

    RUN_TEST(test__public__afwrite);

    return UNITY_END();
}
