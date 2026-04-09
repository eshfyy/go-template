package app

import (
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/usecase"

	"go.uber.org/fx"
)

var UseCaseModule = fx.Module("usecase",
	fx.Provide(fx.Annotate(usecase.NewCreateNotification, fx.As(new(uc.CreateNotification)))),
	fx.Provide(fx.Annotate(usecase.NewGetNotification, fx.As(new(uc.GetNotification)))),
	fx.Provide(fx.Annotate(usecase.NewListNotifications, fx.As(new(uc.ListNotifications)))),
	fx.Provide(fx.Annotate(usecase.NewDeleteNotification, fx.As(new(uc.DeleteNotification)))),
	fx.Provide(fx.Annotate(usecase.NewSendNotification, fx.As(new(uc.SendNotification)))),
	fx.Provide(fx.Annotate(usecase.NewRetryFailed, fx.As(new(uc.RetryFailed)))),
	fx.Provide(fx.Annotate(usecase.NewCreateUser, fx.As(new(uc.CreateUser)))),
	fx.Provide(fx.Annotate(usecase.NewGetUser, fx.As(new(uc.GetUser)))),
	fx.Provide(fx.Annotate(usecase.NewListUsers, fx.As(new(uc.ListUsers)))),
	fx.Provide(fx.Annotate(usecase.NewUpdateUser, fx.As(new(uc.UpdateUser)))),
	fx.Provide(fx.Annotate(usecase.NewDeleteUser, fx.As(new(uc.DeleteUser)))),
)
