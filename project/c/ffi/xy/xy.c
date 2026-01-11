#include "xy.h"
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

Coordinate* create_coordinate(int64_t x, int64_t y) {
    Coordinate* c = (Coordinate*)malloc(sizeof(Coordinate));
    c->x = x;
    c->y = y;
    return c;
}

void destroy_coordinate(Coordinate* c) {
    free(c);
}

void print_coordinate(Coordinate* c) {
    printf("(%d, %d)\n", c->x, c->y);
}
