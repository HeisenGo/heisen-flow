package presenter

import (
	"server/internal/user"
)

type UserRegisterReq struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func UserRegisterToUserDomain(up *UserRegisterReq) *user.User {
	return &user.User{
		FirstName: up.FirstName,
		LastName:  up.LastName,
		Email:     up.Email,
		Password:  up.Password,
	}
}
