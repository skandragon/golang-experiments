package main

//
// Copyright 2021 Michael Graff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
