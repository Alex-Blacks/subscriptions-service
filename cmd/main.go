// @title Subscription Service API
// @version 1.0
// @description REST API for subscriptions service
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Alex-Blacks/subscriptions/docs"
	"github.com/Alex-Blacks/subscriptions/internal/config"
	"github.com/Alex-Blacks/subscriptions/internal/service"
	"github.com/Alex-Blacks/subscriptions/internal/storage"
	"github.com/Alex-Blacks/subscriptions/internal/transport"
	_ "github.com/Alex-Blacks/subscriptions/internal/transport/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	_ = godotenv.Load()

	cfg := config.MustLoad()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("error create db pool")
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("error ping db: %v", err)
	}

	st := storage.NewStorage(pool)
	srv := service.NewService(st, st)

	router := transport.NewRouter(srv)

	server := http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Printf("server started on %s", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()

	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	pool.Close()
	log.Println("server stopped")
}
