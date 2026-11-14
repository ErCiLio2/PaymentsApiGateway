package charges

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ChargeTTL       = 30 * time.Minute
	chargesIndexKey = "kashe:charges"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusAccepted Status = "accepted"
	StatusRejected Status = "rejected"
	StatusExpired  Status = "expired"
)

type Charge struct {
	ID          string    `json:"id"`
	PosID       int       `json:"pos_id"`
	Amount      float64   `json:"amount"`
	CallbackURL string    `json:"callback_url"`
	Status      Status    `json:"status"`
	EntityID    string    `json:"entity_id"`
	OTPMethod   string    `json:"otp_method"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
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

func (s *Store) Close() error {
	return s.rdb.Close()
}

func (s *Store) Save(ctx context.Context, c *Charge) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	pipe := s.rdb.Pipeline()
	pipe.SetEx(ctx, key(c.ID), string(data), ChargeTTL)
	pipe.ZAdd(ctx, chargesIndexKey, redis.Z{
		Score:  float64(c.CreatedAt.Unix()),
		Member: c.ID,
	})
	_, err = pipe.Exec(ctx)
	return err
}

func (s *Store) Get(ctx context.Context, id string) (*Charge, error) {
	raw, err := s.rdb.Get(ctx, key(id)).Result()
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

func (s *Store) UpdateStatus(ctx context.Context, id string, status Status) error {
	c, err := s.Get(ctx, id)
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
	return s.rdb.SetEx(ctx, key(id), string(data), ttl).Err()
}

func (s *Store) List(ctx context.Context) ([]*Charge, error) {
	ids, err := s.rdb.ZRevRange(ctx, chargesIndexKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []*Charge{}, nil
	}

	pipe := s.rdb.Pipeline()
	cmds := make([]*redis.StringCmd, len(ids))
	for i, id := range ids {
		cmds[i] = pipe.Get(ctx, key(id))
	}
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	result := make([]*Charge, 0, len(ids))
	for _, cmd := range cmds {
		raw, err := cmd.Result()
		if err == redis.Nil || err != nil {
			continue
		}
		var c Charge
		if err := json.Unmarshal([]byte(raw), &c); err != nil {
			continue
		}
		if time.Now().After(c.ExpiresAt) {
			c.Status = StatusExpired
		}
		result = append(result, &c)
	}
	return result, nil
}

func (s *Store) ListByEntity(ctx context.Context, entityID string) ([]*Charge, error) {
	all, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*Charge, 0)
	for _, c := range all {
		if c.EntityID == entityID {
			out = append(out, c)
		}
	}
	return out, nil
}

func key(id string) string {
	return "kashe:charge:" + id
}
