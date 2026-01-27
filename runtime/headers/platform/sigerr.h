#ifndef PLATFORM_CONTEXT_H
#define PLATFORM_CONTEXT_H

/* Install fatal signal handlers */
void init_error_handlers(void);

/* Raise a runtime error and print stack trace */
void __public__runtime_error(const char *msg);

#endif