package main

import (
	"fmt"
	"math/rand"
)

//using generics -> so that code is usable for any data type

func generator[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)
	//infinite loop
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

func main() {
	done := make(chan int)
	defer close(done)
	fn := func() int { return rand.Intn(100) }
	stream := generator(done, fn)
	for i := 0; i < 10; i++ {
		fmt.Println(<-stream)
	}
	println("Done")

}
