package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meanii/api.wisper/configs"
	"time"
)

var (
	SecretToken  = configs.GetConfig().SecretToken
	RefreshToken = configs.GetConfig().RefreshToken
)

type JWT struct{}

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func (j *JWT) GenerateAccessToken(payload *Payload, hours time.Duration) (string, error) {
	return j.generateToken(payload, SecretToken, hours)
}

func (j *JWT) GenerateRefreshToken(payload *Payload, hours time.Duration) (string, error) {
	return j.generateToken(payload, RefreshToken, hours)
}

func (j *JWT) ValidateAccessToken(tokenString string) (*Payload, error) {
	return j.validateToken(tokenString, SecretToken)
}

func (j *JWT) ValidateRefreshToken(tokenString string) (*Payload, error) {
	return j.validateToken(tokenString, RefreshToken)
}

func (j *JWT) generateToken(payload *Payload, secretToken string, hours time.Duration) (string, error) {
	payload, err := newPayload(payload.Username, time.Hour*hours)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(secretToken))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) validateToken(tokenString string, secretToken string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretToken), nil
	})
	if err != nil {
		return nil, err
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, err
	}

	return payload, nil
}

func newPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       tokenID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "whisper",
			ID:        tokenID.String(),
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	return payload, nil
}
