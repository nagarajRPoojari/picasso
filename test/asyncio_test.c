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
#include "io.h"
#include "alloc.h"
#include "gc.h"
#include "ggc.h"
#include "initutils.h"

__thread arena_t* __global__arena__;

void setUp(void) {
    init_scheduler();
        init_io();
    gc_init();
}

void tearDown(void) {
    // __global__arena__ = NULL;
    /* @todo: gracefull termination */
}

static void redirect_stdin_pipe(const char *input, int *saved_stdin) {
    int pipefd[2];
    TEST_ASSERT_EQUAL(0, pipe(pipefd));

    /* save real stdin */
    *saved_stdin = dup(STDIN_FILENO);
    TEST_ASSERT_TRUE(*saved_stdin >= 0);

    /* write input into pipe */
    ssize_t len = strlen(input);
    TEST_ASSERT_EQUAL(len, write(pipefd[1], input, len));
    close(pipefd[1]);

    /* redirect stdin -> pipe read end */
    dup2(pipefd[0], STDIN_FILENO);
    close(pipefd[0]);
}

static void restore_stdin(int saved_stdin) {
    dup2(saved_stdin, STDIN_FILENO);
    close(saved_stdin);
}

static int redirect_stdout_pipe(int *saved_stdout) {
    int pipefd[2];
    TEST_ASSERT_EQUAL(0, pipe(pipefd));

    /* save real stdout */
    *saved_stdout = dup(STDOUT_FILENO);
    TEST_ASSERT_TRUE(*saved_stdout >= 0);

    dup2(pipefd[1], STDOUT_FILENO);
    close(pipefd[1]);

    return pipefd[0];
}

static void restore_stout(int saved_stdout) {
    dup2(saved_stdout, STDOUT_FILENO);
    close(saved_stdout);
}

static atomic_int completed;

void submit_task(void*(*fn)(void*), int count, int timeout) {
    atomic_store(&completed, 0);

    for (int i = 0; i < count; i++) {
        thread(fn, NULL);
    }

    /* bounded wait — no deadlock */
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = 1000000; /* 1ms */

    int spins = 0;
    while (atomic_load_explicit(&completed, memory_order_acquire) < count) {
        nanosleep(&ts, NULL);
        spins++;
        if (spins > timeout*1000) { /* ~timeout sec */
            TEST_FAIL_MESSAGE("scheduler did not complete tasks");
        }
    }

    TEST_ASSERT_EQUAL_INT(count, atomic_load(&completed));
}


static void* __public__ascan_thread_func(void* arg) {
    (void)arg;

    int saved_stdin;
    char *buf;

    /* small reads */
    redirect_stdin_pipe("dummy input from user\n", &saved_stdin);
    
    self_yield();
    buf = __public__ascan(11);
    restore_stdin(saved_stdin);

    TEST_ASSERT_NOT_NULL(buf);

    printf("what I read: %s \n ", buf);

    /* @verify: can't expect full n byte read, need to think */
    TEST_ASSERT_EQUAL_STRING("dummy input", buf);

    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    return NULL;
}

static void* __public__aprtinf_thread_func(void* arg) {

    int saved_stdout;
    int readfd = redirect_stdout_pipe(&saved_stdout);

    /* call function */
    self_yield();
    ssize_t ret = __public__aprintf("hello %d %s", 42, "world");

    /* restore stdout */
    restore_stout(saved_stdout);

    /* read captured output */
    char buf[64] = {0};
    ssize_t r = read(readfd, buf, sizeof(buf) - 1);
    close(readfd);

    TEST_ASSERT_EQUAL(14, ret);
    TEST_ASSERT_EQUAL(14, r);
    TEST_ASSERT_EQUAL_STRING("hello 42 world", buf);

    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    return NULL;
}

static void* __public__afread_thread_func(void* arg) {
    FILE* file = fopen("/workspaces/x-language/test/data/test__public__sfread.txt", "r+");

    char* buf = __public__alloc(1024);
    ssize_t r = __public__afread((char*)file, buf, 1024, 0);
    
    TEST_ASSERT_EQUAL(57, r);
    TEST_ASSERT_EQUAL_STRING("Synchronously read n bytes from a file at a given offset.", buf);

    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    return NULL;
}

static void* __public__afwrite_thread_func(void* arg) {
    FILE* file = fopen("/workspaces/x-language/test/data/test__public__swrite.txt", "w");
    
    char buf[10];
    for(int i=0; i<10; i++) buf[i] = 'a' + i%26;
    ssize_t r = __public__sfwrite((char*)file, buf, 10, 0);
    
    TEST_ASSERT_EQUAL(10, r);
    atomic_fetch_add_explicit(&completed, 1, memory_order_release);
    return NULL;
}

void test__public__ascan(void) {
    submit_task(__public__ascan_thread_func, 1, 5);        
}

void test__public__aprintf(void) {
    submit_task(__public__aprtinf_thread_func, 1, 5);        
}

void test__public__afread(void) {
    submit_task(__public__afread_thread_func, 1, 5);
}

void test__public__afwrite(void) {
    submit_task(__public__afwrite_thread_func, 1, 5);
}

int main(void) {
    __global__arena__ = arena_create();
    srand(time(NULL));


    UNITY_BEGIN();

    RUN_TEST(test__public__ascan);
    // RUN_TEST(test__public__aprintf);
    // RUN_TEST(test__public__afread);
    // RUN_TEST(test__public__afwrite);

    return UNITY_END();
}