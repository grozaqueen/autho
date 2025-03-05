package user

import (
	"context"
	"github.com/grozaqueen/julse/internal/usecase/user"
	"github.com/segmentio/kafka-go"
	"log/slog"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/utils"
)

type sessionCreator interface {
	Create(ctx context.Context, userID uint32) (string, error)
}

type UsersDelivery struct {
	userService    user.UsersService
	inputValidator *utils.InputValidator
	sessionService sessionCreator
	errResolver    errs.GetErrorCode
	log            *slog.Logger
}

func NewUsersDelivery(userManager user.UsersService, inputValidator *utils.InputValidator, sessionService sessionCreator, errResolver errs.GetErrorCode, log *slog.Logger) *UsersDelivery {
	return &UsersDelivery{
		userService:    userManager,
		inputValidator: inputValidator,
		sessionService: sessionService,
		errResolver:    errResolver,
		log:            log,
	}
}

type MessageProducer struct {
	writer *kafka.Writer
	log    *slog.Logger
}
