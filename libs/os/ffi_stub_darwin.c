#include "platform/darwin/os.h"

/**
 * Force references to all public OS wrapper functions so Clang/LLVM 
 * emits the necessary declarations and symbols for FFI discovery.
 */
void* __ffi_force[] = {

};