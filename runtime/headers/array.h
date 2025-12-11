#ifndef ARRAY_H
#define ARRAY_H

/**
 * @struct Array
 */
typedef struct {
    int length;
    void* data;
    int* shape; /** int64 block of memory, e.g [3,3,4] */
    int rank;   /** dimension of the array */
} Array;

/**
 * @brief allocate block of memory for array through custom allocator
 *
 * @param count length of array
 * @param elem_size size of each element
 */
Array* __public__alloc_array(int count, int elem_size);
#endif