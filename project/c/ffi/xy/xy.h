#ifndef MATH_H
#define MATH_H

#include <stdint.h>
#include <stdlib.h>

typedef struct Coordinate {
    int64_t x;
    int64_t y;
} Coordinate;

Coordinate* create_coordinate(int64_t x, int64_t y);
void destroy_coordinate(Coordinate* c);
void print_coordinate(Coordinate* c);
#endif