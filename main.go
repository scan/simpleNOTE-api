package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
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

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		logger.Error("datastore.NewClient", zap.Error(err))
		return
	}
	db, err := newDB(client, logger)
	if err != nil {
		logger.Error("newFirestoreDB", zap.Error(err))
	}
	defer db.Close()

	router := chi.NewRouter()

	setupMiddleware(router, logger)

	router.NotFound(handleNotFound)
	router.Get("/health", handleHealth)

	graphqlHandler := handleGraphql(logger, db).ServeHTTP
	router.Get("/graphql", graphqlHandler)
	router.Post("/graphql", graphqlHandler)

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

func handleGraphql(logger *zap.Logger, db *datastoreDB) http.Handler {
	config := Config{
		Resolvers: &Resolver{logger, db},
	}

	config.Directives.IsLoggedIn = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		return next(ctx)
	}

	srv := handler.New(
		NewExecutableSchema(config),
	)

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.MultipartForm{
		MaxUploadSize: 50 * 1024 * 1024,
		MaxMemory:     64 * 1024 * 1024,
	})

	srv.Use(extension.Introspection{})

	return srv
}
