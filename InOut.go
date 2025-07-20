package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

//generator of random no.

//collect into one channel->wait groups arre handled because too many reqs will come to handle them gracefully mutual exclusion is needed
//print output

func repeatFunc(done <-chan int, fn func() int) <-chan int {
	stream := make(chan int)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			default:
				stream <- fn()
			}
		}
	}()
	return stream
}

func take(done <-chan int, stream <-chan int, n int) <-chan int {
	taken := make(chan int)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case v, ok := <-stream:
				if !ok {
					return
				}
				taken <- v
			}
		}
	}()
	return taken
}

func fanIn[T any](done <-chan int, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	fannedInSteam := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInSteam <- i:
			}
		}
	}
	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}
	go func() {
		wg.Wait()
		close(fannedInSteam)
	}()
	return fannedInSteam
}

func primeFinder(done <-chan int, randIntStream <-chan int) <-chan int {
	IsPrime := func(randomInt int) bool {
		if randomInt <= 1 {
			return false
		}
		for i := 2; i*i <= randomInt; i++ {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}

	Primes := make(chan int)
	go func() {
		defer close(Primes)
		for {
			select {
			case <-done:
				return
			case randomInt, ok := <-randIntStream:
				if !ok {
					return
				}
				if IsPrime(randomInt) {
					Primes <- randomInt
				}
			}
		}
	}()
	return Primes
}

func main() {
	start := time.Now()
	done := make(chan int)
	defer close(done)

	randNumFetcher := func() int { return rand.Intn(1000000) }
	randIntStream := repeatFunc(done, randNumFetcher)

	CPUcount := runtime.NumCPU()
	primeFinderChannels := make([]<-chan int, CPUcount)
	for i := 0; i < CPUcount; i++ {
		primeFinderChannels[i] = primeFinder(done, randIntStream)
	}
	//fan in
	fanInStream := fanIn(done, primeFinderChannels...)
	for rando := range take(done, fanInStream, 10) {
		fmt.Println(rando) // Replace with your logic to handle the prime numbers
	}

	fmt.Println("Time taken:", time.Since(start))
}
