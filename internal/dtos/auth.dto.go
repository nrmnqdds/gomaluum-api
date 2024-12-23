package dtos

type LoginProps struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
