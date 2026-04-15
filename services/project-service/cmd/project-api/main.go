package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Daedalus/project-service/internal/adapters/handlers"
	"github.com/Daedalus/project-service/internal/adapters/repositories"
	"github.com/Daedalus/project-service/internal/core/services"
)

func main() {
	// 1. Load .env (Kliops pattern)
	godotenv.Load()

	// 2. Connect to PostgreSQL (pgxpool — Kliops pattern)
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "postgres://daedalus:daedalus@localhost:5432/daedalus?sslmode=disable"
	}

	dbPool, err := pgxpool.New(context.Background(), dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// 3. Wire dependencies (bottom-up: repo → service → handler)
	repo := repositories.NewProjectPostgres(dbPool)
	service := services.NewProjectService(repo)
	handler := handlers.NewProjectHandler(service)

	// 4. Setup routes (http.ServeMux — Kliops pattern)
	mux := http.NewServeMux()

	// Health check (public, Kubernetes probe)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// Project API routes
	handler.RegisterRoutes(mux)

	// Apply middleware stack
	var httpHandler http.Handler = mux
	httpHandler = handlers.CORSMiddleware(httpHandler)
	httpHandler = handlers.RequestLoggingMiddleware(httpHandler)

	// 5. Start HTTP server (Kliops pattern)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      httpHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Starting Daedalus Project Service on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// 6. Graceful shutdown (Kliops pattern)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
