package main

import (
	"context"
	"fmt"

	"github.com/EmilioCliff/e-commerce/db/api"
	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	token "github.com/EmilioCliff/e-commerce/db/token"
	"github.com/EmilioCliff/e-commerce/db/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func main() {
	conn, err := pgxpool.New(context.Background(), "postgresql://root:secret@localhost:5432/e-commerce?sslmode=diable")
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "Database Connection").
			Msg("Cannot Connect to db")

		return
	}
	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: "0.0.0.0:6379",
		// Addr:     config.REDIS_URI,
		// Password: config.REDIS_PASSWORD,
	}

	// TODO: redis client used for caching
	_ = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:6379",
		Password: "",
		DB:       5,
	})

	taskDistributor := worker.NewTaskDistributor(redisOpt)

	maker, err := token.NewPasetoMaker("12345678901234567890123456789012")
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "NewPasetoMaker").
			Msg("Cannot start service")

		return
	}

	go runRedisProcessorServer(redisOpt, store)
	server := api.NewServer(store, maker, &taskDistributor)
	err = server.Start("0.0.0.0:8080")
	if err != nil {
		fmt.Println("error starting the server")
	}
}

func runRedisProcessorServer(redisOpt asynq.RedisClientOpt, store *db.Store) {
	processor := worker.NewTaskProcessor(redisOpt, store)
	if err := processor.Start(); err != nil {
		log.Fatal().
			Str("service", "asynq processor").
			Msg("failed to start asynq server for the processor")
	}
	log.Info().Msg("asynq server started")
}
