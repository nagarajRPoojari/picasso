cmake_minimum_required(VERSION 3.10)
project(GC_Tests C)

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -Wall -Wextra")

set(RUNTIME_HEADERS ${CMAKE_SOURCE_DIR}/../runtime/headers)
set(GC_HEADERS ${CMAKE_SOURCE_DIR}/../gc/headers)
set(UNITY_SRC unity/unity.c)

file(GLOB GC_RUNTIME_SOURCES_ALL
    "${CMAKE_SOURCE_DIR}/../gc/src/*.c"
    "${CMAKE_SOURCE_DIR}/../runtime/src/*.c"
)

set(FILE_TO_EXCLUDE "${CMAKE_SOURCE_DIR}/../runtime/src/main.c")

list(REMOVE_ITEM GC_RUNTIME_SOURCES_ALL "${FILE_TO_EXCLUDE}")

set(GC_RUNTIME_SOURCES ${GC_RUNTIME_SOURCES_ALL})


file(GLOB TEST_SOURCES "*.c")


include_directories(
    src headers test
    ${RUNTIME_HEADERS}
    ${GC_HEADERS}
    /usr/include/liburing 
)

enable_testing()


foreach(test_src ${TEST_SOURCES})
    get_filename_component(test_name ${test_src} NAME_WE)

    if (NOT test_name STREQUAL "unity")
        
        add_executable(${test_name}
            ${test_src}
            ${GC_RUNTIME_SOURCES}
            ${UNITY_SRC}
        )

        target_link_libraries(${test_name} PRIVATE
            uring  # Requires the system to have liburing installed
            pthread
            m      # Math library
            gc     # If your garbage collector uses libgc.so
        )

        add_test(NAME ${test_name} COMMAND ${test_name})
        
    endif()
endforeach()