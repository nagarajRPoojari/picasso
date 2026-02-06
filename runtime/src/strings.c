#include "platform.h"

#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <pthread.h>

#include "str.h"
#include "alloc.h"

extern __thread arena_t* __arena__;

/**
 * @brief Format a string with variable arguments
 * @param fmt Format string
 * @param ... Variable arguments to format
 * @return Pointer to the formatted string
 */
char* __public__strings_format(const char* fmt, ...) {
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

/**
 * @brief Get the length of a string
 * @param str String to measure
 * @return Length of the string
 */
int __public__strings_length(const char* str) {
    return strlen(str);
}

/**
 * @brief Compare two strings
 * @param str1 First string to compare
 * @param str2 Second string to compare
 * @return 0 if equal, negative if str1 < str2, positive if str1 > str2
 */
int __public__strings_compare(const char* str1, const char* str2) {
    return strcmp(str1, str2);
}
