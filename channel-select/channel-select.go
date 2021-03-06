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
// This code will create multiple channels, each of which will send a
// message every few seconds.  The main code will call `select` to read
// from any of the channels.
//
// The senders will send `(id + 1) * 2` messages before closing the channel
//
// An aggregate channel is used in this example.  Each sender is given its
// own channel (of type int) to send its messages on, and the aggregate is
// sent a message that indicates the id, value, and closed status of the
// specific sender.
//

import (
	"log"
	"time"
)

const (
	nThreads = 3 // if this change, change the select below as well.
)

func sender(id int, c chan int) {
	defer close(c)
	for i := 0; i < (id+1)*2; i++ {
		c <- i
		time.Sleep(time.Second * 1)
	}
}

type message struct {
	id     int
	value  int
	closed bool
}

func main() {
	agg := make(chan *message)

	for i := 0; i < nThreads; i++ {
		ch := make(chan int)
		go func(c chan int, id int) {
			for msg := range c {
				agg <- &message{id: id, value: msg, closed: false}
			}
			agg <- &message{id: id, value: 0, closed: true}
		}(ch, i)
		go sender(i, ch)
	}

	activeCount := 3

	for msg := range agg {
		if msg.closed {
			log.Printf("Channel %d closed", msg.id)
			activeCount--
			if activeCount == 0 {
				break
			}
		} else {
			log.Printf("channel %d sent %d", msg.id, msg.value)
		}
	}
}
