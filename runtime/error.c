#include "error.h"
#include <stdio.h>

void runtime_error(const char* msg) {
    fprintf(stderr, "%s\n", msg);
    exit(1);
}
