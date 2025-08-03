package gRPC

import (
	"context"
	"github.com/mhthrh/common_pkg/pkg/logger"
	"github.com/mhthrh/common_pkg/pkg/model/user"
	userGrpc "github.com/mhthrh/common_pkg/pkg/model/user/grpc/v1"
	"github.com/mhthrh/common_pkg/pkg/xErrors"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"restfullApi/pkg/service"
)

type User struct {
	userGrpc.UnimplementedUserServiceServer
	lgr logger.ILogger
}

var (
	srv *service.Service
)

func New(l logger.ILogger) User {
	srv = service.New(l)
	return User{
		UnimplementedUserServiceServer: userGrpc.UnimplementedUserServiceServer{},
		lgr:                            l,
	}
}

func (u User) Create(ctx context.Context, in *userGrpc.UserRequest) (*emptypb.Empty, error) {
	u.lgr.Info(ctx, "start grpc create", zap.Any("in user parameters", in))

	u.lgr.Info(ctx, "start service create calling")

	err := srv.Create(ctx, &user.User{
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		UserName:    in.UserName,
		Password:    in.Password,
	})
	u.lgr.Info(ctx, "get response from service create user", zap.Any("response", err))
	return &emptypb.Empty{}, status.Errorf(xErrors.GetGrpcCode(err), "%s", xErrors.Yaml(err))
}

func (u User) GetByUserName(ctx context.Context, in *userGrpc.UserName) (*userGrpc.UserResponse, error) {
	u.lgr.Info(ctx, "start grpc GetByUserName", zap.Any("in user parameters", in))

	u.lgr.Info(ctx, "start service create calling")
	usr, err := srv.GetByUserName(ctx, in.Username)

	u.lgr.Info(ctx, "get response from service create user", zap.Any("response", err), zap.Any("user info", usr))

	if err.Code != xErrors.SuccessCode {
		return nil, status.Errorf(xErrors.GetGrpcCode(err), "%s", xErrors.Yaml(err))
	}
	return &userGrpc.UserResponse{Usr: &userGrpc.UserRequest{
		FirstName:   usr.FirstName,
		LastName:    usr.LastName,
		Email:       usr.Email,
		PhoneNumber: usr.PhoneNumber,
		UserName:    usr.UserName,
		Password:    usr.Password,
	}}, status.Errorf(xErrors.GetGrpcCode(err), "%s", xErrors.Yaml(err))

}

func (u User) Update(ctx context.Context, in *userGrpc.UserRequest) (*emptypb.Empty, error) {
	u.lgr.Info(ctx, "start grpc Update", zap.Any("in user parameters", in))

	u.lgr.Info(ctx, "start service Update calling")

	err := srv.Update(ctx, &user.User{
		FirstName:   in.FirstName,
		LastName:    in.LastName,
		Email:       in.Email,
		PhoneNumber: in.PhoneNumber,
		UserName:    in.UserName,
		Password:    in.Password,
	})
	u.lgr.Info(ctx, "get response from service update user", zap.Any("response", err))

	return &emptypb.Empty{}, status.Errorf(xErrors.GetGrpcCode(err), "%s", xErrors.Yaml(err))
}

func (u User) Remove(ctx context.Context, in *userGrpc.UserName) (*emptypb.Empty, error) {
	u.lgr.Info(ctx, "start grpc remove", zap.Any("in user parameters", in))

	u.lgr.Info(ctx, "start service remove calling")
	err := srv.Remove(ctx, in.Username)

	u.lgr.Info(ctx, "get response from service Remove user", zap.Any("response", err))

	return &emptypb.Empty{}, status.Errorf(xErrors.GetGrpcCode(err), "%s", xErrors.Yaml(err))
}
