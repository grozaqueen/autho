package user

import (
	"context"
	grpc_gen "github.com/grozaqueen/julse/api/protos/user/gen"
	"log/slog"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/utils"
)

type sessionCreator interface {
	Create(ctx context.Context, userID uint32) (string, error)
}

type UsersDelivery struct {
	userClientGrpc grpc_gen.UserServiceClient
	inputValidator *utils.InputValidator
	sessionService sessionCreator
	errResolver    errs.GetErrorCode
	log            *slog.Logger
}

func NewUsersDelivery(userManager grpc_gen.UserServiceClient, inputValidator *utils.InputValidator, sessionService sessionCreator, errResolver errs.GetErrorCode, log *slog.Logger) *UsersDelivery {
	return &UsersDelivery{
		userClientGrpc: userManager,
		inputValidator: inputValidator,
		sessionService: sessionService,
		errResolver:    errResolver,
		log:            log,
	}
}
