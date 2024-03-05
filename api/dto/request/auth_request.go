package request

type AuthRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
