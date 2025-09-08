package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/dothedada/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

var BAD_WORDS = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w, req)
	})
}

func handlerChirpValidation(w http.ResponseWriter, req *http.Request) {
	text, isValid := isValidRequest(req)
	if isValid.error != nil {
		resWithErr(w, isValid)
		return
	}

	cleanMsg := profanityCleaner(text, BAD_WORDS)
	resJson(w, http.StatusOK, resValid{CleanedBody: cleanMsg})
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
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Cannot connect to DB: %s", err)
	}

	const port = "8080"
	const fileRoot = "."

	conf := apiConfig{
		dbQueries: database.New(db),
	}
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
