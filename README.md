# Picasso Programming Language

Picasso is a modern, compiled programming language designed for myself.

## Overview

Picasso combines the performance of compiled languages with the ease of use of modern high-level languages. It features automatic memory management, built-in concurrency primitives, and a rich standard library while maintaining zero-cost abstractions.

## Key Features

- **Compiled Native Code**: Direct compilation to native machine code without virtual machine overhead.

- **Procedural with Object Support**: Clean procedural programming with full support for classes and objects.

- **Rich Type System**: Signed/unsigned integers (`int8` to `int64`, `uint8` to `uint64`), floating point (`float`, `double`), strings, atomics, dynamic arrays, and user-defined classes.

- **Built-in Concurrency**: Lightweight green threads with `thread()` function - no explicit async/await required. Scale to hundreds of thousands of concurrent tasks.

- **Automatic Memory Management**: Garbage collected runtime - allocate and forget.

- **C Interoperability**: Foreign Function Interface (FFI) for seamless integration with C libraries.

- **Modular Design**: Simple module system with `using` statements and clear namespace separation.

- **Cross-Platform Support**: Linux and macOS on aarch64/arm64 architectures.

- **Comprehensive Standard Library**: Network I/O, file I/O, OS integration, synchronization primitives, string manipulation, and array operations.

## Syntax Examples

### Hello World

```picasso
using "builtin/syncio";

fn start() {
    syncio.printf("Hello, World!\n");
}
```

### Classes and Objects

```python
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

```python
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

```python
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

```python
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

```python
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

```python
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

```python
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

```python
// math.pic
using "builtin/syncio";

class Calculator {
    fn Calculator() {}
    
    fn add(a: int, b: int): int {
        return a + b;
    }
}
```

```python
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

```python
say x: int = 10;
say name: string = "Alice";
say numbers: []int = array.create(int, 5);
say person: start.Person = new start.Person("Bob", 25);
```

## Access Modifiers

- **Public fields/methods**: Use `say` keyword (accessible from other modules)
- **Internal fields/methods**: Use `say internal` keyword (module-private)

```python
class Example {
    say publicField: int;
    say internal privateField: int;
    
    fn Example() {}
    
    fn publicMethod() {}
    
    fn internal privateMethod() {}
}
```

## Built-in Libraries

- **syncio**: Synchronous I/O operations including console output and file operations.
- **net**: Network programming with TCP sockets, client/server support.
- **array**: Dynamic array operations including creation, length, and append.
- **strings**: String manipulation utilities including formatting, comparison, and substring operations.
- **atomics**: Lock-free atomic operations for concurrent programming.
- **types**: Type conversion and type-related utilities.
