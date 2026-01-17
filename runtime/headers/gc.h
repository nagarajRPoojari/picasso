#ifndef GC_H
#define GC_H

#include <pthread.h>
#include <stdatomic.h>
#include <ucontext.h>
#include "scheduler.h" // external: runtime


#define GC_TIMEPERIOD (10 * 1000000)
#define MAX_ARENAS 12
#define MAX_SCHEDULERS 12 

#ifndef GC_PTR_ALIGNMENT
#define GC_PTR_ALIGNMENT 8
#endif
#define GC_ALIGN_MASK (GC_PTR_ALIGNMENT - 1)


typedef struct gc_state {
    atomic_int       world_stopped;      // 0 = running, 1 = requested stop
    atomic_int       stopped_count;      // number of threads that have stopped
    pthread_mutex_t  lock;
    pthread_cond_t   cv_mutators_stopped;
    pthread_cond_t   cv_world_resumed;

    // add_lock protects adding total_threads
    pthread_mutex_t add_lock;
    atomic_int        total_threads;      // mutators only
} gc_state_t;

arena_t* gc_create_global_arena();
arena_t* gc_create_arena();
void gc_register_root(task_t* t);
void gc_unregister_root(task_t* t);

void gc_init(); 
void gc_start();
void gc_stop_the_world();
void gc_resume_world();

#endif