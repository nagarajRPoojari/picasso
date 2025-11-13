#include "error.h"
#include <stdio.h>
#include <stdlib.h>

/**
 * @brief raises runtime error
 * 
 * @param msg message to be printed in error
 */
void runtime_error(const char* msg) {
    fprintf(stderr, "%s\n", msg);
    exit(1);
}
