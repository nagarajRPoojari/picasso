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

    TEST_ASSERT_NULL(arr);
}

void test_alloc_array_large(void) {
    int count = 2000;
    int32_t elem_size = sizeof(float);
    int32_t rank = 2;

    __public__array_t* arr = __public__alloc_array(elem_size, rank, count, count);
    TEST_ASSERT_NOT_NULL(arr);
    TEST_ASSERT_NOT_NULL(arr->data);
    TEST_ASSERT_NOT_NULL(arr->shape);
    TEST_ASSERT_EQUAL_INT(count, arr->length);
    TEST_ASSERT_EQUAL_INT(rank, arr->rank);
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


void test_alloc_array_multidimensional(void) {
    int32_t elem_size = sizeof(int);
    int32_t rank = 3;
    int64_t dim1 = 2, dim2 = 3, dim3 = 4;

    __public__array_t* arr = __public__alloc_array(elem_size, rank, dim1, dim2, dim3);

    TEST_ASSERT_NOT_NULL(arr);
    TEST_ASSERT_NOT_NULL(arr->data);
    TEST_ASSERT_NOT_NULL(arr->shape);
    TEST_ASSERT_EQUAL_INT(dim1, arr->length);
    TEST_ASSERT_EQUAL_INT(rank, arr->rank);
    
    /* Verify shape array */
    TEST_ASSERT_EQUAL_INT64(dim1, arr->shape[0]);
    TEST_ASSERT_EQUAL_INT64(dim2, arr->shape[1]);
    TEST_ASSERT_EQUAL_INT64(dim3, arr->shape[2]);
    
    release(__arena__, arr);
}

void test_alloc_array_write_read(void) {
    int count = 10;
    int32_t elem_size = sizeof(int);
    int32_t rank = 1;

    __public__array_t* arr = __public__alloc_array(elem_size, rank, count);
    
    /* Write pattern */
    int* data = (int*)arr->data;
    for (int i = 0; i < count; i++) {
        data[i] = i * 10;
    }
    
    /* Read and verify */
    for (int i = 0; i < count; i++) {
        TEST_ASSERT_EQUAL_INT(i * 10, data[i]);
    }

    release(__arena__, arr);
}

void test_alloc_array_boundary_sizes(void) {
    int32_t elem_size = sizeof(double);
    int32_t rank = 1;
    int64_t sizes[] = {1, 2, 7, 8, 15, 16, 31, 32, 63, 64, 127, 128, 255, 256};
    
    for (int i = 0; i < sizeof(sizes) / sizeof(sizes[0]); i++) {
        __public__array_t* arr = __public__alloc_array(elem_size, rank, sizes[i]);
        
        TEST_ASSERT_NOT_NULL(arr);
        TEST_ASSERT_NOT_NULL(arr->data);
        TEST_ASSERT_EQUAL_INT64(sizes[i], arr->length);
        TEST_ASSERT_TRUE(is_zeroed(arr->data, sizes[i] * elem_size));
        
        release(__arena__, arr);
    }
}

void test_alloc_array_different_types(void) {
    /* Test different element sizes */
    struct {
        int32_t elem_size;
        int64_t count;
    } test_cases[] = {
        {sizeof(char), 100},
        {sizeof(short), 50},
        {sizeof(int), 25},
        {sizeof(long), 20},
        {sizeof(float), 30},
        {sizeof(double), 15},
    };
    
    for (int i = 0; i < sizeof(test_cases) / sizeof(test_cases[0]); i++) {
        __public__array_t* arr = __public__alloc_array(
            test_cases[i].elem_size, 1, test_cases[i].count);
        
        TEST_ASSERT_NOT_NULL(arr);
        TEST_ASSERT_NOT_NULL(arr->data);
        TEST_ASSERT_EQUAL_INT64(test_cases[i].count, arr->length);
        
        release(__arena__, arr);
    }
}

int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();
    __arena__ = gc_create_global_arena();

    RUN_TEST(test_alloc_array_basic);
    RUN_TEST(test_alloc_array_zero_rank);
    RUN_TEST(test_alloc_array_large);
    RUN_TEST(test_alloc_array_single_element);
    RUN_TEST(test_alloc_array_multidimensional);
    RUN_TEST(test_alloc_array_write_read);
    RUN_TEST(test_alloc_array_boundary_sizes);
    RUN_TEST(test_alloc_array_different_types);

    return UNITY_END();
}