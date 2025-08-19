package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, req)
	})
}

func handlerServerStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerShowPageViews(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	count := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(count))
}

func (cfg *apiConfig) handlerResetPageViews(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

func main() {
	const port = "8080"
	const fileRoot = "."

	var conf apiConfig
	conf.fileserverHits.Store(0)

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		conf.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRoot)))),
	)
	mux.HandleFunc("/metrics", conf.handlerShowPageViews)
	mux.HandleFunc("/reset", conf.handlerResetPageViews)
	mux.HandleFunc("/healthz", handlerServerStatus)

	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: mux,
	}

	fmt.Printf("Serving from port '%s'\n", port)
	log.Fatal(server.ListenAndServe())
}
