#ifndef COO_H
#define COO_H

#include "stdint.h"

typedef struct Coordinate {
    /* data */
    int64_t x;
    int64_t y;
} Coordinate;

Coordinate* __public__coo_create(int64_t x, int64_t y);
void __public__coo_dump(Coordinate* c);

Coordinate __public__coo_create_value(int64_t x, int64_t y);
void __public__coo_get_dimensions(int64_t* width, int64_t* height);

typedef struct Dimensions {
    int64_t width;
    int64_t height;
} Dimensions;

Dimensions __public__coo_get_dimensions_wrapper();

#endif
