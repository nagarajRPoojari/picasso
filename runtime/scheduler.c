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
#include <gc.h>

#include "scheduler.h"
#include "io.h"
#include "queue.h"
#include "task.h"



__thread task_t* current_task;

/**
 * @brief Entry trampoline used by tasks to invoke their function and clean up.
 * 
 * @param t Pointer to the current task.
 * @return Always returns NULL after task exits.
 */
void* task_trampoline(task_t *t) {
    t->fn(t);
    task_destroy(t);
    return NULL;
}

/**
 * @brief Signal handler for SIGSEGV that grows task stacks dynamically.
 * 
 * Triggered when the task overflows its current stack and accesses the guard page.
 * Allocates a new, larger stack, copies old stack data, and updates the context.
 * 
 * @param sig     Signal number (SIGSEGV).
 * @param si      Signal info.
 * @param unused  Pointer to ucontext_t (processor state at time of fault).
 */
void more_stack(int sig, siginfo_t *si, void *unused) {
    ucontext_t *ctx = (ucontext_t *)unused;

    uintptr_t sp = ctx->uc_mcontext.sp;
    uintptr_t old_stack_base = (uintptr_t)current_task->stack;
    uintptr_t old_stack_top = old_stack_base + current_task->stack_size;
    size_t old_stack_size = current_task->stack_size + PAGE_SIZE;

    /** @todo: print only in debug mode */
    printf("[task] Stack overflow detected, growing stack from %zu bytes\n", current_task->stack_size);
    printf("[task] SP = 0x%lx, old_stack = [0x%lx - 0x%lx]\n", sp, old_stack_base, old_stack_top);

    /** allocate new stack of twice size  */
    size_t new_stack_size = current_task->stack_size * 2;
    size_t total_size = new_stack_size + PAGE_SIZE;

    void* mapped = mmap(NULL, total_size, PROT_READ | PROT_WRITE,
                        MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    
    if (mapped == MAP_FAILED) {
        perror("mmap");
        exit(1);
    }

    /** 
    * protect first page (GUARD_PAGE)
    * @fix: what if stack steps ahead of 1 page ?
    */ 
    if (mprotect(mapped, PAGE_SIZE, PROT_NONE) != 0) {
        perror("mprotect");
        exit(1);
    }

    /**
     * Stack layout:
     *
     * Old stack:
     * 
     *   +-------------------------+  <- old_stack_top (SP grows downward)
     *   |       Stack Data        |
     *   |       AAA               |
     *   |       BBB               |  <- sp 
     *   |       ...               |
     *   +-------------------------+
     *   |       Guard Page        |  <- protected by mprotect(PROT_NONE)
     *   +-------------------------+
     *
     * New stack (after growth): @fix: do i need to twice ?, any othey ways ?
     *
     *   +-------------------------+  <- new_stack_top 
     *   |       Stack Data        |
     *   |       AAA               |
     *   |       BBB               |  <- sp (old one)
     *   |       ...               |
     *   |       ...               |
     *   |       ...               |
     *   +-------------------------+
     *   |       Guard Page        |  <- protected by mprotect(PROT_NONE)
     *   +-------------------------+
     *
     */

    /**
     * @fix: older page roots need to be removed, or else gc
     * will continue scanning it.
     * skip guard page
    */
    void* root_start = (char*)mapped + PAGE_SIZE; 
    void* root_end   = (char*)mapped + total_size;
    GC_add_roots(root_start, root_end);

    void* new_stack = (char*)mapped + PAGE_SIZE;

    uintptr_t sp_offset = old_stack_top - sp;
    /** @todo: need to verify this */
    /** 
     * stack pointer at the time of inturrupt could have crossed limit
     * if (sp_offset > current_task->stack_size) {
     *     fprintf(stderr, "SP offset too large\n");
     *     exit(1);
     * }
    */ 

    /** copy current stack content to new stack */
    memcpy((char*)new_stack + new_stack_size - current_task->stack_size,
           current_task->stack,
           current_task->stack_size);

    /** update current_task stack base & size */
    current_task->stack = new_stack;
    current_task->stack_size = new_stack_size;

    /**
     * @important: update new context stack pointer to same offset so that it
     * can resume from where it left off.
    */
    ctx->uc_mcontext.sp = (uintptr_t)new_stack + new_stack_size - sp_offset;
    printf("[task] New stack allocated: %zu bytes\n", new_stack_size);

    /**
     * clear grand old stack memory.
     * @ex: in 3rd iteration it will clear 1st stack memory
     * @fix: fix, may be in scheduler 
    */
    static void* _old_stack = NULL;
    static size_t _old_stack_size = 0;
    /**
     * @warning seems to be error prone, @test
     * if(_old_stack) munmap(_old_stack, _old_stack_size);
    */
    _old_stack = (char*)old_stack_base;
    _old_stack_size = old_stack_size + PAGE_SIZE;
}

/**
 * @brief Initialize alternate stack and install SIGSEGV handler.
 * 
 * This ensures the signal handler has a safe stack to run on if the current
 * task stack is corrupted or overflown.
 */
void init_stack_signal_handler() {
    /**
     * allocating alternate stack from SIG handler
     */
    stack_t altstack;
    altstack.ss_sp = malloc(SIGSTKSZ);
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
 * map of all threads preempt state
 * timer SIGNAL can interrupt any kernel thread,
 * I want it to set own thread preempt state
 */
volatile sig_atomic_t preempt[SCHEDULER_THREAD_POOL_SIZE];

/**
 * @brief Timer callback that forces preemption of a scheduler thread.
 * 
 * Sets the corresponding thread’s preempt flag to 1.
 * 
 * @param sv Signal value passed by the POSIX timer.
 */
void force_preempt(union sigval sv) {
    int tid = *(int *)sv.sival_ptr;
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
    its.it_value.tv_nsec = 50 * 1000;  // 0.05s
    its.it_interval.tv_sec = 0;
    its.it_interval.tv_nsec = 50 * 1000; // 0.05s

    if (timer_settime(tid, 0, &its, NULL) == -1) {
        perror("timer_settime");
        pthread_exit(NULL);
    }
}

/**
 * @brief Create a new task with its own protected stack and context.
 * 
 * @param fn   The function to run in the new task.
 * @param this Pointer argument passed to the task function.
 * @param kt   Pointer to the owning kernel thread (scheduler worker).
 * @return Pointer to the created task structure.
 */
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
    makecontext(&t->ctx, (void(*)(void))fn, 1, this);
    return t;
}

/**
 * @brief Clean up a task and release its resources.
 * 
 * @param t Task to destroy.
 */
void task_destroy(task_t *t) {
    if (!t) return;
    free(t->stack);
    free(t);
}

/**
 * @brief Yield execution from current task back to its scheduler thread.
 * 
 * @param kt The kernel thread running the task.
 */
void task_yield(kernel_thread_t* kt) {
    if (!kt->current) return;
    // swap context back to scheduler
    swapcontext(&kt->current->ctx, &kt->sched_ctx);
}

/**
 * @brief Cooperative preemption check.
 * 
 * Called periodically (e.g., via timer) to allow preemptive multitasking.
 * If the current task’s preempt flag is set, it yields control.
 */
void self_yield() {
    if(preempt[current_task->sched_id]) {
        kernel_thread_t* kt = kernel_thread_map[current_task->sched_id];
        safe_q_push(&kt->ready_q, current_task);
        preempt[current_task->sched_id] = 0;
        task_yield(kt);
    }
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
    // swap context back to task
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
        task_t *t;

        // Run ready tasks
        while ((t = safe_q_pop(&kt->ready_q)) != NULL) {
            current_task = t; // Set TLS pointer
            task_resume(t, kt);
            current_task = NULL; // Clear after task yields
        }

        // Wait for epoll events
        int n = epoll_wait(epfd, events, MAX_EVENTS, 100);
        for (int i = 0; i < n; i++) {
            task_t *t = (task_t *)events[i].data.ptr;
            t->nread = read(t->fd, t->buf, t->readn);
            epoll_ctl(epfd, EPOLL_CTL_DEL, t->fd, NULL);
            safe_q_push(&(kernel_thread_map[t->sched_id]->ready_q), t);
        }
    }
}