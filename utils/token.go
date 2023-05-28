package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/meanii/api.wisper/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	SecretToken  = configs.GetConfig().SecretToken
	RefreshToken = configs.GetConfig().RefreshToken
)

type JWT struct{}

type Payload struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}

type PayloadJwt struct {
	Payload
	jwt.RegisteredClaims
}

func (j *JWT) GenerateAccessToken(payload *Payload, hours time.Duration) (string, error) {
	return j.generateToken(payload, SecretToken, hours)
}

func (j *JWT) GenerateRefreshToken(payload *Payload, hours time.Duration) (string, error) {
	return j.generateToken(payload, RefreshToken, hours)
}

func (j *JWT) ValidateAccessToken(tokenString string) (*PayloadJwt, error) {
	return j.validateToken(tokenString, SecretToken)
}

func (j *JWT) ValidateRefreshToken(tokenString string) (*PayloadJwt, error) {
	return j.validateToken(tokenString, RefreshToken)
}

func (j *JWT) RefreshToken(refreshToken string, hours time.Duration) (string, error) {
	payload, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	return j.GenerateAccessToken(&payload.Payload, time.Hour*hours)
}

func (j *JWT) generateToken(payload *Payload, secretToken string, hours time.Duration) (string, error) {
	jwtPayload, err := newPayload(payload, time.Hour*hours)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)
	tokenString, err := token.SignedString([]byte(secretToken))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) validateToken(tokenString string, secretToken string) (*PayloadJwt, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PayloadJwt{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretToken), nil
	})
	if err != nil {
		return nil, err
	}

	payload, ok := token.Claims.(*PayloadJwt)
	if !ok {
		return nil, err
	}

	return payload, nil
}

func newPayload(payload *Payload, duration time.Duration) (*PayloadJwt, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	jwtPayload := &PayloadJwt{
		Payload: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "whisper",
			ID:        tokenID.String(),
			Subject:   "access_token",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}
	return jwtPayload, nil
}
