# Picasso Programming Language

Picasso is a modern, compiled programming language designed for myself.

## Overview

Picasso combines the performance of compiled languages with the ease of use of modern high-level languages. It features automatic memory management, built-in concurrency primitives, and a rich standard library while maintaining zero-cost abstractions.

## Key Features

### Compiled Native Code
- Direct compilation to native machine code
- No virtual machine overhead
- Optimized for performance-critical applications

### Procedural with Object Support
- Clean procedural programming model
- Full support for classes and objects
- Method-based object orientation

### Rich Type System
- **Signed integers**: `int8`, `int16`, `int32`, `int64`
- **Unsigned integers**: `uint8`, `uint16`, `uint32`, `uint64`
- **Floating point**: `float`, `double`
- **Strings**: First-class string type with built-in operations
- **Atomics**: Lock-free atomic operations for all integer types
- **Arrays**: Dynamic arrays with automatic memory management
- **Classes**: User-defined types with methods and fields

### Built-in Concurrency
- Lightweight green threads with the `thread()` function
- No explicit async/await syntax required
- Automatic scheduling and context switching
- Scale to hundreds of thousands of concurrent tasks

### Automatic Memory Management
- Garbage collected runtime
- Allocate memory without manual cleanup
- No memory leaks or dangling pointers

### C Interoperability
- Foreign Function Interface (FFI) for C libraries
- Easy integration with existing C codebases
- Extend functionality with native C code

### Modular Design
- Simple module system with `using` statements
- Package management built into the language
- Clear namespace separation

### Cross-Platform Support
- Linux (aarch64/arm64)
- macOS (aarch64/arm64)
- Optimized for modern ARM architectures

### Comprehensive Standard Library
- **Network I/O**: TCP sockets, client/server programming
- **File I/O**: Synchronous file operations
- **OS Integration**: Platform-specific system calls
- **Synchronization**: Mutexes and atomic operations
- **String Manipulation**: Rich string processing utilities
- **Array Operations**: Dynamic array management

## Syntax Examples

### Hello World

```picasso
using "builtin/syncio";

fn start() {
    syncio.printf("Hello, World!\n");
}
```

### Classes and Objects

```picasso
using "builtin/syncio";

class Person {
    say name: string;
    say age: int;

    fn Person(name: string, age: int) {
        this.name = name;
        this.age = age;
    }

    fn greet() {
        syncio.printf("Hello, I'm %s and I'm %d years old\n", this.name, this.age);
    }
}

fn start() {
    say person: start.Person = new start.Person("Alice", 30);
    person.greet();
}
```

### Control Flow

```picasso
using "builtin/syncio";

fn start() {
    say x: int = 10;
    
    if (x < 0) {
        syncio.printf("Negative\n");
    } else if (x == 0) {
        syncio.printf("Zero\n");
    } else {
        syncio.printf("Positive\n");
    }
    
    // While loop
    say i: int = 0;
    while (i < 5) {
        syncio.printf("%d ", i);
        i = i + 1;
    }
    
    // Foreach loop
    foreach j in 0..10 {
        syncio.printf("%d ", j);
    }
}
```

### Arrays

```picasso
using "builtin/syncio";
using "builtin/array";

fn start() {
    say numbers: []int = array.create(int, 5);
    
    foreach i in 0..array.len(numbers) {
        numbers[i] = i * 10;
    }
    
    array.append(numbers, 50);
    array.append(numbers, 60);
    
    foreach i in 0..array.len(numbers) {
        syncio.printf("numbers[%d] = %d\n", i, numbers[i]);
    }
}
```

### Concurrency

```picasso
using "builtin/syncio";

class Worker {
    say id: int;
    
    fn Worker(id: int) {
        this.id = id;
    }
    
    fn work() {
        syncio.printf("Worker %d is working\n", this.id);
    }
}

fn start() {
    foreach i in 0..10 {
        say worker: start.Worker = new start.Worker(i);
        thread(worker.work);
    }
}
```

### Atomic Operations

```picasso
using "builtin/syncio";
using "builtin/atomics";

fn start() {
    say counter: atomic int64;
    
    atomics.store_int64(counter, int64(0));
    atomics.add_int64(counter, int64(10));
    atomics.sub_int64(counter, int64(3));
    
    say value: int64 = atomics.load_int64(counter);
    syncio.printf("Counter value: %ld\n", value);
}
```

### Network Programming

```picasso
using "builtin/syncio";
using "builtin/net";
using "builtin/array";

class Server {
    say addr: string;
    say port: int16;
    
    fn Server(addr: string, port: int16) {
        this.addr = addr;
        this.port = port;
    }
    
    fn start() {
        say fd: int = net.listen(this.addr, this.port, 4096, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0);
        
        if (fd < 0) {
            syncio.printf("Failed to start server\n");
            return;
        }
        
        syncio.printf("Server listening on %s:%d\n", this.addr, this.port);
        
        say clientFd: int = net.accept(fd);
        say buf: []uint8 = array.create(uint8, 1024);
        say n: int = net.read(clientFd, buf, 1024);
        
        if (n > 0) {
            net.write(clientFd, buf, n);
        }
    }
}

fn start() {
    say server: start.Server = new start.Server("127.0.0.1", 8080);
    server.start();
}
```

### File I/O

```picasso
using "builtin/syncio";
using "builtin/array";

fn start() {
    say file: string = syncio.fopen("data.txt", "w+");
    
    say data: []uint8 = array.create(uint8, 10);
    foreach i in 0..array.len(data) {
        data[i] = i;
    }
    
    syncio.fwrite(file, data, array.len(data), 0);
    
    say readBuf: []uint8 = array.create(uint8, 10);
    syncio.fread(file, readBuf, 10, 0);
    
    syncio.fclose(file);
}
```

### Module System

```picasso
// math.pic
using "builtin/syncio";

class Calculator {
    fn Calculator() {}
    
    fn add(a: int, b: int): int {
        return a + b;
    }
}
```

```picasso
// start.pic
using "builtin/syncio";
using "math" as m;

fn start() {
    say calc: m.Calculator = new m.Calculator();
    say result: int = calc.add(5, 3);
    syncio.printf("Result: %d\n", result);
}
```

## Variable Declaration

Variables are declared using the `say` keyword:

```picasso
say x: int = 10;
say name: string = "Alice";
say numbers: []int = array.create(int, 5);
say person: start.Person = new start.Person("Bob", 25);
```

## Access Modifiers

- **Public fields/methods**: Use `say` keyword (accessible from other modules)
- **Internal fields/methods**: Use `say internal` keyword (module-private)

```picasso
class Example {
    say publicField: int;
    say internal privateField: int;
    
    fn Example() {}
    
    fn publicMethod() {}
    
    fn internal privateMethod() {}
}
```

## Built-in Libraries

### syncio
Synchronous I/O operations including console output and file operations.

### net
Network programming with TCP sockets, client/server support.

### array
Dynamic array operations including creation, length, and append.

### strings
String manipulation utilities including formatting, comparison, and substring operations.

### atomics
Lock-free atomic operations for concurrent programming.

### types
Type conversion and type-related utilities.

## Building and Running

The project uses Bazel as its build system. Example test cases can be found in the `e2e/` directory.

## Project Structure

```
.
├── cli/              # Command-line interface
├── docs/             # Documentation
├── e2e/              # End-to-end tests
├── irgen/            # IR generation
├── libs/             # Standard library implementations
├── runtime/          # Runtime system
│   ├── headers/      # Runtime headers
│   └── src/          # Runtime source code
└── examples/         # Example programs
```

## Design Philosophy

Picasso is designed with the following principles:

1. **Performance First**: Compiled to native code with minimal runtime overhead
2. **Developer Productivity**: Automatic memory management and built-in concurrency
3. **Simplicity**: Clean syntax without unnecessary complexity
4. **Safety**: Strong type system with compile-time checks
5. **Interoperability**: Easy integration with C libraries
6. **Modern Architecture**: Optimized for ARM64 processors

## License

See LICENSE file for details.

## Contributing

Contributions are welcome. Please see CONTRIBUTING.md for guidelines.
