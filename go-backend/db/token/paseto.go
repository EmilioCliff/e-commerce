package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type Maker interface {
	CreateToken(user_id int64, duration time.Duration, admin string) (string, Payload, error)
	VerifyToken(token string) (*Payload, string, error)
}

type PasetoMaker struct {
	SymmetricalKey []byte     `json:"symmetrical_key"`
	Paseto         *paseto.V2 `json:"paseto"`
}

// Return a PasetoMaker that has all the methods of the Maker interface
func NewPasetoMaker(symmetricalKey string) (Maker, error) {
	if len(symmetricalKey) != chacha20poly1305.KeySize {
		return nil, errors.New("symmetrical key is not the key_size for the chacha20poly1305")
	}

	return &PasetoMaker{
		SymmetricalKey: []byte(symmetricalKey),
		Paseto:         paseto.NewV2(),
	}, nil
}

// Generates a new token
func (maker *PasetoMaker) CreateToken(userID int64, duration time.Duration, admin string) (string, Payload, error) {
	var isAdmin bool
	if admin == "admin" {
		isAdmin = true
	}

	if admin == "user" {
		isAdmin = false
	}

	payload, err := NewPayload(userID, duration, isAdmin)
	if err != nil {
		return "", Payload{}, fmt.Errorf("failed to token payload")
	}

	token, err := maker.Paseto.Encrypt(maker.SymmetricalKey, payload, Footer)

	return token, *payload, err
}

// Decrypts a token and verifys its validity
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, string, error) {
	var payload Payload
	var footer string
	err := maker.Paseto.Decrypt(token, maker.SymmetricalKey, &payload, &footer)
	if err != nil {
		return nil, "", fmt.Errorf("error decrypting the token")
	}

	return &payload, footer, err
}
