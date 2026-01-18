#ifndef STR_H
#define STR_H
#include "platform.h"
char* __public__strings_format(const char* fmt, ...);

int __public__strings_length(const char* str);

int __public__strings_compare(const char* str1, const char* str2);

#endif