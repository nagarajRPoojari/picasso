#ifndef PLATFORM_H
#define PLATFORM_H

#if defined(__linux__)
#include <signal.h>
#include <ucontext.h>
#elif defined(__APPLE__)
#include <sys/ucontext.h>
#include <signal.h>
#else
#error Unsupported platform
#endif

#endif
