package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	database "github.com/ellielle/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment")
	}

	db, err := database.NewDBConnection("database.json")
	if err != nil {
		log.Fatal(err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
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
	// PUT endpoint for user updates
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)
	// POST endpoint for users to login
	mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)
	// POST endpoint for refreshing access tokens
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerUsersRefresh)
	// POST endpoint to revoke access token with refresh token
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerTokensRevoke)
	// DELETE endpoint to remove chirps
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerChirpsDelete)

	// POST endpoint for "Polka" user upgraded events
	mux.HandlerFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhooks)

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
