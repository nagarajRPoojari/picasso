#include <stdio.h>
#include <stdatomic.h>
#include <unistd.h>
#include <string.h>
#include <stddef.h>
#include "alloc.h"
#include "gc.h"

/* expected to be initialized in gc_init */
gc_state_t gc_state;

/* arenas to keep track of. Should be initialized by arena_create */
arena_t* arenas[MAX_ARENAS];
int arenas_count;

/* keep track of kernel/scheduler threads */
kernel_thread_t* kts[MAX_SCHEDULERS];
int kts_count;

/* forward declarations */
static void gc_mark_mem_region(char *start, char *end);

arena_t* gc_create_arena(kernel_thread_t* kt) {
    /* register thread*/
    kts[kts_count++] = kt;

    if(arenas_count + 1 > MAX_ARENAS) {
        perror("failed to create new arena: arena count reached max\n");
        return NULL;
    }

    arena_t* ar = arena_create();

    /* register arena */
    arenas[arenas_count++] = ar;
    return ar;
}

static void mark_chunk_recursive(inuse_chunk_t* ch) {
    if(ch->prev_size & __GC_MARK_FLAG_MASK) return;

    printf("marking chunk of size = %zu \n", ch->size & __CHUNK_SIZE_MASK);

    char* payload_start = (char*)ch + HEADER_SIZE;
    char* payload_end = payload_start + (ch->size & __CHUNK_SIZE_MASK);
    gc_mark_mem_region(payload_start, payload_end);
}

static inuse_chunk_t* find_chunk_in_heap(arena_t* ar, alloced_heap_t* heap, char* ptr){
    char* scan = heap->start;

    while (scan < heap->end) {
        inuse_chunk_t *chunk = (inuse_chunk_t*)scan;
        size_t payload_size = chunk->size & __CHUNK_SIZE_MASK;

        if( !(chunk->size & __CURR_IN_USE_FLAG_MASK) ) {
            scan += HEADER_SIZE + payload_size;
            continue;
        }

        char* payload_start = scan + HEADER_SIZE;
        char* payload_end   = payload_start + payload_size;

        if (ptr >= payload_start && ptr < payload_end)
            return chunk;

        scan += HEADER_SIZE + payload_size;
    }
    return NULL;
}


static inline int is_pointer_aligned(uintptr_t v) {
    return (v & GC_ALIGN_MASK) == 0;
}

/* Try to mark a single candidate pointer value.
 * Returns 1 if it marked something, 0 otherwise.
 */
static int try_mark_pointer(uintptr_t val) {
    if (!is_pointer_aligned(val)) return 0;

    char *p = (char*)val;

    for (int j = 0; j < arenas_count; ++j) {
        arena_t *ar = arenas[j];
        for (int i = 0; i < ar->alloced_heap_count; ++i) {
            char *hs = ar->alloced_heaps[i].start;
            char *he = ar->alloced_heaps[i].end;

            if ((char*)p < hs || (char*)p >= he) continue;

            inuse_chunk_t *ch = find_chunk_in_heap(ar, &ar->alloced_heaps[i], p);
            if (!ch) return 0;

            ch->prev_size |= __GC_MARK_FLAG_MASK;
            mark_chunk_recursive(ch);
            return 1;
        }
    }
    return 0;
}

/* Scan a memory region (stack or other) - read words safely and test each as pointer */
static void gc_mark_mem_region(char *start, char *end) {
    // printf("gc_mark_mem_region - [%p ... %p]\n", (char*)start, (void*)end);
    
    // read as uintptr_t; use memcpy for safe unaligned reads on strict platforms
    for (char *p = start; p + (ptrdiff_t)sizeof(uintptr_t) <= end; p += sizeof(uintptr_t)) {
        uintptr_t val;
        memcpy(&val, p, sizeof(val));   // safe on all architectures

        if (val == 0) continue;
        try_mark_pointer(val);
    }
}

/* Scan saved registers from ucontext (AArch64) */
static void gc_mark_registers(ucontext_t *ctx) {
    // X0-X30
    for (int i = 0; i < 31; ++i) {
        uintptr_t val = (uintptr_t)ctx->uc_mcontext.regs[i];
        if (val == 0) continue;
        try_mark_pointer(val);
    }

    // SP, PC
    try_mark_pointer((uintptr_t)ctx->uc_mcontext.sp);
    try_mark_pointer((uintptr_t)ctx->uc_mcontext.pc);
}

static void gc_mark_task(ucontext_t *ctx) {
    // stack
    char *stack_bottom = (char*)ctx->uc_stack.ss_sp;             // lowest address
    char *stack_top    = stack_bottom + ctx->uc_stack.ss_size;   // highest address

    char *sp = (char*)ctx->uc_mcontext.sp;                       // current SP

    if (sp < stack_bottom || sp > stack_top) {
        perror(" wrong stack \n");
        return;
    }

    gc_mark_mem_region(stack_bottom, stack_top);

    // registers
    gc_mark_registers(ctx);
}

static void gc_mark() {

    for(int kti = 0; kti < kts_count; kti++ ){

        kernel_thread_t* kt = kts[kti];
        if (!kt->current) {
            continue; 
        }

        gc_mark_task(&kt->current->ctx);

        /* It took 9hr to debug the issue. previously scanning only RUNNING task stack, 
        need to scan all task stack. */
        task_node_t* head = kt->ready_q.head;
        while(head) {
            task_t* t = head->t;
            gc_mark_task(&t->ctx);
            head = head->next;
        }
    }
    
}



static void gc_sweep() {

    /*  
        go through all arenas, all heaps 
        if gc not marked, call release()
    */

    for(int i=0; i<arenas_count; i++) {
        arena_t* ar = arenas[i];

        for(int j=0; j<ar->alloced_heap_count; j++) {
            alloced_heap_t* heap = &(ar->alloced_heaps[j]);
            
            char* scan = heap->start;

            while (scan < heap->end) {

                inuse_chunk_t *chunk = (inuse_chunk_t*)scan;
                size_t payload_size = chunk->size & __CHUNK_SIZE_MASK;

                if( chunk->prev_size & __GC_MARK_FLAG_MASK || !(chunk->size & __CURR_IN_USE_FLAG_MASK)) {
                    chunk->prev_size & ~__GC_MARK_FLAG_MASK;
                } else {
                    // printf("[gc] releasing %p [size=%zu] \n", (char*)chunk + HEADER_SIZE, chunk->size & __CHUNK_SIZE_MASK);
                    release(ar, (char*)chunk + HEADER_SIZE);
                }

                scan += HEADER_SIZE + payload_size;
            }
        }
    }
}

void gc_run() {
    printf(" started gc \n");
    while (1) {
        gc_stop_the_world();

        printf(" STOP THE WORLD \n");

        gc_mark();
        gc_sweep();

        gc_resume_world();
        usleep(GC_TIMEPERIOD);
    }
}


void gc_init() {
    // printf("preparing gc ... \n");
    atomic_store(&gc_state.world_stopped, 0);
    atomic_store(&gc_state.stopped_count, 0);
    atomic_store(&gc_state.total_threads, 0);
    pthread_mutex_init(&gc_state.lock, 0);
    pthread_cond_init(&gc_state.cv_mutators_stopped, 0);
    pthread_cond_init(&gc_state.cv_world_resumed, 0);


    pthread_t t;
    pthread_create(&t, NULL, gc_run, 0);
}

void gc_stop_the_world() {
    // printf("[GC-stw] acquire lock \n");
    pthread_mutex_lock(&gc_state.lock);
    // printf("[GC-stw] set world stopped \n");
    atomic_store(&gc_state.world_stopped, 1);
    
    // printf("[GC-stw] wait for stopped_count to become = total_threads \n");
    while (atomic_load(&gc_state.stopped_count) < atomic_load(&gc_state.total_threads)){
        // printf("[GC-stw] still waiting for stopped_count to become = total_threads \n");
        pthread_cond_wait(&gc_state.cv_mutators_stopped, &gc_state.lock);
    }
    
    // printf("[GC-stw] release lock\n");
    pthread_mutex_unlock(&gc_state.lock);
}


void gc_resume_world() {
    // printf("[GC-resume] acquire lock \n");
    pthread_mutex_lock(&gc_state.lock);

    // printf("[GC-resume] reset world stopped \n");
    atomic_store(&gc_state.world_stopped, 0);

    // printf("[GC-resume] reset stopped_count \n");
    atomic_store(&gc_state.stopped_count, 0);

    // printf("[GC-resume] broadcast to mutator \n");
    pthread_cond_broadcast(&gc_state.cv_world_resumed);

    // printf("[GC-resume] release lock\n");
    pthread_mutex_unlock(&gc_state.lock);
}