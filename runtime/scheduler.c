#include <ucontext.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <pthread.h>
#include <sys/epoll.h>
#include <sys/mman.h>
#include <signal.h>
#include <stdint.h> // Include for uintptr_t

#include "scheduler.h"
#include "io.h"
#include "queue.h"
#include "task.h"


__thread task_t* current_task;

// This replaces the dangerous static global variables in more_stack.
typedef struct {
    void* base_ptr;
    size_t total_size;
} old_stack_info_t;

__thread safe_queue_t old_stack_cleanup_q;

void task_yield(kernel_thread_t* kt);
void task_resume(task_t *t, kernel_thread_t* kt);


/**
 * @brief Entry trampoline used by tasks to invoke their function and clean up.
 * @param t Pointer to the current task.
 * @param this Arbitary arg to be passed to func.
 * @return Always returns NULL after task exits.
 */
void* task_trampoline(task_t *t, void *this) {
    t->fn(this);
    return NULL;
}

/**
 * @deprecated currently i am not implementing dynamic stack growth
 * @brief Signal handler for SIGSEGV that grows task stacks dynamically.
 */
void more_stack(int sig, siginfo_t *si, void *unused) {
    ucontext_t *ctx = (ucontext_t *)unused;

    // --- 1. Calculate Old Stack Pointers and Active Data Size ---
    uintptr_t old_stack_base = (uintptr_t)current_task->stack;                  // Base of usable stack
    uintptr_t sp = (uintptr_t)ctx->uc_mcontext.sp;

    if (sp < old_stack_base) {
        sp = old_stack_base;   // clamp to prevent underflow
    }

    uintptr_t old_stack_base_guard = (uintptr_t)current_task->stack - PAGE_SIZE; // Base of mmap region
    size_t    old_stack_size = current_task->stack_size;
    uintptr_t old_stack_top = old_stack_base + old_stack_size;                  // Top of usable stack
    uintptr_t copy_size = old_stack_top - sp;


    printf("[DEBUG] SP (captured): 0x%lx\n", sp);
    printf("[DEBUG] Old Stack Base (usable): 0x%lx\n", old_stack_base);
    printf("[DEBUG] Old Stack Top (usable): 0x%lx\n", old_stack_top);
    printf("[DEBUG] Calculated Copy Size: %zu\n", copy_size);
    printf("[DEBUG] Old Stack Size (allocated): %zu\n", old_stack_size);


    printf("[%lu ////] [%lu ///// %lu ....... %lu]\n",
       (unsigned long)(uintptr_t)old_stack_base_guard,
       (unsigned long)(uintptr_t)old_stack_base,
       (unsigned long)(uintptr_t)sp,
       (unsigned long)(uintptr_t)old_stack_top);

    // Active stack data runs from SP up to old_stack_top
    if (copy_size > old_stack_size) {
        // safe_debug("[FATAL] Stack pointer error: copy size exceeds task stack size.\n");
        _exit(1);
    }

    mprotect(old_stack_base_guard, PAGE_SIZE, PROT_NONE);

    // Ensure we don't try to copy more than the stack size (shouldn't happen 
    // if the guard page is respected, but good for safety)
    if (copy_size > old_stack_size) {
        fprintf(stderr, "[FATAL] Stack pointer error: copy size exceeds task stack size.\n");
        exit(1);
    }
    
    // --- 2. Allocate New Stack Memory ---
    size_t new_stack_size = old_stack_size * 2;
    size_t total_size = new_stack_size + PAGE_SIZE;

    void* mapped = mmap(NULL, total_size, PROT_READ | PROT_WRITE,
                        MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    
    if (mapped == MAP_FAILED) {
        perror("mmap failed during stack growth");
        exit(1);
    }

    // Protect the first page (GUARD_PAGE)
    if (mprotect(mapped, PAGE_SIZE, PROT_NONE) != 0) {
        perror("mprotect failed during stack growth");
        exit(1);
    }

    void* new_stack_base = (char*)mapped + PAGE_SIZE;
    uintptr_t new_stack_top = (uintptr_t)new_stack_base + new_stack_size;

    // The stack grows downward. We copy only the active data (from SP to old_stack_top).
    // The new SP must maintain the same offset from the new stack top.
    uintptr_t new_sp = new_stack_top - copy_size;
    
    // memcpy(destination, source, size)
    memcpy((void*)new_sp, (void*)sp, copy_size); 

    // Track the old stack for asynchronous cleanup
    old_stack_info_t old_stack_info = {
        .base_ptr = (void*)old_stack_base_guard,
        .total_size = old_stack_size + PAGE_SIZE
    };

    /* @danger: can't push whatever i want to safe_q (which expects only task_t*)*/
    // Use a safe queue push if queue operations are not signal-safe (recommended)
    safe_q_push(&old_stack_cleanup_q, &old_stack_info); 

    current_task->stack = new_stack_base;
    current_task->stack_size = new_stack_size;
    
    ctx->uc_mcontext.sp = new_sp;

    printf("Allocated new \n");
    printf("[0x%lx////][0x%lx/////0x%lx.......0x%lx] \n", mapped, new_stack_base, new_sp ,new_stack_top);
}

/**
 * @brief Initialize alternate stack and install SIGSEGV handler.
 */
void init_stack_signal_handler() {
    // Initialize the thread-local cleanup queue
    safe_q_init(&old_stack_cleanup_q, 10);

    // allocating alternate stack for SIG handler
    stack_t altstack;
    altstack.ss_sp = malloc(SIGSTKSZ); // This is safe because it runs on the kernel stack
    if (altstack.ss_sp == NULL) {
        perror("malloc for altstack");
        exit(1);
    }
    altstack.ss_size = SIGSTKSZ;
    altstack.ss_flags = 0;
    if (sigaltstack(&altstack, NULL) < 0) {
        perror("sigaltstack");
        exit(1);
    }

    struct sigaction sa;
    sa.sa_flags = SA_SIGINFO | SA_ONSTACK;
    sa.sa_sigaction = more_stack;
    sigemptyset(&sa.sa_mask);
    if (sigaction(SIGSEGV, &sa, NULL) < 0) {
        perror("sigaction");
        exit(1);
    }
}

/**
 * @brief Clean up any old, munmap-able stacks pending cleanup.
 * @note Must be called from the scheduler thread main loop.
 */
void cleanup_old_stacks() {
    old_stack_info_t* info;
    while ((info = safe_q_pop(&old_stack_cleanup_q)) != NULL) {
        if (munmap(info->base_ptr, info->total_size) != 0) {
            // Note: If munmap fails, it's a critical error/leak.
            perror("munmap failed during stack cleanup");
        }
        // free(info); // Free the info structure if safe_q_push malloced it
    }
}

/**
 * @brief Clean up a task and release its resources.
 * * @param t Task to destroy.
 */
void task_destroy(task_t *t) {
    if (!t) return;

    void* original_base = (char*)t->stack - PAGE_SIZE; 
    size_t total_size = t->stack_size + PAGE_SIZE;
    
    if (munmap(original_base, total_size) != 0) {
        perror("munmap failed in task_destroy");
    }
    // free(t);
}

volatile sig_atomic_t preempt[SCHEDULER_THREAD_POOL_SIZE];

void force_preempt(union sigval sv) {
    int tid = *(int *)sv.sival_ptr;
    preempt[tid] = 1;
}

void init_timer_signal_handler(void *arg) {
    timer_t tid;
    struct sigevent sev;
    struct itimerspec its;

    int id = *(int *)arg;
    sev.sigev_notify = SIGEV_THREAD;
    sev.sigev_notify_function = force_preempt;
    sev.sigev_notify_attributes = NULL;
    sev.sigev_value.sival_ptr = arg;

    if (timer_create(CLOCK_REALTIME, &sev, &tid) == -1) {
        perror("timer_create");
        pthread_exit(NULL);
    }

    its.it_value.tv_sec = 0;
    its.it_value.tv_nsec = 500000;  // Changed to 0.5ms (500,000ns) for less overhead
    its.it_interval.tv_sec = 0;
    its.it_interval.tv_nsec = 500000;

    if (timer_settime(tid, 0, &its, NULL) == -1) {
        perror("timer_settime");
        pthread_exit(NULL);
    }
}

task_t* task_create(void* (*fn)(void *), void* this, kernel_thread_t* kt) {
    task_t *t = calloc(1, sizeof(*t));
    if (!t) { 
        perror("calloc"); 
        exit(1); 
    }

    t->fn = fn;
    t->stack_size = STACK_SIZE;
    t->sched_id = kt->id;

    // allocate stack + guard page
    void* mapped = mmap(NULL, t->stack_size + PAGE_SIZE,
                        PROT_READ | PROT_WRITE,
                        MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    if (mapped == MAP_FAILED) { 
        perror("mmap"); 
        exit(1); 
    }

    // protect guard page
    if (mprotect(mapped, PAGE_SIZE, PROT_NONE) != 0) {
        perror("mprotect");
        exit(1);
    }

    // usable stack starts just above guard page
    t->stack = (char*)mapped + PAGE_SIZE;

    // initialize ucontext
    getcontext(&t->ctx);
    t->ctx.uc_stack.ss_sp = t->stack;
    t->ctx.uc_stack.ss_size = t->stack_size;
    t->ctx.uc_link = &(kt->sched_ctx);

    // make trampoline
    makecontext(&t->ctx, (void(*)(void))task_trampoline, 2, t, this);
    return t;
}

void task_yield(kernel_thread_t* kt) {
    if (!kt->current) return;
    swapcontext(&kt->current->ctx, &kt->sched_ctx);
}

void self_yield() {
    if(preempt[current_task->sched_id]) {
        kernel_thread_t* kt = kernel_thread_map[current_task->sched_id];
        safe_q_push(&kt->ready_q, current_task);
        preempt[current_task->sched_id] = 0;
        task_yield(kt);
    }
}

void task_resume(task_t *t, kernel_thread_t* kt) {
    kt->current = t;
    current_task = t;
    swapcontext(&kt->sched_ctx, &t->ctx);
    kt->current = NULL;
    current_task = NULL;
}

void* scheduler_run(void* arg) {
    kernel_thread_t* kt = (kernel_thread_t*)arg;
    struct epoll_event events[MAX_EVENTS];

    init_stack_signal_handler();
    init_timer_signal_handler(arg);

    while (1) {
        cleanup_old_stacks(); 
        task_t *t;
        while ((t = safe_q_pop(&kt->ready_q)) != NULL) {
            task_resume(t, kt);
        }
    }
    return NULL; 
}