cmake_minimum_required(VERSION 3.10)
project(Scheduler C)

# Where your Go and C sources are
set(GO_MAIN "${CMAKE_SOURCE_DIR}/main.go")
file(GLOB C_SOURCES "${CMAKE_SOURCE_DIR}/runtime/*.c")

# Output object files
set(BIN_DIR "${CMAKE_BINARY_DIR}/bin")
file(MAKE_DIRECTORY ${BIN_DIR})

set(GO_LL "${BIN_DIR}/output.ll")
set(GO_BC "${BIN_DIR}/output.bc")
set(GO_OBJ "${BIN_DIR}/output.o")

# Step 1: Generate LLVM IR from Go
add_custom_command(
    OUTPUT ${GO_LL}
    COMMAND go run ${GO_MAIN}
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    COMMENT "Generating LLVM IR from Go"
)

# Step 2: Assemble LLVM IR to object
add_custom_command(
    OUTPUT ${GO_OBJ}
    DEPENDS ${GO_LL}
    COMMAND llvm-as-16 ${GO_LL} -o ${GO_BC}
    COMMAND llc-16 -filetype=obj ${GO_BC} -o ${GO_OBJ}
    COMMENT "Compiling LLVM IR to object"
)

# Step 4: Add executable including all objects
add_executable(scheduler
    ${C_SOURCES}
    ${GO_OBJ}
)

# Include directories
target_include_directories(scheduler PRIVATE
    ${CMAKE_SOURCE_DIR}/headers
    /usr/local/include
)

# Link libraries
target_link_libraries(scheduler PRIVATE
    /usr/local/lib/libgc.so
    pthread
)

# Compiler and linker options
target_compile_options(scheduler PRIVATE -O2)
target_link_options(scheduler PRIVATE -O2)

message(STATUS "Using Boehm GC include dirs: /usr/local/include")
message(STATUS "Using Boehm GC libraries: /usr/local/lib/libgc.so")
