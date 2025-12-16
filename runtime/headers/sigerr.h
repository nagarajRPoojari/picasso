#ifndef SIGERR_H
#define SIGERR_H

/**
 * @brief raises runtime error
 * 
 * @param msg message to be printed in error
 */
void __public__runtime_error(const char* msg);

/**
 * @brief registers error handlers for common signals.
 */
void init_error_handlers(void);

#endif