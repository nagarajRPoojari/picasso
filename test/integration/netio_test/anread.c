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
void tearDown(void) {
    /* @todo: graceful termination */
}

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
    
    return arr;
}

static atomic_int completed;

#define MESSAGE "hello\n"
#define MESSAGE_LEN 6

static void __public__afread_thread_func(void* arg, int fd) {
    (void)arg;

    int buf_size = 10;

    self_yield();
    __public__array_t* buf = mock_alloc_array(buf_size, sizeof(size_t), 1);

    self_yield();
    int n = __public__net_read(fd, buf, buf_size - 1);

    if (n > 0) {
        buf->data[n] = '\0';
        TEST_ASSERT_EQUAL_STRING_LEN_MESSAGE(MESSAGE, buf->data, MESSAGE_LEN, "received invalid message");
    }

    close(fd);

    atomic_fetch_add(&completed, 1);
}

void __test_netio_read_basic(void* arg, int count, int lld) {
    (void)arg;

    for (int i = 0; i < count; i++) {

        self_yield();
        int fd = __public__net_accept(lld);

        self_yield();
        thread(__public__afread_thread_func, 2, NULL, fd);
    }

}

int connect_to(const char *ip, int port) {
    int fd = socket(AF_INET, SOCK_STREAM, 0);
    if (fd < 0) {
        perror("socket");
        return -1;
    }

    struct sockaddr_in addr = {0};
    addr.sin_family = AF_INET;
    addr.sin_port   = htons(port);
    inet_pton(AF_INET, ip, &addr.sin_addr);

    self_yield();
    if (connect(fd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        perror("connect");
        close(fd);
        return -1;
    }

    return fd;
}

void simulate_client(void* arg, int count, char* addr, int port) {
    (void)arg;
    for (int i = 0; i < count; i++) {

        self_yield();
        int fd = connect_to(addr, port);
        assert(fd >= 0);


        self_yield();
        write(fd, MESSAGE, MESSAGE_LEN);
        close(fd);
    }

}

void test_netio_read_basic(void) {
    atomic_store(&completed, 0);

    int count = 100;
    int timeout_sec = 5;

    char* addr = "127.0.0.1";
    int port = 8000;

    int lld = __public__net_listen(addr, port, 4096);
    TEST_ASSERT(lld >= 0);

    thread(__test_netio_read_basic, 3, NULL, count, lld);
    thread(simulate_client, 4, NULL, count, addr, port);
    
    struct timespec ts = {0, 1000000}; /* 1ms */
    int spins = 0;
    int max_spins = timeout_sec * 1000;

    while (atomic_load_explicit(&completed, memory_order_acquire) < count) {
        nanosleep(&ts, NULL);
        if (++spins > max_spins) {
            TEST_FAIL_MESSAGE("scheduler did not complete tasks");
        }
    }

    TEST_ASSERT_EQUAL_INT(count, atomic_load(&completed));
}


int main(void) {
    srand(time(NULL));

    __global__arena__ = gc_create_global_arena();
    gc_init();

    init_io();
    init_scheduler();

    gc_start();

    UNITY_BEGIN();
    RUN_TEST(test_netio_read_basic);
    return UNITY_END();
}
