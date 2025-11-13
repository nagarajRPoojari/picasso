cmake_minimum_required(VERSION 3.10)
project(Scheduler C)

# ============================================================
#  Build Type & Global Options
# ============================================================
if (NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Debug CACHE STRING
        "Choose the build type: Debug or Release" FORCE)
endif()

message(STATUS "Build type: ${CMAKE_BUILD_TYPE}")

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# ============================================================
#  Paths & Directories
# ============================================================
set(SRC_DIR     "${CMAKE_SOURCE_DIR}/runtime")
set(BIN_DIR     "${CMAKE_SOURCE_DIR}/bin")     # keep outputs under project root
file(MAKE_DIRECTORY ${BIN_DIR})

set(GO_MAIN "${CMAKE_SOURCE_DIR}/main.go")
set(GO_LL   "${BIN_DIR}/output.ll")
set(GO_BC   "${BIN_DIR}/output.bc")
set(GO_OBJ  "${BIN_DIR}/output.o")

file(GLOB C_SOURCES "${SRC_DIR}/*.c")

# ============================================================
#  Step 1: Generate LLVM IR from Go
# ============================================================
add_custom_command(
    OUTPUT ${GO_LL}
    COMMAND go run ${GO_MAIN}
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "→ Generating LLVM IR from Go"
)

# ============================================================
#  Step 2: Assemble LLVM IR to Object
# ============================================================
add_custom_command(
    OUTPUT ${GO_OBJ}
    DEPENDS ${GO_LL}
    COMMAND llvm-as-16 ${GO_LL} -o ${GO_BC}
    COMMAND llc-16 -filetype=obj ${GO_BC} -o ${GO_OBJ}
    COMMENT "→ Compiling LLVM IR to object file"
)

# ============================================================
#  Step 3: Build the C Runtime + Link Everything
# ============================================================
add_executable(scheduler
    ${C_SOURCES}
    ${GO_OBJ}
)

# ============================================================
#  Include Directories
# ============================================================
target_include_directories(scheduler PRIVATE
    ${CMAKE_SOURCE_DIR}/headers
    /usr/local/include
    /usr/include
)

# ============================================================
#  Link Libraries
# ============================================================
target_link_libraries(scheduler PRIVATE
    /usr/local/lib/libgc.so    # Boehm GC
    pthread
    uring
    m
)

# ============================================================
#  Compiler / Linker Options
# ============================================================
target_compile_options(scheduler PRIVATE
    -Wall -Wextra -Wno-unused-parameter
)

if (CMAKE_BUILD_TYPE STREQUAL "Debug")
    message(STATUS "Enabling Debug flags and AddressSanitizer")
    target_compile_options(scheduler PRIVATE -g -O0 -fsanitize=address)
    target_link_options(scheduler PRIVATE -g -O0 -fsanitize=address)
else()
    target_compile_options(scheduler PRIVATE -O2)
    target_link_options(scheduler PRIVATE -O2)
endif()

# ============================================================
#  Output Settings
# ============================================================
set_target_properties(scheduler PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY ${BIN_DIR}
)

# ============================================================
#  Status Summary
# ============================================================
message(STATUS "--------------------------------------------------")
message(STATUS "Scheduler Build Configuration")
message(STATUS "--------------------------------------------------")
message(STATUS "Build Type:        ${CMAKE_BUILD_TYPE}")
message(STATUS "Binary Output Dir: ${BIN_DIR}")
message(STATUS "Using Boehm GC:    /usr/local/lib/libgc.so")
message(STATUS "Linked Libraries:  pthread, uring, m")
message(STATUS "--------------------------------------------------")
