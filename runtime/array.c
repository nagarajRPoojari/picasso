#include "array.h"
#include <gc.h>

/**
 * @brief allocate block of memory for array through GC_MALLOC 
 *
 * @param count length of array
 * @param elem_size size of each element
 */
Array* lang_alloc_array(int count, int elem_size) {
    Array* arr = GC_MALLOC(sizeof(Array));
    arr->length = count;
    arr->data = GC_MALLOC(count * elem_size);
    return arr;
}