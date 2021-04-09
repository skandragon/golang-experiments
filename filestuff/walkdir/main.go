package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
)

func main() {
	err := filepath.WalkDir(".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.Type().IsRegular() {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		log.Printf("Got file %s, %d bytes", info.Name(), len(content))

		return nil
	})
	if err != nil {
		panic(err)
	}
}
