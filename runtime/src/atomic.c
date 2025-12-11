#include "atomic.h"
#include <stdio.h>

/** @bool: */

/**
 * @brief Atomically store a boolean value.
 * @param ptr Pointer to an atomic _Bool variable.
 * @param val Value to store.
 */
void __public__atomic_store_bool(_Atomic _Bool *ptr, _Bool val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load a boolean value.
 * @param ptr Pointer to an atomic _Bool variable.
 * @return The loaded boolean value.
 */
_Bool __public__atomic_load_bool(_Atomic _Bool *ptr) { 
    return atomic_load(ptr); 
}

/** @char: */

/**
 * @brief Atomically store a char value.
 * @param ptr Pointer to an atomic char variable.
 * @param val Value to store.
 */
void __public__atomic_store_char(_Atomic char *ptr, char val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load a char value.
 * @param ptr Pointer to an atomic char variable.
 * @return The loaded char value.
 */
char __public__atomic_load_char(_Atomic char *ptr) { 
    return atomic_load(ptr); 
}
/**
 * @brief Atomically add to a char variable and return the previous value.
 * @param ptr Pointer to an atomic char variable.
 * @param val Value to add.
 * @return The value of *ptr before the addition.
 */
char __public__atomic_add_char(_Atomic char *ptr, char val) { 
    return atomic_fetch_add(ptr, val); 
}
/**
 * @brief Atomically subtract from a char variable and return the previous value.
 * @param ptr Pointer to an atomic char variable.
 * @param val Value to subtract.
 * @return The value of *ptr before the subtraction.
 */
char __public__atomic_sub_char(_Atomic char *ptr, char val) { 
    return atomic_fetch_sub(ptr, val); 
}


/** @short: */

/**
 * @brief Atomically store a short value.
 * @param ptr Pointer to an atomic short variable.
 * @param val Value to store.
 */
void __public__atomic_store_short(_Atomic short *ptr, short val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load a short value.
 * @param ptr Pointer to an atomic short variable.
 * @return The loaded short value.
 */
short __public__atomic_load_short(_Atomic short *ptr) { 
    return atomic_load(ptr); 
}
/**
 * @brief Atomically add to a short variable and return the previous value.
 * @param ptr Pointer to an atomic short variable.
 * @param val Value to add.
 * @return The value of *ptr before the addition.
 */
short __public__atomic_add_short(_Atomic short *ptr, short val) { 
    return atomic_fetch_add(ptr, val); 
}
/**
 * @brief Atomically subtract from a short variable and return the previous value.
 * @param ptr Pointer to an atomic short variable.
 * @param val Value to subtract.
 * @return The value of *ptr before the subtraction.
 */
short __public__atomic_sub_short(_Atomic short *ptr, short val) { 
    return atomic_fetch_sub(ptr, val); 
}

/** @int: */

/**
 * @brief Atomically store an int value.
 * @param ptr Pointer to an atomic int variable.
 * @param val Value to store.
 */
void __public__atomic_store_int(_Atomic int *ptr, int val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load an int value.
 * @param ptr Pointer to an atomic int variable.
 * @return The loaded int value.
 */
int __public__atomic_load_int(_Atomic int *ptr) { 
    return atomic_load(ptr); 
}
/**
 * @brief Atomically add to an int variable and return the previous value.
 * @param ptr Pointer to an atomic int variable.
 * @param val Value to add.
 * @return The value of *ptr before the addition.
 */
int __public__atomic_add_int(_Atomic int *ptr, int val) { 
    return atomic_fetch_add(ptr, val); 
}
/**
 * @brief Atomically subtract from an int variable and return the previous value.
 * @param ptr Pointer to an atomic int variable.
 * @param val Value to subtract.
 * @return The value of *ptr before the subtraction.
 */
int __public__atomic_sub_int(_Atomic int *ptr, int val) { 
    return atomic_fetch_sub(ptr, val); 
}

/** @long: */

/**
 * @brief Atomically store a long value.
 * @param ptr Pointer to an atomic long variable.
 * @param val Value to store.
 */
void __public__atomic_store_long(_Atomic long *ptr, long val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load a long value.
 * @param ptr Pointer to an atomic long variable.
 * @return The loaded long value.
 */
long __public__atomic_load_long(_Atomic long *ptr) { 
    return atomic_load(ptr); 
}
/**
 * @brief Atomically add to a long variable and return the previous value.
 * @param ptr Pointer to an atomic long variable.
 * @param val Value to add.
 * @return The value of *ptr before the addition.
 */
long __public__atomic_add_long(_Atomic long *ptr, long val) { 
    return atomic_fetch_add(ptr, val); 
}
/**
 * @brief Atomically subtract from a long variable and return the previous value.
 * @param ptr Pointer to an atomic long variable.
 * @param val Value to subtract.
 * @return The value of *ptr before the subtraction.
 */
long __public__atomic_sub_long(_Atomic long *ptr, long val) { 
    return atomic_fetch_sub(ptr, val); 
}

/** @longlong: */

/**
 * @brief Atomically store a long long value.
 * @param ptr Pointer to an atomic long long variable.
 * @param val Value to store.
 */
void __public__atomic_store_llong(_Atomic long long *ptr, long long val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load a long long value.
 * @param ptr Pointer to an atomic long long variable.
 * @return The loaded long long value.
 */
long long __public__atomic_load_llong(_Atomic long long *ptr) { 
    return atomic_load(ptr); 
}
/**
 * @brief Atomically add to a long long variable and return the previous value.
 * @param ptr Pointer to an atomic long long variable.
 * @param val Value to add.
 * @return The value of *ptr before the addition.
 */
long long __public__atomic_add_llong(_Atomic long long *ptr, long long val) { 
    return atomic_fetch_add(ptr, val); 
}
/**
 * @brief Atomically subtract from a long long variable and return the previous value.
 * @param ptr Pointer to an atomic long long variable.
 * @param val Value to subtract.
 * @return The value of *ptr before the subtraction.
 */
long long __public__atomic_sub_llong(_Atomic long long *ptr, long long val) { 
    return atomic_fetch_sub(ptr, val); 
}

/** @float: */

/**
 * @brief Atomically store a float value.
 * @param ptr Pointer to an atomic float variable.
 * @param val Value to store.
 */
void __public__atomic_store_float(_Atomic float *ptr, float val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load a float value.
 * @param ptr Pointer to an atomic float variable.
 * @return The loaded float value.
 */
float __public__atomic_load_float(_Atomic float *ptr) { 
    return atomic_load(ptr); 
}

/** @double: */

/**
 * @brief Atomically store a double value.
 * @param ptr Pointer to an atomic double variable.
 * @param val Value to store.
 */
void __public__atomic_store_double(_Atomic double *ptr, double val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load a double value.
 * @param ptr Pointer to an atomic double variable.
 * @return The loaded double value.
 */
double __public__atomic_load_double(_Atomic double *ptr) { 
    return atomic_load(ptr); 
}

/** @pointer: */

/**
 * @brief Atomically store a pointer value.
 * @param ptr Pointer to an atomic pointer variable.
 * @param val Pointer value to store.
 */
void __public__atomic_store_ptr(_Atomic void **ptr, void *val) { 
    atomic_store(ptr, val); 
}
/**
 * @brief Atomically load a pointer value.
 * @param ptr Pointer to an atomic pointer variable.
 * @return The loaded pointer value.
 */
void *__public__atomic_load_ptr(_Atomic void **ptr) { 
    return atomic_load(ptr); 
}