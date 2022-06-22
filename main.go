package main

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	storagePath := pflag.StringP("storage", "s", "storage.json", "path to JSON storing user data")
	pflag.Parse()
	model, err := NewModel(*storagePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalln(err)
		}
		model = NewEmptyModel(*storagePath)
	}
	log.Println(model)
}
