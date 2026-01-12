#include "coo.h"
#include "stdlib.h"
#include "stdio.h"

Coordinate* __public__coo_create(int64_t x, int64_t y) {
    Coordinate* c = (Coordinate*)malloc(sizeof(Coordinate));
    c->x = x;
    c->y = y;

    return c;
}

void __public__coo_dump(Coordinate* c) {
    printf("Coordinate(%d, %d)\n", c->x, c->y);
}