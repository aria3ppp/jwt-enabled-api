package auth

import (
	"crypto/ecdsa"
	"encoding/base64"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func ECPrivateKeyFromBase64(
	keyBase64 []byte,
	logger *zap.Logger,
) (*ecdsa.PrivateKey, error) {
	// decode base64 to bytes
	encoding := base64.StdEncoding
	keyBytes := make([]byte, encoding.DecodedLen(len(keyBase64)))
	if _, err := encoding.Decode(keyBytes, keyBase64); err != nil {
		return nil, err
	}
	// log ecdsa private key
	logger.Info("signing key decoded", zap.ByteString("key", keyBytes))
	// parse ecdsa private key
	return jwt.ParseECPrivateKeyFromPEM(keyBytes)
}
