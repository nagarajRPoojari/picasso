#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <pthread.h>
#include <stdatomic.h>
#include "queue.h"
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

void test_safe_queue_init(void) {
    safe_q_t queue;
    int result = safe_q_init(&queue, 10);
    
    TEST_ASSERT_EQUAL_INT(0, result);
    TEST_ASSERT_EQUAL_INT(10, queue.capacity);
    TEST_ASSERT_EQUAL_INT(0, atomic_load(&queue.head));
    TEST_ASSERT_EQUAL_INT(0, atomic_load(&queue.tail));
}

void test_safe_queue_push_pop_single(void) {
    safe_q_t queue;
    safe_q_init(&queue, 10);
    
    void* test_ptr = (void*)0x1234;
    int push_result = safe_q_push(&queue, test_ptr);
    TEST_ASSERT_EQUAL_INT(0, push_result);
    
    void* popped = safe_q_pop(&queue);
    TEST_ASSERT_EQUAL_PTR(test_ptr, popped);
}

void test_safe_queue_push_pop_multiple(void) {
    safe_q_t queue;
    safe_q_init(&queue, 10);
    
    void* ptrs[] = {(void*)0x1, (void*)0x2, (void*)0x3, (void*)0x4, (void*)0x5};
    int num_ptrs = sizeof(ptrs) / sizeof(ptrs[0]);
    
    /* Push all */
    for (int i = 0; i < num_ptrs; i++) {
        int result = safe_q_push(&queue, ptrs[i]);
        TEST_ASSERT_EQUAL_INT(0, result);
    }
    
    /* Pop all in FIFO order */
    for (int i = 0; i < num_ptrs; i++) {
        void* popped = safe_q_pop(&queue);
        TEST_ASSERT_EQUAL_PTR(ptrs[i], popped);
    }
}

void test_safe_queue_empty_pop(void) {
    safe_q_t queue;
    safe_q_init(&queue, 10);
    
    void* popped = safe_q_pop(&queue);
    TEST_ASSERT_NULL(popped);
}

void test_safe_queue_full_push(void) {
    safe_q_t queue;
    int capacity = 5;
    safe_q_init(&queue, capacity);
    
    /* Fill the queue */
    for (int i = 0; i < capacity; i++) {
        int result = safe_q_push(&queue, (void*)(uintptr_t)i);
        TEST_ASSERT_EQUAL_INT(0, result);
    }
    
    /* Try to push one more - should fail */
    int result = safe_q_push(&queue, (void*)0xFFFF);
    TEST_ASSERT_EQUAL_INT(-1, result);
}

void test_safe_queue_wraparound(void) {
    safe_q_t queue;
    int capacity = 5;
    safe_q_init(&queue, capacity);
    
    /* Fill and empty multiple times to test wraparound */
    for (int cycle = 0; cycle < 3; cycle++) {
        /* Fill */
        for (int i = 0; i < capacity; i++) {
            int val = cycle * 100 + i;
            safe_q_push(&queue, (void*)(uintptr_t)val);
        }
        
        /* Empty */
        for (int i = 0; i < capacity; i++) {
            void* popped = safe_q_pop(&queue);
            int expected = cycle * 100 + i;
            TEST_ASSERT_EQUAL_PTR((void*)(uintptr_t)expected, popped);
        }
    }
}

void test_safe_queue_interleaved_ops(void) {
    safe_q_t queue;
    safe_q_init(&queue, 10);
    
    /* Interleave push and pop operations */
    safe_q_push(&queue, (void*)0x1);
    safe_q_push(&queue, (void*)0x2);
    
    void* p1 = safe_q_pop(&queue);
    TEST_ASSERT_EQUAL_PTR((void*)0x1, p1);
    
    safe_q_push(&queue, (void*)0x3);
    safe_q_push(&queue, (void*)0x4);
    
    void* p2 = safe_q_pop(&queue);
    TEST_ASSERT_EQUAL_PTR((void*)0x2, p2);
    
    void* p3 = safe_q_pop(&queue);
    TEST_ASSERT_EQUAL_PTR((void*)0x3, p3);
    
    void* p4 = safe_q_pop(&queue);
    TEST_ASSERT_EQUAL_PTR((void*)0x4, p4);
    
    void* p5 = safe_q_pop(&queue);
    TEST_ASSERT_NULL(p5);
}

#define NUM_THREADS 4
#define OPS_PER_THREAD 1000

typedef struct {
    safe_q_t* queue;
    int thread_id;
    atomic_int* push_count;
    atomic_int* pop_count;
} thread_arg_t;

void* producer_thread(void* arg) {
    thread_arg_t* targ = (thread_arg_t*)arg;
    
    for (int i = 0; i < OPS_PER_THREAD; i++) {
        void* data = (void*)(uintptr_t)((targ->thread_id << 16) | i);
        
        /* Retry on full queue */
        while (safe_q_push(targ->queue, data) != 0) {
            sched_yield();
        }
        atomic_fetch_add(targ->push_count, 1);
    }
    
    return NULL;
}

void* consumer_thread(void* arg) {
    thread_arg_t* targ = (thread_arg_t*)arg;
    
    for (int i = 0; i < OPS_PER_THREAD; i++) {
        void* data;
        
        /* Retry on empty queue */
        while ((data = safe_q_pop(targ->queue)) == NULL) {
            sched_yield();
        }
        atomic_fetch_add(targ->pop_count, 1);
    }
    
    return NULL;
}

void test_safe_queue_concurrent(void) {
    safe_q_t queue;
    safe_q_init(&queue, 100);
    
    atomic_int push_count = 0;
    atomic_int pop_count = 0;
    
    pthread_t producers[NUM_THREADS];
    pthread_t consumers[NUM_THREADS];
    thread_arg_t args[NUM_THREADS * 2];
    
    /* Start producers */
    for (int i = 0; i < NUM_THREADS; i++) {
        args[i].queue = &queue;
        args[i].thread_id = i;
        args[i].push_count = &push_count;
        args[i].pop_count = &pop_count;
        pthread_create(&producers[i], NULL, producer_thread, &args[i]);
    }
    
    /* Start consumers */
    for (int i = 0; i < NUM_THREADS; i++) {
        args[NUM_THREADS + i].queue = &queue;
        args[NUM_THREADS + i].thread_id = i;
        args[NUM_THREADS + i].push_count = &push_count;
        args[NUM_THREADS + i].pop_count = &pop_count;
        pthread_create(&consumers[i], NULL, consumer_thread, &args[NUM_THREADS + i]);
    }
    
    /* Wait for all threads */
    for (int i = 0; i < NUM_THREADS; i++) {
        pthread_join(producers[i], NULL);
        pthread_join(consumers[i], NULL);
    }
    
    /* Verify counts */
    int expected_total = NUM_THREADS * OPS_PER_THREAD;
    TEST_ASSERT_EQUAL_INT(expected_total, atomic_load(&push_count));
    TEST_ASSERT_EQUAL_INT(expected_total, atomic_load(&pop_count));
}

void test_safe_queue_null_pointer(void) {
    safe_q_t queue;
    safe_q_init(&queue, 10);
    
    /* Push NULL pointer - should be allowed */
    int result = safe_q_push(&queue, NULL);
    TEST_ASSERT_EQUAL_INT(0, result);
    
    /* Pop should return NULL */
    void* popped = safe_q_pop(&queue);
    TEST_ASSERT_NULL(popped);
}

int main(void) {
    UNITY_BEGIN();
    __global__arena__ = gc_create_global_arena();
    __arena__ = gc_create_global_arena();

    RUN_TEST(test_safe_queue_init);
    RUN_TEST(test_safe_queue_push_pop_single);
    RUN_TEST(test_safe_queue_push_pop_multiple);
    RUN_TEST(test_safe_queue_empty_pop);
    RUN_TEST(test_safe_queue_full_push);
    RUN_TEST(test_safe_queue_wraparound);
    RUN_TEST(test_safe_queue_interleaved_ops);
    RUN_TEST(test_safe_queue_concurrent);
    RUN_TEST(test_safe_queue_null_pointer);

    return UNITY_END();
}
