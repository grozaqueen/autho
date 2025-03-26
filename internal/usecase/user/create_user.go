package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/grozaqueen/julse/internal/model"
	"github.com/grozaqueen/julse/internal/utils"
)

func (us *UsersService) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	requestID, err := utils.GetContextRequestID(ctx)
	if err != nil {
		return model.User{}, err
	}

	us.log.Info("[UsersService.CreateUser] Started executing, requestID", slog.Any("request-id", requestID))
	fmt.Println("[UsersService.CreateUser] Started executing, requestI")
	user.Username = us.inputValidator.SanitizeString(user.Username)
	user.Email = us.inputValidator.SanitizeString(user.Email)
	user.Password = us.inputValidator.SanitizeString(user.Password)

	salt, err := utils.GenerateSalt()
	if err != nil {
		us.log.Error("[ UsersService.CreateUser ] Ошибка при генерации соли", slog.String("error", err.Error()))
		return model.User{}, err
	}

	user.Password = utils.HashPassword(user.Password, salt)
	dbUser, err := us.userRepo.CreateUser(ctx, user)
	if err != nil {
		us.log.Error("[ UsersService.CreateUser ] Ошибка при создании юзера на уровне репозитория", slog.String("error", err.Error()))
		return model.User{}, err
	}

	return dbUser, nil
}
