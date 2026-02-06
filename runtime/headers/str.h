#ifndef STR_H
#define STR_H
#include "platform.h"


/**
 * @brief Format a string with variable arguments
 * @param fmt Format string
 * @param ... Variable arguments to format
 * @return Pointer to the formatted string
 */
char* __public__strings_format(const char* fmt, ...);

/**
 * @brief Get the length of a string
 * @param str String to measure
 * @return Length of the string
 */
int __public__strings_length(const char* str);

/**
 * @brief Compare two strings
 * @param str1 First string to compare
 * @param str2 Second string to compare
 * @return 0 if equal, negative if str1 < str2, positive if str1 > str2
 */
int __public__strings_compare(const char* str1, const char* str2);

#endif