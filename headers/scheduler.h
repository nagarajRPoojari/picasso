#ifndef SCHEDULER_H
#define SCHEDULER_H

#include <ucontext.h>
#include "globals.h"
#include "task.h"
#include "queue.h"

typedef struct kernel_thread {
    int id;
    ucontext_t sched_ctx;
    task_t *current;
    safe_queue_t ready_q;

} kernel_thread_t;

extern kernel_thread_t **kernel_thread_map;

// task_create initialises task context with given 
// function and allocates fixed stack memory
task_t *task_create(void* (*fn)(void *), void *this, kernel_thread_t* kt);

// task_destroy deallocates memory
void task_destroy(task_t *t);

// task_yield cooperatively yelds, swapping context
// back to scheduler
void task_yield(kernel_thread_t* kt);

// self_yield voluntarily yields back to scheduler, pushing itself
// to the back of ready queue
void self_yield();

// task_resume resumes READY task by swapping context
// back to function
void task_resume(task_t *t, kernel_thread_t* kt);

// scheduler_run is the main scheduler loop to pop 
// task from ready queue & resume 
void* scheduler_run(void* arg);
#endif