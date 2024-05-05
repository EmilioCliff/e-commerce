package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrTokenExpired = errors.New("token has expired")

const (
	Footer = "E-commerceWebsite"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    int64     `json:"user_id"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	ExpireAt  time.Time `json:"expire_at"`
}

// creates a new payload for the token payload
func NewPayload(userID int64, expire_at time.Duration, admin bool) (*Payload, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid for token payload: %w", err)
	}
	return &Payload{
		ID:        id,
		UserID:    userID,
		IsAdmin:   admin,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(expire_at),
	}, nil
}

// checks the expire_at field of the payload if its expired returns an error
func Valid(payload *Payload) error {
	if time.Now().After(payload.ExpireAt) {
		return ErrTokenExpired
	}

	return nil
}
