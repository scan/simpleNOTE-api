package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"

	"github.com/go-chi/chi"
)

//go:generate go run github.com/99designs/gqlgen

func main() {
	ctx := context.Background()

	logger, err := zapdriver.NewProduction()
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

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		logger.Error("firestore.NewClient", zap.Error(err))
		return
	}
	db, err := newFirestoreDB(client, logger)
	if err != nil {
		logger.Error("newFirestoreDB", zap.Error(err))
	}
	defer db.Close()

	router := chi.NewRouter()

	setupMiddleware(router, logger)

	router.NotFound(handleNotFound)
	router.Get("/health", handleHealth)

	/*
		graphqlHandler := s.handleGraphql().ServeHTTP
		router.Get("/graphql", graphqlHandler)
		router.Post("/graphql", graphqlHandler)
	*/

	logger.Info("application started, listening", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Error("server error", zap.Error(err))
	}
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"not found"}`))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"health":"ok"}`))
}
