package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grozaqueen/merch-service/internal/errs"
	"github.com/grozaqueen/merch-service/internal/utils"
)

func (u *UsersDelivery) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UsersSignUpRequest

	if err := utils.ValidateRegistration(req.Email, req.Password, req.RepeatPassword); err != nil {
		err, code := u.errResolver.Get(err)
		utils.WriteJSON(w, code, errs.HTTPErrorResponse{
			ErrorMessage: err.Error(),
		})
		u.log.Error("[ UsersDelivery.CreateUser ] Валидация регистрации не прошла успешно", slog.String("error", err.Error()))
		return
	}

	usersDefaultResponse, err := u.userService.CreateUser(context.Background(), req.ToModel())

	sessionID, err := u.sessionService.Create(r.Context(), usersDefaultResponse.ID)
	if err != nil {
		err, code := u.errResolver.Get(err)
		utils.WriteJSON(w, code, errs.HTTPErrorResponse{
			ErrorMessage: err.Error(),
		})

		u.log.Error("[ UsersDelivery.CreateUser ] Ошибка при создании сессии ",
			slog.String("error", err.Error()),
			slog.Any("userId", usersDefaultResponse.ID),
		)

		return
	}

	http.SetCookie(w, utils.SetSessionCookie(sessionID))

	utils.WriteJSON(w, http.StatusOK, UsersDefaultResponse{
		UserID: usersDefaultResponse.ID,
		Email:  usersDefaultResponse.Email,
	})
}
