package main

import (
	"log"

	"github.com/spf13/pflag"
)

func main() {
	storagePath := pflag.StringP("storage", "s", "storage.json", "path to JSON storing user data")
	pflag.Parse()
	model, err := NewModel(*storagePath)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(model)
}
