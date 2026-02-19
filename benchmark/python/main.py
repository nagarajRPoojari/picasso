#!/usr/bin/env python3.12


import threading
import sys
from typing import List

# Increase integer string conversion limit for large Fibonacci numbers
sys.set_int_max_str_digits(0)  # 0 = unlimited

# BenchmarkUtils class
class BenchmarkUtils:
    def __init__(self):
        pass
    
    # Fibonacci calculation (recursive)
    def fibonacci(self, n: int) -> int:
        if n <= 1:
            return n
        return self.fibonacci(n - 1) + self.fibonacci(n - 2)
    
    # Fibonacci calculation (iterative)
    def fibonacci_iterative(self, n: int) -> int:
        if n <= 1:
            return n
        
        a, b = 0, 1
        for i in range(2, n + 1):
            a, b = b, a + b
        
        return b
    
    # Array operations benchmark
    def array_operations(self, size: int) -> int:
        arr = []
        
        # Append elements
        for i in range(size):
            arr.append(i)
        
        # Sum all elements
        total = 0
        for i in range(len(arr)):
            total += arr[i]
        
        return total
    
    # String operations benchmark
    def string_operations(self, iterations: int) -> int:
        count = 0
        
        for i in range(iterations):
            temp = f"Iteration {i}"
            count += len(temp)
        
        return count
    
    # Prime number calculation
    def is_prime(self, n: int) -> int:
        if n <= 1:
            return 0
        if n <= 3:
            return 1
        if n % 2 == 0:
            return 0
        if n % 3 == 0:
            return 0
        
        i = 5
        while i * i <= n:
            if n % i == 0:
                return 0
            if n % (i + 2) == 0:
                return 0
            i += 6
        
        return 1
    
    def count_primes(self, limit: int) -> int:
        count = 0
        for i in range(2, limit):
            if self.is_prime(i) == 1:
                count += 1
        return count


# Matrix class
class Matrix:
    def __init__(self, rows: int, cols: int):
        self.rows = rows
        self.cols = cols
        self.data = [[0 for _ in range(cols)] for _ in range(rows)]
    
    def set(self, row: int, col: int, value: int):
        self.data[row][col] = value
    
    def get(self, row: int, col: int) -> int:
        return self.data[row][col]


# MatrixOps class
class MatrixOps:
    def __init__(self):
        pass
    
    def multiply(self, size: int) -> int:
        a = Matrix(size, size)
        b = Matrix(size, size)
        c = Matrix(size, size)
        
        # Initialize matrices
        for i in range(size):
            for j in range(size):
                a.set(i, j, i + j)
                b.set(i, j, i - j)
        
        # Multiply
        for i in range(size):
            for j in range(size):
                total = 0
                for k in range(size):
                    total += a.get(i, k) * b.get(k, j)
                c.set(i, j, total)
        
        return c.get(0, 0)


# Counter class
class Counter:
    def __init__(self):
        self.value = 0
        self.lock = threading.Lock()
    
    def increment(self):
        with self.lock:
            self.value += 1
    
    def get_value(self) -> int:
        with self.lock:
            return self.value


# ConcurrencyBench class
class ConcurrencyBench:
    def __init__(self):
        self.counter = Counter()
    
    def increment_worker(self, iterations: int):
        for i in range(iterations):
            self.counter.increment()
    
    def run(self, threads: int, iterations: int) -> int:
        thread_list = []
        
        for i in range(threads):
            t = threading.Thread(target=self.increment_worker, args=(iterations,))
            thread_list.append(t)
            t.start()
        
        for t in thread_list:
            t.join()
        
        return self.counter.get_value()


def main():
    print("=== Python Performance Benchmark ===\n")
    
    utils = BenchmarkUtils()
    matrix_ops = MatrixOps()
    
    # Benchmark 1: Fibonacci (recursive)
    print("1. Fibonacci (recursive, n=30)...")
    fib_result = utils.fibonacci(30)
    print(f"   Result: {fib_result}\n")
        
    # Benchmark 2: Fibonacci (iterative)
    print("2. Fibonacci (iterative, n=1000000)...")
    fib_iter_result = utils.fibonacci_iterative(1000000)
    # Calculate number of digits without full string conversion
    import math
    if fib_iter_result > 0:
        num_digits = int(math.log10(fib_iter_result)) + 1
        print(f"   Result: <large number with {num_digits} digits>\n")
    else:
        print(f"   Result: {fib_iter_result}\n")
    
    # Benchmark 3: Array operations
    print("3. Array operations (size=100000)...")
    array_result = utils.array_operations(100000)
    print(f"   Sum: {array_result}\n")
    
    # Benchmark 4: String operations
    print("4. String operations (iterations=10000)...")
    string_result = utils.string_operations(10000)
    print(f"   Total length: {string_result}\n")
    
    # Benchmark 5: Matrix multiplication
    print("5. Matrix multiplication (50x50)...")
    matrix_result = matrix_ops.multiply(50)
    print(f"   Result[0][0]: {matrix_result}\n")
    
    # Benchmark 6: Prime counting
    print("6. Prime counting (limit=100000)...")
    prime_count = utils.count_primes(100000)
    print(f"   Primes found: {prime_count}\n")
    
    # Benchmark 7: Concurrency
    print("7. Concurrent counter (10 threads, 10000 iterations each)...")
    conc_bench = ConcurrencyBench()
    concurrent_result = conc_bench.run(10, 10000)
    print(f"   Final count: {concurrent_result}\n")
    
    print("=== Benchmark Complete ===")


if __name__ == "__main__":
    main()

