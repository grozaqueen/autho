package sessions

import (
	"context"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/model"
)

type sessionManager interface {
	Delete(ctx context.Context, session model.Session) error
	Get(ctx context.Context, sessionID string) (model.Session, error)
}

type SessionHandler struct {
	sessionManager sessionManager
	errResolver    errs.GetErrorCode
}

func NewSessionDelivery(sessionManager sessionManager, errResolver errs.GetErrorCode) *SessionHandler {
	return &SessionHandler{
		sessionManager: sessionManager,
		errResolver:    errResolver,
	}
}
