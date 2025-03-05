package user

import (
	"github.com/grozaqueen/julse/internal/model"
)

type UsersSignUpRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type UsersLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UsersDefaultResponse struct {
	UserID uint32 `json:"user_id"`
	Email  string `json:"username"`
}

func (us *UsersSignUpRequest) ToModel() model.User {
	return model.User{
		Email:    us.Email,
		Password: us.Password,
	}
}

func (ul *UsersLoginRequest) ToModel() model.User {
	return model.User{
		Email:    ul.Email,
		Password: ul.Password,
	}
}
