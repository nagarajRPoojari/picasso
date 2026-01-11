#include "xy.h"

/* Force references so Clang emits declarations */
void* __ffi_force[] = {
    (void*)create_coordinate,
    (void*)print_coordinate,
    (void*)destroy_coordinate,
};
