#ifndef PLATFORM_CONTEXT_H
#define PLATFORM_CONTEXT_H

#include <stdint.h>
#include <stddef.h>

typedef struct platform_ctx {
#if defined(PLATFORM_UCONTEXT)
    void *uc;          // ucontext_t*
#else
    uintptr_t sp;      // saved stack pointer (asm)
#endif
} platform_ctx_t;

int  platform_ctx_init(platform_ctx_t *ctx);

int  platform_ctx_make(platform_ctx_t *ctx,
                       void (*entry)(uintptr_t, uintptr_t),
                       uintptr_t a,
                       uintptr_t b,
                       void *stack,
                       size_t stack_size);

void platform_ctx_switch(platform_ctx_t *from,
                         platform_ctx_t *to);

#endif
