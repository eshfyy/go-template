package app

import (
	iservice "go-template/internal/contracts/service"
	"go-template/internal/service"

	"go.uber.org/fx"
)

var ServiceModule = fx.Module("service",
	fx.Provide(
		fx.Annotate(
			service.NewNotificationSenderService,
			fx.As(new(iservice.NotificationSenderService)),
		),
	),
)
