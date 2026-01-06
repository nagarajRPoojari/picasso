#ifndef SYNC_H
#define SYNC_H

// #include "scheduler.h"
// #include "queue.h"

// #define WRITER_BIT   1
// #define READER_INC   2

// typedef struct __public__rwmutex_t {
//     atomic_int state;
//     pthread_mutex_t lock;      // per-mutex lock to protect enqueue + sleep
//     safe_queue_t readers;
//     safe_queue_t writers;
// } __public__rwmutex_t;


// __public__rwmutex_t* __public__create_rwmutex();

// void __public__mutex_rlock(__public__rwmutex_t* mux);

// void __public__mutex_rwlock(__public__rwmutex_t* mux);

// void __public__mutex_runlock(__public__rwmutex_t* mux);

// void __public__mutex_rwunlock(__public__rwmutex_t* mux);


// typedef struct __public__mutex_t {
//     atomic_int state;       // 0 = unlocked, 1 = locked
//     pthread_mutex_t lock;   // per-mutex lock for enqueue + sleep
//     safe_queue_t waiters;   // queue of waiting tasks
// } __public__mutex_t;


// __public__mutex_t* __public__create_mutex();

// void __public__mutex_lock(__public__mutex_t* mtx);

// void __public__mutex_unlock(__public__mutex_t* mtx);

#endif