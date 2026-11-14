package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const ChargeTTL = 30 * time.Minute

type ChargeStatus string

const (
	StatusPending  ChargeStatus = "pending"
	StatusAccepted ChargeStatus = "accepted"
	StatusRejected ChargeStatus = "rejected"
	StatusExpired  ChargeStatus = "expired"
)

type Charge struct {
	ID          string       `json:"id"`
	PosID       int          `json:"pos_id"`
	Amount      float64      `json:"amount"`
	CallbackURL string       `json:"callback_url"`
	Status      ChargeStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	ExpiresAt   time.Time    `json:"expires_at"`
}

type Store struct {
	rdb *redis.Client
}

func New(redisURL string) (*Store, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Store{rdb: rdb}, nil
}

func (s *Store) SaveCharge(ctx context.Context, c *Charge) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return s.rdb.SetEx(ctx, chargeKey(c.ID), string(data), ChargeTTL).Err()
}

func (s *Store) GetCharge(ctx context.Context, id string) (*Charge, error) {
	raw, err := s.rdb.Get(ctx, chargeKey(id)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var c Charge
	if err := json.Unmarshal([]byte(raw), &c); err != nil {
		return nil, err
	}
	if time.Now().After(c.ExpiresAt) {
		c.Status = StatusExpired
	}
	return &c, nil
}

func (s *Store) UpdateChargeStatus(ctx context.Context, id string, status ChargeStatus) error {
	c, err := s.GetCharge(ctx, id)
	if err != nil || c == nil {
		return err
	}
	c.Status = status
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	ttl := time.Until(c.ExpiresAt)
	if ttl <= 0 {
		return nil
	}
	return s.rdb.SetEx(ctx, chargeKey(id), string(data), ttl).Err()
}

func (s *Store) Close() error {
	return s.rdb.Close()
}

func chargeKey(id string) string {
	return "proxypay-sandbox:charge:" + id
}
