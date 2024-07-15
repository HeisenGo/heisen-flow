package presenter

import (
	"server/internal/user"
)

type UserRegisterReq struct {
	FirstName string `json:"first_name" validate:"required" example:"yourname"`
	LastName  string `json:"last_name" validate:"required" example:"yourlastname"`
	Email     string `json:"email" validate:"required" example:"abc@gmail.com"`
	Password  string `json:"password" validate:"required" example:"Abc@123"`
}

type UserLoginReq struct {
	Email    string `json:"email" validate:"required" example:"valid_email@folan.com"`
	Password string `json:"password" validate:"required" example:"Abc@123"`
}

func UserRegisterToUserDomain(up *UserRegisterReq) *user.User {
	return &user.User{
		FirstName: up.FirstName,
		LastName:  up.LastName,
		Email:     up.Email,
		Password:  up.Password,
	}
}
