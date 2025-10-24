#ifndef IO_H
#define IO_H

#include "globals.h"
#include "queue.h"
#include "task.h"

// io_worker is the main loop which pops
// blocking io task & waits for io
void *io_worker(void *arg);

// utility io tasks
void* async_file_read(int, char*, int);
void _async_file_read();

void* async_stdin_read();
void _async_stdin_read();

#endif