cmake_minimum_required(VERSION 3.16)
project(niyama C)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

find_library(LIBURING_LIB uring REQUIRED)
find_library(LIBGC_LIB gc REQUIRED)
find_package(Threads REQUIRED)

set(BIN_DIR ${CMAKE_BINARY_DIR}/bin)
file(MAKE_DIRECTORY ${BIN_DIR})

add_custom_command(
    OUTPUT ${BIN_DIR}/irgen
    COMMAND go build -o ${BIN_DIR}/irgen ${CMAKE_SOURCE_DIR}/irgen/main.go
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/irgen
    DEPENDS
        ${CMAKE_SOURCE_DIR}/irgen/main.go
        ${CMAKE_SOURCE_DIR}/irgen/go.mod
)

add_custom_target(irgen ALL
    DEPENDS ${BIN_DIR}/irgen
)


file(GLOB RUNTIME_SOURCES
    runtime/src/*.c
    gc/src/*.c
)

add_library(runtime STATIC
    ${RUNTIME_SOURCES}
)

target_include_directories(runtime PUBLIC
    runtime/headers
    gc/headers
    # Note: If liburing headers are in standard paths (/usr/include), 
    # we don't strictly need to add them, but it's safe to keep.
)

# Even though it's static, we tell CMake that 'runtime' needs uring.
# This helps if you ever link runtime into a CMake target, 
# but it DOES NOT help when you manually invoke ld in your CLI.
target_link_libraries(runtime PRIVATE 
    ${LIBURING_LIB}
)

set_target_properties(runtime PROPERTIES
    ARCHIVE_OUTPUT_DIRECTORY ${BIN_DIR}
)

add_executable(niyama
    cli/main.c
)

target_link_libraries(niyama PRIVATE
    ${LIBGC_LIB}
    Threads::Threads
    m
    ${LIBURING_LIB}
)

add_dependencies(niyama irgen)

set_target_properties(niyama PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY ${BIN_DIR}
)

install(TARGETS niyama
    RUNTIME DESTINATION bin
)

install(TARGETS runtime
    ARCHIVE DESTINATION lib
)