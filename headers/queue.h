#ifndef QUEUE_H 
#define QUEUE_H

#include <pthread.h>
#include "task.h"

/* Node in the task queue */
typedef struct task_node { 
    task_t *t;                 /* Pointer to the task */
    struct task_node *next;    /* Next node in the queue */
} task_node_t;

/* thread safe FIFO queue for tasks  */
typedef struct {
    int size_limit;            /* Maximum number of tasks in the queue */
    
    task_node_t *head, *tail;  /* Queue head and tail pointers */
    pthread_mutex_t lock;       /* Mutex protecting queue operations */
    pthread_cond_t cond;        /* Condition variable for waiting threads */
} safe_queue_t;

/**
 * @brief Initialize a thread-safe queue.
 * 
 * Sets head/tail to NULL, initializes mutex and condition variable,
 * and sets the maximum size of the queue.
 * 
 * @param q Pointer to the queue to initialize.
 * @param size Maximum number of elements allowed in the queue.
 */
void safe_q_init(safe_queue_t *q, int size);

/**
 * @brief Push a task onto the queue in a thread-safe manner.
 * 
 * Wakes up any threads waiting for a task. If the queue has a size limit,
 * the caller should handle blocking or dropping tasks as needed.
 * 
 * @param q Pointer to the queue.
 * @param t Task to push.
 */
void safe_q_push(safe_queue_t *q, task_t *t);

/**
 * @brief Pop a task from the queue in a non-blocking way.
 * 
 * If the queue is empty, returns NULL immediately.
 * 
 * @param q Pointer to the queue.
 * @return Pointer to the task, or NULL if queue is empty.
 */
task_t *safe_q_pop(safe_queue_t *q);

/**
 * @brief Pop a task from the queue in a blocking way.
 * 
 * If the queue is empty, the calling thread waits until a task becomes available.
 * Never returns NULL.
 * 
 * @param q Pointer to the queue.
 * @return Pointer to the next task in the queue.
 */
task_t *safe_q_pop_wait(safe_queue_t *q);


#endif