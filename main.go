package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	port := pflag.StringP("port", "p", "8080", "port to listen on")
	storagePath := pflag.StringP("storage", "s", "storage.json", "path to JSON storing user data")
	pflag.Parse()
	model, err := NewModel(*storagePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalln(err)
		}
		model = NewEmptyModel(*storagePath)
	}
	log.Println("Listening on port", *port)
	log.Fatalln(http.ListenAndServe(":"+*port, NewServer(model)))
}
