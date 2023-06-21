package main

import "log"

func main() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan interface{}, inputStream <-chan interface{}, limit int) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)

			for i := 0; i < limit; i++ {
				select {
				case <-done:
					return
				case valueStream <- <-inputStream:
				}
			}
		}()
		return valueStream
	}

	// orDone := func(done, c <-chan interface{}) <-chan interface{} {
	// 	valStream := make(chan interface{})
	// 	go func() {
	// 		defer close(valStream)
	// 		for {
	// 			select {
	// 			case <-done:
	// 				return
	// 			case v, ok := <-c:
	// 				if !ok {
	// 					return
	// 				}
	// 				select {
	// 				case valStream <- v:
	// 				case <-done:
	// 				}
	// 			}
	// 		}
	// 	}()
	// 	return valStream
	// }

	tee := func(done <-chan interface{}, inputStream <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})

		go func() {
			defer close(out1)
			defer close(out2)

			// for v := range inputStream {
			for v := range inputStream {
				var out1, out2 = out1, out2

				// Note while reading read one reac from out1, out2 else we will endup in deadlock
				for i := 0; i < 2; i++ {
					select {
					case <-done:
					case out1 <- v:
						out1 = nil
					case out2 <- v:
						out2 = nil
					}
				}
			}
		}()

		return out1, out2
	}

	done := make(chan interface{})
	defer close(done)

	resultStream := take(done, repeat(done, 1, 2, 3, 4), 6)
	for v := range resultStream {
		log.Println(v)
	}

	outStream1, outStream2 := tee(done, take(done, repeat(done, 1, 2, 3, 4), 6))
	for val1 := range outStream1 {
		log.Printf("out1: %v, out2: %v\n", val1, <-outStream2)
	}

}
