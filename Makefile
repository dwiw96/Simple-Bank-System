dockerStart:
	docker run --rm --name postgres_bank -p 5432:5432 -e POSTGRES_USER=db -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=bank postgres
dockerExec:
	docker exec -it postgres_bank psql -U db bank

newMigration:
	migrate create -ext sql -dir db/migration -seq add_users
migrateUp:
	migrate -path db/migration -database "postgresql://db:secret@localhost:5432/bank?sslmode=disable" -verbose up
# The number 1 here means that we only want to run 1 next migration, or more precisely, just run the next up migration version that was applied current one.	
migrateUp1:
	migrate -path db/migration -database "postgresql://db:secret@localhost:5432/bank?sslmode=disable" -verbose up 1
migrateDown:
	migrate -path db/migration -database "postgresql://db:secret@localhost:5432/bank?sslmode=disable" -verbose down
# The number 1 here means that we only want to rollback 1 last migration, or more precisely, just run the last down migration version that was applied before.	
migrateDown1:
	migrate -path db/migration -database "postgresql://db:secret@localhost:5432/bank?sslmode=disable" -verbose down 1
migrateForce:
	migrate -path db/migration -database "postgresql://db:secret@localhost:5432/bank?sslmode=disable" force $(v)

server:
	go run main.go

test:
	go test -v -cover ./...
.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock