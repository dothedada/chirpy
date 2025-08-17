package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const fileRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(fileRoot)))

	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: mux,
	}

	fmt.Printf("Serving from port '%s'\n", port)
	log.Fatal(server.ListenAndServe())

}
