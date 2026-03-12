#include "coo.h"
#include "stdlib.h"
#include "stdio.h"

// Original functions - struct pointer based
Coordinate* __public__coo_create(int64_t x, int64_t y) {
    Coordinate* c = (Coordinate*)malloc(sizeof(Coordinate));
    c->x = x;
    c->y = y;

    return c;
}

void __public__coo_dump(Coordinate* c) {
    printf("Coordinate(%lld, %lld)\n", (long long)c->x, (long long)c->y);
}

// New function - returns bare struct by value
Coordinate __public__coo_create_value(int64_t x, int64_t y) {
    Coordinate c;
    c.x = x;
    c.y = y;
    return c;
}

// Simulated C API that uses pointer parameters for "out" values
void __public__coo_get_dimensions(int64_t* width, int64_t* height) {
    *width = 1920;
    *height = 1080;
}

// Wrapper function that converts pointer-based API to struct return
Dimensions __public__coo_get_dimensions_wrapper() {
    Dimensions d;
    __public__coo_get_dimensions(&d.width, &d.height);
    return d;
}