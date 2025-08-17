package main

import (
	"fmt"
	"log"
	"net/http"
)

func handlerServerStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	const port = "8080"
	const fileRoot = "."

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot))),
	)
	mux.HandleFunc("/healthz", handlerServerStatus)

	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: mux,
	}

	fmt.Printf("Serving from port '%s'\n", port)
	log.Fatal(server.ListenAndServe())

}
