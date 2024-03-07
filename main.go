package main

import (
	"flag"
	"log"
	"net/http"

	database "github.com/ellielle/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	db, err := database.NewDBConnection("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	// Wipe test database in debug mode
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		db.DebugWipeTestDatabase()
		log.Print("Database deleted successfully...")
	}

	// Create new request multiplexer
	mux := http.NewServeMux()
	// Fileserver for handling static pages
	fileseverHandler := apiCfg.middelwareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fileseverHandler)

	// API endpoints under the api subroute
	// Health check endpoing
	mux.HandleFunc("GET /api/healthz", healthzResponseHandler)
	// Page hit count reset endpoint
	mux.HandleFunc("GET /api/reset", apiCfg.handlerMetricsReset)
	// GET endpoint for retrieving all chirps
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGetAll)
	// GET endpoint for retrieving a single chirp
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	// POST endpoint to submit "Chirps". Chrips must be 140 chars or less, and should be in JSON
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	// POST endpoint to submit an email and create a new User
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	// POST endpoint for users to login
	mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)

	// Admin route, which only contains the metrics endpoint for now
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricsResponse)

	// Wrap mux in CORS headers and serve
	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
