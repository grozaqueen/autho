package user

import (
	grpc_gen "github.com/grozaqueen/julse/api/protos/user/gen"
	"google.golang.org/grpc"
)

func (u *UsersGrpc) Register(grpcServer *grpc.Server) {
	grpc_gen.RegisterUserServiceServer(grpcServer, u)
}
