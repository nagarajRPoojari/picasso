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
    #define TRANSFER_SIZE (64 * 1024)        /* 64 KB */
    #define NUM_TRANSFERS 5
    #define STRESS_LEVEL "1 (Light)"
#elif defined(STRESS_LEVEL_2)
    #define TRANSFER_SIZE (256 * 1024)       /* 256 KB */
    #define NUM_TRANSFERS 10
    #define STRESS_LEVEL "2 (Medium)"
#elif defined(STRESS_LEVEL_4)
    #define TRANSFER_SIZE (4 * 1024 * 1024)  /* 4 MB */
    #define NUM_TRANSFERS 20
    #define STRESS_LEVEL "4 (Extreme)"
#else
    #define TRANSFER_SIZE (1024 * 1024)      /* 1 MB */
    #define NUM_TRANSFERS 10
    #define STRESS_LEVEL "3 (Heavy - Default)"
#endif

#define TEST_PORT 9878
#define MAX_TIMEOUT ((NUM_TRANSFERS * TRANSFER_SIZE / (1024 * 1024)) * 100 + 300)

static atomic_int transfers_completed;
static atomic_int server_ready;
static atomic_long total_bytes_sent;
static atomic_long total_bytes_received;

void large_transfer_client(void* arg) {
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
    
    char* send_data = (char*)malloc(TRANSFER_SIZE);
    if (!send_data) {
        close(client_fd);
        return;
    }
    
    for (size_t i = 0; i < TRANSFER_SIZE; i++) {
        send_data[i] = (char)((i + client_id) % 256);
    }
    
    __public__array_t send_buf;
    send_buf.data = send_data;
    send_buf.length = TRANSFER_SIZE;
    send_buf.capacity = TRANSFER_SIZE;
    send_buf.elem_size = 1;
    
    size_t total_sent = 0;
    while (total_sent < TRANSFER_SIZE) {
        ssize_t sent = __public__net_write(client_fd, &send_buf, TRANSFER_SIZE - total_sent);
        if (sent <= 0) {
            break;
        }
        total_sent += sent;
        send_buf.data = send_data + total_sent;
    }
    
    atomic_fetch_add(&total_bytes_sent, total_sent);
    
    /* Read acknowledgment */
    char ack[16] = {0};
    __public__array_t ack_buf;
    ack_buf.data = ack;
    ack_buf.length = 0;
    ack_buf.capacity = sizeof(ack);
    ack_buf.elem_size = 1;
    
    ssize_t ack_bytes = __public__net_read(client_fd, &ack_buf, sizeof(ack) - 1);
    
    free(send_data);
    close(client_fd);
    
    if (total_sent == TRANSFER_SIZE && ack_bytes > 0) {
        atomic_fetch_add(&transfers_completed, 1);
    }
}

void large_transfer_handler(void* arg) {
    int client_fd = (int)(uintptr_t)arg;
    
    /* Allocate receive buffer */
    char* recv_data = (char*)malloc(TRANSFER_SIZE);
    if (!recv_data) {
        close(client_fd);
        return;
    }
    
    __public__array_t recv_buf;
    recv_buf.data = recv_data;
    recv_buf.length = 0;
    recv_buf.capacity = TRANSFER_SIZE;
    recv_buf.elem_size = 1;
    
    /* Read large data */
    size_t total_received = 0;
    while (total_received < TRANSFER_SIZE) {
        ssize_t received = __public__net_read(client_fd, &recv_buf, TRANSFER_SIZE - total_received);
        if (received <= 0) {
            break;
        }
        total_received += received;
        recv_buf.data = recv_data + total_received;
    }
    
    atomic_fetch_add(&total_bytes_received, total_received);
    
    /* Send acknowledgment */
    char ack[16];
    snprintf(ack, sizeof(ack), "OK:%zu", total_received);
    
    __public__array_t ack_buf;
    ack_buf.data = ack;
    ack_buf.length = strlen(ack);
    ack_buf.capacity = strlen(ack);
    ack_buf.elem_size = 1;
    
    __public__net_write(client_fd, &ack_buf, strlen(ack));
    
    free(recv_data);
    close(client_fd);
}

void large_transfer_server(void* arg) {
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
        262144,     /* rcvbuf - 256KB */
        262144,     /* sndbuf - 256KB */
        0           /* ipv6_only */
    );
    
    if (listen_fd < 0) {
        return;
    }
    
    atomic_store(&server_ready, 1);
    
    /* Accept connections */
    for (int i = 0; i < NUM_TRANSFERS; i++) {
        ssize_t client_fd = __public__net_accept(listen_fd);
        
        if (client_fd >= 0) {
            thread(large_transfer_handler, 1, (void*)(uintptr_t)client_fd);
        }
    }
    
    close(listen_fd);
}

void test_stress_large_data_transfer(void) {
    atomic_store(&transfers_completed, 0);
    atomic_store(&server_ready, 0);
    atomic_store(&total_bytes_sent, 0);
    atomic_store(&total_bytes_received, 0);
    
    printf("Testing %d large transfers of %zu bytes each (Stress Level: %s)...\n", 
           NUM_TRANSFERS, (size_t)TRANSFER_SIZE, STRESS_LEVEL);
    printf("Total data to transfer: %.2f MB\n", 
           (NUM_TRANSFERS * TRANSFER_SIZE) / (1024.0 * 1024.0));
    
    thread(large_transfer_server, 0);
    
    struct timespec ts = {0, 100000000}; /* 100ms */
    nanosleep(&ts, NULL);
    
    for (int i = 0; i < NUM_TRANSFERS; i++) {
        thread(large_transfer_client, 1, (void*)(uintptr_t)i);
    }
    
    ts.tv_sec = 0;
    ts.tv_nsec = 200000000; /* 200ms */
    int timeout = 0;
    
    while (atomic_load(&transfers_completed) < NUM_TRANSFERS && timeout < MAX_TIMEOUT) {
        nanosleep(&ts, NULL);
        timeout++;
    }
    
    int completed = atomic_load(&transfers_completed);
    long sent = atomic_load(&total_bytes_sent);
    long received = atomic_load(&total_bytes_received);
    
    printf("Completed: %d/%d transfers\n", completed, NUM_TRANSFERS);
    printf("Total sent: %.2f MB, Total received: %.2f MB\n", 
           sent / (1024.0 * 1024.0), received / (1024.0 * 1024.0));
    
    TEST_ASSERT_GREATER_OR_EQUAL(NUM_TRANSFERS * 0.95, completed);
    TEST_ASSERT_EQUAL_INT64(sent, received);
}

int main(void) {
    UNITY_BEGIN();
    
    srand(time(NULL));
    __global__arena__ = gc_create_global_arena();
    gc_init();
    
    init_io();
    init_scheduler();
    gc_start();
    
    RUN_TEST(test_stress_large_data_transfer);
    
    wait_for_schedulers();
    
    return UNITY_END();
}
