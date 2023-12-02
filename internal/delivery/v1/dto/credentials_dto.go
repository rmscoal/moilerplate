package dto

type SignUpRequest struct {
	FirstName   string `json:"firstName,omitempty" example:"FirstName"`
	LastName    string `json:"lastName,omitempty" example:"LastName"`
	Username    string `json:"username,omitempty" example:"Username"`
	Password    string `json:"password,omitempty" example:"verystrongpassword"`
	Email       string `json:"email,omitempty" example:"email@email.com"`
	PhoneNumber string `json:"phoneNumber,omitempty" example:"628123456789"`
}

type LoginRequest struct {
	Username string `json:"username,omitempty" example:"Username"`
	Password string `json:"password,omitempty" example:"verystrongpassword"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Username     string `json:"username,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken,omitempty" example:"refreshTokenHere"`
}

type RefreResponse TokenResponse
