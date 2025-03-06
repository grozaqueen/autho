package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/grozaqueen/julse/internal/model"
	"github.com/grozaqueen/julse/internal/utils"
	"net/http"
)

type sessionGetter interface {
	Get(ctx context.Context, sessionID string) (model.Session, error)
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New()

		ctx := context.WithValue(r.Context(), utils.RequestIDName, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
