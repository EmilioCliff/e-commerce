# Create a postgres docker container 
postgres:
	docker run --name postgres2 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 5432:5432 -d postgres:alpine3.19

# Create or drop the e-commerce db
createdb:
	docker exec -it postgres2 createdb --username=root --owner=root e-commerce

dropdb:
	docker exec -it postgres2 dropdb e-commerce

# run migrations
migrateup:
	migrate -path go-backend/db/migrations -database postgresql://root:secret@localhost:5432/e-commerce?sslmode=disable -verbose up

migratedown:
	migrate -path go-backend/db/migrations -database postgresql://root:secret@localhost:5432/e-commerce?sslmode=disable -verbose down

# Generate sqlc
sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc