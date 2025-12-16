#include "sigerr.h"
#include <stdio.h>
#include <stdlib.h>
#include <signal.h>
#include <libunwind.h>
#include <string.h>
#include <pthread.h>
#include <unistd.h>
#include <sys/mman.h>


static void print_stacktrace(void) {
    unw_cursor_t cursor;
    unw_context_t context;

    flockfile(stderr);

    if(unw_getcontext(&context) < 0) {
        fprintf(stderr, "unw_getcontext failed\n");
        funlockfile(stderr);
        return;
    }
    
    if (unw_init_local(&cursor, &context) < 0) {
        fprintf(stderr, "unw_init_local failed\n");
        funlockfile(stderr);
        return;
    }
    
    // Skip the first 2 frames (crash_handler and signal trampoline)
    for (int i = 0; i < 2; i++) {
        if (unw_step(&cursor) <= 0) {
            fprintf(stderr, "Stack too shallow\n");
            funlockfile(stderr);
            return;
        }
    }
    
    fprintf(stderr, "===== STACK TRACE (Thread: %lu) =====\n", pthread_self());
    int frame = 0;
    while (unw_step(&cursor) > 0) {
        unw_word_t ip, sp;
        char func_name[256];
        unw_get_reg(&cursor, UNW_REG_IP, &ip);
        unw_get_reg(&cursor, UNW_REG_SP, &sp);
        
        if (unw_get_proc_name(&cursor, func_name, sizeof(func_name), NULL) == 0) {
            fprintf(stderr, "\t#%02d  %p  %s (SP=%p)\n", frame++, (void*)ip, func_name, (void*)sp);
        } else {
            fprintf(stderr, "\t#%02d  %p  <unknown> (SP=%p)\n", frame++, (void*)ip, (void*)sp);
        }
    }
    fprintf(stderr, "=======================\n");
    funlockfile(stderr);
}

static void crash_handler(int sig, siginfo_t *info, void *ucontext) {
    (void)info;
    (void)ucontext;

    /* unsafe but acceptable for crash diagnostics */
    print_stacktrace();

    _exit(128 + sig);
}


static void *sigstack_mem = NULL;

static void setup_altstack(void) {
    size_t size = SIGSTKSZ * 2;

    sigstack_mem = mmap(NULL, size,
                        PROT_READ | PROT_WRITE,
                        MAP_PRIVATE | MAP_ANONYMOUS | MAP_STACK,
                        -1, 0);
    if (sigstack_mem == MAP_FAILED) {
        perror("mmap sigaltstack");
        _exit(1);
    }

    stack_t ss;
    ss.ss_sp = sigstack_mem;
    ss.ss_size = size;
    ss.ss_flags = 0;

    if (sigaltstack(&ss, NULL) != 0) {
        perror("sigaltstack");
        _exit(1);
    }
}

/**
 * @brief raises runtime error
 * 
 * @param msg message to be printed in error
 */
void __public__runtime_error(const char* msg) {
    fprintf(stderr, "%s", msg);
    print_stacktrace();

    exit(1);
}


/**
 * @brief registers error handlers for common signals.
 */
void init_error_handlers(void) {
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));

    sa.sa_sigaction = crash_handler;
    sigemptyset(&sa.sa_mask);
    sa.sa_flags = SA_SIGINFO | SA_ONSTACK | SA_RESTART;

    int sigs[] = { SIGSEGV, SIGBUS, SIGILL, SIGFPE, SIGABRT };
    for (size_t i = 0; i < sizeof(sigs)/sizeof(sigs[0]); i++) {
        if (sigaction(sigs[i], &sa, NULL) != 0) {
            perror("sigaction");
            _exit(1);
        }
    }
}
