package dto

type CreateUserRequestDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
