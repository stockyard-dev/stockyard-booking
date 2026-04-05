package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stockyard-dev/stockyard-booking/internal/server"
	"github.com/stockyard-dev/stockyard-booking/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9800"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./booking-data"
	}

	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("booking: %v", err)
	}
	defer db.Close()

	srv := server.New(db, server.DefaultLimits())

	fmt.Printf("\n  Booking — Self-hosted appointment booking and scheduling\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Questions? hello@stockyard.dev — I read every message\n\n", port, port)
	log.Printf("booking: listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
