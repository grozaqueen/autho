package user

import (
	"context"
	"fmt"
	proto "github.com/grozaqueen/julse/api/protos/user/gen"
	"github.com/grozaqueen/julse/internal/configs"
	postgres "github.com/grozaqueen/julse/internal/configs/postgresql"
	user2 "github.com/grozaqueen/julse/internal/grpc_api/user"
	userRepoLib "github.com/grozaqueen/julse/internal/repository/user"
	userServiceLib "github.com/grozaqueen/julse/internal/usecase/user"
	"github.com/grozaqueen/julse/internal/utils"
	"google.golang.org/grpc"
	"log/slog"
)

type usersDelivery interface {
	Register(grpcServer *grpc.Server)
	LoginUser(ctx context.Context, in *proto.UsersLoginRequest) (*proto.UsersDefaultResponse, error)
	GetUserById(ctx context.Context, in *proto.GetUserByIdRequest) (*proto.UsersDefaultResponse, error)
	CreateUser(ctx context.Context, in *proto.UsersSignUpRequest) (*proto.UsersDefaultResponse, error)
}

type UsersApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	delivery   usersDelivery
	config     configs.ServiceViperConfig
}

func NewUsersApp(log *slog.Logger, grpcServer *grpc.Server,
	serviceConf map[string]any) (*UsersApp, error) {
	c, err := configs.ParseServiceViperConfig(serviceConf)
	if err != nil {
		slog.Error("UsersApp [NewUsersApp] Failed to parse service cfg")

		return nil, err
	}

	dbPool, err := postgres.LoadPgxPool()
	if err != nil {
		return nil, fmt.Errorf("не инициализируется бд %v", err)
	}

	// todo добавить
	inputValidator := utils.NewInputValidator()

	userRepo := userRepoLib.NewUsersStore(dbPool, log)
	userService := userServiceLib.NewUserService(userRepo, inputValidator, log)

	delivery := user2.NewUsersGrpc(userService, userRepo, log)

	return &UsersApp{
		log:        log,
		gRPCServer: grpcServer,
		delivery:   delivery,
		config:     c,
	}, nil
}
