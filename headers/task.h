#ifndef TASK_H
#define TASK_H

#include <ucontext.h>
#include <signal.h>
#define TASK_IO_BUFFER 256

typedef struct task {

    // id to uniquely identify task
    int id;

    // own CPU context
    ucontext_t ctx;

    size_t stack_size;
    // sched_id
    int sched_id;

    // function to executue
    void* (*fn)(void *);

    // private stack
    char *stack;

    // file descriptor if doing any io
    int fd;

    // buffer pool for io
    char *buf;
    ssize_t readn;

    ssize_t nread;
    int use_epoll; 
} task_t;

#endif