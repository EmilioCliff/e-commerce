package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type Maker interface {
	CreateToken(username string, duration time.Duration, admin bool) (string, error)
	VerifyToken(token string) (*Payload, string, error)
}

type PasetoMaker struct {
	Symmetrical_key []byte     `json:"symmetrical_key"`
	Paseto          *paseto.V2 `json:"paseto"`
}

func NewPasetoMaker(symmetrical_key string) (Maker, error) {
	if len(symmetrical_key) != chacha20poly1305.KeySize {
		return nil, errors.New("symmetrical key is not the key_size for the chacha20poly1305")
	}

	return &PasetoMaker{
		Symmetrical_key: []byte(symmetrical_key),
		Paseto:          paseto.NewV2(),
	}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration, admin bool) (string, error) {
	payload, err := NewPayload(username, duration, admin)
	if err != nil {
		return "", fmt.Errorf("failed to token payload")
	}

	return maker.Paseto.Encrypt(maker.Symmetrical_key, payload, Footer)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, string, error) {
	var payload Payload
	var footer string
	err := maker.Paseto.Decrypt(token, maker.Symmetrical_key, &payload, &footer)
	if err != nil {
		return nil, "", fmt.Errorf("error decrypting the token")
	}

	return &payload, footer, err
}
