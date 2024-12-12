package model

type Auth struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type AuthResponse struct {
	User         UserResponse   `json:"user"`
	Wallet       WalletResponse `json:"wallet"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
}

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
