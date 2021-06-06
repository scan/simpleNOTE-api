package main

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"

	chilogger "github.com/766b/chi-logger"
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth/v6/limiter"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func corsMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	corsConfig := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Accept-Language"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		AllowOriginFunc: func(origin string) bool {
			logger.Info("Testing origin", zap.String("origin", origin))
			if allowedOrigin, ok := os.LookupEnv("ALLOWED_CORS_ORIGIN"); ok {
				switch origin {
				case allowedOrigin, "http://localhost:8080", "http://localhost:3000":
					return true
				}
				logger.Warn("Origin denied", zap.String("origin", origin))
				return false
			}
			return true
		},
	})

	return corsConfig.Handler
}

func setupMiddleware(router chi.Router, logger *zap.Logger) {
	router.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Heartbeat("/health"),
		chilogger.NewZapMiddleware("router", logger),
		middleware.Compress(6, "test/plain", "text/html", "application/json"),
		corsMiddleware(logger),
	)

	limit := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Second})
	router.Use(tollbooth_chi.LimitHandler(limit))
}
