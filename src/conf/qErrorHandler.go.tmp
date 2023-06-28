package conf

import (
	"context"

	"github.com/hibiken/asynq"
)

type QueueErrorHandler struct{}

func (er *QueueErrorHandler) HandleError(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	Logger.Errorw("QErr", "task", task.Type(), "Retries", retried, "err", err)
}
