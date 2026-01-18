#ifndef PLATFORM_H
#define PLATFORM_H

#if defined(__linux__)

#ifndef _GNU_SOURCE
#define _GNU_SOURCE 1
#endif

#include <signal.h>
#include <ucontext.h>


#elif defined(__APPLE__)
#include <sys/ucontext.h>
#include <signal.h>
#else
#error Unsupported platform
#endif

#endif
