#ifndef ARRAY_H
#define ARRAY_H

typedef struct {
    int length;
    void* data;
    int* shape;
    int rank;
} Array;

Array* lang_alloc_array(int count, int elem_size);

#endif