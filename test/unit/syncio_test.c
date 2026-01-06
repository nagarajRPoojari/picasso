#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include "diskio.h"
#include "alloc.h"
#include "gc.h"
#include "array.h"

/* for local usage : use here */
arena_t* __test__global__arena__;

/* per thread __arena__ : don't use here */
extern __thread arena_t* __arena__; 

/* global arena for internal runtime usage : don't use here */
extern arena_t* __global__arena__;

void setUp(void) {
    __test__global__arena__ = arena_create();
}

void tearDown(void) {
    __test__global__arena__ = NULL;
}


__public__array_t* mock_alloc_array(int count, int elem_size, int rank) {
    size_t data_size = (size_t)count * elem_size;
    size_t shape_size = (size_t)rank * sizeof(int64_t);
    size_t total_size = sizeof(__public__array_t) + data_size + shape_size;

    __public__array_t* arr = (__public__array_t*)allocate(__test__global__arena__, total_size);

    
    arr->data = (int8_t*)(arr + 1); 
    
    if (rank > 0) {
        arr->shape = (int64_t*)(arr->data + data_size);
    } else {
        arr->shape = NULL;
    }
    
    arr->length = count;
    arr->rank = rank;
    
    return arr;
}

/* blocking io test */

/* utility func to simulate stdin input */
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


void test__public__sscan(void) {
    int saved_stdin;
    __public__array_t *buf;

    /* small reads */
    redirect_stdin_pipe("dummy input from user\n", &saved_stdin);
    buf = __public__sscan(11);
    restore_stdin(saved_stdin);

    TEST_ASSERT_NOT_NULL(buf);
    TEST_ASSERT_EQUAL_STRING("dummy input", buf->data);

    /* large reads */
    int n = 1000;
    char input[n];
    for (int i = 0; i < n - 1; i++) {
        input[i] = 'a' + (i % 26);
    }
    input[n - 1] = '\0';

    redirect_stdin_pipe(input, &saved_stdin);
    buf = __public__sscan(n - 1);
    restore_stdin(saved_stdin);

    TEST_ASSERT_NOT_NULL(buf);
    TEST_ASSERT_EQUAL_STRING(input, buf->data);

    /* input < requested length (read-some semantics) */
    redirect_stdin_pipe("input is only", &saved_stdin);
    buf = __public__sscan(20);
    restore_stdin(saved_stdin);

    TEST_ASSERT_NOT_NULL(buf);
    TEST_ASSERT_EQUAL_STRING("input is only", buf->data);
}


void test__public__sprintf(void) {
    int saved_stdout;
    int readfd = redirect_stdout_pipe(&saved_stdout);

    /* call function */
    ssize_t ret = __public__sprintf("hello %d %s", 42, "world");

    /* restore stdout */
    restore_stout(saved_stdout);

    /* read captured output */
    char buf[64] = {0};
    ssize_t r = read(readfd, buf, sizeof(buf) - 1);
    close(readfd);

    TEST_ASSERT_EQUAL(14, ret);
    TEST_ASSERT_EQUAL(14, r);
    TEST_ASSERT_EQUAL_STRING("hello 42 world", buf);
}

void test__public__sfread(void) {
    FILE* file = fopen("test/data/test__public__sfread.txt", "r+");

    __public__array_t* buf = mock_alloc_array(1024, sizeof(size_t), 1);
    ssize_t r = __public__sfread((char*)file, buf, 1024, 0);
    
    TEST_ASSERT_EQUAL(57, r);
    TEST_ASSERT_EQUAL_STRING("Synchronously read n bytes from a file at a given offset.", buf->data);
}

void test__public__sfwrite(void) {
    FILE* file = fopen("test/data/test__public__swrite.txt", "w");
    
    __public__array_t* buf = mock_alloc_array(10, sizeof(size_t), 1);
    for(int i=0; i<10; i++) buf->data[i] = 'a' + i%26;
    ssize_t r = __public__sfwrite((char*)file, buf, 10, 0);
    
    TEST_ASSERT_EQUAL(10, r);
}

int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();
    __arena__ = gc_create_global_arena();

    RUN_TEST(test__public__sscan);
    RUN_TEST(test__public__sprintf);
    RUN_TEST(test__public__sfread);
    RUN_TEST(test__public__sfwrite);
    return UNITY_END();
}