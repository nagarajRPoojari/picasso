#include "str.h"

/* Force references so Clang emits declarations */
void* __ffi_force[] = {
    (void*)__public__strings_format,
    (void*)__public__strings_length,
    (void*)__public__strings_alloc,
    (void*)__public__strings_get,
    (void*)__public__strings_compare,
};
