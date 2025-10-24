#include "array.h"
#include <gc.h>

Array* lang_alloc_array(int count, int elem_size) {
    Array* arr = GC_MALLOC(sizeof(Array));
    arr->length = count;
    arr->data = GC_MALLOC(count * elem_size);
    return arr;
}