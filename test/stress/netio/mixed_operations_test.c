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
    #define NUM_CLIENTS 10
    #define MESSAGES_PER_CLIENT 3
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define NUM_CLIENTS 25
    #define MESSAGES_PER_CLIENT 5
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define NUM_CLIENTS 100
    #define MESSAGES_PER_CLIENT 10
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define NUM_CLIENTS 50
    #define MESSAGES_PER_CLIENT 5
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define TEST_PORT 9879
#define MAX_TIMEOUT ((NUM_CLIENTS * MESSAGES_PER_CLIENT / 10) * 50 + 300)

static atomic_int messages_completed;
static atomic_int server_ready;

void simple_client_task(void* arg) {
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
    
    for (int i = 0; i < MESSAGES_PER_CLIENT; i++) {
        char msg[64];
        snprintf(msg, sizeof(msg), "MSG:%d:%d", client_id, i);
        
        __public__array_t send_buf;
        send_buf.data = msg;
        send_buf.length = strlen(msg);
        send_buf.capacity = strlen(msg);
        send_buf.elem_size = 1;
        
        ssize_t written = __public__net_write(client_fd, &send_buf, strlen(msg));
        
        /* Read response */
        char resp[64] = {0};
        __public__array_t recv_buf;
        recv_buf.data = resp;
        recv_buf.length = 0;
        recv_buf.capacity = sizeof(resp);
        recv_buf.elem_size = 1;
        
        ssize_t read_bytes = __public__net_read(client_fd, &recv_buf, sizeof(resp) - 1);
        
        if (written > 0 && read_bytes > 0) {
            atomic_fetch_add(&messages_completed, 1);
        }
    }
    
    close(client_fd);
}

void simple_handler(void* arg) {
    int client_fd = (int)(uintptr_t)arg;
    
    char buffer[128];
    
    for (int i = 0; i < MESSAGES_PER_CLIENT; i++) {
        __public__array_t recv_buf;
        recv_buf.data = buffer;
        recv_buf.length = 0;
        recv_buf.capacity = sizeof(buffer);
        recv_buf.elem_size = 1;
        
        ssize_t received = __public__net_read(client_fd, &recv_buf, sizeof(buffer) - 1);
        
        if (received <= 0) {
            break;
        }
        
        buffer[received] = '\0';
        
        __public__array_t send_buf;
        send_buf.data = buffer;
        send_buf.length = received;
        send_buf.capacity = received;
        send_buf.elem_size = 1;
        
        __public__net_write(client_fd, &send_buf, received);
    }
    
    close(client_fd);
}

void simple_server_task(void* arg) {
    /* Create listening socket */
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
    
    /* Accept NUM_CLIENTS connections */
    for (int i = 0; i < NUM_CLIENTS; i++) {
        ssize_t client_fd = __public__net_accept(listen_fd);
        
        if (client_fd >= 0) {
            thread(simple_handler, 1, (void*)(uintptr_t)client_fd);
        }
    }
    
    close(listen_fd);
}

void test_stress_multiple_messages(void) {
    atomic_store(&messages_completed, 0);
    atomic_store(&server_ready, 0);
    
    printf("Testing %d clients with %d messages each (Stress Level: %s)...\n",
           NUM_CLIENTS, MESSAGES_PER_CLIENT, STRESS_LEVEL);
    printf("Total messages: %d\n", NUM_CLIENTS * MESSAGES_PER_CLIENT);
    
    thread(simple_server_task, 0);
    
    struct timespec ts = {0, 100000000}; /* 100ms */
    nanosleep(&ts, NULL);
    
    for (int i = 0; i < NUM_CLIENTS; i++) {
        thread(simple_client_task, 1, (void*)(uintptr_t)i);
    }
    
    ts.tv_sec = 0;
    ts.tv_nsec = 100000000; /* 100ms */
    int timeout = 0;
    int expected_msgs = NUM_CLIENTS * MESSAGES_PER_CLIENT;
    
    while (atomic_load(&messages_completed) < expected_msgs && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
        
        if (timeout % 10 == 0) {
            int completed = atomic_load(&messages_completed);
            printf("Progress: %d/%d messages\n", completed, expected_msgs);
        }
    }
    
    int completed = atomic_load(&messages_completed);
    printf("Completed: %d/%d messages\n", completed, expected_msgs);
    
    TEST_ASSERT_GREATER_OR_EQUAL(expected_msgs * 0.90, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_multiple_messages);
    
    wait_for_schedulers();
    
    return UNITY_END();
}
