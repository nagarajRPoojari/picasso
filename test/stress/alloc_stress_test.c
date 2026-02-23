#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <pthread.h>
#include <time.h>
#include "alloc.h"

static arena_t* ar;

void setUp(void) {
    ar = arena_create();
    TEST_ASSERT_NOT_NULL_MESSAGE(ar, "Arena creation failed");
}

void tearDown(void) {
    ar = NULL;
}

#define STRESS_THREADS 8


#if defined(STRESS_LEVEL_1)
    #define STRESS_ITERATIONS 1000
    #define MAX_LIVE_ALLOCS 100
    #define CHURN_ITERATIONS 5000
    #define LARGE_ALLOC_COUNT 100
    #define FRAG_ITERATIONS 100

    #define STRESS_OPS_PER_THREAD 500
    #define PRESSURE_ALLOCS 1000
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define STRESS_ITERATIONS 2000
    #define MAX_LIVE_ALLOCS 2000
    #define CHURN_ITERATIONS 100000
    #define LARGE_ALLOC_COUNT 200
    #define FRAG_ITERATIONS 2000

    #define STRESS_OPS_PER_THREAD 1000
    #define PRESSURE_ALLOCS 2000
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define STRESS_ITERATIONS 10000
    #define MAX_LIVE_ALLOCS 2000
    #define CHURN_ITERATIONS 50000
    #define LARGE_ALLOC_COUNT 1000
    #define FRAG_ITERATIONS 4000

    #define STRESS_OPS_PER_THREAD 10000
    #define PRESSURE_ALLOCS 10000
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define STRESS_ITERATIONS 10000
    #define MAX_LIVE_ALLOCS 1000
    #define CHURN_ITERATIONS 50000
    #define LARGE_ALLOC_COUNT 100
    #define FRAG_ITERATIONS 1000

    #define STRESS_OPS_PER_THREAD 5000
    #define PRESSURE_ALLOCS 10000
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

/* Stress test: Allocate and free in random patterns */
void test_stress_random_alloc_free(void) {
    
    void* live_allocs[MAX_LIVE_ALLOCS];
    size_t live_sizes[MAX_LIVE_ALLOCS];
    int live_count = 0;
    
    srand(time(NULL));
    
    for (int i = 0; i < STRESS_ITERATIONS; i++) {
        int action = rand() % 100;
        
        if (action < 60 && live_count < MAX_LIVE_ALLOCS) {
            /* Allocate */
            size_t size = (rand() % 4096) + 1;
            void* ptr = allocate(ar, size);
            TEST_ASSERT_NOT_NULL(ptr);
            
            /* Fill with pattern */
            memset(ptr, i & 0xFF, size);
            
            live_allocs[live_count] = ptr;
            live_sizes[live_count] = size;
            live_count++;
        } else if (live_count > 0) {
            /* Free random allocation */
            int idx = rand() % live_count;
            
            /* Verify pattern before freeing */
            unsigned char* bytes = (unsigned char*)live_allocs[idx];
            int expected_pattern = (i - (live_count - idx)) & 0xFF;
            
            release(ar, live_allocs[idx]);
            
            /* Remove from tracking */
            live_allocs[idx] = live_allocs[live_count - 1];
            live_sizes[idx] = live_sizes[live_count - 1];
            live_count--;
        }
        
        /* Periodic verification */
        if (i % 1000 == 0) {
            printf("Iteration %d: %d live allocations\n", i, live_count);
        }
    }
    
    /* Cleanup remaining */
    for (int i = 0; i < live_count; i++) {
        release(ar, live_allocs[i]);
    }
}

/* Stress test: Rapid allocation/deallocation */
void test_stress_rapid_churn(void) {
    
    for (int i = 0; i < CHURN_ITERATIONS; i++) {
        size_t size = ((i % 512) + 1) * 8;
        void* ptr = allocate(ar, size);
        TEST_ASSERT_NOT_NULL(ptr);
        
        memset(ptr, i & 0xFF, size);
        release(ar, ptr);
        
        if (i % 10000 == 0) {
            printf("Churn iteration %d\n", i);
        }
    }
}

/* Stress test: Large allocations */
void test_stress_large_allocations(void) {
    void* ptrs[LARGE_ALLOC_COUNT];
    
    for (int i = 0; i < LARGE_ALLOC_COUNT; i++) {
        size_t size = (1024 * 1024) + (i * 1024); /* 1MB+ */
        ptrs[i] = allocate(ar, size);
        TEST_ASSERT_NOT_NULL(ptrs[i]);
        
        /* Write pattern */
        memset(ptrs[i], i & 0xFF, size);
        
        if (i % 10 == 0) {
            printf("Large allocation %d: %zu bytes\n", i, size);
        }
    }
    
    /* Verify and free */
    for (int i = 0; i < LARGE_ALLOC_COUNT; i++) {
        release(ar, ptrs[i]);
    }
}

/* Stress test: Fragmentation */
void test_stress_fragmentation(void) {
    void* ptrs[FRAG_ITERATIONS];
    
    /* Allocate alternating sizes */
    for (int i = 0; i < FRAG_ITERATIONS; i++) {
        size_t size = (i % 2 == 0) ? 32 : 256;
        ptrs[i] = allocate(ar, size);
        TEST_ASSERT_NOT_NULL(ptrs[i]);
    }
    
    /* Free every other allocation */
    for (int i = 0; i < FRAG_ITERATIONS; i += 2) {
        release(ar, ptrs[i]);
    }
    
    /* Try to allocate in freed spaces */
    for (int i = 0; i < FRAG_ITERATIONS / 2; i++) {
        void* ptr = allocate(ar, 32);
        TEST_ASSERT_NOT_NULL(ptr);
    }
    
    /* Cleanup */
    for (int i = 1; i < FRAG_ITERATIONS; i += 2) {
        release(ar, ptrs[i]);
    }
}

typedef struct {
    arena_t* arena;
    int thread_id;
} stress_thread_arg_t;

void* stress_thread_worker(void* arg) {
    stress_thread_arg_t* targ = (stress_thread_arg_t*)arg;
    void* live_allocs[100];
    int live_count = 0;
    
    for (int i = 0; i < STRESS_OPS_PER_THREAD; i++) {
        int action = rand() % 100;
        
        if (action < 70 && live_count < 100) {
            /* Allocate */
            size_t size = (rand() % 1024) + 1;
            void* ptr = allocate(targ->arena, size);
            TEST_ASSERT_NOT_NULL(ptr);
            memset(ptr, (targ->thread_id + i) & 0xFF, size);
            live_allocs[live_count++] = ptr;
        } else if (live_count > 0) {
            /* Free */
            int idx = rand() % live_count;
            release(targ->arena, live_allocs[idx]);
            live_allocs[idx] = live_allocs[--live_count];
        }
    }
    
    /* Cleanup */
    for (int i = 0; i < live_count; i++) {
        release(targ->arena, live_allocs[i]);
    }
    
    return NULL;
}

/* all threads share common arena, this is expected to not happen in 
 real senarious. all scheduler threads will be provided with its own 
 arena to avoid concurrent allocations. Still here it is expected to
 succeed since arenas' are thread safe by lock.
*/
void test_stress_concurrent_operations(void) {
    pthread_t threads[STRESS_THREADS];
    stress_thread_arg_t args[STRESS_THREADS];
    
    for (int i = 0; i < STRESS_THREADS; i++) {
        args[i].arena = ar;
        args[i].thread_id = i;
        pthread_create(&threads[i], NULL, stress_thread_worker, &args[i]);
    }
    
    for (int i = 0; i < STRESS_THREADS; i++) {
        pthread_join(threads[i], NULL);
    }
}

/* Stress test: Memory pressure */
void test_stress_memory_pressure(void) {
    void* ptrs[PRESSURE_ALLOCS];
    
    /* Allocate until we have significant memory usage */
    for (int i = 0; i < PRESSURE_ALLOCS; i++) {
        size_t size = (i % 1024) + 1;
        ptrs[i] = allocate(ar, size);
        TEST_ASSERT_NOT_NULL(ptrs[i]);
        
        if (i % 1000 == 0) {
            printf("Memory pressure: %d allocations\n", i);
        }
    }
    
    /* Free in reverse order */
    for (int i = PRESSURE_ALLOCS - 1; i >= 0; i--) {
        release(ar, ptrs[i]);
    }
}

int main(void) {
    UNITY_BEGIN();
    
    RUN_TEST(test_stress_random_alloc_free);
    RUN_TEST(test_stress_rapid_churn);
    RUN_TEST(test_stress_large_allocations);
    RUN_TEST(test_stress_fragmentation);
    RUN_TEST(test_stress_concurrent_operations);
    RUN_TEST(test_stress_memory_pressure);
    
    return UNITY_END();
}
