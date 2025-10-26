#ifndef IO_H
#define IO_H

#include "queue.h"
#include "task.h"


/* Number of I/O worker threads in the pool */
#define IO_THREAD_POOL_SIZE 4

/* Maximum number of events returned by epoll_wait at once */
#define MAX_EVENTS 16

/* Maximum number of tasks in the I/O queue */
#define IO_QUEUE_SIZE 256


/**
 * Thread-safe queue for pending I/O tasks 
 * @owner: io.c
 */
extern safe_queue_t io_queue;

/**
 * epoll file descriptor used by all I/O workers
 * @owner: io.c
 */
extern int epfd;

/** @todo: move this a userlib or somewhere else */
/**
 * @brief Main loop for an I/O worker thread.
 * 
 * Continuously pops tasks from the io_queue, performs blocking I/O,
 * and waits for I/O readiness using epoll. Reschedules tasks when ready.
 * 
 * @param arg Pointer to worker thread ID or context (implementation-specific)
 * @return Never returns under normal operation.
 */
void *io_worker(void *arg);


/**
 * @brief Initiate an asynchronous file read task.
 * 
 * Schedules a read operation on the I/O queue. The actual read happens
 * in the I/O worker thread.
 * 
 * @param fd     File descriptor to read from.
 * @param buf    Buffer to store data.
 * @param count  Number of bytes to read.
 * @return Pointer to a task representing this async operation.
 */
void* async_file_read(int fd, char *buf, int count);

/* Internal helper for file read (used by I/O worker) */
void _async_file_read();

/**
 * @brief Initiate an asynchronous read from stdin.
 * 
 * Schedules a read operation on stdin via the I/O queue.
 * 
 * @return Pointer to a task representing this async operation.
 */
void* async_stdin_read();

/* Internal helper for stdin read (used by I/O worker) */
void _async_stdin_read();

#endif