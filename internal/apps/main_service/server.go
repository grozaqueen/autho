package main_service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	usergrpc "github.com/grozaqueen/julse/api/protos/user/gen"
	"github.com/grozaqueen/julse/internal/configs"
	"github.com/grozaqueen/julse/internal/configs/logger"
	postgres "github.com/grozaqueen/julse/internal/configs/postgresql"
	"github.com/grozaqueen/julse/internal/configs/redis"
	sessionsDeliveryLib "github.com/grozaqueen/julse/internal/delivery/sessions"
	userDeliveryLib "github.com/grozaqueen/julse/internal/delivery/user"
	errResolveLib "github.com/grozaqueen/julse/internal/errs"
	middlewares "github.com/grozaqueen/julse/internal/middleware"
	"github.com/grozaqueen/julse/internal/model"
	sessionsRepoLib "github.com/grozaqueen/julse/internal/repository/sessions"
	"github.com/grozaqueen/julse/internal/usecase/sessions"
	"github.com/grozaqueen/julse/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	mainService = "main_service"
	userService = "user_go"
)

type usersDelivery interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}

type SessionDelivery interface {
	Delete(w http.ResponseWriter, r *http.Request)
	Get(ctx context.Context, sessionID string) (model.Session, error)
}

type Server struct {
	r        *mux.Router
	sessions SessionDelivery
	auth     usersDelivery
	cfg      configs.ServiceViperConfig
	log      *slog.Logger
}

type csrfDelivery interface {
	GetCsrf(w http.ResponseWriter, r *http.Request)
}

func NewServer() (*Server, error) {
	log := logger.InitLogger()
	router := mux.NewRouter()
	v, err := configs.SetupViper()
	if err != nil {
		return nil, err
	}

	dbPool, err := postgres.LoadPgxPool()
	if err != nil {
		return nil, err
	}

	inputValidator := utils.NewInputValidator()

	errResolver := errResolveLib.NewErrorStore()

	redisClient, err := redis.LoadRedisClient()
	if err != nil {
		return nil, err
	}

	sessionsRepo := sessionsRepoLib.NewSessionRepo(redisClient, log)
	sessionService := sessions.NewSessionService(sessionsRepo, log)
	sessionsDelivery := sessionsDeliveryLib.NewSessionDelivery(sessionsRepo, errResolver)

	userGRPCCfg := v.GetStringMap(userService)
	userCfg, err := configs.ParseServiceViperConfig(userGRPCCfg)
	if err != nil {
		return nil, err
	}

	userConn, err := grpc.NewClient(fmt.Sprintf("%s:%s", userCfg.Domain, userCfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	userClient := usergrpc.NewUserServiceClient(userConn)

	userHandler := userDeliveryLib.NewUsersDelivery(userClient, inputValidator, sessionService, errResolver, log)

	cfg := v.GetStringMap(mainService)
	mainCfg, err := configs.ParseServiceViperConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		r:        router,
		auth:     userHandler,
		cfg:      mainCfg,
		log:      log,
		sessions: sessionsDelivery,
	}, nil
}

func (s *Server) setupRoutes() {
	errResolver := errResolveLib.NewErrorStore()

	s.r.HandleFunc("/api/v1/login", s.auth.LoginUser).Methods(http.MethodPost)
	s.r.HandleFunc("/api/v1/logout", s.sessions.Delete).Methods(http.MethodPost)
	s.r.HandleFunc("/api/v1/signup", s.auth.CreateUser).Methods(http.MethodPost)

	authSub := s.r.Methods(http.MethodGet, http.MethodPost, http.MethodPut).Subrouter()
	authSub.HandleFunc("/api/v1/", s.auth.GetUserById).Methods(http.MethodGet)

	authSub.Use(middlewares.RequestIDMiddleware)
	authSub.Use(middlewares.AuthMiddleware(s.sessions, errResolver))
}

func (s *Server) Run() error {
	s.setupRoutes()

	handler := middlewares.CorsMiddleware(s.r)

	s.log.Info("starting  server", slog.String("address:", s.cfg.Port))
	return http.ListenAndServe(fmt.Sprintf(":%s", s.cfg.Port), handler)
}
