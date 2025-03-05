package middleware

import (
	"errors"
	"net/http"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/utils"
)

func AuthMiddleware(sessionGetter sessionGetter, errResolver errs.GetErrorCode) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(utils.SessionName)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					next.ServeHTTP(w, r)
					return
				}

				utils.WriteErrorJSON(w, http.StatusInternalServerError, errs.InternalServerError)
				return
			}

			session, err := sessionGetter.Get(r.Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, errs.SessionNotFound) {
					next.ServeHTTP(w, r)
					return
				}

				err, code := errResolver.Get(err)
				utils.WriteErrorJSON(w, code, err)
				return
			}

			ctx := utils.SetContextSessionUserID(r.Context(), session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
