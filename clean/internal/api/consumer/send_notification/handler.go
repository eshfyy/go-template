package send_notification

import (
	"context"
	"encoding/json"
	"fmt"
	uc "go-template/internal/contracts/usecase"
	"go-template/internal/domain"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

type payload struct {
	NotificationID uuid.UUID `json:"NotificationID"`
}

func New(useCase uc.SendNotification) func(ctx context.Context, record *kgo.Record) error {
	return func(ctx context.Context, record *kgo.Record) error {
		var p payload
		if err := json.Unmarshal(record.Value, &p); err != nil {
			return fmt.Errorf("unmarshal %s: %w", domain.NotificationCreated, err)
		}

		return useCase.Execute(ctx, p.NotificationID)
	}
}
