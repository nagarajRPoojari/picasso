#include "context.h"
#include <string.h>
#include <stdlib.h>
#include <stdio.h>

extern void platform_ctx_entry(void);

void task_completed_callback(platform_ctx_t* current, platform_ctx_t* parent) {
    fflush(stdout);
    
    if (parent) {
        platform_ctx_switch(current, parent);
    } else {
        exit(0);
    }
    
    while(1) {}
}

int platform_ctx_init(platform_ctx_t *ctx) {
    if (!ctx) return -1;
    memset(ctx, 0, sizeof(platform_ctx_t));
    return 0;
}

int platform_ctx_make(platform_ctx_t *ctx, void (*entry)(uintptr_t, uintptr_t), 
                      uintptr_t a, uintptr_t b, void *stack, size_t stack_size, 
                      platform_ctx_t* back_link) {
    if (!ctx || !stack) return -1;
    
    uintptr_t top = (uintptr_t)stack + stack_size;
    top &= ~0x0FL; 
    top -= 128;
    
    memset(&ctx->reg, 0, sizeof(registers_t));
    ctx->stack = stack;
    ctx->stack_size = stack_size;
    ctx->back_link = back_link;
    
    // Store BOTH current context and back_link on stack
    // Stack layout from top: [current_ctx] [back_link] [rest of stack]
    top -= 16;
    *(platform_ctx_t**)top = ctx;        // Current context pointer
    top -= 16;
    *(platform_ctx_t**)top = back_link;  // Parent context pointer
    
    ctx->reg.regs[0] = (uintptr_t)entry;           // x19: entry function
    ctx->reg.regs[1] = a;                          // x20: arg a
    ctx->reg.regs[2] = b;                          // x21: arg b
    ctx->reg.regs[10] = 0;                         // x29: frame pointer
    ctx->reg.regs[11] = (uintptr_t)platform_ctx_entry; // x30: link register
    ctx->reg.sp = (uintptr_t)top;
    
    return 0;
}

void platform_ctx_destroy(platform_ctx_t* ctx) {
    if (ctx) memset(ctx, 0, sizeof(platform_ctx_t));
}

uintptr_t platform_ctx_get_stack_base(platform_ctx_t* ctx) {
    return (uintptr_t)ctx->stack;
}

uintptr_t platform_ctx_get_stack_pointer(platform_ctx_t* ctx) {
    return (uintptr_t)ctx->reg.sp;
}

uintptr_t platform_ctx_get_pc(platform_ctx_t* ctx) {
    return (uintptr_t)ctx->reg.regs[11];
}

int64_t platform_ctx_get_stack_size(platform_ctx_t* ctx) {
    return (int64_t)ctx->stack_size;
}

uint64_t platform_ctx_get_reg(platform_ctx_t *ctx, int i) {
    if (i >= 19 && i <= 30) return ctx->reg.regs[i - 19];
    if (i == 31) return (uint64_t)ctx->reg.sp;
    return 0; 
}