package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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

func handlerChirpValidation(w http.ResponseWriter, req *http.Request) {
	type resError struct {
		Error string `json:"error"`
	}
	type resBody struct {
		Valid bool `json:"valid"`
	}

	w.Header().Set("Content-Type", "application/json")

	isValid := isValidRequest(req)
	if isValid.error != nil {
		errMsg := resError{
			Error: isValid.error.Error(),
		}

		errBytes, err := json.Marshal(errMsg)
		if err != nil {
			log.Printf("Error marshaling JSON response %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(isValid.status)
		w.Write(errBytes)
		return
	}

	resMsg := resBody{
		Valid: true,
	}
	resBytes, err := json.Marshal(resMsg)
	if err != nil {
		log.Printf("Error marshaling JSON response %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}

func handlerServerStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerShowPageViews(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

	count := fmt.Sprintf(html, cfg.fileserverHits.Load())
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
	mux.HandleFunc("GET /admin/metrics", conf.handlerShowPageViews)
	mux.HandleFunc("POST /admin/reset", conf.handlerResetPageViews)
	mux.HandleFunc("GET /api/healthz", handlerServerStatus)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValidation)

	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: mux,
	}

	fmt.Printf("Serving from port '%s'\n", port)
	log.Fatal(server.ListenAndServe())
}
