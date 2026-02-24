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
    #define NUM_CYCLES 50
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define NUM_CYCLES 200
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define NUM_CYCLES 1000
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define NUM_CYCLES 500
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define TEST_PORT 9877
#define MAX_TIMEOUT ((NUM_CYCLES / 10) * 50 + 200)

static atomic_int cycles_completed;
static atomic_int server_ready;

void rapid_client_task(void* arg) {
    int cycle_id = (int)(uintptr_t)arg;
    
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
    
    char msg[16];
    snprintf(msg, sizeof(msg), "%d", cycle_id);
    
    __public__array_t send_buf;
    send_buf.data = msg;
    send_buf.length = strlen(msg);
    send_buf.capacity = strlen(msg);
    send_buf.elem_size = 1;
    
    ssize_t written = __public__net_write(client_fd, &send_buf, strlen(msg));
    
    close(client_fd);
    
    if (written > 0) {
        atomic_fetch_add(&cycles_completed, 1);
    }
}

void rapid_handler(void* arg) {
    int client_fd = (int)(uintptr_t)arg;
    
    char recv_buf[16] = {0};
    __public__array_t read_buf;
    read_buf.data = recv_buf;
    read_buf.length = 0;
    read_buf.capacity = sizeof(recv_buf);
    read_buf.elem_size = 1;
    
    __public__net_read(client_fd, &read_buf, sizeof(recv_buf) - 1);
    
    close(client_fd);
}

void rapid_server_task(void* arg) {
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
    
    for (int i = 0; i < NUM_CYCLES; i++) {
        ssize_t client_fd = __public__net_accept(listen_fd);
        
        if (client_fd >= 0) {
            thread(rapid_handler, 1, (void*)(uintptr_t)client_fd);
        }
    }
    
    close(listen_fd);
}

void test_stress_rapid_connect_disconnect(void) {
    atomic_store(&cycles_completed, 0);
    atomic_store(&server_ready, 0);
    
    printf("Testing %d rapid connect/disconnect cycles (Stress Level: %s)...\n", 
           NUM_CYCLES, STRESS_LEVEL);
    
    thread(rapid_server_task, 0);
    
    struct timespec ts = {0, 50000000}; /* 50ms */
    nanosleep(&ts, NULL);
    
    for (int i = 0; i < NUM_CYCLES; i++) {
        thread(rapid_client_task, 1, (void*)(uintptr_t)i);
        if (i % 10 == 0) {
            struct timespec tiny = {0, 1000000}; /* 1ms */
            nanosleep(&tiny, NULL);
        }
    }
    
    ts.tv_sec = 0;
    ts.tv_nsec = 100000000; /* 100ms */
    int timeout = 0;
    
    while (atomic_load(&cycles_completed) < NUM_CYCLES && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int completed = atomic_load(&cycles_completed);
    printf("Completed: %d/%d rapid cycles\n", completed, NUM_CYCLES);
    
    TEST_ASSERT_GREATER_OR_EQUAL(NUM_CYCLES * 0.90, completed);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_rapid_connect_disconnect);
    
    struct timespec ts = {0, 100000000}; /* 100ms */
    nanosleep(&ts, NULL);
    
    wait_for_schedulers();
    
    return UNITY_END();
}
