#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include "str.h"
#include "alloc.h"
#include "gc.h"

arena_t* __test__global__arena__;
extern __thread arena_t* __arena__;
extern arena_t* __global__arena__;

void setUp(void) {
    __test__global__arena__ = arena_create();
}

void tearDown(void) {
    __test__global__arena__ = NULL;
}

__public__string_t* create_test_string(const char* data, size_t size) {
    __public__string_t* s = allocate(__test__global__arena__, sizeof(__public__string_t));
    s->data = allocate(__test__global__arena__, size + 1);
    memcpy(s->data, data, size);
    s->data[size] = '\0';
    s->size = size;
    return s;
}

void test_string_alloc_basic(void) {
    const char* test_data = "Hello, World!";
    size_t len = strlen(test_data);
    
    __public__string_t* str = create_test_string(test_data, len);
    
    TEST_ASSERT_NOT_NULL(str);
    TEST_ASSERT_NOT_NULL(str->data);
    TEST_ASSERT_EQUAL_INT(len, str->size);
    TEST_ASSERT_EQUAL_STRING(test_data, str->data);
}

void test_string_alloc_empty(void) {
    __public__string_t* str = create_test_string("", 0);
    
    TEST_ASSERT_NOT_NULL(str);
    TEST_ASSERT_NOT_NULL(str->data);
    TEST_ASSERT_EQUAL_INT(0, str->size);
    TEST_ASSERT_EQUAL_STRING("", str->data);
}

void test_string_alloc_large(void) {
    size_t large_size = 10000;
    char* large_data = malloc(large_size + 1);
    memset(large_data, 'A', large_size);
    large_data[large_size] = '\0';
    
    __public__string_t* str = create_test_string(large_data, large_size);
    
    TEST_ASSERT_NOT_NULL(str);
    TEST_ASSERT_NOT_NULL(str->data);
    TEST_ASSERT_EQUAL_INT(large_size, str->size);
    TEST_ASSERT_EQUAL_MEMORY(large_data, str->data, large_size);
    
    free(large_data);
}

void test_string_with_null_bytes(void) {
    char data[] = {'H', 'e', 'l', '\0', 'l', 'o'};
    size_t size = sizeof(data);
    
    __public__string_t* str = create_test_string(data, size);
    
    TEST_ASSERT_NOT_NULL(str);
    TEST_ASSERT_EQUAL_INT(size, str->size);
    TEST_ASSERT_EQUAL_MEMORY(data, str->data, size);
}

void test_string_unicode(void) {
    const char* unicode_str = "Hello 世界 🌍";
    size_t len = strlen(unicode_str);
    
    __public__string_t* str = create_test_string(unicode_str, len);
    
    TEST_ASSERT_NOT_NULL(str);
    TEST_ASSERT_EQUAL_INT(len, str->size);
    TEST_ASSERT_EQUAL_STRING(unicode_str, str->data);
}

void test_string_special_chars(void) {
    const char* special = "Tab:\t Newline:\n Quote:\" Backslash:\\";
    size_t len = strlen(special);
    
    __public__string_t* str = create_test_string(special, len);
    
    TEST_ASSERT_NOT_NULL(str);
    TEST_ASSERT_EQUAL_INT(len, str->size);
    TEST_ASSERT_EQUAL_STRING(special, str->data);
}

void test_string_boundary_sizes(void) {
    /* Test various boundary sizes */
    size_t sizes[] = {1, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127, 128, 129, 255, 256, 257};
    
    for (int i = 0; i < sizeof(sizes) / sizeof(sizes[0]); i++) {
        char* data = malloc(sizes[i] + 1);
        memset(data, 'X', sizes[i]);
        data[sizes[i]] = '\0';
        
        __public__string_t* str = create_test_string(data, sizes[i]);
        
        TEST_ASSERT_NOT_NULL(str);
        TEST_ASSERT_EQUAL_INT(sizes[i], str->size);
        TEST_ASSERT_EQUAL_MEMORY(data, str->data, sizes[i]);
        
        free(data);
    }
}

void test_string_multiple_allocations(void) {
    #define NUM_STRINGS 100
    __public__string_t* strings[NUM_STRINGS];
    
    for (int i = 0; i < NUM_STRINGS; i++) {
        char buffer[32];
        snprintf(buffer, sizeof(buffer), "String %d", i);
        strings[i] = create_test_string(buffer, strlen(buffer));
        TEST_ASSERT_NOT_NULL(strings[i]);
    }
    
    /* Verify all strings are still valid */
    for (int i = 0; i < NUM_STRINGS; i++) {
        char expected[32];
        snprintf(expected, sizeof(expected), "String %d", i);
        TEST_ASSERT_EQUAL_STRING(expected, strings[i]->data);
    }
}

int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();
    __arena__ = gc_create_global_arena();

    RUN_TEST(test_string_alloc_basic);
    RUN_TEST(test_string_alloc_empty);
    RUN_TEST(test_string_alloc_large);
    RUN_TEST(test_string_with_null_bytes);
    RUN_TEST(test_string_unicode);
    RUN_TEST(test_string_special_chars);
    RUN_TEST(test_string_boundary_sizes);
    RUN_TEST(test_string_multiple_allocations);

    return UNITY_END();
}
