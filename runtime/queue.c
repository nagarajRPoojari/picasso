#include <stdlib.h>
#include <pthread.h>

#include "queue.h"
#include "task.h"

void safe_q_init(safe_queue_t *q, int size) {
    q->head = q->tail = NULL;
    q->size_limit = size;
    // initialize lock
    pthread_mutex_init(&q->lock, NULL);
    pthread_cond_init(&q->cond, NULL);
}

void safe_q_push(safe_queue_t *q, task_t *t) {
    task_node_t *n = malloc(sizeof(*n));
    n->t = t; n->next = NULL;

    pthread_mutex_lock(&q->lock);
    if(q->size_limit <= 0) {
        pthread_mutex_unlock(&q->lock);
        return;
    }
    q->size_limit--;
    if (!q->tail) 
        q->head = q->tail = n;
    else {
        q->tail->next = n; 
        q->tail = n; 
    }

    // signal only one consumer thread about the availability
    pthread_cond_signal(&q->cond);
    pthread_mutex_unlock(&q->lock);
}

task_t *safe_q_pop(safe_queue_t *q) {   
    pthread_mutex_lock(&q->lock);

    task_node_t *n = q->head;
    if (!n) { 
        pthread_mutex_unlock(&q->lock); 
        return NULL; 
    }
    q->head = n->next;
    if (!q->head) 
        q->tail = NULL;
    q->size_limit++;
    pthread_mutex_unlock(&q->lock);

    task_t *t = n->t;

    free(n);
    return t;
}

task_t *safe_q_pop_wait(safe_queue_t *q) {
    pthread_mutex_lock(&q->lock);

    // suspend if queue is empty
    // release lock & wait for signal from push
    while (!q->head) 
        pthread_cond_wait(&q->cond, &q->lock);


    task_node_t *n = q->head;
    q->head = n->next;
    if (!q->head) 
        q->tail = NULL;

    q->size_limit++;
    pthread_mutex_unlock(&q->lock);
    task_t *t = n->t;

    free(n);
    return t;
}