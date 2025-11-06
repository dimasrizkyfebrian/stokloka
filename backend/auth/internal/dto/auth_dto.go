package dto

// RegisterUserDTO adalah struct untuk input registrasi
type RegisterUserDTO struct {
	Nama     string `json:"nama" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	RoleID   uint   `json:"role_id" validate:"required,numeric"`
}

// LoginUserDTO adalah struct untuk input login
type LoginUserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse adalah struct untuk output JSON
type AuthResponse struct {
	Token string `json:"token"`
}
