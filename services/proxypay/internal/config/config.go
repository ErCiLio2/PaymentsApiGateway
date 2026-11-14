package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	RedisURL       string
	APIKey         string // chave que os clientes usam para criar charges
	PosID          int
	WebhookSecret  string // segredo para assinar os callbacks enviados ao cliente
	BaseURL        string // URL pública desta instância (para deeplinks e qrcode_url)
	WebURL         string // URL do painel web (SvelteKit)
	SandboxEnabled bool
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		APIKey:         getEnv("API_KEY", "dev-key"),
		WebhookSecret:  getEnv("WEBHOOK_SECRET", ""),
		BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
		WebURL:         getEnv("WEB_URL", "http://localhost:5173"),
		SandboxEnabled: getEnv("SANDBOX_ENABLED", "true") == "true",
	}

	posIDStr := getEnv("POS_ID", "1")
	posID, err := strconv.Atoi(posIDStr)
	if err != nil {
		return nil, fmt.Errorf("POS_ID inválido: %w", err)
	}
	cfg.PosID = posID

	if cfg.WebhookSecret == "" {
		return nil, fmt.Errorf("WEBHOOK_SECRET é obrigatório (mínimo 32 caracteres)")
	}
	if len(cfg.WebhookSecret) < 32 {
		return nil, fmt.Errorf("WEBHOOK_SECRET demasiado curto (mínimo 32 caracteres)")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
