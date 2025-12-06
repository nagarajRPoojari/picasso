cmake_minimum_required(VERSION 3.10)
project(Scheduler C)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# ------------------------------------------------------------
# Paths
# ------------------------------------------------------------
set(SRC_DIR "${CMAKE_SOURCE_DIR}/runtime")
set(BIN_DIR "${CMAKE_SOURCE_DIR}/bin")
file(MAKE_DIRECTORY ${BIN_DIR})

file(GLOB C_SOURCES "${SRC_DIR}/*.c")

set(GO_MAIN "${CMAKE_SOURCE_DIR}/main.go")
set(GO_LL   "${BIN_DIR}/output.ll")
set(GO_BC   "${BIN_DIR}/output.bc")
set(GO_OBJ  "${BIN_DIR}/output.o")

# ------------------------------------------------------------
# Force rebuild + clean bin/ before anything
# ------------------------------------------------------------
add_custom_target(always ALL
    COMMAND ${CMAKE_COMMAND} -E echo "Cleaning bin/ ..."
    COMMAND ${CMAKE_COMMAND} -E rm -rf ${BIN_DIR}/*
    COMMAND ${CMAKE_COMMAND} -E echo "Forcing rebuild..."
)

# ------------------------------------------------------------
# Go → LLVM → object (forced always)
# ------------------------------------------------------------
add_custom_command(
    OUTPUT ${GO_OBJ}
    COMMAND ${CMAKE_COMMAND} -E echo "Building Go → LLVM → Object"
    COMMAND go run ${GO_MAIN}
    COMMAND llvm-as-16 ${GO_LL} -o ${GO_BC}
    COMMAND llc-16 -filetype=obj ${GO_BC} -o ${GO_OBJ}
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    DEPENDS ${GO_MAIN} always
)

# Mark as symbolic → run every build
set_property(SOURCE ${GO_OBJ} PROPERTY SYMBOLIC TRUE)

add_custom_target(build_go_obj ALL
    DEPENDS ${GO_OBJ}
)

# ------------------------------------------------------------
# Executable
# ------------------------------------------------------------
add_executable(scheduler
    ${C_SOURCES}
    ${GO_OBJ}
)

add_dependencies(scheduler always)
add_dependencies(scheduler build_go_obj)

target_include_directories(scheduler PRIVATE
    ${CMAKE_SOURCE_DIR}/headers
)

# ------------------------------------------------------------
# Link libs
# ------------------------------------------------------------
target_link_libraries(scheduler PRIVATE
    /usr/local/lib/libgc.so   # Boehm GC
    pthread
    uring
    m
)

# ------------------------------------------------------------
# Flags
# ------------------------------------------------------------
if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Debug CACHE STRING "" FORCE)
endif()

if(CMAKE_BUILD_TYPE STREQUAL "Debug")
    target_compile_options(scheduler PRIVATE -g -O0 -fsanitize=address -Wall -Wextra)
    target_link_options(scheduler PRIVATE -g -O0 -fsanitize=address)
else()
    target_compile_options(scheduler PRIVATE -O2 -Wall -Wextra)
    target_link_options(scheduler PRIVATE -O2)
endif()

# ------------------------------------------------------------
# Output dir
# ------------------------------------------------------------
set_target_properties(scheduler PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY ${BIN_DIR}
)
