package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type RedisTaskDistributor interface {
	DistributeVerifyEmail(ctx context.Context, payload SendVerifyEmailPayload, opts ...asynq.Option) error
}

type TaskDistributor struct {
	client *asynq.Client
}

func NewTaskDistributor(redisOption asynq.RedisClientOpt) RedisTaskDistributor {
	client := asynq.NewClient(redisOption)

	return &TaskDistributor{
		client: client,
	}
}
