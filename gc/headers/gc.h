#ifndef GC_H
#define GC_H

#include <pthread.h>
#include <stdatomic.h>


typedef struct gc_state {
    atomic_int       world_stopped;      // 0 = running, 1 = requested stop
    atomic_int       stopped_count;      // number of threads that have stopped
    pthread_mutex_t  lock;
    pthread_cond_t   cv_mutators_stopped;
    pthread_cond_t   cv_world_resumed;

    atomic_int        total_threads;      // mutators only
} gc_state_t;

void gc_init();
void gc_stop_the_world();
void gc_resume_world();

#endif