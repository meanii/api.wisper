package utils

import (
	"errors"
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

const (
	RefreshTokenDuration = 24 * 2 // 2 days
	AccessTokenDuration  = 6      // 6 hours
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AccessTokenRawPayload create a custom type for the claims
type AccessTokenRawPayload struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Scopes   []string           `json:"scopes"`
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
	Payload *T `json:"payload"`
	jwt.RegisteredClaims
}

type JwtInstance struct {
	token     string
	tokenType string
	expiredAt time.Duration
}

// GetInstance get the instance of the JWT
func (j *JWT[T]) GetInstance() *JwtInstance {
	accessToken := &JWT[AccessTokenRawPayload]{}
	if reflect.TypeOf(j) == reflect.TypeOf(accessToken) {
		return &JwtInstance{token: SecretToken, tokenType: "access_token", expiredAt: AccessTokenDuration}
	}
	return &JwtInstance{token: RefreshToken, tokenType: "refresh_token", expiredAt: RefreshTokenDuration}
}

// GenerateToken generate the payload for the jwt
func (j *JWT[T]) GenerateToken(payload T) (string, error) {

	// Generate the token using the payload and other parameters
	generatedPayload, err := j.generatePayload(&payload, j.GetInstance().expiredAt)
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

func RefreshTokens(tokens Tokens) (*Tokens, error) {
	refreshJwt := JWT[RefreshTokenRawPayload]{} // initialize the refresh jwt
	accessJwt := JWT[AccessTokenRawPayload]{}   // initialize the access jwt

	refreshJwtPayload, err := refreshJwt.ValidateToken(tokens.RefreshToken)
	if err != nil {
		return nil, err
	}

	// check if the access token is valid
	if (refreshJwtPayload.Payload).AccessToken != tokens.AccessToken {
		return nil, errors.New("invalid access token")
	}

	// check if the access token is valid and not expired
	accessJwtPayload, err := accessJwt.ValidateToken(tokens.AccessToken)
	if err != nil {
		return nil, err
	}

	// generate the new access token
	newAccessToken, err := accessJwt.GenerateToken(*accessJwtPayload.Payload)
	if err != nil {
		return nil, err
	}

	// generate the new refresh token
	newRefreshToken, err := refreshJwt.GenerateToken(*refreshJwtPayload.Payload)
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: newAccessToken, RefreshToken: newRefreshToken}, nil
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
