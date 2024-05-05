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

func NewTaskDistributor(redisOpts asynq.RedisClientOpt) RedisTaskDistributor {
	client := asynq.NewClient(redisOpts)

	return &TaskDistributor{
		client: client,
	}
}
