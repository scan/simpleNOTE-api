package main

import (
	"log"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot initialise zap logger: %v", err)
	}
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		logger.Fatal("GOOGLE_CLOUD_PROJECT must be set")
	}
}
