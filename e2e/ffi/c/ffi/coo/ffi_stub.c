#include "coo.h"

/* Force references so Clang emits declarations */
void* __ffi_force[] = {
    (void*)__public__coo_create,
    (void*)__public__coo_dump,
    (void*)__public__coo_create_value,
    (void*)__public__coo_get_dimensions,
    (void*)__public__coo_get_dimensions_wrapper,
};
