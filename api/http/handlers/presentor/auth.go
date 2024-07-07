package presenter

import (
	"server/internal/user"
)

type UserRegisterReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func UserRegisterToUserDomain(up *UserRegisterReq) *user.User {
	return &user.User{
		FirstName: up.FirstName,
		LastName:  up.LastName,
		Email:     up.Email,
		Password:  up.Password,
	}
}
