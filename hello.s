.global _main             ; Export the entry point for the linker
.align 4                  ; Align instructions to 4-byte boundaries

_main:
    ; write(1, message, 13)
    mov x0, #1            ; Arg 0: File descriptor (1 = stdout)
    adrp x1, message@PAGE ; Arg 1: Load page address of the string
    add x1, x1, message@PAGEOFF ; Add the offset within that page
    mov x2, #13           ; Arg 2: Length of string
    mov x16, #4           ; System call number for 'write'
    svc #0                ; Supervisor Call (invokes the kernel)

    ; exit(0)
    mov x0, #0            ; Return code 0
    mov x16, #1           ; System call number for 'exit'
    svc #0                ; Invoke kernel
    
.data
message:
    .ascii "Hello, ARM64\n"