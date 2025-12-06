#include "array.h"
#include <stdlib.h>
/**
 * @brief allocate block of memory for array through GC_MALLOC 
 *
 * @param count length of array
 * @param elem_size size of each element
 */
Array* lang_alloc_array(int count, int elem_size) {
    Array* arr = (Array*)malloc(sizeof(Array));
    arr->length = count;
    arr->data = calloc(count, elem_size); /* should init with zero */
    return arr;
}