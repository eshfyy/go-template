package retry_failed

import (
	"context"
	uc "go-template/internal/contracts/usecase"

	"github.com/hibiken/asynq"
)

const TaskType = "notifications:retry_failed"

func NewTask() *asynq.Task {
	return asynq.NewTask(TaskType, nil)
}

func New(retryUC uc.RetryFailed) asynq.HandlerFunc {
	return func(ctx context.Context, _ *asynq.Task) error {
		return retryUC.Execute(ctx)
	}
}
