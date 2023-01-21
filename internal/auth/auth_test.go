package auth_test

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"
	"testing"
	"time"

	"jwt-enabled-api/internal/auth"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var privateKey *ecdsa.PrivateKey

func init() {
	var err error
	keyBytes, err := os.ReadFile(filepath.Join("testdata", "key.pem"))
	if err != nil {
		panic(err)
	}
	privateKey, err = jwt.ParseECPrivateKeyFromPEM(keyBytes)
	if err != nil {
		panic(err)
	}
}

func TestAuth_GenerateJwtToken(t *testing.T) {
	require := require.New(t)

	type fields struct {
		signingKey                *ecdsa.PrivateKey
		jwtExpiresInSecs          int
		refreshTokenExpiresInSecs int
	}
	type args struct {
		payload *auth.Payload
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				signingKey:                privateKey,
				jwtExpiresInSecs:          10,
				refreshTokenExpiresInSecs: 100,
			},
			args: args{
				payload: &auth.Payload{UserID: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := auth.NewAuth(
				tt.fields.signingKey,
				tt.fields.jwtExpiresInSecs,
				tt.fields.refreshTokenExpiresInSecs,
			)
			t0 := time.Now()
			gotToken, gotExpiresAt, err := auth.GenerateJwtToken(
				tt.args.payload,
			)
			t1 := time.Now()
			if tt.wantErr {
				require.Error(err)
				require.Empty(gotToken)
				require.Empty(gotExpiresAt)
			} else {
				require.NoError(err)
				expJwtToken, err := jwt.Parse(gotToken, func(t *jwt.Token) (any, error) { return &privateKey.PublicKey, nil })
				require.NoError(err)
				expClaims, ok := expJwtToken.Claims.(jwt.MapClaims)
				require.True(ok)
				expExpiresAtFloat64, ok := expClaims["exp"].(float64)
				require.True(ok)
				expExpiresAt := time.Unix(int64(expExpiresAtFloat64), 0)
				require.GreaterOrEqual(expExpiresAt, t0)
				require.LessOrEqual(expExpiresAt, t1.Add(time.Second*time.Duration(tt.fields.jwtExpiresInSecs)))
				require.Equal(int64(expExpiresAtFloat64), gotExpiresAt.Unix())
				expUserIDFloat64, ok := expClaims["user_id"].(float64)
				require.True(ok)
				require.Equal(int(expUserIDFloat64), tt.args.payload.UserID)
			}
		})
	}
}

func TestAuth_GenerateRefreshToken(t *testing.T) {
	require := require.New(t)

	type fields struct {
		signingKey                *ecdsa.PrivateKey
		jwtExpiresInSecs          int
		refreshTokenExpiresInSecs int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				signingKey:                privateKey,
				jwtExpiresInSecs:          10,
				refreshTokenExpiresInSecs: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := auth.NewAuth(
				tt.fields.signingKey,
				tt.fields.jwtExpiresInSecs,
				tt.fields.refreshTokenExpiresInSecs,
			)
			t0 := time.Now()
			gotToken, gotExpiresAt, err := auth.GenerateRefreshToken()
			t1 := time.Now()
			if tt.wantErr {
				require.Error(err)
				require.Empty(gotToken)
				require.Empty(gotExpiresAt)
			} else {
				require.NoError(err)
				require.NotEmpty(gotToken)
				_, err := uuid.Parse(gotToken)
				require.NoError(err)
				require.GreaterOrEqual(gotExpiresAt, t0)
				require.LessOrEqual(gotExpiresAt, t1.Add(time.Second*time.Duration(tt.fields.refreshTokenExpiresInSecs)))
			}
		})
	}
}

func TestAuth_ParseJwtToken(t *testing.T) {
	require := require.New(t)

	payload := &auth.Payload{UserID: 1}

	type fields struct {
		signingKey                *ecdsa.PrivateKey
		jwtExpiresInSecs          int
		refreshTokenExpiresInSecs int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				signingKey:                privateKey,
				jwtExpiresInSecs:          10,
				refreshTokenExpiresInSecs: 100,
			},
		},
		{
			name: "expired token",
			fields: fields{
				signingKey:                privateKey,
				jwtExpiresInSecs:          -10,
				refreshTokenExpiresInSecs: 100,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := auth.NewAuth(
				tt.fields.signingKey,
				tt.fields.jwtExpiresInSecs,
				tt.fields.refreshTokenExpiresInSecs,
			)

			tokenString, _, err := auth.GenerateJwtToken(payload)
			require.NoError(err)
			gotPayload, err := auth.ParseJwtToken(tokenString)
			if tt.wantErr {
				require.Error(err)
				require.Nil(gotPayload)
			} else {
				require.NoError(err)
				require.Equal(gotPayload, payload)
			}
		})
	}
}
