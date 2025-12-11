#include "array.h"
#include "alloc.h"
#include <stdlib.h>
/**
 * @brief allocate block of memory for array through GC_MALLOC 
 *
 * @param count length of array
 * @param elem_size size of each element
 */
extern __thread arena_t* __arena__;

Array* lang_alloc_array(int count, int elem_size) {
    Array* arr = (Array*)allocate(__arena__, sizeof(Array));
    arr->length = count;
    arr->data = allocate(__arena__, count*elem_size); /* should init with zero */
    return arr;
}