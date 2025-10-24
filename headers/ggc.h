#ifndef GGC_H
#define GGC_H

// Initialize GC at program startup
void runtime_init();

// Allocate memory managed by GC
void *lang_alloc(long size) ;

// Allocate memory without scanning (useful for raw byte buffers/strings)
void *lang_alloc_atomic(long size);

// Debug helper
void runtime_collect();

#endif