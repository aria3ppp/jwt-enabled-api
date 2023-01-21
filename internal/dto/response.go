package dto

type JwtTokenResponse struct {
	JwtToken     string `json:"jwt_token"`
	JwtExpiresAt int64  `json:"jwt_expires_at"`
}

type TokenResponse struct {
	JwtTokenResponse
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
	UserID           int    `json:"user_id"`
}
