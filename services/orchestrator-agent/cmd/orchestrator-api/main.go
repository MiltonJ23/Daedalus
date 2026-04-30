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

	"github.com/Daedalus/orchestrator-agent/internal/adapters/handlers"
	"github.com/Daedalus/orchestrator-agent/internal/adapters/publishers"
	"github.com/Daedalus/orchestrator-agent/internal/adapters/repositories"
	"github.com/Daedalus/orchestrator-agent/internal/core/services"
)

func main() {
	godotenv.Load()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://daedalus:daedalus@localhost:5434/daedalus_orchestrator?sslmode=disable"
	}
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer dbPool.Close()
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("postgres ping: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	goalRepo := repositories.NewGoalPostgres(dbPool)
	taskRepo := repositories.NewTaskPostgres(dbPool)

	// PB-025: TaskPublisher port — wire a Redis Streams adapter here when
	// the module cache provides go-redis. The InMemory publisher is the
	// default fallback so the service is fully functional out of the box.
	publisher := publishers.NewInMemory()
	log.Println("Using in-memory TaskPublisher (Redis Streams adapter pending)")

	svc := services.NewOrchestratorService(goalRepo, taskRepo, publisher)
	handler := handlers.NewOrchestratorHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})
	mux.Handle("GET /metrics", handlers.MetricsHandler())
	handler.RegisterRoutes(mux)

	var h http.Handler = mux
	h = handlers.MetricsMiddleware(h)
	h = handlers.CORSMiddleware(h)
	h = handlers.RequestLoggingMiddleware(h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go func() {
		log.Printf("Starting Daedalus Orchestrator Agent on port %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
