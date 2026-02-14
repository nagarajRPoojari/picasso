#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include "array.h"
#include "alloc.h"
#include "gc.h"

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

/* Helper: check that data is zero-initialized */
static int is_zeroed(char* data, size_t size) {
    for (size_t i = 0; i < size; i++) {
        if (data[i] != 0) return 0;
    }
    return 1;
}

void test_alloc_array_basic(void) {
    int32_t elem_size = sizeof(int);
    int32_t rank = 3;
    int64_t dim1 = 2, dim2 = 3, dim3 = 4;

    __public__array_t* arr = __public__alloc_array(elem_size, rank, dim1, dim2, dim3);

    TEST_ASSERT_NOT_NULL(arr);
    TEST_ASSERT_NOT_NULL(arr->data);
    TEST_ASSERT_NOT_NULL(arr->shape);
    TEST_ASSERT_EQUAL_INT(dim1, arr->length);
    TEST_ASSERT_EQUAL_INT(rank, arr->rank);

    /* array must be allocated in __arena__ */
    release(__arena__, arr);
}

void test_alloc_array_zero_rank(void) {
    int count = 5;
    int32_t elem_size = sizeof(double);
    int32_t rank = 0;

    __public__array_t* arr = __public__alloc_array(elem_size, rank, count);

    TEST_ASSERT_NOT_NULL(arr);
    TEST_ASSERT_NOT_NULL(arr->data);
    TEST_ASSERT_NULL(arr->shape); // shape should be NULL for rank 0
    TEST_ASSERT_EQUAL_INT(count, arr->length);
    TEST_ASSERT_EQUAL_INT(rank, arr->rank);
    TEST_ASSERT_TRUE(is_zeroed(arr->data, count * elem_size));

    /* array must be allocated in __arena__ */
    release(__arena__, arr);
}

void test_alloc_array_zero_count(void) {
    int count = 0;
    int32_t elem_size = sizeof(int64_t);
    int32_t rank = 2;

    __public__array_t* arr = __public__alloc_array(elem_size, rank, count);

    TEST_ASSERT_NOT_NULL(arr);
    TEST_ASSERT_NOT_NULL(arr->data); // still allocated, but size 0
    TEST_ASSERT_NOT_NULL(arr->shape); 
    TEST_ASSERT_EQUAL_INT(count, arr->length);
    TEST_ASSERT_EQUAL_INT(rank, arr->rank);

    /* array must be allocated in __arena__ */
    release(__arena__, arr);
}

void test_alloc_array_large(void) {
    int count = 1000;
    int32_t elem_size = sizeof(float);
    int32_t rank = 5;

    __public__array_t* arr = __public__alloc_array(elem_size, rank, count);

    TEST_ASSERT_NOT_NULL(arr);
    TEST_ASSERT_NOT_NULL(arr->data);
    TEST_ASSERT_NOT_NULL(arr->shape);
    TEST_ASSERT_EQUAL_INT(count, arr->length);
    TEST_ASSERT_EQUAL_INT(rank, arr->rank);
    TEST_ASSERT_TRUE(is_zeroed(arr->data, count * elem_size));

    /* array must be allocated in __arena__ */
    release(__arena__, arr);
}

void test_alloc_array_single_element(void) {
    int count = 1;
    int32_t elem_size = sizeof(char);
    int32_t rank = 1;

    __public__array_t* arr = __public__alloc_array(elem_size, rank, count);

    TEST_ASSERT_NOT_NULL(arr);
    TEST_ASSERT_NOT_NULL(arr->data);
    TEST_ASSERT_NOT_NULL(arr->shape);
    TEST_ASSERT_EQUAL_INT(1, arr->length);
    TEST_ASSERT_EQUAL_INT(1, arr->rank);
    TEST_ASSERT_TRUE(is_zeroed(arr->data, count * elem_size));

    /* array must be allocated in __arena__ */
    release(__arena__, arr);
}


int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();
    __arena__ = gc_create_global_arena();

    RUN_TEST(test_alloc_array_basic);
    // RUN_TEST(test_alloc_array_zero_rank);
    // RUN_TEST(test_alloc_array_zero_count);
    // RUN_TEST(test_alloc_array_large);
    // RUN_TEST(test_alloc_array_single_element);

    return UNITY_END();
}