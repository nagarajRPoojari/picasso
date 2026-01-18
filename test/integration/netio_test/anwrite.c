#include "test/unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>
#include <pthread.h>
#include <stdatomic.h>
#include <time.h>
#include <stdio.h>
#include <errno.h>

#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include "netio.h"
#include "alloc.h"
#include "gc.h"
#include "ggc.h"
#include "initutils.h"

extern arena_t* __global__arena__;
arena_t* __test__global__arena__;

void setUp(void) {
    __test__global__arena__ = arena_create();
}
void tearDown(void) { }

__public__array_t* mock_alloc_array(int count, int elem_size, int rank) {
    size_t data_size = (size_t)count * elem_size;
    size_t shape_size = (size_t)rank * sizeof(int64_t);
    size_t total_size = sizeof(__public__array_t) + data_size + shape_size;

    __public__array_t* arr = (__public__array_t*)allocate(__test__global__arena__, total_size);
    
    arr->data = (char*)(arr + 1); 
    if (rank > 0) {
        arr->shape = (int64_t*)(arr->data + data_size);
    } else {
        arr->shape = NULL;
    }
    
    arr->length = count;
    arr->rank = rank;
    
    memset(arr->data, 0, data_size); 
    return arr;
}

#define MESSAGE "hello\n"
#define MESSAGE_LEN 6

// Atomic counter to track completed writes
static atomic_int completed;

static void* writer_thread(void* arg, int fd) {
    (void)arg;

    __public__array_t* buf = mock_alloc_array(MESSAGE_LEN, sizeof(char), 1);
    memcpy(buf->data, MESSAGE, MESSAGE_LEN);

    // Ensure full write
    int written = 0;
    while (written < MESSAGE_LEN) {
        int n = __public__net_write(fd, buf + written, MESSAGE_LEN - written);
        if (n > 0) written += n;
        else self_yield(); // yield if async write didn't progress
    }

    close(fd);
    atomic_fetch_add(&completed, 1);
    return NULL;
}

static void* server_thread(void* arg, int count, int listen_fd) {
    (void)arg;
    for (int i = 0; i < count; i++) {
        self_yield();
        int fd = __public__net_accept(listen_fd);
        if (fd >= 0) {
            thread(writer_thread, 2, NULL, fd); // schedule writer in runtime
        } else {
            atomic_fetch_add(&completed, 1);
        }
    }
    return NULL;
}

typedef struct { int count; int port; const char* addr; } client_args_t;

static void* native_client(void* varg) {
    client_args_t* args = (client_args_t*)varg;

    for (int i = 0; i < args->count; i++) {
        int fd = socket(AF_INET, SOCK_STREAM, 0);
        struct sockaddr_in addr = {0};
        addr.sin_family = AF_INET;
        addr.sin_port = htons(args->port);
        inet_pton(AF_INET, args->addr, &addr.sin_addr);

        // Retry connect briefly
        int connected = 0;
        for (int retry = 0; retry < 50; retry++) { // 50ms total
            if (connect(fd, (struct sockaddr*)&addr, sizeof(addr)) == 0) {
                connected = 1;
                break;
            }
            usleep(1000);
        }
        if (!connected) {
            close(fd);
            continue;
        }

        char read_buf[MESSAGE_LEN + 1] = {0};
        int total_read = 0;
        while (total_read < MESSAGE_LEN) {
            int n = read(fd, read_buf + total_read, MESSAGE_LEN - total_read);
            if (n > 0) total_read += n;
            else if (n == 0) break;
            else if (errno == EINTR || errno == EAGAIN) continue;
            else break;
        }

        TEST_ASSERT_EQUAL_STRING_LEN_MESSAGE(MESSAGE, read_buf, total_read, "Data mismatch");

        close(fd);
        usleep(500);
    }

    return NULL;
}

void test_netio_write_basic(void) {
    atomic_store(&completed, 0);

    int count = 100;
    const char* addr = "127.0.0.1";
    int port = 8002;

    int listen_fd = __public__net_listen(addr, port, 4096);
    TEST_ASSERT_MESSAGE(listen_fd >= 0, "Failed to listen");

    // Start server thread inside runtime scheduler
    thread(server_thread, 3, NULL, count, listen_fd);

    // Start native client pthread
    pthread_t client;
    client_args_t cargs = {count, port, addr};
    pthread_create(&client, NULL, native_client, &cargs);

    // Wait until all writes complete
    while (atomic_load_explicit(&completed, memory_order_acquire) < count) {
        struct timespec ts = {0, 5000000}; // 5ms
        nanosleep(&ts, NULL);
    }

    pthread_join(client, NULL);
    TEST_ASSERT_EQUAL_INT(count, atomic_load(&completed));

    close(listen_fd);
}

int main(void) {
    __global__arena__ = gc_create_global_arena();
    gc_init();
    init_io();
    init_scheduler();
    gc_start();

    UNITY_BEGIN();
    RUN_TEST(test_netio_write_basic);
    return UNITY_END();
}
