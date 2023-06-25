package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/meanii/api.wisper/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type JwtInstance struct {
	token     string
	tokenType string
}

// GetInstance get the instance of the JWT
func (j *JWT[T]) GetInstance() *JwtInstance {
	accessToken := &JWT[AccessTokenRawPayload]{}
	if reflect.TypeOf(j) == reflect.TypeOf(accessToken) {
		return &JwtInstance{token: SecretToken, tokenType: "access_token"}
	}
	return &JwtInstance{token: RefreshToken, tokenType: "refresh_token"}
}

// GenerateToken generate the payload for the jwt
func (j *JWT[T]) GenerateToken(payload T, hours time.Duration) (string, error) {
	// Generate the token using the payload and other parameters
	generatedPayload, err := j.generatePayload(&payload, hours)
	if err != nil {
		fmt.Println("Something went wrong while generating payload")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, generatedPayload)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(j.GetInstance().token))
	if err != nil {
		fmt.Println("Something went wrong while generating tokenString")
	}
	return tokenString, nil
}

// ValidateToken validate the token
func (j *JWT[T]) ValidateToken(tokenString string) (*JwtPayload[T], error) {

	token, err := jwt.ParseWithClaims(tokenString, &JwtPayload[T]{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.GetInstance().token), nil
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

// generatePayload generate the payload for the jwt
func (j *JWT[T]) generatePayload(payload *T, duration time.Duration) (*JwtPayload[T], error) {
	jp := JwtPayload[T]{} // initialize the jwtPayload
	jp.Payload = payload
	// set the default claims
	jp.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "wisper",
		Subject:   j.GetInstance().tokenType,
	}
	return &jp, nil
}
