package user

import (
	proto "github.com/grozaqueen/julse/api/protos/user/gen"
	"github.com/grozaqueen/julse/internal/model"
)

func toModel(us *proto.UsersSignUpRequest) model.User {
	return model.User{
		Email:    us.GetEmail(),
		Username: us.GetUsername(),
		Password: us.GetPassword(),
	}
}

func toUserModel(us *proto.UsersLoginRequest) model.User {
	return model.User{
		Email:    us.GetEmail(),
		Password: us.GetPassword(),
	}
}
