package main

import "log"

func main() {

	bridge := func(done <-chan interface{}, channelStream <-chan <-chan interface{}) <-chan interface{} {
		valueStream := make(chan interface{})

		go func() {
			defer close(valueStream)

			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-channelStream:
					if !ok {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}

				for v := range stream {
					select {
					case <-done:
					case valueStream <- v:
					}
				}

			}
		}()

		return valueStream
	}

	genVals := func() <-chan <-chan interface{} {
		channelStream := make(chan (<-chan interface{}))
		go func() {
			defer close(channelStream)

			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				channelStream <- stream
			}
		}()
		return channelStream

	}

	done := make(chan interface{})
	defer close(done)

	for v := range bridge(done, genVals()) {
		log.Println(v)
	}

}
