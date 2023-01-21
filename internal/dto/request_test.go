package dto_test

import (
	"testing"

	"jwt-enabled-api/internal/dto"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func TestTokenBody_Validate(t *testing.T) {
	uuid, err := uuid.NewV4()
	require.NoError(t, err)

	testCases := []struct {
		name     string
		params   dto.UserRefreshTokenRequest
		expError error
	}{
		{
			name:   "tc1",
			params: dto.UserRefreshTokenRequest{},
			expError: validation.Errors{
				"token": validation.ErrRequired,
			},
		},
		{
			name: "tc2",
			params: dto.UserRefreshTokenRequest{
				Token: "invlid-uuid-string",
			},
			expError: validation.Errors{
				"token": is.ErrUUID,
			},
		},
		{
			name: "tc3",
			params: dto.UserRefreshTokenRequest{
				Token: uuid.String(),
			},
			expError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			require.Equal(tc.expError, tc.params.Validate())
		})
	}
}
