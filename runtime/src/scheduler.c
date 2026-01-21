#include "platform.h"
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <sys/mman.h>
#include <semaphore.h>
#include <stdatomic.h>
#include <ffi.h>
#include <stdarg.h>

#include "scheduler.h"
#include "platform/context.h"
#include "diskio.h"
#include "queue.h"
#include "task.h"
#include "alloc.h"
#include "gc.h"
#include "initutils.h"
#include "sigerr.h"

__thread task_t* current_task;

// This replaces the dangerous static global variables in more_stack.
typedef struct {
    void* base_ptr;
    size_t total_size;
} old_stack_info_t;

__thread safe_queue_t old_stack_cleanup_q;
__thread arena_t* __arena__;

extern gc_state_t gc_state;
extern arena_t* __global__arena__;

void task_yield(kernel_thread_t* kt);
void task_resume(task_t *t, kernel_thread_t* kt);


/**
 * @brief Entry trampoline used by tasks to invoke their function and clean up.
 * @param t Pointer to the current task.
 * @param this Arbitary arg to be passed to func.
 * @return Always returns NULL after task exits.
 */
void task_trampoline(uintptr_t _t, uintptr_t _p) {
    task_t* t = (task_t*)_t;
    task_payload_t* payload = (task_payload_t*)_p;

    void *retval;

    // Dynamically invoke the function with the prepared registers
    ffi_call(&payload->cif, FFI_FN(payload->fn), &retval, payload->arg_values);

    t->state = TASK_FINISHED;

    // Clean up payload memory
    for (int i = 0; i < payload->nargs; i++) {
        release(__global__arena__, payload->arg_values[i]);
    }
    release(__global__arena__, payload->arg_types);
    release(__global__arena__, payload->arg_values);
    release(__global__arena__, payload);
}

/**
 * @deprecated currently i am not implementing dynamic stack growth
 * @brief Signal handler for SIGSEGV that grows task stacks dynamically.
 */
void more_stack(int sig, siginfo_t *si, void *unused) {
    platform_ctx_t *ctx = (platform_ctx_t *)unused;

    uintptr_t old_stack_base = (uintptr_t)current_task->stack;                  
    uintptr_t sp = platform_ctx_get_stack_pointer(ctx);

    if (sp < old_stack_base) {
        sp = old_stack_base;   // clamp to prevent underflow
    }

    uintptr_t old_stack_base_guard = (uintptr_t)current_task->stack - PAGE_SIZE; 
    size_t    old_stack_size = current_task->stack_size;
    uintptr_t old_stack_top = old_stack_base + old_stack_size;                  
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

    if (sp > old_stack_base) {
        return;
    }

    /* abort for now */
    printf("[FATAL] ================= stack over flow ================= \n" );
    exit(1);
}

/**
 * @brief Initialize alternate stack and install SIGSEGV handler.
 * 
 * This ensures the signal handler has a safe stack to run on if the current
 * task stack is corrupted or overflown.
 */
void init_stack_signal_handler() {
    // Initialize the thread-local cleanup queue
    safe_q_init(&old_stack_cleanup_q, 10);

    // allocating alternate stack for SIG handler
    stack_t altstack;
    altstack.ss_sp = allocate(__arena__, SIGSTKSZ); // This is safe because it runs on the kernel stack
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
 * @brief Clean up a task and release its resources.
 * * @param t Task to destroy.
 */
void task_destroy(task_t *t) {
    if (!t) return;
    gc_unregister_root(t);

    void* original_base = (char*)t->stack - PAGE_SIZE; 
    size_t total_size = t->stack_size + PAGE_SIZE;
    
    if (munmap(original_base, total_size) != 0) {
        perror("munmap failed in task_destroy");
    }
    release(__global__arena__, t);
}

volatile sig_atomic_t preempt[SCHEDULER_THREAD_POOL_SIZE];
/**
 * @brief Timer callback that forces preemption of a scheduler thread.
 * 
 * Sets the corresponding thread’s preempt flag to 1.
 * 
 * @param sv Signal value passed by the POSIX timer.
 */
void force_preempt(int sig, siginfo_t *si, void *uc) {
    int tid = *(int *)si->si_value.sival_ptr; 
    preempt[tid] = 1;
}

/**
 * @brief Initialize per-thread timer signal handler.
 * 
 * Creates a periodic POSIX timer (SIGEV_THREAD) that triggers preemption
 * at fixed intervals for the current scheduler thread.
 * 
 * @param arg Pointer to the scheduler thread ID (int*).
 */
void init_timer_signal_handler(void *arg) {
    // struct sigaction sa;
    // sa.sa_flags = SA_SIGINFO; 
    // sa.sa_sigaction = force_preempt; 
    
    // sigemptyset(&sa.sa_mask); 
    
    // if (sigaction(SIGRTMIN, &sa, NULL) == -1) {
    //     perror("sigaction failed");
    // }

    // timer_t tid;
    // struct sigevent sev;
    // struct itimerspec its;

    // sev.sigev_notify = SIGEV_SIGNAL;
    // sev.sigev_signo = SIGRTMIN; 
    // sev.sigev_value.sival_ptr = arg; 


    // if (timer_create(CLOCK_REALTIME, &sev, &tid) == -1) {
    //     perror("timer_create");
    //     pthread_exit(NULL);
    // }

    // its.it_value.tv_sec = 0;
    // its.it_value.tv_nsec = 50000000;
    // its.it_interval.tv_sec = 0;
    // its.it_interval.tv_nsec = 50000000;

    // if (timer_settime(tid, 0, &its, NULL) == -1) {
    //     perror("timer_settime");
    //     pthread_exit(NULL);
    // }
}

task_t* task_create(void* (*fn)(), void* payload, kernel_thread_t* kt) {
    task_t *t = allocate(__global__arena__, sizeof(task_t));
    platform_ctx_init(&t->ctx);
    
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

    platform_ctx_make(&t->ctx, task_trampoline, (uintptr_t)t, (uintptr_t)payload, t->stack, t->stack_size, &(kt->sched_ctx));

    gc_register_root(t);
    return t;
}

/**
 * @brief Yield execution from current task back to its scheduler thread.
 * 
 * @param kt The kernel thread running the task.
 */
void task_yield(kernel_thread_t* kt) {
    if (!kt->current) return;
    platform_ctx_switch(&kt->current->ctx, &kt->sched_ctx);
}

/**
 * @brief Cooperative preemption check.
 * 
 * Called periodically (e.g., via timer) to allow preemptive multitasking.
 * If the current task’s preempt flag is set, it yields control.
 */
void self_yield() {
    if (!atomic_load_explicit(&gc_state.world_stopped, memory_order_acquire)) {
        if(preempt[current_task->sched_id]) {
            kernel_thread_t* kt = kernel_thread_map[current_task->sched_id];
            safe_q_push(&kt->ready_q, current_task);
            preempt[current_task->sched_id] = 0;
            task_yield(kt);
        }
        return;
    }

    pthread_mutex_lock(&gc_state.lock);

    if (atomic_fetch_add(&gc_state.stopped_count, 1) + 1 == atomic_load(&gc_state.total_threads)){
        pthread_cond_signal(&gc_state.cv_mutators_stopped);
    }

    while (atomic_load(&gc_state.world_stopped)){
        pthread_cond_wait(&gc_state.cv_world_resumed, &gc_state.lock);
    }

    pthread_mutex_unlock(&gc_state.lock);
}

/**
 * @brief Resume a specific task on the given scheduler thread.
 * 
 * @param t  Task to resume.
 * @param kt The kernel thread executing the task.
 */
void task_resume(task_t *t, kernel_thread_t* kt) {
    kt->current = t;
    current_task = t;
    platform_ctx_switch(&kt->sched_ctx, &t->ctx);
    kt->current = NULL;
    current_task = NULL;
}

/**
 * @brief Main run loop for a scheduler worker thread.
 * 
 * Continuously executes ready tasks from the queue, handles preemption,
 * and waits for epoll I/O events to resume blocked tasks.
 * 
 * @param arg Pointer to the kernel_thread_t structure for this thread.
 * @return Never returns under normal operation.
 */
void* scheduler_run(void* arg) {
    kernel_thread_t* kt = (kernel_thread_t*)arg;

    __arena__ = gc_create_arena(kt);

    // init_stack_signal_handler();
    init_timer_signal_handler(arg);

    #if defined(__linux__)
        init_error_handlers();
    #endif
    while (1) {
        task_t *t;
        while (1) {
            
            t = safe_q_pop_wait(&kt->ready_q);
            if(!t) {
                return NULL;
            }
            kt->current = t;

            pthread_mutex_lock(&gc_state.add_lock);
            atomic_fetch_add(&gc_state.total_threads, 1);
            pthread_mutex_unlock(&gc_state.add_lock);

            unsafe_ioq_remove(&kt->wait_q, t);

            t->state = TASK_RUNNING;
            task_resume(t, kt);

            if (t->state == TASK_FINISHED) {
                /* @verify: this doesn't seems to be efficient way */
                task_destroy(t);
                atomic_fetch_sub(&task_count, 1);
                if(!atomic_load(&task_count)) {
                    for(int i=0; i<SCHEDULER_THREAD_POOL_SIZE; i++) {
                        safe_q_push(&kernel_thread_map[i]->ready_q, NULL);
                    }
                    return NULL;
                }
            }
            atomic_fetch_sub(&gc_state.total_threads, 1);
        }
    }
    return NULL; 
}