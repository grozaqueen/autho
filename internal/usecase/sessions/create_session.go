package sessions

import (
	"context"

	"github.com/google/uuid"
	"github.com/grozaqueen/julse/internal/model"
)

func (ss *SessionService) Create(ctx context.Context, userID uint32) (string, error) {
	id := uuid.New()

	session := model.Session{
		SessionID: id.String(),
		UserID:    userID,
	}

	return ss.SessionRepo.Create(ctx, session)
}
