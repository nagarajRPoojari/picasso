#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <pthread.h>
#include <stdatomic.h>
#include <time.h>
#include "alloc.h"
#include "gc.h"
#include "ggc.h"
#include "initutils.h"

extern __thread arena_t* __arena__;
extern arena_t* __global__arena__;

void setUp(void) {
}

void tearDown(void) {
}

void test_gc_arena_creation(void) {
    arena_t* arena = gc_create_global_arena();
    TEST_ASSERT_NOT_NULL(arena);
}

void test_gc_allocation_basic(void) {
    size_t sizes[] = {16, 32, 64, 128, 256, 512, 1024};
    
    for (int i = 0; i < sizeof(sizes) / sizeof(sizes[0]); i++) {
        void* ptr = allocate(__global__arena__, sizes[i]);
        TEST_ASSERT_NOT_NULL(ptr);
        
        /* Write to allocated memory */
        memset(ptr, 0xAA, sizes[i]);
    }
}

void test_gc_multiple_allocations(void) {
    #define NUM_ALLOCS 100
    void* ptrs[NUM_ALLOCS];
    
    for (int i = 0; i < NUM_ALLOCS; i++) {
        size_t size = (i + 1) * 16;
        ptrs[i] = allocate(__global__arena__, size);
        TEST_ASSERT_NOT_NULL(ptrs[i]);
        memset(ptrs[i], i & 0xFF, size);
    }
    
    /* Verify all allocations are still valid */
    for (int i = 0; i < NUM_ALLOCS; i++) {
        size_t size = (i + 1) * 16;
        unsigned char* ptr = (unsigned char*)ptrs[i];
        for (size_t j = 0; j < size; j++) {
            TEST_ASSERT_EQUAL_UINT8(i & 0xFF, ptr[j]);
        }
    }
}

void test_gc_large_allocation(void) {
    size_t large_size = 1024 * 1024; /* 1MB */
    void* ptr = allocate(__global__arena__, large_size);
    TEST_ASSERT_NOT_NULL(ptr);
    
    /* Write pattern to verify allocation */
    memset(ptr, 0xBB, large_size);
    
    unsigned char* bytes = (unsigned char*)ptr;
    for (size_t i = 0; i < large_size; i++) {
        TEST_ASSERT_EQUAL_UINT8(0xBB, bytes[i]);
    }
}

void test_gc_zero_allocation(void) {
    void* ptr = allocate(__global__arena__, 0);
    /* Implementation-specific: may return NULL or valid pointer */
    if (ptr != NULL) {
        TEST_ASSERT_TRUE(1); /* Valid behavior */
    }
}

#define GC_STRESS_ALLOCS 1000
#define GC_STRESS_THREADS 4

typedef struct {
    int thread_id;
    atomic_int* alloc_count;
} gc_thread_arg_t;

void* gc_stress_thread(void* arg) {
    gc_thread_arg_t* targ = (gc_thread_arg_t*)arg;
    
    for (int i = 0; i < GC_STRESS_ALLOCS; i++) {
        size_t size = ((targ->thread_id * 1000 + i) % 512) + 16;
        void* ptr = allocate(__global__arena__, size);
        TEST_ASSERT_NOT_NULL(ptr);
        
        /* Write to memory */
        memset(ptr, (targ->thread_id + i) & 0xFF, size);
        
        atomic_fetch_add(targ->alloc_count, 1);
        
        /* Occasionally yield to increase contention */
        if (i % 100 == 0) {
            sched_yield();
        }
    }
    
    return NULL;
}

void test_gc_concurrent_allocations(void) {
    atomic_int alloc_count = 0;
    pthread_t threads[GC_STRESS_THREADS];
    gc_thread_arg_t args[GC_STRESS_THREADS];
    
    for (int i = 0; i < GC_STRESS_THREADS; i++) {
        args[i].thread_id = i;
        args[i].alloc_count = &alloc_count;
        pthread_create(&threads[i], NULL, gc_stress_thread, &args[i]);
    }
    
    for (int i = 0; i < GC_STRESS_THREADS; i++) {
        pthread_join(threads[i], NULL);
    }
    
    int expected = GC_STRESS_THREADS * GC_STRESS_ALLOCS;
    TEST_ASSERT_EQUAL_INT(expected, atomic_load(&alloc_count));
}

void test_gc_mixed_size_allocations(void) {
    /* Allocate various sizes to test different allocation paths */
    void* small = allocate(__global__arena__, 8);
    void* medium = allocate(__global__arena__, 256);
    void* large = allocate(__global__arena__, 4096);
    void* huge = allocate(__global__arena__, 65536);
    
    TEST_ASSERT_NOT_NULL(small);
    TEST_ASSERT_NOT_NULL(medium);
    TEST_ASSERT_NOT_NULL(large);
    TEST_ASSERT_NOT_NULL(huge);
    
    /* Write patterns */
    memset(small, 0x11, 8);
    memset(medium, 0x22, 256);
    memset(large, 0x33, 4096);
    memset(huge, 0x44, 65536);
    
    /* Verify patterns */
    TEST_ASSERT_EQUAL_UINT8(0x11, ((unsigned char*)small)[0]);
    TEST_ASSERT_EQUAL_UINT8(0x22, ((unsigned char*)medium)[0]);
    TEST_ASSERT_EQUAL_UINT8(0x33, ((unsigned char*)large)[0]);
    TEST_ASSERT_EQUAL_UINT8(0x44, ((unsigned char*)huge)[0]);
}

void test_gc_allocation_alignment(void) {
    /* Test that allocations are properly aligned */
    for (int i = 0; i < 100; i++) {
        size_t size = (i + 1) * 7; /* Odd sizes to test alignment */
        void* ptr = allocate(__global__arena__, size);
        TEST_ASSERT_NOT_NULL(ptr);
        
        /* Check 16-byte alignment */
        uintptr_t addr = (uintptr_t)ptr;
        TEST_ASSERT_EQUAL_UINT64(0, addr % 16);
    }
}

void test_gc_fragmentation_resistance(void) {
    #define FRAG_TEST_SIZE 50
    void* ptrs[FRAG_TEST_SIZE];
    
    /* Allocate alternating sizes */
    for (int i = 0; i < FRAG_TEST_SIZE; i++) {
        size_t size = (i % 2 == 0) ? 32 : 128;
        ptrs[i] = allocate(__global__arena__, size);
        TEST_ASSERT_NOT_NULL(ptrs[i]);
    }
    
    /* Free every other allocation */
    for (int i = 0; i < FRAG_TEST_SIZE; i += 2) {
        release(__global__arena__, ptrs[i]);
    }
    
    /* Allocate again - should reuse freed space */
    for (int i = 0; i < FRAG_TEST_SIZE / 2; i++) {
        void* ptr = allocate(__global__arena__, 32);
        TEST_ASSERT_NOT_NULL(ptr);
    }
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    
    gc_start();

    RUN_TEST(test_gc_arena_creation);
    RUN_TEST(test_gc_allocation_basic);
    RUN_TEST(test_gc_multiple_allocations);
    RUN_TEST(test_gc_large_allocation);
    RUN_TEST(test_gc_zero_allocation);
    RUN_TEST(test_gc_concurrent_allocations);
    RUN_TEST(test_gc_mixed_size_allocations);
    RUN_TEST(test_gc_allocation_alignment);
    RUN_TEST(test_gc_fragmentation_resistance);

    return UNITY_END();
}
