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
// This code will run a command (program arguments become the commnd and its arguments)
// and display stdout and stderr received.
//
// It handles this by using a goroutine to monitor stdout and another to monitor
// stderr.  Each sends to an aggregate channel, with the output stream contents
// wrapped in a `message` which includes the id (1 or 2 for stdout and stderr),
// the data, and a 'closed' boolean.  If closed is set, the data will be empty.
//

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func sender(id int, c chan *message, in io.Reader) {
	buffer := make([]byte, 10240)
	for {
		n, err := in.Read(buffer)
		if n > 0 {
			log.Printf("%d read %d bytes", id, n)
			c <- &message{id: id, value: string(buffer[:n]), closed: false}
		}
		if err == io.EOF {
			c <- &message{id: id, value: "", closed: true}
		}
		if err != nil {
			log.Printf("Got %v in read", err)
			c <- &message{id: id, value: "", closed: true}
		}
	}
}

type message struct {
	id     int
	value  string
	closed bool
}

func main() {
	agg := make(chan *message)

	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	log.Printf("Running %s...", os.Args[1])

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go sender(1, agg, stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go sender(2, agg, stderr)

	log.Printf("Starting command")
	err = cmd.Start()
	if err != nil {
		log.Fatalf("running cmd: %v", err)
	}

	log.Printf("Reading channels")
	activeCount := 2
	for msg := range agg {
		if msg.closed {
			log.Printf("Channel %d closed", msg.id)
			activeCount--
			if activeCount == 0 {
				return
			}
		} else {
			log.Printf("channel %d sent %s", msg.id, msg.value)
		}
	}

	log.Printf("Waiting for command to exit...")
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
	log.Printf("Command exited with status code 0")
}
