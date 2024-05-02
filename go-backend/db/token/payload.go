package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrTokenExpired = errors.New("Token has expired")

const (
	Footer = "E-commerce Website"
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	ExpireAt  time.Time `json:"expire_at"`
}

func NewPayload(username string, expire_at time.Duration, admin bool) (*Payload, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid for token payload: %w", err)
	}
	return &Payload{
		Id:        id,
		Username:  username,
		IsAdmin:   admin,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(expire_at),
	}, nil
}

func Valid(payload *Payload) error {
	if time.Now().After(payload.ExpireAt) {
		return ErrTokenExpired
	}

	return nil
}
