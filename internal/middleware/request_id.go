package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/grozaqueen/julse/internal/utils"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New()

		ctx := context.WithValue(r.Context(), utils.RequestIDName, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
