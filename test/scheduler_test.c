#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include <pthread.h>
#include "io.h"
#include "alloc.h"
#include "initutils.h"


extern kernel_thread_t **kernel_thread_map;
extern struct io_uring **io_ring_map;

extern pthread_t sched_threads[SCHEDULER_THREAD_POOL_SIZE];


__thread arena_t* __global__arena__;

void setUp(void) {
    __global__arena__ = arena_create();

    init_io();
    init_scheduler();
}

void tearDown(void) {
    __global__arena__ = NULL;
}

void* test_func(void*) {
    printf("hello world\n");

    return NULL;
}

void test_scheduler(void) {
    thread(test_func, NULL);

    for (int i = 0; i < SCHEDULER_THREAD_POOL_SIZE; i++) {
        pthread_join(sched_threads[i], NULL);
    }
}


int main(void) {
    UNITY_BEGIN();

    RUN_TEST(test_scheduler);

    return UNITY_END();
}