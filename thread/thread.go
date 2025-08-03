package main

import (
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
	for i := 2; i <= int(math.Sqrt(float64(num))); i++ {
		if num%int(i) == 0 {
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
	var count int64 = 0
	bunch := RANGE / CONCURRENCY
	startTime := time.Now()
	for i := 0; i < CONCURRENCY; i++ {
		start := i * bunch
		end := start + bunch
		if i == CONCURRENCY-1 {
			end = RANGE
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			localCount := int64(0)
			for j := start; j < end; j++ {
				if prime(j) {
					localCount++
				}
			}
			atomic.AddInt64(&count, localCount)
			println("Thread finished:", start, "to", end, "found", localCount, "primes",time.Now().Sub(startTime).Seconds(), "seconds")
		}(start, end)
		
	}
	wg.Wait()
	println("Total prime numbers found:", count)
	duration := time.Since(startTime)
	println("Time taken:", duration.Seconds(), "seconds")
	
}
