
#define _GNU_SOURCE
#include <ucontext.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <pthread.h>
#include <sys/epoll.h>
#include <signal.h>

#include <gc.h>
#include <gc/gc.h>       // optional, some versions

#include "start.h"
#include "array.h"
#include "ggc.h"
#include "globals.h"
#include "io.h"
#include "queue.h"
#include "scheduler.h"
#include "task.h"



kernel_thread_t **kernel_thread_map;


void thread(void*(*fn)(void*), void *this) {
    int kernel_thread_id = rand() % SCHEDULER_THREAD_POOL_SIZE;
    task_t *t1 = task_create(fn, this, kernel_thread_map[kernel_thread_id]);
    t1->id = rand();
    safe_q_push(&(kernel_thread_map[kernel_thread_id]->ready_q), t1);
}


int init_io() {
    safe_q_init(&io_queue, IO_QUEUE_SIZE);
    
    epfd = epoll_create1(0);
    if (epfd == -1) { perror("epoll_create1"); return 1; }
    
    pthread_t io_threads[IO_THREAD_POOL_SIZE];
    // Thread pool
    for (int i=0;i<IO_THREAD_POOL_SIZE;i++){
        pthread_create(&io_threads[i], NULL, io_worker, NULL);
    }

    return 0;
}

pthread_t sched_threads[SCHEDULER_THREAD_POOL_SIZE];
int init_scheduler() {
    kernel_thread_map = calloc(4, sizeof(kernel_thread_t*));
    for (int i=0;i<SCHEDULER_THREAD_POOL_SIZE;i++) {
        kernel_thread_map[i] = calloc(1, sizeof(kernel_thread_t));
        kernel_thread_map[i]->id = i;
        kernel_thread_map[i]->current = NULL;
        safe_queue_t ready_q;
        safe_q_init(&ready_q, SCHEDULER_LOCAL_QUEUE_SIZE);

        kernel_thread_map[i]->ready_q = ready_q;
        pthread_create(&sched_threads[i], NULL, scheduler_run, kernel_thread_map[i]);
    }

    return 0;
}

void clean_scheduler() {
    free(kernel_thread_map[0]);
}


int main() {
    srand(time(NULL));

    GC_INIT();
    GC_allow_register_threads(); 

    init_io();
    init_scheduler();

    thread(start, NULL);


    for (int i = 0; i < SCHEDULER_THREAD_POOL_SIZE; i++) {
        pthread_join(sched_threads[i], NULL);
    }

    clean_scheduler();
    return 0;
}
