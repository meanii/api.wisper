package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/meanii/api.wisper/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"reflect"
	"time"
)

var (
	SecretToken  = configs.GetConfig().SecretToken
	RefreshToken = configs.GetConfig().RefreshToken
)

// AccessTokenRawPayload create a custom type for the claims
type AccessTokenRawPayload struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
}

// RefreshTokenRawPayload create a custom type for the claims
type RefreshTokenRawPayload struct {
	AccessToken string `json:"access_token"`
}

type RawPayload interface {
	AccessTokenRawPayload | RefreshTokenRawPayload
}

// JWT create a generic JWT struct
type JWT[T RawPayload] struct{}

type JwtPayload[T RawPayload] struct {
	Payload *T
	jwt.RegisteredClaims
}

func (j *JWT[T]) GetInstance() string {
	accessToken := &JWT[AccessTokenRawPayload]{}
	if reflect.TypeOf(j) == reflect.TypeOf(accessToken) {
		return SecretToken
	}
	return RefreshToken
}

func (j *JWT[T]) GenerateToken(payload T, hours time.Duration) (string, error) {
	// Generate the token using the payload and other parameters
	generatedPayload, err := generatePayload(&payload, hours)
	if err != nil {
		log.Fatalln("Something went wrong while generating payload")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, generatedPayload)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(j.GetInstance()))
	if err != nil {
		log.Fatalln("Something went wrong while generating tokenString")
	}
	return tokenString, nil
}

func (j *JWT[T]) ValidateToken(tokenString string) (*JwtPayload[T], error) {

	token, err := jwt.ParseWithClaims(tokenString, &JwtPayload[T]{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.GetInstance()), nil
	})

	if err != nil {
		return nil, err
	}
	payload, ok := token.Claims.(*JwtPayload[T])
	if !ok {
		return nil, err
	}
	return payload, nil
}

func generatePayload[T RawPayload](payload *T, duration time.Duration) (*JwtPayload[T], error) {
	jp := JwtPayload[T]{} // initialize the jwtPayload
	jp.Payload = payload
	// set the default claims
	jp.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "wisper",
		Subject:   "access_token",
	}
	return &jp, nil
}
