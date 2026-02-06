#ifndef PLATFORM_CONTEXT_H
#define PLATFORM_CONTEXT_H

/**
 * @brief Install fatal signal handlers
 */
void init_error_handlers(void);

/**
 * @brief Raise a runtime error and print stack trace
 * @param msg Error message to display
 */
void __public__runtime_error(const char *msg);

#endif