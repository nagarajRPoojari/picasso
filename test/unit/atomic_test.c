#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include <stdatomic.h>
#include <pthread.h>
#include <stdatomic.h>
#include <stdint.h>
#include <stdbool.h>

#include "atomics.h"
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
    _Atomic int64_t* p = (_Atomic int64_t*)t->ptr;
    for (int i = 0; i < ITERATIONS; i++) {
        __public__atomics_add_int(p, 1);
    }
    return NULL;
}

void* thread_sub_int(void* arg) {
    thread_arg_t* t = (thread_arg_t*)arg;
    _Atomic int64_t* p = (_Atomic int64_t*)t->ptr;
    for (int i = 0; i < ITERATIONS; i++) {
        __public__atomics_sub_int(p, 1);
    }
    return NULL;
}

void test_atomic_bool_basic(void) {
    _Atomic _Bool b = 0;
    __public__atomics_store_boolean(&b, 1);
    TEST_ASSERT_TRUE(__public__atomics_load_boolean(&b));
    __public__atomics_store_boolean(&b, 0);
    TEST_ASSERT_FALSE(__public__atomics_load_boolean(&b));
}

void test_atomic_int8_basic(void) {
    _Atomic int8_t c = 0;
    __public__atomics_store_int8(&c, 42);
    TEST_ASSERT_EQUAL_CHAR(42, __public__atomics_load_int8(&c));
    TEST_ASSERT_EQUAL_CHAR(42, __public__atomics_add_int8(&c, 1)); // returns old value
    TEST_ASSERT_EQUAL_CHAR(43, __public__atomics_load_int8(&c));
    TEST_ASSERT_EQUAL_CHAR(43, __public__atomics_sub_int8(&c, 2)); // returns old value
    TEST_ASSERT_EQUAL_CHAR(41, __public__atomics_load_int8(&c));
}

void test_atomic_int16_t_basic(void) {
    _Atomic int16_t s = 0;
    __public__atomics_store_int16(&s, 100);
    TEST_ASSERT_EQUAL_INT16(100, __public__atomics_load_int16(&s));
    __public__atomics_add_int16(&s, 10);
    TEST_ASSERT_EQUAL_INT16(110, __public__atomics_load_int16(&s));
    __public__atomics_sub_int16(&s, 20);
    TEST_ASSERT_EQUAL_INT16(90, __public__atomics_load_int16(&s));
}

void test_atomic_int32_basic(void) {
    _Atomic int32_t i = 0;
    __public__atomics_store_int32(&i, 1000);
    TEST_ASSERT_EQUAL_INT(1000, __public__atomics_load_int32(&i));
    __public__atomics_add_int32(&i, 50);
    TEST_ASSERT_EQUAL_INT(1050, __public__atomics_load_int32(&i));
    __public__atomics_sub_int32(&i, 25);
    TEST_ASSERT_EQUAL_INT(1025, __public__atomics_load_int32(&i));
}

void test_atomic_int64_basic(void) {
    _Atomic int64_t i = 0;
    __public__atomics_store_int64(&i, 1000);
    TEST_ASSERT_EQUAL_INT(1000, __public__atomics_load_int64(&i));
    __public__atomics_add_int(&i, 50);
    TEST_ASSERT_EQUAL_INT(1050, __public__atomics_load_int64(&i));
    __public__atomics_sub_int(&i, 25);
    TEST_ASSERT_EQUAL_INT(1025, __public__atomics_load_int64(&i));
}

void test_atomic_float_basic(void) {
    _Atomic float f = 0.0f;
    __public__atomics_store_float32(&f, 3.14f);
    TEST_ASSERT_FLOAT_WITHIN(0.0001f, 3.14f, __public__atomics_load_float32(&f));
}

void test_atomic_double_basic(void) {
    _Atomic double d = 0.0;
    __public__atomics_store_double(&d, 2.718);
    #define DOUBLE_EPS 0.000001
    TEST_ASSERT_FLOAT_WITHIN(DOUBLE_EPS, 2.718, (float)__public__atomics_load_double(&d));
}

void test_atomic_ptr_basic(void) {
    int x = 42;
    int y = 100;
    _Atomic uintptr_t p = 0;
    __public__atomics_store_ptr(&p, &x);
    TEST_ASSERT_EQUAL_PTR(&x, __public__atomics_load_ptr(&p));
    __public__atomics_store_ptr(&p, &y);
    TEST_ASSERT_EQUAL_PTR(&y, __public__atomics_load_ptr(&p));
}

void test_atomic_int_concurrent(void) {
    _Atomic int32_t counter = 0;
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

    TEST_ASSERT_EQUAL_INT(0, __public__atomics_load_int32(&counter));
}

int main(void) {
    UNITY_BEGIN();

    RUN_TEST(test_atomic_bool_basic);
    RUN_TEST(test_atomic_int8_basic);
    RUN_TEST(test_atomic_int64_basic);
    RUN_TEST(test_atomic_float_basic);
    RUN_TEST(test_atomic_double_basic);
    RUN_TEST(test_atomic_ptr_basic);

    RUN_TEST(test_atomic_int_concurrent);

    return UNITY_END();
}
