package main

import (
	"log"

	"gopkg.in/yaml.v3"
)

func main() {

	s := `
a:
  b: 5
`

	var x map[interface{}]interface{}
	err := yaml.Unmarshal([]byte(s), &x)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", x)
	for k, v := range x {
		log.Printf("keyType %T, valueType %T", k, v)
	}

	out, err := yaml.Marshal(x)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Out:\n%s", out)
}
