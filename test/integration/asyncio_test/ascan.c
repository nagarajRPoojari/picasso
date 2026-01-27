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
arena_t* __test__global__arena__;

void setUp(void) {
    __test__global__arena__ = arena_create();
}
void tearDown(void) {
    /* @todo: graceful termination */
}

__public__array_t* mock_alloc_array(int count, int elem_size, int rank) {
    size_t data_size = (size_t)count * elem_size;
    size_t shape_size = (size_t)rank * sizeof(int64_t);
    size_t total_size = sizeof(__public__array_t) + data_size + shape_size;

    __public__array_t* arr = (__public__array_t*)allocate(__test__global__arena__, total_size);

    
    arr->data = (char*)(arr + 1); 
    
    if (rank > 0) {
        arr->shape = (int64_t*)(arr->data + data_size);
    } else {
        arr->shape = NULL;
    }
    
    arr->length = count;
    arr->rank = rank;
    
    return arr;
}

static void redirect_stdin_pipe(const char *input, int *saved_stdin) {
    int pipefd[2];
    TEST_ASSERT_EQUAL(0, pipe(pipefd));

    *saved_stdin = dup(STDIN_FILENO);
    TEST_ASSERT_TRUE(*saved_stdin >= 0);

    ssize_t len = strlen(input);
    TEST_ASSERT_EQUAL(len, write(pipefd[1], input, len));
    close(pipefd[1]);

    dup2(pipefd[0], STDIN_FILENO);
    close(pipefd[0]);
}

static void restore_stdin(int saved_stdin) {
    dup2(saved_stdin, STDIN_FILENO);
    close(saved_stdin);
}

static atomic_int completed;

static void submit_task(void(*fn)(), int count, int timeout_sec) {
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

static void __public__ascan_thread_func(void* arg) {
    (void)arg;

    int saved_stdin;
    __public__array_t *buf;

    redirect_stdin_pipe("dummy input from user\n", &saved_stdin);
    
    self_yield();
    buf = __public__syncio_scan(11);
    restore_stdin(saved_stdin);

    TEST_ASSERT_NOT_NULL(buf);
    TEST_ASSERT_EQUAL_STRING("dummy input", buf->data);

    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
}

void test__public__ascan(void) {
    submit_task(__public__ascan_thread_func, 1, 5);
}

int main(void) {
    srand(time(NULL));

    
    __global__arena__ = gc_create_global_arena();
    gc_init();

    init_io();
    init_scheduler();

    gc_start();

    UNITY_BEGIN();

    RUN_TEST(test__public__ascan);

    return UNITY_END();
}
