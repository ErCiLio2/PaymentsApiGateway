package sandbox

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"github.com/somatec/payments-gateway/proxypay/internal/config"
	"github.com/somatec/payments-gateway/proxypay/internal/store"

	"crypto/rand"
	"encoding/hex"
)

func newID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type Handler struct {
	cfg   *config.Config
	store *store.Store
}

func NewHandler(cfg *config.Config, s *store.Store) *Handler {
	return &Handler{cfg: cfg, store: s}
}

// POST /opg/v1/charges — endpoint compatível com a API real do Proxypay OPG
func (h *Handler) CreateCharge(w http.ResponseWriter, r *http.Request) {
	if !h.validateBasicAuth(r) {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		PosID       int     `json:"pos_id"`
		Amount      float64 `json:"amount"`
		CallbackURL string  `json:"callback_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid_body")
		return
	}
	if body.PosID <= 0 || body.Amount <= 0 || body.CallbackURL == "" {
		jsonError(w, http.StatusUnprocessableEntity, "missing_fields")
		return
	}

	id := newID()
	now := time.Now().UTC()
	expiresAt := now.Add(store.ChargeTTL)
	deeplink := fmt.Sprintf("%s/%s", h.cfg.WebURL, id)

	c := &store.Charge{
		ID:          id,
		PosID:       body.PosID,
		Amount:      body.Amount,
		CallbackURL: body.CallbackURL,
		Status:      store.StatusPending,
		CreatedAt:   now,
		ExpiresAt:   expiresAt,
	}
	if err := h.store.SaveCharge(r.Context(), c); err != nil {
		log.Printf("[sandbox] erro ao guardar charge: %v", err)
		jsonError(w, http.StatusInternalServerError, "store_error")
		return
	}

	qrcodeURL := fmt.Sprintf("%s/api/charges/%s/qrcode.png", h.cfg.BaseURL, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":                id,
		"qrcode_url":        qrcodeURL,
		"deeplink_redirect": deeplink,
		"expires_at":        expiresAt.Format(time.RFC3339),
	})
}

// GET /api/charges/{id} — detalhes do charge (usado pela página web de aprovação)
func (h *Handler) GetCharge(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := h.store.GetCharge(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "store_error")
		return
	}
	if c == nil {
		jsonError(w, http.StatusNotFound, "charge_not_found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":         c.ID,
		"amount":     c.Amount,
		"status":     string(c.Status),
		"expires_at": c.ExpiresAt.Format(time.RFC3339),
	})
}

// GET /api/charges/{id}/qrcode.png — gera QR code PNG do deeplink
func (h *Handler) GetQRCode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := h.store.GetCharge(r.Context(), id)
	if err != nil || c == nil {
		http.NotFound(w, r)
		return
	}
	deeplink := fmt.Sprintf("%s/%s", h.cfg.WebURL, id)
	png, err := qrcode.Encode(deeplink, qrcode.Medium, 300)
	if err != nil {
		http.Error(w, "qrcode error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=1800")
	w.Write(png)
}

// POST /api/charges/{id}/pay — simula pagamento (aceitar ou recusar)
func (h *Handler) PayCharge(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Decision string `json:"decision"` // "accepted" | "rejected"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid_body")
		return
	}
	if body.Decision != "accepted" && body.Decision != "rejected" {
		jsonError(w, http.StatusUnprocessableEntity, "invalid_decision")
		return
	}

	c, err := h.store.GetCharge(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "store_error")
		return
	}
	if c == nil {
		jsonError(w, http.StatusNotFound, "charge_not_found")
		return
	}
	if c.Status != store.StatusPending {
		jsonError(w, http.StatusConflict, fmt.Sprintf("charge_%s", c.Status))
		return
	}

	status := store.StatusAccepted
	if body.Decision == "rejected" {
		status = store.StatusRejected
	}
	if err := h.store.UpdateChargeStatus(r.Context(), id, status); err != nil {
		jsonError(w, http.StatusInternalServerError, "store_error")
		return
	}

	go h.sendCallback(c, body.Decision)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

// sendCallback envia o webhook ao callbackUrl do charge (simula o Proxypay)
func (h *Handler) sendCallback(c *store.Charge, decision string) {
	txID := newID()
	payload, _ := json.Marshal(map[string]any{
		"id":        txID,     // ID da transacção (eventId para idempotência no eDuca)
		"charge_id": c.ID,    // liga ao PaymentTransaction.reference no eDuca
		"status":    decision, // "accepted" | "rejected"
		"amount":    c.Amount,
		"currency":  "AOA",
	})

	sig := hmacHex(payload, []byte(h.cfg.WebhookSecret))
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, c.CallbackURL, bytes.NewReader(payload))
	if err != nil {
		log.Printf("[sandbox] erro ao criar pedido callback: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-proxypay-signature", sig)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[sandbox] erro ao enviar callback para %s: %v", c.CallbackURL, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("[sandbox] callback %s → %s: %d", c.ID, c.CallbackURL, resp.StatusCode)
}

// GET /api/admin/charges — lista todas as charges (protegido por Bearer token)
func (h *Handler) ListCharges(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token != h.cfg.AdminKey {
		jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	charges, err := h.store.ListCharges(r.Context())
	if err != nil {
		log.Printf("[sandbox] erro ao listar charges: %v", err)
		jsonError(w, http.StatusInternalServerError, "store_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(charges)
}

func (h *Handler) validateBasicAuth(r *http.Request) bool {
	user, _, ok := r.BasicAuth()
	return ok && user == h.cfg.APIKey
}

func hmacHex(data, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

func jsonError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": code})
}
