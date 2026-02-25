package auth

type RegisterRequest struct {
	UserName string `json:"userName" validate:"required,min=3,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role"     validate:"omitempty,oneof=admin user"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Role     string `json:"eRole"`
}
