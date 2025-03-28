package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vlone310/bss/internal/adapter/token/maker"
)

const minSecretKeySize = 32

var ErrInvalidSecretKey = fmt.Errorf("secret is too short, must be at least %d characters", minSecretKeySize)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (maker.Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrInvalidSecretKey
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (m *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := maker.NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, getClaims(payload))

	return jwtToken.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) VerifyToken(token string) (*maker.Payload, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, maker.ErrInvalidToken
		}

		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok || !jwtToken.Valid {
		return nil, maker.ErrInvalidToken
	}

	return &maker.Payload{
		ID:        uuid.MustParse(claims.ID),
		Username:  claims.Subject,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiredAt: claims.ExpiresAt.Time,
	}, nil
}

func getClaims(payload *maker.Payload) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		Subject:   payload.Username,
		ID:        payload.ID.String(),
	}
}
