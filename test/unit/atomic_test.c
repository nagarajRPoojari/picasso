#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include <stdatomic.h>
#include <pthread.h>
#include "atomic.h"
#include "alloc.h"
#include "gc.h"

#define NUM_THREADS 4
#define ITERATIONS 100000

/* for local usage : use here */
arena_t* __test__global__arena__;

/* per thread __arena__ : don't use here */
extern __thread arena_t* __arena__; 

/* global arena for internal runtime usage : don't use here */
extern arena_t* __global__arena__;
/* Setup / teardown for Unity */
void setUp(void) {
    __test__global__arena__ = (arena_t*)1; // dummy non-null
}

void tearDown(void) {
    __test__global__arena__ = NULL;
}

#define NUM_THREADS 4
#define ITERATIONS 100000

typedef struct {
    void* ptr;
    int id;
} thread_arg_t;

/* Generic thread writer for integer types */
void* thread_add_int(void* arg) {
    thread_arg_t* t = (thread_arg_t*)arg;
    _Atomic int* p = (_Atomic int*)t->ptr;
    for (int i = 0; i < ITERATIONS; i++) {
        __public__atomic_add_int(p, 1);
    }
    return NULL;
}

void* thread_sub_int(void* arg) {
    thread_arg_t* t = (thread_arg_t*)arg;
    _Atomic int* p = (_Atomic int*)t->ptr;
    for (int i = 0; i < ITERATIONS; i++) {
        __public__atomic_sub_int(p, 1);
    }
    return NULL;
}

void test_atomic_bool_basic(void) {
    _Atomic _Bool b = 0;
    __public__atomic_store_bool(&b, 1);
    TEST_ASSERT_TRUE(__public__atomic_load_bool(&b));
    __public__atomic_store_bool(&b, 0);
    TEST_ASSERT_FALSE(__public__atomic_load_bool(&b));
}

void test_atomic_char_basic(void) {
    _Atomic char c = 0;
    __public__atomic_store_char(&c, 42);
    TEST_ASSERT_EQUAL_CHAR(42, __public__atomic_load_char(&c));
    TEST_ASSERT_EQUAL_CHAR(42, __public__atomic_add_char(&c, 1)); // returns old value
    TEST_ASSERT_EQUAL_CHAR(43, __public__atomic_load_char(&c));
    TEST_ASSERT_EQUAL_CHAR(43, __public__atomic_sub_char(&c, 2)); // returns old value
    TEST_ASSERT_EQUAL_CHAR(41, __public__atomic_load_char(&c));
}

void test_atomic_short_basic(void) {
    _Atomic short s = 0;
    __public__atomic_store_short(&s, 100);
    TEST_ASSERT_EQUAL_INT16(100, __public__atomic_load_short(&s));
    __public__atomic_add_short(&s, 10);
    TEST_ASSERT_EQUAL_INT16(110, __public__atomic_load_short(&s));
    __public__atomic_sub_short(&s, 20);
    TEST_ASSERT_EQUAL_INT16(90, __public__atomic_load_short(&s));
}

void test_atomic_int_basic(void) {
    _Atomic int i = 0;
    __public__atomic_store_int(&i, 1000);
    TEST_ASSERT_EQUAL_INT(1000, __public__atomic_load_int(&i));
    __public__atomic_add_int(&i, 50);
    TEST_ASSERT_EQUAL_INT(1050, __public__atomic_load_int(&i));
    __public__atomic_sub_int(&i, 25);
    TEST_ASSERT_EQUAL_INT(1025, __public__atomic_load_int(&i));
}

void test_atomic_long_basic(void) {
    _Atomic long l = 0;
    __public__atomic_store_long(&l, 123456);
    TEST_ASSERT_EQUAL_INT64(123456, __public__atomic_load_long(&l));
    __public__atomic_add_long(&l, 44);
    TEST_ASSERT_EQUAL_INT64(123500, __public__atomic_load_long(&l));
}

void test_atomic_llong_basic(void) {
    _Atomic long long ll = 0;
    __public__atomic_store_llong(&ll, 123456789);
    TEST_ASSERT_EQUAL_INT64(123456789, __public__atomic_load_llong(&ll));
    __public__atomic_add_llong(&ll, 11);
    TEST_ASSERT_EQUAL_INT64(123456800, __public__atomic_load_llong(&ll));
}

void test_atomic_float_basic(void) {
    _Atomic float f = 0.0f;
    __public__atomic_store_float(&f, 3.14f);
    TEST_ASSERT_FLOAT_WITHIN(0.0001f, 3.14f, __public__atomic_load_float(&f));
}

void test_atomic_double_basic(void) {
    _Atomic double d = 0.0;
    __public__atomic_store_double(&d, 2.718);
    #define DOUBLE_EPS 0.000001
    TEST_ASSERT_FLOAT_WITHIN(DOUBLE_EPS, 2.718, (float)__public__atomic_load_double(&d));
}

void test_atomic_ptr_basic(void) {
    int x = 42;
    int y = 100;
    _Atomic void* p = NULL;
    __public__atomic_store_ptr(&p, &x);
    TEST_ASSERT_EQUAL_PTR(&x, __public__atomic_load_ptr(&p));
    __public__atomic_store_ptr(&p, &y);
    TEST_ASSERT_EQUAL_PTR(&y, __public__atomic_load_ptr(&p));
}

void test_atomic_int_concurrent(void) {
    _Atomic int counter = 0;
    pthread_t writers[NUM_THREADS];
    pthread_t subtracters[NUM_THREADS];
    thread_arg_t args[NUM_THREADS];

    for (int i = 0; i < NUM_THREADS; i++) {
        args[i].ptr = &counter;
        args[i].id = i;
        pthread_create(&writers[i], NULL, thread_add_int, &args[i]);
        pthread_create(&subtracters[i], NULL, thread_sub_int, &args[i]);
    }

    for (int i = 0; i < NUM_THREADS; i++) {
        pthread_join(writers[i], NULL);
        pthread_join(subtracters[i], NULL);
    }

    TEST_ASSERT_EQUAL_INT(0, __public__atomic_load_int(&counter));
}

int main(void) {
    UNITY_BEGIN();

    RUN_TEST(test_atomic_bool_basic);
    RUN_TEST(test_atomic_char_basic);
    RUN_TEST(test_atomic_short_basic);
    RUN_TEST(test_atomic_int_basic);
    RUN_TEST(test_atomic_long_basic);
    RUN_TEST(test_atomic_llong_basic);
    RUN_TEST(test_atomic_float_basic);
    RUN_TEST(test_atomic_double_basic);
    RUN_TEST(test_atomic_ptr_basic);

    RUN_TEST(test_atomic_int_concurrent);

    return UNITY_END();
}
