#ifndef ARRAY_H
#define ARRAY_H
#include <stdint.h> 

/**
 * @struct Array
 */

typedef struct {
    int8_t* data;
    int64_t* shape; 
    int64_t length; 
    int64_t rank;  
} Array;

/**
 * @brief allocate block of memory for array through custom allocator
 *
 * @param count length of array
 * @param elem_size size of each element
 */
Array* __public__alloc_array(int count, int elem_size, int rank);


/**
 * @brief utility func to print array information
 * 
 * @param arr array struct instance
 */
void __public__debug_array_info(Array* arr);
#endif