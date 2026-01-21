#ifndef NETPOLL_H
#define NETPOLL_H

#include <stdint.h>
#include <stddef.h>

#if defined(__linux__) 
    #include <sys/epoll.h>
    
    typedef struct netpoll {
        int epfd;
    } netpoll_t;

#elif defined(__APPLE__)
    #include <sys/event.h> // macOS uses kqueue instead of epoll

    typedef struct netpoll {
        int kq;
    } netpoll_t;
#endif


typedef void* netpoll_ud_t;

typedef enum {
    NETPOLL_IN   = 1 << 0,   /* readable */
    NETPOLL_OUT  = 1 << 1,   /* writable */
    NETPOLL_ERR  = 1 << 2,   /* error */
    NETPOLL_HUP  = 1 << 3,   /* hangup / eof */
    NETPOLL_ONESHOT = 1 << 4
} netpoll_eventflag_t;

typedef struct {
    int              fd;
    netpoll_eventflag_t events;
    netpoll_ud_t     ud;     /* task_t* */
} netpoll_event_t;

netpoll_t* netpoll_create(void);
void       netpoll_destroy(netpoll_t* p);

int netpoll_add( netpoll_t* p, int fd, netpoll_eventflag_t events, netpoll_ud_t ud );

int netpoll_mod( netpoll_t* p, int fd, netpoll_eventflag_t events, netpoll_ud_t ud);

int netpoll_del( netpoll_t* p, int fd );

int netpoll_wait( netpoll_t* p, netpoll_event_t* events, int max_events, int timeout_ms);

#endif 
