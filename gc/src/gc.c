#include <stdio.h>
#include <stdatomic.h>
#include <unistd.h>
#include "alloc.h"
#include "gc.h"

/* expected to be initialized in gc_init */
gc_state_t gc_state;

void gc_run() {
    printf(" started gc \n");
    for(int i=0;i<10000;i++){
        gc_stop_the_world();

        printf(" STOP THE WORLD \n");

        gc_resume_world();

        usleep(1200 * 1000);
    }
}


void gc_init() {
    printf("preparing gc ... \n");


    atomic_store(&gc_state.world_stopped, 0);
    atomic_store(&gc_state.stopped_count, 0);
    atomic_store(&gc_state.total_threads, 0);
    pthread_mutex_init(&gc_state.lock, 0);
    pthread_cond_init(&gc_state.cv_mutators_stopped, 0);
    pthread_cond_init(&gc_state.cv_world_resumed, 0);


    pthread_t t;
    // pthread_create(&t, NULL, gc_run, 0);
}

void gc_stop_the_world() {
    printf("[GC-stw] acquire lock \n");
    pthread_mutex_lock(&gc_state.lock);
    printf("[GC-stw] set world stopped \n");
    atomic_store(&gc_state.world_stopped, 1);
    
    printf("[GC-stw] wait for stopped_count to become = total_threads \n");
    while (atomic_load(&gc_state.stopped_count) < atomic_load(&gc_state.total_threads)){
        printf("[GC-stw] still waiting for stopped_count to become = total_threads \n");
        pthread_cond_wait(&gc_state.cv_mutators_stopped, &gc_state.lock);
    }
    
    printf("[GC-stw] release lock\n");
    pthread_mutex_unlock(&gc_state.lock);
}


void gc_resume_world() {
    printf("[GC-resume] acquire lock \n");
    pthread_mutex_lock(&gc_state.lock);

    printf("[GC-resume] reset world stopped \n");
    atomic_store(&gc_state.world_stopped, 0);

    printf("[GC-resume] reset stopped_count \n");
    atomic_store(&gc_state.stopped_count, 0);

    printf("[GC-resume] broadcast to mutator \n");
    pthread_cond_broadcast(&gc_state.cv_world_resumed);

    printf("[GC-resume] release lock\n");
    pthread_mutex_unlock(&gc_state.lock);
}
