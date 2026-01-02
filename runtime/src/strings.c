#include "str.h"
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <pthread.h>
#include "alloc.h"

extern __thread arena_t* __arena__;

char* __public__format(const char* fmt, ...) {
    va_list args;
    va_start(args, fmt);

    va_list args_copy;
    va_copy(args_copy, args);

    // First pass: compute length
    int len = vsnprintf(NULL, 0, fmt, args_copy);
    va_end(args_copy);

    if (len < 0) {
        va_end(args);
        return NULL;
    }

    char* buf = allocate(__arena__, len + 1);
    if (!buf) {
        va_end(args);
        return NULL;
    }

    // Second pass: actual formatting
    vsnprintf(buf, len + 1, fmt, args);
    va_end(args);

    return buf;
}


int __public__len(const char* str) {
    return strlen(str);
}

int __public__compare(const char* str1, const char* str2) {
    return strcmp(str1, str2);
}
