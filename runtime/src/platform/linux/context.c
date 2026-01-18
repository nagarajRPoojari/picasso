#include "platform.h"
#include "platform/context.h"
#include <ucontext.h>
#include <stdlib.h>

int  platform_ctx_init(platform_ctx_t *ctx) {
    ctx->uc = (ucontext_t*)malloc(sizeof(ucontext_t));
    getcontext(ctx->uc);

    return 0;
}

int  platform_ctx_make(platform_ctx_t *ctx, void (*entry)(uintptr_t, uintptr_t), uintptr_t a, uintptr_t b, void *stack, size_t stack_size, platform_ctx_t* back_link) {
    ctx->uc->uc_stack.ss_sp = stack;
    ctx->uc->uc_stack.ss_size = stack_size;
    ctx->uc->uc_link = back_link->uc;
    makecontext(ctx->uc, (void(*)(void))entry, 2, a, b);

    return 0;
}

void platform_ctx_switch(platform_ctx_t *from, platform_ctx_t *to) {
    swapcontext(from->uc, to->uc);
}

void plarform_ctx_destroy(platform_ctx_t* ctx) {
    free(ctx->uc);
}

uintptr_t platform_ctx_get_stack_base(platform_ctx_t* ctx) {
    return (uintptr_t)ctx->uc->uc_stack.ss_sp;
}

uintptr_t platform_ctx_get_stack_pointer(platform_ctx_t* ctx) {
    return (uintptr_t)ctx->uc->uc_mcontext.sp;
}

uintptr_t platform_ctx_get_pc(platform_ctx_t* ctx) {
    return (uintptr_t)ctx->uc->uc_mcontext.pc;
}

int64_t platform_ctx_get_stack_size(
    platform_ctx_t* ctx) {
    return (uintptr_t)ctx->uc->uc_stack.ss_size;
}

uint64_t platform_ctx_get_reg(platform_ctx_t *ctx, int i) {
    return ctx->uc->uc_mcontext.regs[i];
}
