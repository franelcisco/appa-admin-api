package models

type LoginRequest struct {
	TokeID string `json:"tokenID"`
}

type LoginResponse struct {
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
	Token   string `json:"token"`
}
