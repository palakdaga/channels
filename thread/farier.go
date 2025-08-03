package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func prime(num int) bool {
	if num < 2 {
		return false
	}
	limit := int(math.Sqrt(float64(num)))
	for i := 2; i <= limit; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	CONCURRENCY := runtime.NumCPU()
	if CONCURRENCY < 1 {
		CONCURRENCY = 10
	}
	RANGE := 100000000

	var wg sync.WaitGroup
	var count int64

	jobs := make(chan int, 1000) // buffered channel for work units (numbers)

	// Start worker goroutines
	for i := 0; i < CONCURRENCY; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			localCount := int64(0)
			for num := range jobs {
				if prime(num) {
					localCount++
				}
			}
			atomic.AddInt64(&count, localCount)
			// Optionally print for each worker when finished
			fmt.Println("Worker", id, "done, found", localCount, "primes")
		}(i)
	}

	startTime := time.Now()

	// Feed jobs channel with numbers to check
	for num := 2; num < RANGE; num++ {
		jobs <- num
	}
	close(jobs) // close channel to signal no more jobs

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Println("Total prime numbers found:", count)
	fmt.Println("Time taken:", duration.Seconds(), "seconds")
}
