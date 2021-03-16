package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {
	home := os.Getenv("HOME")
	pubkeyFilename := fmt.Sprintf("%s/.ssh/id_rsa", home)

	config := &ssh.ClientConfig{
		User: "explorer",
		Auth: []ssh.AuthMethod{
			publicKey(pubkeyFilename),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	runCommand(conn, "env", []string{"LC_FOO=foo", "LC_BAR=bar", "LC_BAZ"})
}

func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Panic(err)
	}
	return ssh.PublicKeys(signer)
}

func runCommand(conn *ssh.Client, cmd string, env []string) {
	sess, err := conn.NewSession()
	if err != nil {
		log.Panic(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	go io.Copy(os.Stdout, sessStdOut)

	sessStderr, err := sess.StderrPipe()
	if err != nil {
		log.Panic(err)
	}
	go io.Copy(os.Stderr, sessStderr)

	// For each envar passed in, if there is no =, use the running process's value.
	// If it has an = (even if empty on the right hand side) use it.
	// else, silently ignore it.
	// Note many SSH servers reject all but a limited subset of envars for
	// security reasons.  LC_* are usually allowed.
	for _, envar := range env {
		items := strings.SplitN(envar, "=", 2)
		if len(items) == 1 {
			val, ok := os.LookupEnv(items[0])
			if ok {
				err = sess.Setenv(items[0], val)
				if err != nil {
					log.Panic(err)
				}
			}
		} else if len(items) == 2 {
			err = sess.Setenv(items[0], items[1])
			if err != nil {
				log.Panic(err)
			}
		} else {
			continue
		}
	}

	err = sess.Run(cmd) // cmd is processed by a shell.
	if err != nil {
		log.Panic(err)
	}
}
