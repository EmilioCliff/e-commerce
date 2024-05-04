package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type Maker interface {
	CreateToken(user_id int64, duration time.Duration, admin bool) (string, error)
	VerifyToken(token string) (*Payload, string, error)
}

type PasetoMaker struct {
	Symmetrical_key []byte     `json:"symmetrical_key"`
	Paseto          *paseto.V2 `json:"paseto"`
}

// Return a PasetoMaker that has all the methods of the Maker interface
func NewPasetoMaker(symmetrical_key string) (Maker, error) {
	if len(symmetrical_key) != chacha20poly1305.KeySize {
		return nil, errors.New("symmetrical key is not the key_size for the chacha20poly1305")
	}

	return &PasetoMaker{
		Symmetrical_key: []byte(symmetrical_key),
		Paseto:          paseto.NewV2(),
	}, nil
}

// Generates a new token
func (maker *PasetoMaker) CreateToken(user_id int64, duration time.Duration, admin bool) (string, error) {
	payload, err := NewPayload(user_id, duration, admin)
	if err != nil {
		return "", fmt.Errorf("failed to token payload")
	}

	return maker.Paseto.Encrypt(maker.Symmetrical_key, payload, Footer)
}

// Decrypts a token and verifys its validity
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, string, error) {
	var payload Payload
	var footer string
	err := maker.Paseto.Decrypt(token, maker.Symmetrical_key, &payload, &footer)
	if err != nil {
		return nil, "", fmt.Errorf("error decrypting the token")
	}

	return &payload, footer, err
}
