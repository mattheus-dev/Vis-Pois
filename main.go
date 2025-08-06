package main

import (
	"log"
	"net/http"
	"os"

	httpHandler "vis-pois/internal/adapters/http"
	"vis-pois/internal/application/services"
)

func main() {
	logger := log.New(os.Stdout, "[CSV-PROCESSOR] ", log.LstdFlags)
	logger.Println("Starting CSV processor application...")

	csvService := services.NewCSVService(logger)

	csvHandler := httpHandler.NewCSVHandler(csvService, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/leitura/teste", csvHandler.HandleUpload)

	port := "8080"
	logger.Printf("Server starting on port %s...", port)
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}