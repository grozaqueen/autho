package sessions

import (
	"errors"
	"net/http"

	"github.com/grozaqueen/julse/internal/errs"
	"github.com/grozaqueen/julse/internal/model"
	"github.com/grozaqueen/julse/internal/utils"
)

func (sd *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(utils.SessionName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			utils.WriteJSON(w, http.StatusUnauthorized, errs.HTTPErrorResponse{
				ErrorMessage: errs.UserNotAuthorized.Error(),
			})

			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, errs.HTTPErrorResponse{
			ErrorMessage: errs.InternalServerError.Error(),
		})

		return
	}

	err = sd.sessionManager.Delete(r.Context(), model.Session{SessionID: cookie.Value})
	if err != nil {
		err, code := sd.errResolver.Get(err)
		utils.WriteJSON(w, code, errs.HTTPErrorResponse{
			ErrorMessage: err.Error(),
		})

		return
	}

	http.SetCookie(w, utils.RemoveSessionCookie())

	utils.WriteJSON(w, http.StatusNoContent, nil)
}
