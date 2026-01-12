#include "coo.h"

/* Force references so Clang emits declarations */
void* __ffi_force[] = {
    (void*)__public__coo_create,
    (void*)__public__coo_dump,
};

