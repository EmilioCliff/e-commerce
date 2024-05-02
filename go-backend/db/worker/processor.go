package worker

import (
	"context"

	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	"github.com/EmilioCliff/e-commerce/db/token"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type RedisTaskProcessor interface {
	Start() error
	ProcessVerifyEmail(ctx context.Context, task asynq.Task) error
}

type TaskProcessor struct {
	server *asynq.Server
	maker  token.Maker
	store  db.Store
}

func NewTaskProcessor(redisOpts asynq.RedisClientOpt, store db.Store) RedisTaskProcessor {
	server := asynq.NewServer(redisOpts, asynq.Config{
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	maker, err := token.NewPasetoMaker("12345678901234567890123456789012")
	if err != nil {
		log.Fatal().
			Str("service", "NewTaskProcessor").
			Msg("failed to load paseto maker to task processor")
	}
	return &TaskProcessor{
		maker:  maker,
		server: server,
		store:  store,
	}
}

func (processor *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	return processor.server.Start(mux)
}
