package main

import (
	"log"
	"time"
)

func main() {
	stop := time.After(3 * time.Second)
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			log.Println("Helo World")
			log.Println(time.Now())
		case <-stop:
			return
		}
	}
}
