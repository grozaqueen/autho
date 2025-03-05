package sessions

import (
	"context"
	"errors"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/model"
	"github.com/redis/go-redis/v9"
)

func (sr *SessionStore) Delete(ctx context.Context, session model.Session) error {
	err := sr.redis.Del(ctx, session.SessionID).Err()
	if errors.Is(err, redis.Nil) {
		return errs.SessionNotFound
	}

	return nil
}
