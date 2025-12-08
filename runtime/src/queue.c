#include <stdlib.h>
#include <pthread.h>

#include "queue.h"
#include "task.h"

/**
 * @brief Initialize a thread-safe queue.
 * 
 * Sets head/tail to NULL, initializes mutex and condition variable,
 * and sets the maximum size of the queue.
 * 
 * @param q Pointer to the queue to initialize.
 * @param size Maximum number of elements allowed in the queue.
 */
void safe_q_init(safe_queue_t *q, int size) {
    q->head = q->tail = NULL;
    q->size_limit = size;
    // initialize lock
    pthread_mutex_init(&q->lock, NULL);
    pthread_cond_init(&q->cond, NULL);
}

/**
 * @brief Push a task onto the queue in a thread-safe manner.
 * 
 * Wakes up any threads waiting for a task. If the queue has a size limit,
 * the caller should handle blocking or dropping tasks as needed.
 * 
 * @param q Pointer to the queue.
 * @param t Task to push.
 */
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

/**
 * @brief Pop a task from the queue in a non-blocking way.
 * 
 * If the queue is empty, returns NULL immediately.
 * 
 * @param q Pointer to the queue.
 * @return Pointer to the task, or NULL if queue is empty.
 */
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

    // free(n);
    return t;
}

/**
 * @brief Pop a task from the queue in a blocking way.
 * 
 * If the queue is empty, the calling thread waits until a task becomes available.
 * Never returns NULL.
 * 
 * @param q Pointer to the queue.
 * @return Pointer to the next task in the queue.
 */
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

    // free(n);
    return t;
}