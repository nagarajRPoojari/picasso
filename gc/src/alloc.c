#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/mman.h>
#include <assert.h>
#include <string.h>
#include "alloc.h"


/* utils */
static inline size_t align16(size_t size) {
    return (size + 15) & ~(size_t)15;
}

static inline size_t align_page(size_t size) {
    static size_t page_size = 0;
    if (!page_size) page_size = (size_t)sysconf(_SC_PAGESIZE);
    return (size + page_size - 1) & ~(page_size - 1);
}

static inline size_t chunk_size(size_t size) {
    return (size & __CHUNK_SIZE_MASK);
}

static inline size_t get_size(free_chunk_t* fc) {
    return chunk_size(fc->size);
}

static inline int get_smallbin_index(size_t size) {
    return (size >> 4) - 1; 
}

static inline void set_curr_inuse(free_chunk_t* fs){
    fs->size |= __CURR_IN_USE_FLAG_MASK;
}

static inline void unset_curr_inuse(free_chunk_t* fs){
    fs->size &= ~__CURR_IN_USE_FLAG_MASK;
}

static int get_largebin_index(size_t size) {
    // Sanity check: If size is too small, it belongs to small/fast bins
    if(size <= SMALLBIN_MAX_SIZE) return -1;
    
    /*  range 0-31: 64-byte steps (from 512 B to 64 KB) */
    if (size <= (64 * 1024)) { 
        /* 
            the first 32 large bins cover 512 bytes up to 64KB. 
            e.g, 
                512 - 575
                576 - 639
                ...
        */
        int index = (int)(((size - 1) >> 6) - 8); /* shift by 8 indices */ 
        
        return (index > 31) ? 31 : index; 
    }

    /*
        range 32-39: 256-byte steps (64 KB to 256 KB)
        total 24 bins (8 bins * 3 groups)        
    */
    if (size <= (256 * 1024)) { 
        /** index 32 is the first bin in this group. Step size is 2^8 (256). */
        return (int)(((size - (64 * 1024 + 1)) >> 8) + 32);
    }

    /* range 40-47: 1 KB steps (256 KB to 1 MB) */
    if (size <= (1024 * 1024)) {
        /* index 40 is the first bin in this group. Step size is 2^10 (1024). */
        return (int)(((size - (256 * 1024 + 1)) >> 10) + 40);
    }

    /* range 48-55: 4 KB steps (1 MB to 4 MB) */
    if (size <= (4 * 1024 * 1024)) { 
        /* index 48 is the first bin in this group. Step size is 2^12 (4096). */
        return (int)(((size - (1024 * 1024 + 1)) >> 12) + 48);
    }

    /*
        range 56-63: Exponential growth (for sizes > 4 MB) 
        these last 8 bins cover vast size ranges, typically doubling the range with each bin.
        we start indexing from 56 (4MB) and rely on the logarithm of the size.
        the following uses a simplified fixed-step approach for clarity:
    */
    
    if (size <= (8 * 1024 * 1024)) return 56;
    if (size <= (16 * 1024 * 1024)) return 57;
    if (size <= (32 * 1024 * 1024)) return 58;
    if (size <= (64 * 1024 * 1024)) return 59;
    if (size <= (128 * 1024 * 1024)) return 60;
    if (size <= (256 * 1024 * 1024)) return 61;
    if (size <= (512 * 1024 * 1024)) return 62;

    /* max bin 63 */
    /* all sizes greater than 512 MB go into the final bin. */
    return LARGEBIN_MAX_INDEX;
}


static void unlink_chunk(free_chunk_t* p) {
    if (p->bk && p->fd) {
        p->fd->bk = p->bk;
        p->bk->fd = p->fd;
        p->fd = p->bk = NULL;
    }
}

/* for inserting in sentinal node case */
static void insert_chunk_head(free_chunk_t* head, free_chunk_t* p) {
    p->fd = head->fd;
    p->bk = head;
    head->fd->bk = p;
    head->fd = p;
}


static void* alloc_by_mmap(size_t size) {
    void* p = mmap(NULL, size, PROT_READ | PROT_WRITE, MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    if(p == MAP_FAILED) return NULL;
    return p; 
}

static free_chunk_t* request_chunk_by_mmap(size_t size){
    size_t total_size = align_page(size + HEADER_SIZE);
    void* p = alloc_by_mmap(total_size);
    if (!p) return NULL;

    free_chunk_t* fc = (free_chunk_t*) p;
    fc->size = total_size - HEADER_SIZE;
    fc->size |= __MMAP_ALLOCATED_FLAG_MASK;
    fc->fd = fc->bk = NULL;

    Debug("Mapped new chunk: %zu bytes\n", fc->size & __CHUNK_SIZE_MASK);
    return fc;
}


/* split_top_chunk: Carve a piece from the main arena top_chunk */
static free_chunk_t* split_top_chunk(arena_t* ar, size_t payload_size) {
    if(!ar->top_chunk) return NULL;

    const size_t required_total = payload_size + HEADER_SIZE;
    size_t available_total = get_size(ar->top_chunk) + HEADER_SIZE;
    free_chunk_t* victim = ar->top_chunk;
    
    if (available_total >= required_total) {
        size_t remaining_total = available_total - required_total;

        victim->size = (payload_size | (victim->size & __PREV_IN_USE_FLAG_MASK));
        set_curr_inuse(victim);

        if (remaining_total == 0) {
            ar->top_chunk = NULL;
            return victim;
        }
        
        if (remaining_total < MIN_PAYLOAD_SIZE + HEADER_SIZE) {
            size_t new_payload_size = payload_size + remaining_total;
            victim->size = (new_payload_size | (victim->size & __PREV_IN_USE_FLAG_MASK));
            
            ar->top_chunk = NULL;
            return victim;
        }

        free_chunk_t* new_top = (free_chunk_t*)((char*)victim + required_total);
        size_t new_top_payload_size = remaining_total - HEADER_SIZE;
        new_top->size = new_top_payload_size | __PREV_IN_USE_FLAG_MASK;
        
        ar->top_chunk = new_top;
        Debug("Split top chunk. New top payload size: %zu\n", new_top_payload_size);
        return victim;
    }
    
    return NULL;
}

/* find_in_fastbins: Check LIFO single-linked lists for exact size match */
static free_chunk_t* find_in_fastbins(arena_t* ar, size_t payload_size) {

    size_t chunk_sz = payload_size;
    int idx = (chunk_sz >> 4) - 1; // Assuming min chunk size 32, or payload of 16

    if (idx >= 0 && idx < FASTBINS_COUNT && ar->fastbins[idx] != NULL) {
        free_chunk_t* victim = ar->fastbins[idx];
        ar->fastbins[idx] = victim->fd; // Pop head
        victim->fd = NULL;
        Debug("Found in Fastbin[%d]\n", idx);

        set_curr_inuse(victim);
        return victim;
    }
    return NULL;
}

/* insert_into_smallbin: helper to move chunks from unsorted to small */
static void insert_into_smallbin(arena_t* ar, free_chunk_t* chunk) {
    size_t sz = get_size(chunk);
    int idx = get_smallbin_index(sz);
    if (idx < SMALLBINS_COUNT) {
        insert_chunk_head(ar->smallbins[idx], chunk);
        ar->smallbinmap |= (1U << idx);
    }
}

/* find_in_smallbins: Check FIFO double-linked lists for exact size match */
static free_chunk_t* find_in_smallbins(arena_t* ar, size_t payload_size) {
    int idx = get_smallbin_index(payload_size);

    if (idx < SMALLBINS_COUNT && (ar->smallbinmap & 1<<idx)) {
        free_chunk_t* head = ar->smallbins[idx];
        
        // If head->fd == head, it's empty
        if (head->fd != head) {
            free_chunk_t* victim = head->fd; // FIFO: take the first one
            unlink_chunk(victim);
            
            if (head->fd == head) {
                ar->smallbinmap &= ~(1U << idx);
            }
            Debug("Found in Smallbin[%d]\n", idx);
            set_curr_inuse(victim);
            return victim;
        }
    }
    return NULL;
}

/* insert_into_largebin: helper to move chunks from unsortedbin to largebin */
static void insert_into_largebin(arena_t* ar, free_chunk_t* chunk) {
    size_t sz = get_size(chunk);
    int idx = get_largebin_index(sz);
    
    assert(idx < LARGEBINS_COUNT);
    free_chunk_t* head = ar->largebins[idx];
    free_chunk_t* next_size;
    free_chunk_t* prev_size;
    
    next_size = head->next_sizeptr; 

    Debug("head = %p next-size = %p \n", head, next_size);
    
    if (!next_size) {
        next_size = head;
    }
    
    while (next_size != head && get_size(next_size) < sz) {
        next_size = next_size->next_sizeptr;
    }
    
    prev_size = next_size->prev_sizeptr;
    
    chunk->next_sizeptr = next_size;
    chunk->prev_sizeptr = prev_size;
    
    
    next_size->prev_sizeptr = chunk;
    prev_size->next_sizeptr = chunk;
    
    insert_chunk_head(head, chunk);
    
    ar->largebinmap |= (1U << idx);
    
}

// Assuming payload_size is already aligned (e.g., align16(requested_size))
static free_chunk_t* find_in_largebin(arena_t* ar, size_t payload_size) {
    // Required size includes the header
    size_t required_size = payload_size;
    int idx = get_largebin_index(required_size);
    if(idx < 0) return NULL;

    /* start search from the calculated bin and check all subsequent bins (best fit) */
    for (int current_idx = idx; current_idx < LARGEBINS_COUNT; current_idx++) {

        /* check empty? */
        if (!(ar->largebinmap & (1U << current_idx))) {
            continue; 
        }

        free_chunk_t* head = ar->largebins[current_idx];
        free_chunk_t* victim = NULL;
        
        /*
        search the size-sorted list (next_sizeptr) for the best fit
        find the smallest chunk (ceil/victim) that is large enough
        */
        for (free_chunk_t* ceil = head->next_sizeptr; ceil != head; ceil = ceil->next_sizeptr) {
            size_t ceil_size = get_size(ceil);
            
            if (ceil_size >= required_size) {
                victim = ceil;
                break;
            }
        }
        
        if (victim) {
            /* unlink from fd/bk list */
            unlink_chunk(victim); 

            /* unlink from size-sorted list (next_sizeptr/prev_sizeptr) */
            victim->next_sizeptr->prev_sizeptr = victim->prev_sizeptr;
            victim->prev_sizeptr->next_sizeptr = victim->next_sizeptr;
            victim->next_sizeptr = victim->prev_sizeptr = NULL;

            /* update bin map if the list is now empty */
            if (head->fd == head) {
                ar->largebinmap &= ~(1U << current_idx);
            }
            
            /* split */

            size_t victim_size = get_size(victim);
            size_t remainder_size = victim_size - required_size;

            if (remainder_size >= MIN_PAYLOAD_SIZE) {
                /* split the chunk and put the remainder back */
                
                /* prepare victim: Set size to requested size*/
                victim->size = required_size | (victim->size & __SIZE_BITS);

                /* setup remainder chunk */
                free_chunk_t* remainder = (free_chunk_t*)((char*)victim + required_size);
                remainder->size = remainder_size;
                remainder->size |= __PREV_IN_USE_FLAG_MASK; 

                /* link the remainder into the unsorted bin (to be sorted/merged later) */
                insert_chunk_head(ar->unsortedbin, remainder);
            }
            /* else: remainder is too small, victim takes the full original size. (internal fragmentation) */

            Debug("Found and returned chunk from largebin index %d\n", current_idx);
            Debug("== > victim %p \n", victim);

            set_curr_inuse(victim);
            return victim;
        }
    }

    return NULL;
}

/* find_in_unsortedbin: Iterate unsorted bin. Return fit if found, else sort into bins. */
static free_chunk_t* find_in_unsortedbin(arena_t* ar, size_t payload_size) {
    free_chunk_t* curr = ar->unsortedbin->fd;
    free_chunk_t* victim = NULL;
    free_chunk_t* remainder = NULL;

    
    while (curr != ar->unsortedbin) {
        free_chunk_t* next = curr->fd; // Save next
        size_t size = get_size(curr);
        size_t required = payload_size;

        curr->fd = curr->bk = NULL;
        if (ar->unsortedbin->fd == curr) ar->unsortedbin->fd = next;

        if (!victim && size == required) {
            victim = curr;
        } 
        else if (!victim && size > required + MIN_PAYLOAD_SIZE) {
            // Split fit
            size_t remainder_size = size - required;
            
            // Setup victim
            victim = curr;
            victim->size = payload_size | (victim->size & __SIZE_BITS);

            // Setup remainder
            remainder = (free_chunk_t*)((char*)victim + required + HEADER_SIZE);
            remainder->size = remainder_size | __PREV_IN_USE_FLAG_MASK;

        } 
        else {
            /* no fit, move to appropriate bin */
            if (size <= 16 * SMALLBINS_COUNT) {
                Debug("inserting %zu into smallbin \n", curr->size & __CHUNK_SIZE_MASK);
                insert_into_smallbin(ar, curr);
            } else {
                /* push it to largebins */
                Debug("inserting %zu into largebin: %p \n", curr->size & __CHUNK_SIZE_MASK, curr);
                insert_into_largebin(ar, curr);
            }
        }
        curr = next;
    }

    if(victim) Debug("Found in unsorted bin: %zu, \n", victim->size & __CHUNK_SIZE_MASK);
    if(remainder) {        
        insert_chunk_head(ar->unsortedbin, remainder);
    }

    if(victim) set_curr_inuse(victim);
    return victim;
}

static void grow_heap(arena_t* ar, size_t size_needed) {

    size_t request = (size_needed > HEAP_MIN_SIZE) ? size_needed : HEAP_MIN_SIZE;
    free_chunk_t* new_block = request_chunk_by_mmap(request + HEAP_BOUNDARY_SIZE);
    
    if (new_block) {
        free_chunk_t* boundary = (free_chunk_t*)((char*)new_block + request);
        boundary->size = 0;
        boundary->fd = boundary->bk = NULL; /* important */

        ar->top_chunk = new_block;
        ar->top_chunk->size = request | __PREV_IN_USE_FLAG_MASK;

        /* need to unset mmap flag */
        ar->top_chunk->size &= ~__MMAP_ALLOCATED_FLAG_MASK;

    }

}

static free_chunk_t* coalesce(arena_t* ar, free_chunk_t* chunk){
    size_t chunk_size = get_size(chunk);

    if (!(chunk->size & __PREV_IN_USE_FLAG_MASK)) {
        free_chunk_t* prev_chunk = (free_chunk_t*)((char*)chunk - chunk->prev_size - HEADER_SIZE);
        // unlink_chunk(prev_chunk);
        size_t prev_size = get_size(prev_chunk);
        size_t merged_size = prev_size + HEADER_SIZE + chunk_size;
        prev_chunk->size = (prev_chunk->size & __SIZE_BITS) | merged_size;
        chunk = prev_chunk;
        chunk_size = merged_size;
    }

    free_chunk_t* next_chunk = (free_chunk_t*)((char*)chunk + HEADER_SIZE + chunk_size);

    if (next_chunk == ar->top_chunk) {
        size_t top_size = get_size(next_chunk);
        size_t merged_size = chunk_size + HEADER_SIZE + top_size;
        chunk->size = (chunk->size & __SIZE_BITS) | merged_size;
        ar->top_chunk = chunk;
        return NULL;
    }

    if (!(next_chunk->size & __CURR_IN_USE_FLAG_MASK)) {
        unlink_chunk(next_chunk);
        size_t next_size = get_size(next_chunk);
        size_t merged_size = chunk_size + HEADER_SIZE + next_size;
        chunk->size = (chunk->size & __SIZE_BITS) | merged_size;
        chunk_size = merged_size;
    }

    free_chunk_t* next_next_chunk = (free_chunk_t*)((char*)chunk + HEADER_SIZE + chunk_size);
    next_next_chunk->prev_size = chunk_size;
    next_next_chunk->size &= ~__PREV_IN_USE_FLAG_MASK;

    return chunk;
}


static free_chunk_t* fastbin_consolidate(arena_t* ar) {

}

static void* allocate_unsafe(arena_t* ar, size_t requested_size) {
    if(!requested_size) return NULL;


    size_t payload_size = align16(requested_size);
    if (payload_size < MIN_PAYLOAD_SIZE) payload_size = MIN_PAYLOAD_SIZE;
    
    size_t total_size = payload_size + HEADER_SIZE;
    free_chunk_t* fc = NULL;

    /* scan fastbins */
    if (payload_size <= 16 * FASTBINS_COUNT) { 
        fc = find_in_fastbins(ar, payload_size);
        if (fc) return (void*)((char*)fc + HEADER_SIZE);
    }
    
    /* scan smallbins */
    if (payload_size < 16 * SMALLBINS_COUNT) {
        fc = find_in_smallbins(ar, payload_size);
        if (fc) return (void*)((char*)fc + HEADER_SIZE);
    }
    
    /* scan unsorted bins */

    fc = find_in_unsortedbin(ar, payload_size);
    if (fc) return (void*)((char*)fc + HEADER_SIZE);

    
    /* scan largebins */
    fc = find_in_largebin(ar, payload_size);
    if (fc) return (void*)((char*)fc + HEADER_SIZE);

    /* allocate with mmap for huge chunks */
    if (payload_size >= MMAP_THRESHOLD) {
        fc = request_chunk_by_mmap(payload_size);
        set_curr_inuse(fc);
        return (void*)((char*)fc + HEADER_SIZE);
    }

    /* carve from top_chunk */
    fc = split_top_chunk(ar, payload_size);
    if (fc) return (void*)((char*)fc + HEADER_SIZE);
    /* if top_chunk is not NULL push it to unsortedbin */
    // if(ar->top_chunk) {
    //     printf("7.ar->unsortedbin->fd->fd: %p \n", ar->unsortedbin->fd->fd);
    //     insert_chunk_head(ar->unsortedbin, ar->top_chunk);
    //     printf("6.ar->unsortedbin->fd->fd: %p \n", ar->unsortedbin->fd->fd);
    //     ar->top_chunk = NULL;
    // }

    /* grow heap & try to carve again */
    grow_heap(ar, total_size);
    fc = split_top_chunk(ar, payload_size);

    return (fc) ? (void*)((char*)fc + HEADER_SIZE) : NULL;
}

void* allocate(arena_t* ar, size_t requested_size) {
    assert(ar != NULL);
    pthread_rwlock_wrlock(&ar->mu);
    void* res = allocate_unsafe(ar, requested_size);
    pthread_rwlock_unlock(&ar->mu);
    return res;
}

static void release_unsafe(arena_t* ar, void* ptr) {
    if (!ptr) return;
    
    free_chunk_t* fc = (free_chunk_t*)((char*)ptr - HEADER_SIZE);
    if(fc->size & __CURR_IN_USE_FLAG_MASK) {
        unset_curr_inuse(fc);
    }else {
        return;
    }


    /* check whether allocation is from mmap, then unmap immediately*/
    if(fc->size & __MMAP_ALLOCATED_FLAG_MASK) {
        munmap(fc, get_size(fc) + HEADER_SIZE);

        return;
    }

    size_t size = get_size(fc);
    
    /* add to fastbin */
    if (size <= 16 * FASTBINS_COUNT) { 
        int idx = (size >> 4) - 1;
        if (idx >= 0 && idx < FASTBINS_COUNT) {
            Debug("Releasing chunk size: %zu to fastbin \n", size);
            fc->fd = ar->fastbins[idx];
            ar->fastbins[idx] = fc;
            return;
        }
    }
    
    /* check whether we can coalesce */
    free_chunk_t* merged_chunk = coalesce(ar, fc);
    /* merged_chunk = NULL indicates, fc has been merged with top_chunk */
    if(!merged_chunk) return;
    
    size = get_size(merged_chunk);
    Debug("Releasing chunk size (after coalescing): %zu. addr= %p \n", size, merged_chunk);

    /* set next physical chunk prev_in_use bit to no*/
    free_chunk_t* next_chunk = (free_chunk_t*)((char*)merged_chunk + size + HEADER_SIZE);
    next_chunk->prev_size = size;
    next_chunk->size &= ~__PREV_IN_USE_FLAG_MASK;

    /* insert to unsortedbins only if it is not coalesced i.e, independent chunk */
    if(fc == merged_chunk) {
        insert_chunk_head(ar->unsortedbin, fc);
    }
}

void release(arena_t* ar, void* ptr) {
    assert(ar != NULL);
    pthread_rwlock_wrlock(&ar->mu);
    release_unsafe(ar, ptr);
    pthread_rwlock_unlock(&ar->mu);
}

arena_t* arena_create(void) {

    size_t needed = sizeof(arena_t) +
                    (SMALLBINS_COUNT * sizeof(free_chunk_t)) + 
                    (LARGEBINS_COUNT * sizeof(free_chunk_t)) + sizeof(free_chunk_t);
    
    void* block = alloc_by_mmap(align_page(needed));
    if (!block) return NULL;

    arena_t* a = (arena_t*)block;
    char* curs = (char*)block + sizeof(arena_t);


    /* initialize mutext */
    pthread_rwlock_init(&a->mu, NULL);

    /* initialize fastbins */
    for (size_t i = 0; i < FASTBINS_COUNT; i++) {
        a->fastbins[i] = NULL;
    }

    /* initialize smallbins */
    for (size_t i = 0; i < SMALLBINS_COUNT; i++) {
        free_chunk_t* sentinel = (free_chunk_t*)curs;
        sentinel->fd = sentinel;
        sentinel->bk = sentinel;
        sentinel->size = 0; 
        a->smallbins[i] = sentinel;
        curs += sizeof(free_chunk_t);
    }
    
    /* initialize largebins */
    for (size_t i = 0; i < LARGEBINS_COUNT; i++) {
        free_chunk_t* sentinel = (free_chunk_t*)curs;
        sentinel->fd = sentinel;
        sentinel->bk = sentinel;
        sentinel->size = 0; 
        sentinel->next_sizeptr = sentinel->prev_sizeptr = sentinel;
        a->largebins[i] = sentinel;
        curs += sizeof(free_chunk_t);
    }

    /* initialize top chunk */
    a->top_chunk = request_chunk_by_mmap(HEAP_MIN_SIZE + HEAP_BOUNDARY_SIZE);
    a->top_chunk->size = HEAP_MIN_SIZE | __PREV_IN_USE_FLAG_MASK;
    /* unset mmap flag */
    a->top_chunk->size &= ~__MMAP_ALLOCATED_FLAG_MASK;

    free_chunk_t* sentinel = (free_chunk_t*)curs;
    sentinel->fd = sentinel;
    sentinel->bk = sentinel;
    sentinel->size = 0; 
    a->unsortedbin = sentinel;

    a->smallbinmap = 0;
    a->largebinmap = 0;

    return a;
}

