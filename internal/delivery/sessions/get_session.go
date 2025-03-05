package sessions

import (
	"context"

	"github.com/grozaqueen/julse/internal/model"
)

func (sd *SessionHandler) Get(ctx context.Context, sessionID string) (model.Session, error) {
	return sd.sessionManager.Get(ctx, sessionID)
}
