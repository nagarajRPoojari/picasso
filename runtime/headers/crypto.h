#ifndef CRYPTO_H
#define CRYPTO_H

#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdint.h>
/**
 * @brief returns hash of given data
 * 
 * @param data data, e.g, struct pointer, string
 * @param len length of datatype
 */
uint64_t hash(const void *data, size_t len);
#endif