#include "unity/unity.h"
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdlib.h>
#include <assert.h>
#include "io.h"
#include "alloc.h"

extern __thread arena_t* __arena__;

void setUp(void) {
    __arena__ = arena_create();
}

void tearDown(void) {
    __arena__ = NULL;
}


/* blocking io test */

/* utility func to simulate stdin input */
static void redirect_stdin(const char *input) {
    FILE *f = fopen("test_input.txt", "w");
    assert(f);

    fputs(input, f);
    fclose(f);

    freopen("test_input.txt", "r", stdin);
}


void test__public__sscan(void) {
    /* small reads */
    redirect_stdin("dummy input from user\n");
    char* buf = (char*)__public__sscan(11);
    TEST_ASSERT_NOT_NULL(buf);
    TEST_ASSERT_EQUAL_STRING("dummy input", buf);
    
    
    /* large reads */
    int n = 1000;
    char input[n];
    for(int i=0; i<n-1; i++){
        input[i] = 'a' + i%26;
    }
    input[n-1] = '\0';
    redirect_stdin(input);
    buf = (char*)__public__sscan(n-1);
    TEST_ASSERT_EQUAL_STRING(input, buf);
    
    
    /* input < required length */
    /* __public__sscan designed for tty, shouldn't wait for all n bytes */
    redirect_stdin("input is only");
    buf = (char*)__public__sscan(20);
    TEST_ASSERT_EQUAL_STRING("input is only", buf);    
}

void test__public__sprint(void) {

}

int main(void) {
    UNITY_BEGIN();

    RUN_TEST(test__public__sscan);

    return UNITY_END();
}