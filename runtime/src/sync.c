// #include "scheduler.h"
// #include "sync.h"
// #include "alloc.h"

// extern __thread arena_t* __arena__;

// __public__rwmutex_t* __public__create_rwmutex() {
//     __public__rwmutex_t* mux = (__public__rwmutex_t*)allocate(__arena__, sizeof(__public__rwmutex_t));
//     safe_q_init(&mux->readers, SCHEDULER_LOCAL_QUEUE_SIZE);
//     safe_q_init(&mux->writers, SCHEDULER_LOCAL_QUEUE_SIZE);

//     pthread_mutex_init(&mux->lock, 0);
//     atomic_store_explicit(&mux->state, 0, memory_order_relaxed);

//     return mux;
// }

// void __public__mutex_rlock(__public__rwmutex_t* mux) {
//     int old;

//     for (;;) {
//         old = atomic_load_explicit(&mux->state, memory_order_seq_cst);

//         if (old & WRITER_BIT) {
//             // Atomic enqueue + sleep
//             pthread_mutex_lock(&mux->lock);
//             // Recheck to avoid race
//             if (atomic_load_explicit(&mux->state, memory_order_seq_cst) & WRITER_BIT) {
//                 safe_q_push(&mux->readers, current_task);
//                 // Mark task as "sleeping" here implicitly
//                 pthread_mutex_unlock(&mux->lock);
//                 task_yield(kernel_thread_map[current_task->sched_id]);
//                 continue;
//             }
//             pthread_mutex_unlock(&mux->lock);
//             continue;
//         }

//         if (atomic_compare_exchange_weak_explicit(
//                 &mux->state,
//                 &old,
//                 old + READER_INC,
//                 memory_order_seq_cst,
//                 memory_order_relaxed)) {
//             return;
//         }
//     }
// }

// void __public__mutex_runlock(__public__rwmutex_t* mux) {
//     int prev = atomic_fetch_sub_explicit(
//         &mux->state,
//         READER_INC,
//         memory_order_release
//     );

//     // Last reader and a writer waiting
//     if ((prev - READER_INC) == WRITER_BIT) {
//         pthread_mutex_lock(&mux->lock);
//         task_t *w = safe_q_pop(&mux->writers);
//         if (w) {
//             safe_q_push(&kernel_thread_map[current_task->sched_id]->ready_q, w);
//         }
//         pthread_mutex_unlock(&mux->lock);
//     }
// }

// void __public__mutex_rwlock(__public__rwmutex_t* mux) {
//     int old;

//     for (;;) {
//         old = atomic_load_explicit(&mux->state, memory_order_seq_cst);

//         if (old & WRITER_BIT) {
//             pthread_mutex_lock(&mux->lock);
//             if (atomic_load_explicit(&mux->state, memory_order_seq_cst) & WRITER_BIT) {
//                 safe_q_push(&mux->writers, current_task);
//                 pthread_mutex_unlock(&mux->lock);
//                 task_yield(kernel_thread_map[current_task->sched_id]);
//                 continue;
//             }
//             pthread_mutex_unlock(&mux->lock);
//             continue;
//         }

//         if (atomic_compare_exchange_weak_explicit(
//                 &mux->state,
//                 &old,
//                 old | WRITER_BIT,
//                 memory_order_seq_cst,
//                 memory_order_relaxed)) {
//             break;
//         }
//     }

//     while (atomic_load_explicit(&mux->state, memory_order_seq_cst) != WRITER_BIT) {
//         task_yield(kernel_thread_map[current_task->sched_id]);
//     }
// }

// void __public__mutex_rwunlock(__public__rwmutex_t* mux) {
//     pthread_mutex_lock(&mux->lock);

//     // Clear the WRITER_BIT
//     atomic_store_explicit(&mux->state, 0, memory_order_release);

//     // Priority: wake a waiting writer first
//     task_t *w = safe_q_pop(&mux->writers);
//     if (w) {
//         safe_q_push(&kernel_thread_map[current_task->sched_id]->ready_q, w);
//         pthread_mutex_unlock(&mux->lock);
//         return;
//     }

//     // Otherwise wake all waiting readers
//     task_t *r;
//     while ((r = safe_q_pop(&mux->readers)) != NULL) {
//         safe_q_push(&kernel_thread_map[current_task->sched_id]->ready_q, r);
//     }

//     pthread_mutex_unlock(&mux->lock);
// }


// __public__mutex_t* __public__create_mutex() {
//     __public__mutex_t* mux = (__public__mutex_t*)allocate(__arena__, sizeof(__public__mutex_t));
//     pthread_mutex_init(&mux->lock, 0);
//     atomic_store_explicit(&mux->state, 0, memory_order_relaxed);

//     return mux;
// }

// void __public__mutex_lock(__public__mutex_t* mtx) {
//     int expected = 0;

//     // Fast path: try to grab lock immediately
//     if (atomic_compare_exchange_strong_explicit(
//             &mtx->state,
//             &expected,
//             1,
//             memory_order_seq_cst,
//             memory_order_relaxed)) {
//         return; // acquired lock
//     }

//     // Slow path: lock is held, must enqueue and sleep
//     for (;;) {
//         // Mark ourselves waiting and sleep atomically
//         pthread_mutex_lock(&mtx->lock);
//         expected = atomic_load_explicit(&mtx->state, memory_order_seq_cst);
//         if (expected == 0) {
//             // Lock became free while we acquired mtx->lock
//             atomic_store_explicit(&mtx->state, 1, memory_order_seq_cst);
//             pthread_mutex_unlock(&mtx->lock);
//             return;
//         }

//         // Enqueue and yield to scheduler
//         safe_q_push(&mtx->waiters, current_task);
//         pthread_mutex_unlock(&mtx->lock);
//         task_yield(kernel_thread_map[current_task->sched_id]);
//     }
// }

// void __public__mutex_unlock(__public__mutex_t* mtx) {
//     pthread_mutex_lock(&mtx->lock);

//     // Mark the mutex as free
//     atomic_store_explicit(&mtx->state, 0, memory_order_release);

//     // Wake one waiting task, if any
//     task_t *t = safe_q_pop(&mtx->waiters);
//     if (t) {
//         // Give the lock to this task immediately
//         atomic_store_explicit(&mtx->state, 1, memory_order_release);
//         safe_q_push(&kernel_thread_map[t->sched_id]->ready_q, t);
//     }

//     pthread_mutex_unlock(&mtx->lock);
// }
