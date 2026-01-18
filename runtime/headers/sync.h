#ifndef SYNC_H
#define SYNC_H

#include "platform.h"
#include <stdalign.h>

#include "stdatomic.h"
#include "stdint.h"
#include "scheduler.h"
#include "queue.h"

#define WRITER_BIT   1
#define READER_INC   2

typedef struct __public__rwmutex {
    _Atomic int64_t state;
    pthread_mutex_t lock;      // per-mutex lock to protect enqueue + sleep
    safe_queue_t readers;
    safe_queue_t writers;
}__public__rwmutex_t;


__public__rwmutex_t* __public__sync_rwmutex_create();
void __public__sync_rwmutex_rlock(__public__rwmutex_t* mux);
void __public__sync_rwmutex_rwlock(__public__rwmutex_t* mux);
void __public__sync_rwmutex_runlock(__public__rwmutex_t* mux);
void __public__sync_rwmutex_rwunlock(__public__rwmutex_t* mux);



typedef struct __public__mutex {
    _Atomic int64_t state;
    pthread_mutex_t lock;
    safe_queue_t waiters;
}__public__mutex_t;

__public__mutex_t* __public__sync_mutex_create(void);
void __public__sync_mutex_lock(__public__mutex_t* mtx);
void __public__sync_mutex_unlock(__public__mutex_t* mtx);

#endif