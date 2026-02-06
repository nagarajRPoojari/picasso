#ifndef ARRAY_H
#define ARRAY_H

#include "platform.h"
#include <stdint.h> 

typedef struct {
    char* data;
    int64_t* shape; 
    int64_t length; 
    int64_t rank;  
} __public__array_t;

/**
 * @brief allocate block of memory for array through custom allocator
 * @param count length of array
 * @param elem_size size of each element
 */
__public__array_t* __public__alloc_array(int count, int elem_size, int rank);


/**
 * @brief utility func to print array information
 * @param arr array struct instance
 */
void __public__debug_array_info(__public__array_t* arr);
#endif