#ifndef GLOBALS_H
#define GLOBALS_H

#include <ucontext.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <pthread.h>
#include <sys/epoll.h>

#include "queue.h"
#include "task.h"
#include "scheduler.h"


// scheduler.h
#define STACK_SIZE (4*1024)
#define SCHEDULER_THREAD_POOL_SIZE 1
#define SCHEDULER_LOCAL_QUEUE_SIZE 256
#define GUARD_SIZE (4096) 
extern __thread task_t* current_task;

// io.h
#define IO_THREAD_POOL_SIZE 4
#define MAX_EVENTS 16
#define IO_QUEUE_SIZE 256
extern safe_queue_t io_queue;
extern int epfd;

#endif 