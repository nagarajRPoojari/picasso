#ifndef STR_H
#define STR_H
#include "platform.h"
#include <stdint.h>

typedef struct {
    char* data;
    int64_t size; 
} __public__string_t;


__public__string_t* __public__strings_alloc(const char* fmt, size_t size);

/**
 * @brief Format a string with variable arguments
 * @param fmt Format string
 * @param ... Variable arguments to format
 * @return Pointer to the formatted string
 */
__public__string_t* __public__strings_format(__public__string_t* fmt, ...) ;

/**
 * @brief Get the length of a string
 * @param str String to measure
 * @return Length of the string
 */
int __public__strings_length(__public__string_t* str);

/**
 * @brief Get the ith byte of a string
 * @param str String to measure
 * @return Length of the string
 */
uint8_t __public__strings_get(__public__string_t* str, int i);

/**
 * @brief Compare two strings
 * @param str1 First string to compare
 * @param str2 Second string to compare
 * @return 0 if equal, negative if str1 < str2, positive if str1 > str2
 */
int __public__strings_compare(__public__string_t* str1, __public__string_t* str2);

/* utils */
void buf_append(char **buf, size_t *cap, size_t *len, const char *src, size_t n);

size_t u64_to_dec(uint64_t v, char tmp[32]);

size_t i64_to_dec(int64_t v, char tmp[32]);

size_t ptr_to_hex(const void *p, char tmp[32]);

size_t f64_to_dec(double v, char tmp[64]);

#endif