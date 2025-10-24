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

#define  PAGE_SIZE  sysconf(_SC_PAGESIZE)


__thread task_t* current_task;

void* task_trampoline(task_t *t) {
    t->fn(t);
    task_destroy(t);
    return NULL;
}

// sigsegv_handler handles corrupt memory access, i.e accessing GUARD_PAGE
// allocates new stack of double size & updates cpu context accordingly.
// clears earlier stack if has been allocated.
void sigsegv_handler(int sig, siginfo_t *si, void *unused) {
    ucontext_t *ctx = (ucontext_t *)unused;

    uintptr_t sp = ctx->uc_mcontext.sp;
    uintptr_t old_stack_base = (uintptr_t)current_task->stack;
    uintptr_t old_stack_top = old_stack_base + current_task->stack_size;
    size_t old_stack_size = current_task->stack_size + PAGE_SIZE;

    // @todo: print only in debug mode
    printf("[task] Stack overflow detected, growing stack from %zu bytes\n", current_task->stack_size);
    printf("[task] SP = 0x%lx, old_stack = [0x%lx - 0x%lx]\n", sp, old_stack_base, old_stack_top);

    // allocate new stack of twice size
    size_t new_stack_size = current_task->stack_size * 2;
    size_t total_size = new_stack_size + PAGE_SIZE;

    void* mapped = mmap(NULL, total_size, PROT_READ | PROT_WRITE,
                        MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    
    if (mapped == MAP_FAILED) {
        perror("mmap");
        exit(1);
    }

    // protect first page (GUARD_PAGE)
    if (mprotect(mapped, PAGE_SIZE, PROT_NONE) != 0) {
        perror("mprotect");
        exit(1);
    }


    // @todo: older page roots need to be removed, or else gc
    // will continue scanning it.
    // skip guard page
    void* root_start = (char*)mapped + PAGE_SIZE; 
    void* root_end   = (char*)mapped + total_size;
    GC_add_roots(root_start, root_end);

    void* new_stack = (char*)mapped + PAGE_SIZE;

    uintptr_t sp_offset = old_stack_top - sp;
    if (sp_offset > current_task->stack_size) {
        fprintf(stderr, "SP offset too large\n");
        exit(1);
    }

    // copy current stack content to new stack
    memcpy((char*)new_stack + new_stack_size - current_task->stack_size,
           current_task->stack,
           current_task->stack_size);

    // update current_task stack base & size
    current_task->stack = new_stack;
    current_task->stack_size = new_stack_size;

    // !IMP: update new context stack pointer to same offset so that it
    // can resume from where it left off.
    ctx->uc_mcontext.sp = (uintptr_t)new_stack + new_stack_size - sp_offset;
    printf("[task] New stack allocated: %zu bytes\n", new_stack_size);

    // clear grand old stack memory.
    // ex: in 3rd iteration it will clear 1st stack memory
    // @todo: fix, may be in scheduler 
    static void* _old_stack = NULL;
    static size_t _old_stack_size = 0;
    if(_old_stack) munmap(_old_stack, _old_stack_size);
    _old_stack = (char*)old_stack_base;
    _old_stack_size = old_stack_size + PAGE_SIZE;
}


void init_sigsegv_handler() {
    // allocating alternate stack from SIG handler
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
    sa.sa_sigaction = sigsegv_handler;
    sigemptyset(&sa.sa_mask);
    if (sigaction(SIGSEGV, &sa, NULL) < 0) {
        perror("sigaction");
        exit(1);
    }
}

task_t* task_create(void* (*fn)(void *), kernel_thread_t* kt) {
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
    makecontext(&t->ctx, (void(*)(void))fn, 1, t);
    return t;
}

void task_destroy(task_t *t) {
    if (!t) return;
    free(t->stack);
    free(t);
}

void task_yield(kernel_thread_t* kt) {
    if (!kt->current) return;
    // swap context back to scheduler
    swapcontext(&kt->current->ctx, &kt->sched_ctx);
}


volatile sig_atomic_t preempt[SCHEDULER_THREAD_POOL_SIZE];

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
    // swap context back to task
    swapcontext(&kt->sched_ctx, &t->ctx);
    kt->current = NULL;
    current_task = NULL;
}


void timer_callback(union sigval sv) {
    int tid = *(int *)sv.sival_ptr;
    preempt[tid] = 1;
}


void* scheduler_run(void* arg) {
    kernel_thread_t* kt = (kernel_thread_t*)arg;
    struct epoll_event events[MAX_EVENTS];

    int id = *(int *)arg;
    timer_t tid;
    struct sigevent sev;
    struct itimerspec its;

    sev.sigev_notify = SIGEV_THREAD;
    sev.sigev_notify_function = timer_callback;
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