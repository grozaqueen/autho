package middleware

import (
	"context"
	"github.com/grozaqueen/merch-service/internal/model"
)

type sessionGetter interface {
	Get(ctx context.Context, sessionID string) (model.Session, error)
}

type csrfValidator interface {
	IsValidCSRFToken(session model.Session, token string) (bool, error)
}
