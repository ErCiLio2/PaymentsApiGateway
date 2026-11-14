package main

import (
	"log"
	"net/http"

	"github.com/somatec/payments-gateway/proxypay/internal/config"
	"github.com/somatec/payments-gateway/proxypay/internal/sandbox"
	"github.com/somatec/payments-gateway/proxypay/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[proxypay] configuração inválida: %v", err)
	}

	s, err := store.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("[proxypay] falha ao ligar ao Redis: %v", err)
	}
	defer s.Close()

	h := sandbox.NewHandler(cfg, s)

	mux := http.NewServeMux()

	// Proxypay OPG API (autenticada com Basic Auth — usada pelo eDuca ou outro cliente)
	mux.HandleFunc("POST /opg/v1/charges", h.CreateCharge)

	// API interna usada pela página web de aprovação
	mux.HandleFunc("GET /api/charges/{id}", h.GetCharge)
	mux.HandleFunc("GET /api/charges/{id}/qrcode.png", h.GetQRCode)
	mux.HandleFunc("POST /api/charges/{id}/pay", h.PayCharge)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	handler := corsMiddleware(mux, cfg.WebURL)

	log.Printf("[proxypay] a escutar em :%s (web: %s)", cfg.Port, cfg.WebURL)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("[proxypay] servidor falhou: %v", err)
	}
}

func corsMiddleware(next http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
