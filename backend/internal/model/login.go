package model

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}


type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}