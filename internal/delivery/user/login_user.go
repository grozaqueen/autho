package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *UsersDelivery) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req UsersLoginRequest

	usersDefaultResponse, err := u.userService.LoginUser(context.Background(), req.ToModel())

	if err != nil {
		u.log.Error("[ UsersDelivery.LoginUser ] До вывода кода", slog.String("error", err.Error()), slog.String("func", "u.userClientGrpc.LoginUser"))

		grpcErr, ok := status.FromError(err)
		if ok {
			u.log.Error("[ UsersDelivery.LoginUser ] ",
				slog.String("error", err.Error()), slog.String("func", "u.userClientGrpc.LoginUser"),
				slog.Any("code", grpcErr.Code()))

			switch grpcErr.Code() {
			case codes.NotFound:
				err, code := u.errResolver.Get(errs.UserDoesNotExist)
				utils.WriteJSON(w, code, errs.HTTPErrorResponse{
					ErrorMessage: err.Error(),
				})
				return

			case codes.InvalidArgument:
				err, code := u.errResolver.Get(errs.WrongCredentials)
				utils.WriteJSON(w, code, errs.HTTPErrorResponse{
					ErrorMessage: err.Error(),
				})
				return

			default:
				err, code := u.errResolver.Get(errs.InternalServerError)
				utils.WriteJSON(w, code, errs.HTTPErrorResponse{
					ErrorMessage: err.Error(),
				})
				return
			}
		}

		err, code := u.errResolver.Get(err)
		utils.WriteJSON(w, http.StatusInternalServerError, errs.HTTPErrorResponse{
			ErrorMessage: errs.InternalServerError.Error(),
		})
		u.log.Error("[ UsersDelivery.LoginUser ] ", slog.String("error", err.Error()), slog.String("func", "u.userClientGrpc.LoginUser"), slog.Int("code", code))
		return
	}

	sessionID, err := u.sessionService.Create(r.Context(), usersDefaultResponse.ID)
	if err != nil {
		err, code := u.errResolver.Get(err)
		utils.WriteJSON(w, code, errs.HTTPErrorResponse{
			ErrorMessage: err.Error(),
		})
		u.log.Error("[ UsersDelivery.LoginUser ] Ошибка при получении сессии", slog.String("error", err.Error()))
		return
	}

	http.SetCookie(w, utils.SetSessionCookie(sessionID))

	utils.WriteJSON(w, http.StatusOK, UsersDefaultResponse{
		UserID: usersDefaultResponse.ID,
		Email:  usersDefaultResponse.Email,
	})
}
