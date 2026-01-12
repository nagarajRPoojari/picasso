#include "platform.h"
#include "array.h"
#include "alloc.h"
#include <stdlib.h>
#include <string.h>
#include <stddef.h>
#include <stdio.h>

/**
 * @brief allocate block of memory for array through GC_MALLOC 
 *
 * @param count length of array
 * @param elem_size size of each element
 */
extern __thread arena_t* __arena__;

__public__array_t* __public__alloc_array(int count, int elem_size, int rank) {
    size_t data_size = (size_t)count * elem_size;
    size_t shape_size = (size_t)rank * sizeof(int64_t);
    size_t total_size = sizeof(__public__array_t) + data_size + shape_size;

    __public__array_t* arr = (__public__array_t*)allocate(__arena__, total_size);

    
    arr->data = (int8_t*)(arr + 1); 
    
    if (rank > 0) {
        arr->shape = (int64_t*)(arr->data + data_size);
    } else {
        arr->shape = NULL;
    }
    
    arr->length = count;
    arr->rank = rank;
    
    memset(arr->data, 0, data_size); 
    return arr;
}

int64_t __public__len(__public__array_t* arr) {
    return arr->length;
}

/** @deprecated */
void __public__debug_array_info(__public__array_t* arr) {
    if (arr == NULL) {
        printf(" == DEBUG INFO: __public__array_t pointer is NULL ==\n");
        return;
    }

    printf(" == DEBUG INFO: Final __public__array_t State ==\n");
    printf("     Base Address (a.Ptr): %p\n", (void*)arr);
    printf("     data (Offset 0):      %p\n", arr->data);
    printf("     shape (Offset 8):     %p\n", arr->shape);
    printf("     length (Offset 16):   %lld\n", arr->length);
    printf("     rank (Offset 24):     %lld\n", arr->rank);
    
    // Print shape elements if rank > 0
    if (arr->rank > 0 && arr->shape != NULL) {
        printf("     Shape Dims: [");
        for (int i = 0; i < arr->rank; i++) {
            printf("%lld", arr->shape[i]);
            if (i < arr->rank - 1) {
                printf(", ");
            }
        }
        printf("]\n");
    } else {
        printf("     Shape Dims: (None or NULL)\n");
    }
}