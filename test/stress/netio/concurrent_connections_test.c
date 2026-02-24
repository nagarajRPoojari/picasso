#include "unity/unity.h"
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdatomic.h>
#include <time.h>
#include <string.h>
#include "alloc.h"
#include "gc.h"
#include "initutils.h"
#include "scheduler.h"
#include "netio.h"
#include "str.h"
#include "array.h"

extern __thread arena_t* __arena__;
extern arena_t* __global__arena__;
extern atomic_int task_count;
extern kernel_thread_t **kernel_thread_map;

void setUp(void) {}
void tearDown(void) {}

/* Stress levels - configurable via compile-time macros
 * Define one of: STRESS_LEVEL_1, STRESS_LEVEL_2, STRESS_LEVEL_3, STRESS_LEVEL_4
 * Or use default (STRESS_LEVEL_3)
 */
#if defined(STRESS_LEVEL_1)
    #define NUM_CONNECTIONS 10
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define NUM_CONNECTIONS 50
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define NUM_CONNECTIONS 2000
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define NUM_CONNECTIONS 1000
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define TEST_PORT 9876
#define MAX_TIMEOUT ((NUM_CONNECTIONS / 10) * 50 + 200)

static atomic_int connections_accepted;
static atomic_int connections_completed;
static atomic_int server_ready;

void client_task(void* arg) {
    int client_id = (int)(uintptr_t)arg;
    
    while (atomic_load(&server_ready) == 0) {
        struct timespec ts = {0, 1000000}; /* 1ms */
        nanosleep(&ts, NULL);
    }
    
    __public__string_t addr;
    addr.data = "127.0.0.1";
    addr.size = 9;
    
    ssize_t client_fd = __public__net_dial(&addr, TEST_PORT);
    if (client_fd < 0) {
        return;
    }
    
    char msg[64];
    snprintf(msg, sizeof(msg), "Client %d", client_id);
    
    __public__array_t send_buf;
    send_buf.data = msg;
    send_buf.length = strlen(msg);
    send_buf.capacity = strlen(msg);
    send_buf.elem_size = 1;
    
    ssize_t written = __public__net_write(client_fd, &send_buf, strlen(msg));
    
    char recv_buf[64] = {0};
    __public__array_t read_buf;
    read_buf.data = recv_buf;
    read_buf.length = 0;
    read_buf.capacity = sizeof(recv_buf);
    read_buf.elem_size = 1;
    
    ssize_t read_bytes = __public__net_read(client_fd, &read_buf, sizeof(recv_buf) - 1);
    
    close(client_fd);
    
    if (written > 0 && read_bytes > 0) {
        atomic_fetch_add(&connections_completed, 1);
    }
}

void connection_handler(void* arg) {
    int client_fd = (int)(uintptr_t)arg;
    
    char recv_buf[64] = {0};
    __public__array_t read_buf;
    read_buf.data = recv_buf;
    read_buf.length = 0;
    read_buf.capacity = sizeof(recv_buf);
    read_buf.elem_size = 1;
    
    ssize_t read_bytes = __public__net_read(client_fd, &read_buf, sizeof(recv_buf) - 1);
    
    if (read_bytes > 0) {
        /* Echo back the message */
        __public__array_t send_buf;
        send_buf.data = recv_buf;
        send_buf.length = read_bytes;
        send_buf.capacity = read_bytes;
        send_buf.elem_size = 1;
        
        __public__net_write(client_fd, &send_buf, read_bytes);
    }
    
    close(client_fd);
}

void server_task(void* arg) {
    __public__string_t addr;
    addr.data = "127.0.0.1";
    addr.size = 9;
    
    ssize_t listen_fd = __public__net_listen(
        &addr,
        TEST_PORT,
        128,        /* backlog */
        1,          /* close_on_exec */
        1,          /* reuse_addr */
        1,          /* reuse_port */
        0,          /* tcp_nodelay */
        0,          /* tcp_defer_accept */
        0,          /* tcp_fastopen */
        0,          /* keepalive */
        0,          /* rcvbuf */
        0,          /* sndbuf */
        0           /* ipv6_only */
    );
    
    if (listen_fd < 0) {
        return;
    }
    
    atomic_store(&server_ready, 1);
    
    /* Accept connections */
    for (int i = 0; i < NUM_CONNECTIONS; i++) {
        ssize_t client_fd = __public__net_accept(listen_fd);
        
        if (client_fd >= 0) {
            atomic_fetch_add(&connections_accepted, 1);
            thread(connection_handler, 1, (void*)(uintptr_t)client_fd);
        }
    }
    
    close(listen_fd);
}

void test_stress_concurrent_connections(void) {
    atomic_store(&connections_accepted, 0);
    atomic_store(&connections_completed, 0);
    atomic_store(&server_ready, 0);
    
    printf("Testing %d concurrent connections (Stress Level: %s)...\n", 
           NUM_CONNECTIONS, STRESS_LEVEL);
    
    thread(server_task, 0);
    
    struct timespec ts = {0, 50000000}; /* 50ms */
    nanosleep(&ts, NULL);
    
    for (int i = 0; i < NUM_CONNECTIONS; i++) {
        thread(client_task, 1, (void*)(uintptr_t)i);
    }
    
    ts.tv_sec = 0;
    ts.tv_nsec = 100000000; /* 100ms */
    int timeout = 0;
    
    while (atomic_load(&connections_completed) < NUM_CONNECTIONS && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int accepted = atomic_load(&connections_accepted);
    int completed = atomic_load(&connections_completed);
    
    printf("Accepted: %d/%d, Completed: %d/%d\n", 
           accepted, NUM_CONNECTIONS, completed, NUM_CONNECTIONS);
    
    /* Allow some margin for network timing issues */
    TEST_ASSERT_GREATER_OR_EQUAL(NUM_CONNECTIONS * 0.95, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_concurrent_connections);
    
    wait_for_schedulers();
    
    return UNITY_END();
}
