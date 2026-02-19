package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// BenchmarkUtils class equivalent
type BenchmarkUtils struct{}

func NewBenchmarkUtils() *BenchmarkUtils {
	return &BenchmarkUtils{}
}

// Fibonacci calculation (recursive)
func (bu *BenchmarkUtils) Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return bu.Fibonacci(n-1) + bu.Fibonacci(n-2)
}

// Fibonacci calculation (iterative)
func (bu *BenchmarkUtils) FibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}

	return b
}

// Array operations benchmark
func (bu *BenchmarkUtils) ArrayOperations(size int) int {
	arr := make([]int, 0)

	// Append elements
	for i := 0; i < size; i++ {
		arr = append(arr, i)
	}

	// Sum all elements
	sum := 0
	for i := 0; i < len(arr); i++ {
		sum += arr[i]
	}

	return sum
}

// String operations benchmark
func (bu *BenchmarkUtils) StringOperations(iterations int) int {
	count := 0

	for i := 0; i < iterations; i++ {
		temp := fmt.Sprintf("Iteration %d", i)
		count += len(temp)
	}

	return count
}

// Prime number calculation
func (bu *BenchmarkUtils) IsPrime(n int) int {
	if n <= 1 {
		return 0
	}
	if n <= 3 {
		return 1
	}
	if n%2 == 0 {
		return 0
	}
	if n%3 == 0 {
		return 0
	}

	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 {
			return 0
		}
		if n%(i+2) == 0 {
			return 0
		}
	}

	return 1
}

func (bu *BenchmarkUtils) CountPrimes(limit int) int {
	count := 0
	for i := 2; i < limit; i++ {
		if bu.IsPrime(i) == 1 {
			count++
		}
	}
	return count
}

// Matrix class
type Matrix struct {
	data [][]int
	rows int
	cols int
}

func NewMatrix(rows, cols int) *Matrix {
	data := make([][]int, rows)
	for i := range data {
		data[i] = make([]int, cols)
	}

	return &Matrix{
		data: data,
		rows: rows,
		cols: cols,
	}
}

func (m *Matrix) Set(row, col, value int) {
	m.data[row][col] = value
}

func (m *Matrix) Get(row, col int) int {
	return m.data[row][col]
}

// MatrixOps class equivalent
type MatrixOps struct{}

func NewMatrixOps() *MatrixOps {
	return &MatrixOps{}
}

func (mo *MatrixOps) Multiply(size int) int {
	a := NewMatrix(size, size)
	b := NewMatrix(size, size)
	c := NewMatrix(size, size)

	// Initialize matrices
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			a.Set(i, j, i+j)
			b.Set(i, j, i-j)
		}
	}

	// Multiply
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			sum := 0
			for k := 0; k < size; k++ {
				sum += a.Get(i, k) * b.Get(k, j)
			}
			c.Set(i, j, sum)
		}
	}

	return c.Get(0, 0)
}

// Counter class
type Counter struct {
	value int64
}

func NewCounter() *Counter {
	return &Counter{value: 0}
}

func (c *Counter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *Counter) GetValue() int64 {
	return atomic.LoadInt64(&c.value)
}

// ConcurrencyBench class equivalent
type ConcurrencyBench struct {
	counter *Counter
}

func NewConcurrencyBench() *ConcurrencyBench {
	return &ConcurrencyBench{
		counter: NewCounter(),
	}
}

func (cb *ConcurrencyBench) IncrementWorker(iterations int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < iterations; i++ {
		cb.counter.Increment()
	}
}

func (cb *ConcurrencyBench) Run(threads, iterations int) int64 {
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go cb.IncrementWorker(iterations, &wg)
	}

	wg.Wait()
	return cb.counter.GetValue()
}

func main() {
	fmt.Println("=== Go Performance Benchmark ===\n")

	utils := NewBenchmarkUtils()
	matrixOps := NewMatrixOps()

	// Benchmark 1: Fibonacci (recursive)
	fmt.Println("1. Fibonacci (recursive, n=30)...")
	fibResult := utils.Fibonacci(30)
	fmt.Printf("   Result: %d\n\n", fibResult)

	// Benchmark 2: Fibonacci (iterative)
	fmt.Println("2. Fibonacci (iterative, n=1000000)...")
	fibIterResult := utils.FibonacciIterative(1000000)
	fmt.Printf("   Result: %d\n\n", fibIterResult)

	// Benchmark 3: Array operations
	fmt.Println("3. Array operations (size=100000)...")
	arrayResult := utils.ArrayOperations(100000)
	fmt.Printf("   Sum: %d\n\n", arrayResult)

	// Benchmark 4: String operations
	fmt.Println("4. String operations (iterations=10000)...")
	stringResult := utils.StringOperations(10000)
	fmt.Printf("   Total length: %d\n\n", stringResult)

	// Benchmark 5: Matrix multiplication
	fmt.Println("5. Matrix multiplication (50x50)...")
	matrixResult := matrixOps.Multiply(50)
	fmt.Printf("   Result[0][0]: %d\n\n", matrixResult)

	// Benchmark 6: Prime counting
	fmt.Println("6. Prime counting (limit=100000)...")
	primeCount := utils.CountPrimes(100000)
	fmt.Printf("   Primes found: %d\n\n", primeCount)

	// Benchmark 7: Concurrency
	fmt.Println("7. Concurrent counter (10 threads, 10000 iterations each)...")
	concBench := NewConcurrencyBench()
	concurrentResult := concBench.Run(10, 10000)
	fmt.Printf("   Final count: %d\n\n", concurrentResult)

	fmt.Println("=== Benchmark Complete ===")
}
