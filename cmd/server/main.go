package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/middleware"
	"github.com/mytheresa/go-hiring-challenge/models"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting application")

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize repositories
	prodRepo := models.NewProductsRepository(db)
	catRepo := models.NewCategoriesRepository(db)

	// Initialize handlers
	catalogHandler := catalog.NewCatalogHandler(prodRepo)
	categoriesHandler := categories.NewCategoriesHandler(catRepo)

	// Set up routing
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", catalogHandler.HandleGet)
	mux.HandleFunc("GET /catalog/{code}", catalogHandler.HandleGetByCode)
	mux.HandleFunc("GET /categories", categoriesHandler.HandleGet)
	mux.HandleFunc("POST /categories", categoriesHandler.HandlePost)

	// Apply middleware
	handler := middleware.Recovery(logger)(middleware.Logger(logger)(mux))

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: handler,
	}

	// Start the server
	go func() {
		logger.Info("server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", slog.String("error", err.Error()))
			log.Fatalf("Server failed: %s", err)
		}

		logger.Info("server stopped gracefully")
	}()

	<-ctx.Done()
	logger.Info("shutting down server")
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("server shutdown failed", slog.String("error", err.Error()))
	}
	stop()
	logger.Info("server shutdown complete")
}
