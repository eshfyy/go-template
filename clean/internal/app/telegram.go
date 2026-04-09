package app

import (
	"go-template/internal/config"
	"go-template/internal/contracts/infra"
	"go-template/internal/infra/telegram"

	"go.uber.org/fx"
)

var TelegramModule = fx.Module("telegram",
	fx.Provide(
		fx.Annotate(
			func(cfg config.Telegram) (*telegram.Sender, error) {
				return telegram.NewSender(cfg.Token)
			},
			fx.As(new(infra.NotificationSender)),
		),
	),
	fx.Provide(func(sender infra.NotificationSender) []infra.NotificationSender {
		return []infra.NotificationSender{sender}
	}),
)
