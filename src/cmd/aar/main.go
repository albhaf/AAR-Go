package main

import (
	"aar"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

func main() {
	port := os.Getenv("PORT")
	databaseUrl := os.Getenv("DATABASE_URL")
	if port == "" {
		port = "8080"
	}

	var err error
	pgConfig, err := pgx.ParseURI(databaseUrl)
	aar.DB, err = pgx.Connect(pgConfig)
	if err != nil {
		log.Fatal("Error opening database: %q", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	r.HandleFunc("/missions", aar.MissionsHandler)
	r.HandleFunc("/missions/{missionId}", aar.MissionHandler)
	r.HandleFunc("/missions/{missionId}/events", aar.EventsHandler)

	// Bind to a port and pass our router in
	http.ListenAndServe(":"+port, handlers.CORS()(r))
}
