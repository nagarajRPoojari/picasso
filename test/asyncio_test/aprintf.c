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

static int redirect_stdout_pipe(int *saved_stdout, int *write_end) {
    int pipefd[2];
    TEST_ASSERT_EQUAL(0, pipe(pipefd));

    *saved_stdout = dup(STDOUT_FILENO);
    TEST_ASSERT_TRUE(*saved_stdout >= 0);

    dup2(pipefd[1], STDOUT_FILENO);
    *write_end = pipefd[1]; // save write end
    return pipefd[0];        // return read end
}

static void restore_stdout(int saved_stdout, int write_end) {
    dup2(saved_stdout, STDOUT_FILENO);
    close(saved_stdout);
    close(write_end);
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

static void* __public__aprintf_thread_func(void* arg) {
    // (void)arg;

    // int saved_stdout, write_end;
    // // int readfd = redirect_stdout_pipe(&saved_stdout, &write_end);

    // self_yield();
    // printf("[test] => before \n");
    ssize_t ret = __public__aprintf("hello %d %s", 42, "world");

    // // restore_stdout(saved_stdout, write_end);

    // char buf[64] = {0};
    // ssize_t r = read(STDOUT_FILENO, buf, sizeof(buf)-1);
    // close(readfd);

    // TEST_ASSERT_EQUAL(14, ret);
    // TEST_ASSERT_EQUAL(14, r);
    // TEST_ASSERT_EQUAL_STRING("hello 42 world", buf);

    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    return NULL;
}

void test__public__aprintf(void) {
    submit_task(__public__aprintf_thread_func, 1, 5);
}

int main(void) {
    srand(time(NULL));

    __global__arena__ = gc_create_global_arena();

    printf("=========== staring IO ========== \n");
    init_io();
    init_scheduler();
    gc_init();

    UNITY_BEGIN();

    RUN_TEST(test__public__aprintf);

    return UNITY_END();
}
