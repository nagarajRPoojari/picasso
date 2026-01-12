#include "coo.h"

/* Force references so Clang emits declarations */
__attribute__((used))
void* __ffi_force[2] = {
    (void*)__public__coo_create,
    (void*)__public__coo_dump,
};

