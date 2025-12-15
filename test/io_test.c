#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include "io.h"
#include "alloc.h"

extern __thread arena_t* __arena__;

void setUp(void) {
    __arena__ = arena_create();
}

void tearDown(void) {
    __arena__ = NULL;
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
    char *buf;

    /* small reads */
    redirect_stdin_pipe("dummy input from user\n", &saved_stdin);
    buf = __public__sscan(11);
    restore_stdin(saved_stdin);

    TEST_ASSERT_NOT_NULL(buf);
    TEST_ASSERT_EQUAL_STRING("dummy input", buf);

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
    TEST_ASSERT_EQUAL_STRING(input, buf);

    /* input < requested length (read-some semantics) */
    redirect_stdin_pipe("input is only", &saved_stdin);
    buf = __public__sscan(20);
    restore_stdin(saved_stdin);

    TEST_ASSERT_NOT_NULL(buf);
    TEST_ASSERT_EQUAL_STRING("input is only", buf);
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

int main(void) {
    UNITY_BEGIN();

    RUN_TEST(test__public__sscan);

    return UNITY_END();
}