#include "str.h"
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>

char* format(const char* fmt, ...) {
    va_list args;
    va_start(args, fmt);

    // Find required buffer size
    int len = vsnprintf(NULL, 0, fmt, args);
    va_end(args);
    if (len < 0) return NULL;

    // Allocate memory
    char* buf = malloc(len + 1);
    if (!buf) return NULL;

    // Format the string
    va_start(args, fmt);
    vsnprintf(buf, len + 1, fmt, args);
    va_end(args);

    return buf;
}

int len(const char* str) {
    return strlen(str);
}

int compare(const char* str1, const char* str2) {
    return strcmp(str1, str2);
}
