package sessions

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/model"
	"github.com/redis/go-redis/v9"
)

func (sr *SessionStore) Get(ctx context.Context, sessionID string) (model.Session, error) {
	value := sr.redis.Get(ctx, sessionID)
	if errors.Is(value.Err(), redis.Nil) {
		return model.Session{}, errs.SessionNotFound
	}

	userID, err := value.Uint64()
	if err != nil {
		log.Println(fmt.Errorf("[SesionRepo.Get] An error occured: %w", err))
		return model.Session{}, errs.InternalServerError
	}

	return model.Session{
		UserID:    uint32(userID),
		SessionID: sessionID,
	}, nil
}
