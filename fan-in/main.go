package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	generate := func(done <-chan interface{}, numbers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)

			for _, n := range numbers {
				select {
				case <-done:
					return
				case intStream <- n:
				}
			}
		}()
		return intStream
	}

	sq := func(done <-chan interface{}, inputStream <-chan int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)

			for v := range inputStream {
				select {
				case <-done:
					return
				case intStream <- v * v:
				}
			}
		}()
		return intStream
	}

	merge := func(done <-chan interface{}, inputChannels ...<-chan int) <-chan int {
		var wg sync.WaitGroup
		intStream := make(chan int)

		multiplex := func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				select {
				case <-done:
					return
				case intStream <- v:
				}
			}
		}

		wg.Add(len(inputChannels))
		for _, ch := range inputChannels {
			go multiplex(ch)
		}

		go func() {
			wg.Wait()
			close(intStream)
		}()

		return intStream
	}

	start := time.Now()

	done := make(chan interface{})
	defer close(done)

	inputStream := generate(done, 1, 2, 3, 4, 5)
	sq1 := sq(done, inputStream)
	sq2 := sq(done, inputStream)

	outputStream := merge(done, sq1, sq2)
	for v := range outputStream {
		log.Println(v)
	}

	log.Printf("Done in %v\n", time.Since(start))
}
