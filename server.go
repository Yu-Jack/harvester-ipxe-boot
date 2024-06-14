package main

//go:generate go run pkg/template/generate.go

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./public"))

	http.Handle("/", fs)

	port := ":3333"

	log.Printf("Serving on HTTP port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
