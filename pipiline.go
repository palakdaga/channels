// here we are making a pipeline
package main

import "math/rand"

// Generator function that yields integers
func genInt[T any, K any](done <-chan T, fn func() K) <-chan K {
	ch := make(chan K)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			default:
				ch <- fn()
			}
		}
	}()
	return ch
}

func findOdd(done <-chan int, ch <-chan int, num int ) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case n, ok := <-ch:
				if !ok {
					return
				}
				if n%2 != 0 {
					select {
					case out <- n:
					case <-done:
						return
					}
				}
			}
		}

	}()
	return out
}

func main() {
	done := make(chan int)
	defer close(done)

	// Generator function to produce integers
	intGen := genInt(done, func() int {
		return rand.Intn(100) // Replace with your logic to generate integers
	})

	// Find odd integers from the generator
	oddInts := findOdd(done, intGen, 10)

	// Consume the odd integers
	for oddInt := range oddInts {
		println(oddInt) // Replace with your logic to process odd integers
	}
}
