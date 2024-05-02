package main

import (
	"context"
	"fmt"

	"github.com/EmilioCliff/e-commerce/db/api"
	db "github.com/EmilioCliff/e-commerce/db/sqlc"
	token "github.com/EmilioCliff/e-commerce/db/token"
	"github.com/jackc/pgx/v5/pgxpool"
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

	maker, err := token.NewPasetoMaker("12345678901234567890123456789012")
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "NewPasetoMaker").
			Msg("Cannot start service")

		return
	}

	server := api.NewServer(store, maker)
	err = server.Start("0.0.0.0:8080")
	if err != nil {
		fmt.Println("error starting the server")
	}
}
