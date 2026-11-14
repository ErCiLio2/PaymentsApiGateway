package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/somatec/kashe/internal/config"
	"github.com/somatec/kashe/internal/db"
	"github.com/somatec/kashe/internal/payments/proxypay"
	"github.com/somatec/kashe/internal/store/charges"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[kashe] configuração inválida: %v", err)
	}

	// PostgreSQL
	pool, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[kashe] falha ao ligar ao PostgreSQL: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(context.Background(), pool); err != nil {
		log.Fatalf("[kashe] falha nas migrações: %v", err)
	}

	// Redis
	chargeStore, err := charges.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("[kashe] falha ao ligar ao Redis: %v", err)
	}
	defer chargeStore.Close()

	ppHandler := proxypay.NewHandler(cfg, chargeStore)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(corsMiddleware(cfg.WebURL))

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","version":"2.0.0"}`))
	})

	// Proxypay OPG API (Basic Auth)
	r.Post("/opg/v1/charges", ppHandler.CreateCharge)

	// API pública de charges
	r.Get("/api/charges/{id}", ppHandler.GetCharge)
	r.Get("/api/charges/{id}/qrcode.png", ppHandler.GetQRCode)
	r.Post("/api/charges/{id}/pay", ppHandler.PayCharge)

	// Admin API (Bearer ADMIN_KEY)
	r.Get("/api/admin/charges", ppHandler.ListCharges)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[kashe] a escutar em :%s (web: %s)", cfg.Port, cfg.WebURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[kashe] servidor falhou: %v", err)
		}
	}()

	<-quit
	log.Println("[kashe] a encerrar...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
