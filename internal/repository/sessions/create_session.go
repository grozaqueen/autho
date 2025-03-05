package sessions

import (
	"context"
	"fmt"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/model"
	"github.com/grozaqueen/julse/internal/utils"
)

func (sr *SessionStore) Create(ctx context.Context, session model.Session) (string, error) {
	userIDStr := fmt.Sprintf("%d", session.UserID)

	err := sr.redis.Set(ctx, session.SessionID, userIDStr, utils.DefaultSessionLifetime).Err()
	if err != nil {
		return "", errs.SessionCreationError
	}

	return session.SessionID, nil
}
