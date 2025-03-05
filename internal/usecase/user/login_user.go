package user

import (
	"context"
	"log/slog"

	"github.com/grozaqueen/merch-service/internal/errs"
	"github.com/grozaqueen/merch-service/internal/model"
	"github.com/grozaqueen/merch-service/internal/utils"
)

func (us *UsersService) LoginUser(ctx context.Context, user model.User) (model.User, error) {
	dbUser, err := us.userRepo.GetUserByEmail(ctx, user)
	if err != nil {
		us.log.Error("[ UsersService.LoginUser ] Не найден юзер", slog.String("error", err.Error()))
		return model.User{}, err
	}

	if !utils.VerifyPassword(dbUser.Password, user.Password) {
		us.log.Info("[ UsersService.LoginUser ] Не прошла валидация паролей")
		return model.User{}, errs.WrongCredentials
	}

	return dbUser, nil
}
