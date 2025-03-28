package paseto

import (
	"time"

	"github.com/o1egl/paseto"
	"github.com/vlone310/bss/internal/adapter/token/maker"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (maker.Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, maker.ErrInvalidKey
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (m *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := maker.NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return m.paseto.Encrypt(m.symmetricKey, payload, nil)
}
func (m *PasetoMaker) VerifyToken(token string) (*maker.Payload, error) {
	payload := &maker.Payload{}
	err := m.paseto.Decrypt(token, m.symmetricKey, payload, nil)
	if err != nil {
		return nil, maker.ErrInvalidToken
	}

	if time.Now().After(payload.ExpiredAt) {
		return nil, maker.ErrExpiredToken
	}

	return payload, nil
}
