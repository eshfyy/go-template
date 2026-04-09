package grpc

import (
	"go-template/gen/go/notificationpb"
	"go-template/gen/go/userpb"
	"go-template/internal/api/grpc/interceptor"
	uc "go-template/internal/contracts/usecase"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type UseCases struct {
	CreateNotification uc.CreateNotification
	GetNotification    uc.GetNotification
	ListNotifications  uc.ListNotifications
	DeleteNotification uc.DeleteNotification
	CreateUser         uc.CreateUser
	GetUser            uc.GetUser
	ListUsers          uc.ListUsers
	UpdateUser         uc.UpdateUser
	DeleteUser         uc.DeleteUser
}

func NewServer(ucs UseCases, log *zap.Logger) *grpc.Server {
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(interceptor.ErrorUnaryInterceptor(log)),
	)

	notificationpb.RegisterNotificationServiceServer(srv, NewNotificationServer(
		ucs.CreateNotification,
		ucs.GetNotification,
		ucs.ListNotifications,
		ucs.DeleteNotification,
	))

	userpb.RegisterUserServiceServer(srv, NewUserServer(
		ucs.CreateUser,
		ucs.GetUser,
		ucs.ListUsers,
		ucs.UpdateUser,
		ucs.DeleteUser,
	))

	return srv
}
