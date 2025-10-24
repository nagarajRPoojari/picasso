#ifndef QUEUE_H 
#define QUEUE_H

#include <pthread.h>
#include "task.h"

typedef struct task_node { 
    task_t *t; 
    struct task_node *next; 
} task_node_t;

// thread safe queue
typedef struct {
    int size_limit;
    
    task_node_t *head, *tail;
    pthread_mutex_t lock;
    pthread_cond_t cond;
} safe_queue_t;

// safe_q_init initializes data structure
void safe_q_init(safe_queue_t *q, int size);

void safe_q_push(safe_queue_t *q, task_t *t);

// safe_q_pop is non blocking, will be returned 
// with NULL if there is no elements pop
task_t *safe_q_pop(safe_queue_t *q);

// safe_q_pop_wait is blocking in nature, waits for 
// element if empty, never returns NULL
task_t *safe_q_pop_wait(safe_queue_t *q);

#endif