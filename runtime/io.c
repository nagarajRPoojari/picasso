#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>
#include <sys/epoll.h>

#include "io.h"
#include "queue.h"
#include "task.h"
#include "scheduler.h"

safe_queue_t io_queue;
int epfd;

void *io_worker(void *arg) {
    (void)arg;
    while (1) {
        task_t *t = safe_q_pop_wait(&io_queue);
        if (!t) continue;
        // @todo: IO could be read or write, genralize
        t->nread = read(t->fd, t->buf, t->readn);

        // once IO is done, push back to ready queue for 
        // scheduler to pick up
        // @todo: either push back to global queue or same 
        // local queue
        safe_q_push(&(kernel_thread_map[t->sched_id]->ready_q), t);
    }
    return NULL;
}

void _async_stdin_read() {
    struct epoll_event ev;
    ev.events = EPOLLIN | EPOLLET;
    ev.data.ptr = current_task;
    fcntl(current_task->fd, F_SETFL, O_NONBLOCK);

    if (epoll_ctl(epfd, EPOLL_CTL_ADD, current_task->fd, &ev) == -1) {
        if (errno != EEXIST) perror("epoll_ctl ADD");
    }
    // yield cooperatively
    task_yield(kernel_thread_map[current_task->sched_id]);
}

void _async_file_read() {
    safe_q_push(&io_queue, current_task);
    // yield cooperatively
    task_yield(kernel_thread_map[current_task->sched_id]);
}

void* async_file_read(int fd, char* buf, int n) {
    current_task->fd = fd;
    current_task->buf = buf;
    current_task->readn = n;
    _async_file_read();
    return &(current_task->nread);
}

void* async_stdin_read(char* buf, int n) {
    current_task->fd = STDIN_FILENO;
    current_task->buf = buf;
    current_task->readn = n;
    _async_stdin_read();
}
