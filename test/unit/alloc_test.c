#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include "alloc.h"

static arena_t* ar;


static int is_aligned(void* p) {
    return ((uintptr_t)p % 16) == 0;
}

/* write specific pattern to check overlap later */
static void fill_pattern(void* p, size_t size, uint8_t pattern) {
    if (p) memset(p, pattern, size);
}

static int check_pattern(void* p, size_t size, uint8_t pattern) {
    uint8_t* byte_ptr = (uint8_t*)p;
    for (size_t i = 0; i < size; i++) {
        if (byte_ptr[i] != pattern) return 0;
    }
    return 1;
}

void setUp(void) {
    ar = arena_create();
    TEST_ASSERT_NOT_NULL_MESSAGE(ar, "Arena creation failed");
}

void tearDown(void) {
    ar = NULL;
}


void test_basic_allocation_and_alignment(void) {
    size_t sizes[] = {1, 8, 16, 24, 32, 128, 1024};
    
    for (int i = 0; i < 7; i++) {
        void* p = allocate(ar, sizes[i]);
        TEST_ASSERT_NOT_NULL(p);
        TEST_ASSERT_TRUE_MESSAGE(is_aligned(p), "Pointer not 16-byte aligned");
        
        /* ensure we can write to the full extent without segfaulting */
        memset(p, 0xAA, sizes[i]);
    }
}

void test_zero_allocation(void) {
    /* Standard behavior varies, but usually returns NULL or a specific unique pointer */
    /* adjust assertion based on your specific implementation requirements */
    void* p = allocate(ar, 0);
    if (p != NULL) {
        TEST_ASSERT_TRUE(is_aligned(p));
    }
}

void test_free_null_is_safe(void) {
    // release(ar, NULL); 
}


/* test fastbin */
void test_fastbin_reuse(void) {
    /* alloc 10 chunks of same size */
    int n = 10;
    void* ptrs[n];
    for(int i=0; i<n ;i++)
        ptrs[i] = allocate(ar, 32);

    /* free in specific order */
    for(int i=0; i<n; i++)
        release(ar, ptrs[i]);

    /* re-alloc should return reverse order of free */
    for(int i=0; i<n; i++){
        void* r = allocate(ar, 32);
        TEST_ASSERT_EQUAL_PTR_MESSAGE(ptrs[n-i-1], r, "allocated invalid pointer");
    }
}

/* test smallbin */
void test_smallbin_reuse(void) {
    int n = 7;
    void* ptrs[n];
    /* [64][128][][]...[top_chunk] */
    for(int i=0; i<n; i++)
        ptrs[i] = allocate(ar, (i+1)*64);

    /* don't release last chunk to avoid merge with top_chunk */
    /* [64-free][128][192-free]....[top_chunk] */
    for(int i=0; i<n - 1; i++) {
        if(i%2==0) release(ar, ptrs[i]);
    }
    /* ar->unsortedbin->[576]->[448]->[320]->[192] */

    /* allocate something large to force all bins in unsortedbin to smallbin */
    allocate(ar, 576); /* slightly smaller than [576] biggest chunk */

    /* check smallbins */
    for(int i=1; i<n-1; i++){
        if(i%2==0){
            size_t size = (i+1)*64;
            int smallbin_index = (size >> 4) - 1;

            void* smallbin_head = ar->smallbins[smallbin_index]->fd;
            void* original_ptr = (char*)ptrs[i] - HEADER_SIZE;
            
            TEST_ASSERT_EQUAL_PTR(smallbin_head, original_ptr);
        }
    }


    /* consume chunk from smallbin */
    free_chunk_t* x= allocate(ar, 192);
    TEST_ASSERT_EQUAL_PTR(x,  ptrs[2]); /* ptrs[0] is in fastbin */
}


/* test unsortedbin */
void test_unsortedbin_reuse(void) {
    int n = 10;
    void* ptrs[n];
    /* [64][128]....[640][top_chunk] */
    for(int i=0; i<n; i++) {
        ptrs[i] = allocate(ar, (i+1)*64);
    }

    /* clear few chunks in left to right order */
    /* [64][........][][] .... [top_chunk]*/
    size_t expected_coalesced_sz = 0;
    for(int i=0; i<n/2; i++) {
        release(ar, ptrs[i]);
        expected_coalesced_sz +=  (i+1)*64 + HEADER_SIZE;
    }
    /* [64][......][][]...[top_chunk] */
    /* ptrs[0](64byte) fastbin not coalesced */
    expected_coalesced_sz -= (HEADER_SIZE + 64 + HEADER_SIZE);

    /* unsortedbin should contain single big chunk */
    TEST_ASSERT_NOT_NULL(ar->unsortedbin->fd);
    TEST_ASSERT_EQUAL_PTR(ar->unsortedbin->fd->fd, ar->unsortedbin);
    TEST_ASSERT_TRUE((ar->unsortedbin->fd->size & __CHUNK_SIZE_MASK ) == expected_coalesced_sz);

    /* allocate a big chunk */
    void* big = allocate(ar, expected_coalesced_sz);
    TEST_ASSERT_EQUAL_PTR(big, ptrs[1]);

    /* allocate a small chunk */
    void* small = allocate(ar, 64);
    TEST_ASSERT_EQUAL_PTR(small, ptrs[0]);
}

/* test largebins */
void test_largebin_reuse(void) {
    size_t sizes[] = {1024, 2048, 4096, 8192};
    int n = 4;
    
    void* ptrs[n];
    void* barriers[n];

    /* Layout: [1024][barrier][2048][barrier][4096][barrier][8192][barrier]... */
    for(int i=0; i<n; i++){
        ptrs[i] = allocate(ar, sizes[i]);
        /* Allocate a small chunk in between to prevent coalescing when we free them */
        barriers[i] = allocate(ar, 64); 
    }

    /* We free them in reverse order, or random order, it shouldn't matter */
    /* They enter the Unsorted Bin. */
    for(int i=0; i<n; i++) {
        release(ar, ptrs[i]);
    }

    /* any allocation should trigger unsortedbin to largebin flow*/
    /* using large size greater than all chunks available to avoid spliting */
    allocate(ar, 10000); 

    
    /* exact Fit */
    /* Request 4096. Should get exactly ptrs[2]. */
    void* fit_exact = allocate(ar, 4096);
    TEST_ASSERT_EQUAL_PTR(fit_exact, ptrs[2]);

    /* best fit (splitting) */
    /* request 1500. */
    /* 1024 (ptrs[0]) is too small. */
    /* 2048 (ptrs[1]) is the best fit. */
    /* 8192 (ptrs[3]) is valid but wasteful (allocator should skip it). */
    void* fit_split = allocate(ar, 1500);
    
    TEST_ASSERT_EQUAL_PTR(fit_split, ptrs[1]);

    /* remaining large chunk */
    /* Request 8000. Should get ptrs[3] (8192). */
    void* fit_large = allocate(ar, 8000);
    TEST_ASSERT_EQUAL_PTR(fit_large, ptrs[3]);
}
void test_coalesce_forward(void) {
    void* p1 = allocate(ar, 128);
    void* p2 = allocate(ar, 128); 
    void* barrier = allocate(ar, 16); 

    release(ar, p2); // Free higher address first
    release(ar, p1); // Free lower address, should merge with p2

    // Request size equal to p1 + p2 + header_overhead
    // If not coalesced, this alloc would fail or take from top chunk
    size_t combined_size = 128 + 128 + HEADER_SIZE; 
    void* big = allocate(ar, combined_size);

    TEST_ASSERT_EQUAL_PTR_MESSAGE(p1, big, "Chunks did not coalesce forward");
}

void test_coalesce_backward(void) {
    void* p1 = allocate(ar, 128);
    void* p2 = allocate(ar, 128);
    void* barrier = allocate(ar, 16);

    release(ar, p1); // Free lower address first
    release(ar, p2); // Free higher address, should merge back into p1

    size_t combined_size = 128 + 128 + HEADER_SIZE;
    void* big = allocate(ar, combined_size);

    TEST_ASSERT_EQUAL_PTR_MESSAGE(p1, big, "Chunks did not coalesce backward");
}

void test_coalesce_sandwich(void) {
    void* p1 = allocate(ar, 4096);
    void* p2 = allocate(ar, 4096);
    void* p3 = allocate(ar, 4096);
    void* barrier = allocate(ar, 16);

    release(ar, p1);
    release(ar, p3);
    release(ar, p2); 

    size_t combined_size = 4096 * 3 + HEADER_SIZE * 2;
    void* big = allocate(ar, combined_size);

    TEST_ASSERT_EQUAL_PTR_MESSAGE(p1, big, "Sandwich coalesce failed");
}

void test_mmap(void) {
    /* request size above MMAP_THRESHOLD */
    void* p = allocate(ar, MMAP_THRESHOLD);
    TEST_ASSERT_NOT_NULL(p);

    free_chunk_t* chunk = (free_chunk_t*)((char*)p - HEADER_SIZE);

    TEST_ASSERT_TRUE(chunk->size & __MMAP_ALLOCATED_FLAG_MASK);
}   

void test_splitting_large_chunk(void) {
    void* big = allocate(ar, 512);

    release(ar, big);

    void* small = allocate(ar, 128);
    TEST_ASSERT_EQUAL_PTR(big, small);

    void* remainder = allocate(ar, 128);
    uintptr_t expected_addr = (uintptr_t)small + 128 + HEADER_SIZE;
    
    TEST_ASSERT_EQUAL_PTR((void*)expected_addr, remainder);
}

void test_no_overlap_random_allocs(void) {
    #define NUM_PTRS 50
    void* ptrs[NUM_PTRS];
    size_t sizes[NUM_PTRS];

    // 1. Allocate varying sizes and fill with unique patterns
    for (int i = 0; i < NUM_PTRS; i++) {
        sizes[i] = (i + 1) * 8;
        ptrs[i] = allocate(ar, sizes[i]);
        TEST_ASSERT_NOT_NULL(ptrs[i]);
        fill_pattern(ptrs[i], sizes[i], (uint8_t)i);
    }

    // 2. Verify patterns (ensure alloc N didn't overwrite alloc N-1)
    for (int i = 0; i < NUM_PTRS; i++) {
        TEST_ASSERT_TRUE_MESSAGE(check_pattern(ptrs[i], sizes[i], (uint8_t)i), 
                                 "Memory Corruption: Overlap detected");
    }

    // 3. Release all
    for (int i = 0; i < NUM_PTRS; i++) {
        release(ar, ptrs[i]);
    }
}

#define N_THREADS        8
#define ITERS_PER_THREAD 10000
#define MAX_LIVE_ALLOCS  256
#define MAX_ALLOC_SIZE  4096
#define CANARY          0xAB

static arena_t *global_arena;

/* simple xorshift rng (thread-local) */
static inline uint32_t rng(uint32_t *state) {
    uint32_t x = *state;
    x ^= x << 13;
    x ^= x >> 17;
    x ^= x << 5;
    *state = x;
    return x;
}

typedef struct {
    void  *ptr;
    size_t size;
} alloc_rec_t;

void* worker(void *arg) {
    uint32_t seed = (uintptr_t)arg ^ (uintptr_t)&seed;

    alloc_rec_t live[MAX_LIVE_ALLOCS] = {0};
    size_t live_count = 0;

    for (int i = 0; i < ITERS_PER_THREAD; i++) {
        uint32_t r = rng(&seed);

        /* 60% alloc, 40% free */
        if ((r % 10) < 6 && live_count < MAX_LIVE_ALLOCS) {
            size_t size = (rng(&seed) % MAX_ALLOC_SIZE) + 1;

            void *p = allocate(global_arena, size);
            TEST_ASSERT_NOT_NULL(p);
            TEST_ASSERT_TRUE_MESSAGE(is_aligned(p),
                                     "Pointer not 16-byte aligned");

            /* fill memory */
            memset(p, CANARY, size);

            /* store */
            live[live_count++] = (alloc_rec_t){ p, size };
        } else if (live_count > 0) {
            /* free a random live allocation */
            size_t idx = rng(&seed) % live_count;
            alloc_rec_t rec = live[idx];

            /* verify canary before free */
            unsigned char *c = rec.ptr;
            for (size_t j = 0; j < rec.size; j++) {
                TEST_ASSERT_EQUAL_HEX8_MESSAGE(
                    CANARY, c[j], "Memory corruption detected");
            }

            release(global_arena, rec.ptr);

            /* remove from table (swap-delete) */
            live[idx] = live[--live_count];
        }
    }

    /* cleanup remaining allocations */
    for (size_t i = 0; i < live_count; i++) {
        release(global_arena, live[i].ptr);
    }

    return NULL;
}

void test_multithread_shared_arena(void) {
    global_arena = arena_create();
    TEST_ASSERT_NOT_NULL(global_arena);

    pthread_t threads[N_THREADS];

    for (int i = 0; i < N_THREADS; i++) {
        pthread_create(&threads[i], NULL, worker, (void *)(uintptr_t)i);
    }

    for (int i = 0; i < N_THREADS; i++) {
        pthread_join(threads[i], NULL);
    }
}

int main(void) {
    UNITY_BEGIN();

    RUN_TEST(test_basic_allocation_and_alignment);
    RUN_TEST(test_zero_allocation);
    RUN_TEST(test_free_null_is_safe);
    
    RUN_TEST(test_fastbin_reuse);
    RUN_TEST(test_unsortedbin_reuse);
    RUN_TEST(test_smallbin_reuse);
    RUN_TEST(test_largebin_reuse);
    
    RUN_TEST(test_coalesce_forward);
    RUN_TEST(test_coalesce_backward);
    RUN_TEST(test_coalesce_sandwich);
    
    RUN_TEST(test_splitting_large_chunk);
    RUN_TEST(test_no_overlap_random_allocs);

    RUN_TEST(test_mmap);

    RUN_TEST(test_multithread_shared_arena);

    return UNITY_END();
}