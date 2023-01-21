package auth

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

type Interface interface {
	GenerateJwtToken(*Payload) (token string, expiresAt time.Time, err error)
	GenerateRefreshToken() (token string, expiresAt time.Time, err error)
}

type Payload struct {
	UserID int `json:"user_id"`
	// more fields?
}

type Auth struct {
	signingKey                *ecdsa.PrivateKey
	signingMethod             jwt.SigningMethod
	jwtExpiresInSecs          int
	refreshTokenExpiresInSecs int
}

var _ Interface = &Auth{}

func NewAuth(
	signingKey *ecdsa.PrivateKey,
	jwtExpiresInSecs int,
	refreshTokenExpiresInSecs int,
) *Auth {
	return &Auth{
		signingKey:                signingKey,
		signingMethod:             jwt.SigningMethodES256,
		jwtExpiresInSecs:          jwtExpiresInSecs,
		refreshTokenExpiresInSecs: refreshTokenExpiresInSecs,
	}
}

type jwtClaims struct {
	*Payload
	jwt.RegisteredClaims
}

func (auth *Auth) GenerateJwtToken(
	payload *Payload,
) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().
		Add(time.Second * time.Duration(auth.jwtExpiresInSecs))
	// set payload
	claims := jwtClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	// sign token
	token, err = jwt.NewWithClaims(auth.signingMethod, claims).
		SignedString(auth.signingKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

func (auth *Auth) GenerateRefreshToken() (token string, expiresAt time.Time, err error) {
	// generate a UUIDv4
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt = time.Now().
		Add(time.Second * time.Duration(auth.refreshTokenExpiresInSecs))
	return uuid.String(), expiresAt, nil
}

func (auth *Auth) ParseJwtToken(tokenString string) (*Payload, error) {
	// prepare key provider function
	keyFunc := func(t *jwt.Token) (any, error) {
		if auth.signingMethod.Alg() != t.Method.Alg() {
			return nil, fmt.Errorf(
				"unexpected jwt signing method=%v",
				t.Method.Alg(),
			)
		}
		return &auth.signingKey.PublicKey, nil
	}
	// parse token string
	var claims jwtClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims.Payload, nil
}
