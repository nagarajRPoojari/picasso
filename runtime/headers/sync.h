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

/**
 * @brief Create a new read-write mutex
 * @return Pointer to the created read-write mutex
 */
__public__rwmutex_t* __public__sync_rwmutex_create();

/**
 * @brief Acquire a read lock on the read-write mutex
 * @param mux Pointer to the read-write mutex
 */
void __public__sync_rwmutex_rlock(__public__rwmutex_t* mux);

/**
 * @brief Acquire a write lock on the read-write mutex
 * @param mux Pointer to the read-write mutex
 */
void __public__sync_rwmutex_rwlock(__public__rwmutex_t* mux);

/**
 * @brief Release a read lock on the read-write mutex
 * @param mux Pointer to the read-write mutex
 */
void __public__sync_rwmutex_runlock(__public__rwmutex_t* mux);

/**
 * @brief Release a write lock on the read-write mutex
 * @param mux Pointer to the read-write mutex
 */
void __public__sync_rwmutex_rwunlock(__public__rwmutex_t* mux);



typedef struct __public__mutex {
    _Atomic int64_t state;
    pthread_mutex_t lock;
    safe_queue_t waiters;
}__public__mutex_t;

/**
 * @brief Create a new mutex
 * @return Pointer to the created mutex
 */
__public__mutex_t* __public__sync_mutex_create(void);

/**
 * @brief Acquire a lock on the mutex
 * @param mtx Pointer to the mutex
 */
void __public__sync_mutex_lock(__public__mutex_t* mtx);

/**
 * @brief Release a lock on the mutex
 * @param mtx Pointer to the mutex
 */
void __public__sync_mutex_unlock(__public__mutex_t* mtx);

#endif