package main

//
// This code will spin up `nThreads` goroutines, each listening to a
// channel.  It will loop until the channel is closed.
//
// This shows that if there are multiple listeners, each message is
// consumed by one listener only, so this could be a simple work
// queue.
//

import (
	"log"
	"sync"
)

const (
	nThreads = 3
)

func listener(id int, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range c {
		log.Printf("goroutine %d got value %d", id, value)
	}
}

func main() {
	var wg sync.WaitGroup
	c := make(chan int)

	wg.Add(nThreads)
	for i := 0; i < nThreads; i++ {
		go listener(i, c, &wg)
	}

	for i := 0; i < 20; i++ {
		c <- i
	}
	close(c)

	wg.Wait()
}
