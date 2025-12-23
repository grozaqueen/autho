package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *UsersDelivery) CreateUser(w http.ResponseWriter, r *http.Request) {
	requestID, err := utils.GetContextRequestID(r.Context())
	if err != nil {
		u.log.Error("[UsersDelivery.CreateUser] No request ID")
		utils.WriteErrorJSONByError(w, err, u.errResolver)
		return
	}

	u.log.Info("[UsersDelivery.CreateUser] Started executing", slog.Any("request-id", requestID))

	var req UsersSignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.log.Error("[UsersDelivery.CreateUser] Failed to decode body", slog.String("error", err.Error()))
		utils.WriteErrorJSONByError(w, errs.InvalidJSONFormat, u.errResolver) // если у тебя нет errs.BadRequest — замени на подходящую ошибку
		return
	}

	if err = utils.ValidateRegistration(req.Email, req.Username, req.Password, req.RepeatPassword); err != nil {
		err, code := u.errResolver.Get(err)
		utils.WriteJSON(w, code, errs.HTTPErrorResponse{ErrorMessage: err.Error()})
		u.log.Error("[UsersDelivery.CreateUser] Validation failed", slog.String("error", err.Error()))
		return
	}

	newCtx, err := utils.AddMetadataRequestID(r.Context())
	if err != nil {
		err, code := u.errResolver.Get(err)
		utils.WriteJSON(w, code, errs.HTTPErrorResponse{
			ErrorMessage: err.Error(),
		})
		return 
	}

	usersDefaultResponse, err := u.userClientGrpc.CreateUser(newCtx, req.ToGrpcSignupRequest())

	grpcErr, ok := status.FromError(err)
	if err != nil {
		if ok {
			switch grpcErr.Code() {
			case codes.InvalidArgument:
				u.log.Error("[UsersDelivery.CreateUser] User already exists", slog.Any("error", err.Error()))
				utils.WriteErrorJSONByError(w, errs.UserAlreadyExists, u.errResolver)
				return
			default:
				u.log.Error("[UsersDelivery.CreateUser] Unknown error", slog.String("error", err.Error()))
				utils.WriteErrorJSONByError(w, errs.InternalServerError, u.errResolver)
				return
			}
		}

		u.log.Error("[UsersDelivery.CreateUser] Failed to parse grpc error code")
		utils.WriteErrorJSONByError(w, errs.InternalServerError, u.errResolver)
		return
	}

	sessionID, err := u.sessionService.Create(r.Context(), usersDefaultResponse.GetUserId())
	if err != nil {
		err, code := u.errResolver.Get(err)
		utils.WriteJSON(w, code, errs.HTTPErrorResponse{
			ErrorMessage: err.Error(),
		})

		u.log.Error("[UsersDelivery.CreateUser] Session create error",
			slog.String("error", err.Error()),
			slog.Any("userId", usersDefaultResponse.UserId),
		)
		return
	}

	http.SetCookie(w, utils.SetSessionCookie(sessionID))

	utils.WriteJSON(w, http.StatusOK, UsersDefaultResponse{
		UserID:   usersDefaultResponse.UserId,
		Username: usersDefaultResponse.Username,
		City:     usersDefaultResponse.City,
	})
}
