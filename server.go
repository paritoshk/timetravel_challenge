package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/api"
	"github.com/rainbowmga/timetravel/service"
)

func logError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func main() {
	router := mux.NewRouter()

	sqliteService, err := service.NewSQLiteRecordService("./records.db")
	if err != nil {
		log.Fatalf("Failed to create SQLite service: %v", err)
	}

	apiHandler := api.NewAPI(sqliteService)

	apiHandler.CreateRoutes(router)

	router.Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		logError(err)
	})

	address := "127.0.0.1:8000"
	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is running on %s", address)
	log.Fatal(srv.ListenAndServe())
}
