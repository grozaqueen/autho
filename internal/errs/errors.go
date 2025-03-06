package errs

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type GetErrorCode interface {
	Get(error) (error, int)
}

var (
	InvalidJSONFormat     = errors.New("Неверный формат JSON")
	InvalidUsernameFormat = errors.New("Неверный формат имени пользователя")
	InvalidEmailFormat    = errors.New("Неверный формат email")
	InvalidPasswordFormat = errors.New("Неверный формат пароля")
	PasswordsDoNotMatch   = errors.New("Пароли не совпадают")
	UserAlreadyExists     = errors.New("Пользователь уже существует")
	InternalServerError   = errors.New("Внутренняя ошибка сервера")
	SessionCreationError  = errors.New("Ошибка при создании сессии")
	SessionSaveError      = errors.New("Ошибка при сохранении сессии")
	SessionNotFound       = errors.New("Сессия не найдена")
	WrongCredentials      = errors.New("Неверный email или пароль")
	UserNotAuthorized     = errors.New("Пользователь не авторизован")
	LogoutError           = errors.New("Ошибка при завершении сессии")
	UserDoesNotExist      = errors.New("Пользователь не существует")
	WrongPassword         = errors.New("Неправильный пароль")
	RequestIDNotFound     = errors.New("Не удалось получить request-id")
	FailedToParseConfig   = errors.New("Ошибка парсинга конфигурации")
)

type ErrorStore struct {
	mux        sync.RWMutex
	errorCodes map[error]int
}

func (e *ErrorStore) Get(err error) (error, int) {
	e.mux.RLock()
	defer e.mux.RUnlock()
	errCode, present := e.errorCodes[err]
	if !present {
		log.Println(fmt.Errorf("unexpected error occured: %w", err))
		return InternalServerError, http.StatusInternalServerError
	}

	return err, errCode
}

func NewErrorStore() *ErrorStore {
	return &ErrorStore{
		mux: sync.RWMutex{},
		errorCodes: map[error]int{
			InvalidJSONFormat:     http.StatusBadRequest,
			InvalidUsernameFormat: http.StatusBadRequest,
			InvalidEmailFormat:    http.StatusBadRequest,
			InvalidPasswordFormat: http.StatusBadRequest,
			UserAlreadyExists:     http.StatusConflict,
			InternalServerError:   http.StatusInternalServerError,
			SessionCreationError:  http.StatusInternalServerError,
			SessionSaveError:      http.StatusInternalServerError,
			SessionNotFound:       http.StatusUnauthorized,
			WrongCredentials:      http.StatusUnauthorized,
			UserNotAuthorized:     http.StatusUnauthorized,
			LogoutError:           http.StatusInternalServerError,
			PasswordsDoNotMatch:   http.StatusBadRequest,
			UserDoesNotExist:      http.StatusNotFound,
			WrongPassword:         http.StatusBadRequest,
			RequestIDNotFound:     http.StatusBadRequest,
			FailedToParseConfig:   http.StatusInternalServerError,
		},
	}
}

type HTTPErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}
