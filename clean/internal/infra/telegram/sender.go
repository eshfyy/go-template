package telegram

import (
	"context"
	"fmt"
	"go-template/internal/contracts/infra"
	"go-template/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	bot *tgbotapi.BotAPI
}

func NewSender(token string) (*Sender, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("telegram bot: %w", err)
	}
	return &Sender{bot: bot}, nil
}

func (s *Sender) Send(_ context.Context, req infra.SendRequest) error {
	text := fmt.Sprintf("*%s*\n\n%s", req.Notification.Title, req.Notification.Text)

	msg := tgbotapi.NewMessage(req.Recipient.TelegramID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := s.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("telegram send: %w", err)
	}
	return nil
}

func (s *Sender) Channel() domain.NotificationChannel {
	return domain.NotificationChannelTelegram
}
