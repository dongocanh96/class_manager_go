package token

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTMaker struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewJWTMaker(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (*JWTMaker, error) {
	return &JWTMaker{privateKey, publicKey}, nil
}

func (maker *JWTMaker) CreateToken(userid int64, username string, isteacher bool, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userid, username, isteacher, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	token, err := jwtToken.SignedString(maker.privateKey)
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, ErrInvalidToken
		}
		return maker.publicKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
