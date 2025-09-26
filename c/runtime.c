#include <gc.h>
#include <stdio.h>
#include <stdlib.h>

// Initialize GC at program startup
void runtime_init() {
    GC_INIT();
}

// Allocate memory managed by GC
void *lang_alloc(long size) {
    return GC_MALLOC(size);
}

// Allocate memory without scanning (useful for raw byte buffers/strings)
void *lang_alloc_atomic(long size) {
    return GC_MALLOC_ATOMIC(size);
}

// Debug helper
void runtime_collect() {
    GC_gcollect();
}

void runtime_error(const char* msg) {
    fprintf(stderr, "%s\n", msg);
    exit(1);
}