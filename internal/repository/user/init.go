package user

import (
	"github.com/grozaqueen/julse/internal/repository/pool"
	"log/slog"
)

type UsersStore struct {
	db  pool.DBPool
	log *slog.Logger
}

func NewUsersStore(db pool.DBPool, log *slog.Logger) *UsersStore {
	return &UsersStore{db: db,
		log: log}
}
