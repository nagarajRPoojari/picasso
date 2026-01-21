#ifndef PLATFORM_CONTEXT_H
#define PLATFORM_CONTEXT_H

#include "platform.h"
#include <stdint.h>
#include <stddef.h>
#if defined(__linux__)
    #include <ucontext.h>
#endif

#if defined(__APPLE__)
    typedef struct registers_t {
        // x19-28
        uint64_t regs[12];
        uintptr_t sp;
    } registers_t;
#endif


typedef struct platform_ctx {
#if defined(__linux__)
    ucontext_t *uc;          // ucontext_t*
#elif defined(__APPLE__)
    registers_t reg;
    void* stack;
    size_t stack_size;
    struct platform_ctx* back_link;
#endif
} platform_ctx_t;

int  platform_ctx_init(platform_ctx_t *ctx);

int  platform_ctx_make(platform_ctx_t *ctx, void (*entry)(uintptr_t, uintptr_t), uintptr_t a, uintptr_t b, void *stack, size_t stack_size, platform_ctx_t* back_link);

void platform_ctx_switch(platform_ctx_t *from, platform_ctx_t *to);

void plarform_ctx_destroy(platform_ctx_t* ctx);

uintptr_t platform_ctx_get_stack_base(platform_ctx_t* ctx);

uintptr_t platform_ctx_get_stack_pointer(platform_ctx_t* ctx);

uintptr_t platform_ctx_get_pc(platform_ctx_t* ctx);

int64_t platform_ctx_get_stack_size(platform_ctx_t* ctx);

uint64_t platform_ctx_get_reg(platform_ctx_t *ctx, int i) ;

#endif
