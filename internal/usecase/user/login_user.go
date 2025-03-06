package user

import (
	"context"
	"log/slog"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/model"
	"github.com/grozaqueen/julse/internal/utils"
)

func (us *UsersService) LoginUser(ctx context.Context, user model.User) (model.User, error) {
	requestID, err := utils.GetContextRequestID(ctx)
	if err != nil {
		return model.User{}, err
	}

	us.log.Info("[UsersService.LoginUser] Started executing, requestID", slog.Any("request-id", requestID))

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
