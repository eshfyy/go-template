package main

import (
	"context"
	"errors"
	"net"
	"net/http"

	grpcapi "go-template/internal/api/grpc"
	"go-template/internal/api/rest"
	"go-template/internal/app"
	"go-template/internal/config"
	uc "go-template/internal/contracts/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	fx.New(
		app.ConfigModule,
		app.LoggerModule,
		app.OTelModule("api"),
		app.PostgresModule,
		app.KafkaProducerModule,
		app.UseCaseModule,

		fx.Provide(newRESTRouter),
		fx.Provide(newGRPCServer),

		fx.Invoke(startREST),
		fx.Invoke(startGRPC),
	).Run()
}

type RouterParams struct {
	fx.In
	Log                *zap.Logger
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

func newRESTRouter(p RouterParams) *gin.Engine {
	return rest.NewRouter(rest.UseCases{
		CreateNotification: p.CreateNotification,
		GetNotification:    p.GetNotification,
		ListNotifications:  p.ListNotifications,
		DeleteNotification: p.DeleteNotification,
		CreateUser:         p.CreateUser,
		GetUser:            p.GetUser,
		ListUsers:          p.ListUsers,
		UpdateUser:         p.UpdateUser,
		DeleteUser:         p.DeleteUser,
	}, p.Log)
}

func newGRPCServer(p RouterParams) *grpc.Server {
	return grpcapi.NewServer(grpcapi.UseCases{
		CreateNotification: p.CreateNotification,
		GetNotification:    p.GetNotification,
		ListNotifications:  p.ListNotifications,
		DeleteNotification: p.DeleteNotification,
		CreateUser:         p.CreateUser,
		GetUser:            p.GetUser,
		ListUsers:          p.ListUsers,
		UpdateUser:         p.UpdateUser,
		DeleteUser:         p.DeleteUser,
	}, p.Log)
}

func startREST(lc fx.Lifecycle, router *gin.Engine, cfg config.HTTP, log *zap.Logger) {
	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			log.Info("starting rest", zap.String("addr", cfg.Addr()))
			go func() {
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error("rest server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("stopping rest")
			return srv.Shutdown(ctx)
		},
	})
}

func startGRPC(lc fx.Lifecycle, srv *grpc.Server, cfg config.GRPC, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			lis, err := net.Listen("tcp", cfg.Addr())
			if err != nil {
				return err
			}
			log.Info("starting grpc", zap.String("addr", cfg.Addr()))
			go srv.Serve(lis)
			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("stopping grpc")
			srv.GracefulStop()
			return nil
		},
	})
}
