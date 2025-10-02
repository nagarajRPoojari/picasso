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

typedef struct {
    long length;
    void* data;
} Array;

// Allocates an array of `count` elements of size `elem_size`
Array* lang_alloc_array(long count, long elem_size) {
    Array* arr = GC_MALLOC(sizeof(Array));
    arr->length = count;
    arr->data = GC_MALLOC(count * elem_size);
    return arr;
}


// Allocate array of count elements, each of elem_size bytes
void* lang_alloc_array(long count, long elem_size) {
    return GC_MALLOC(count * elem_size);
}
