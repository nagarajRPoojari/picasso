#ifndef COO_H
#define COO_H

#include "stdint.h"

const int64_t __public__coo_nagaraj = 100;

typedef struct Coordinate {
    /* data */
    int64_t x;
    int64_t y;
} Coordinate;

extern const int64_t nagaraj;

Coordinate* __public__coo_create(int64_t x, int64_t y);
void __public__coo_dump(Coordinate* c);
#endif