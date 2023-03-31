dockerStart:
	docker run --rm --name postgres_bank -p 5432:5432 -e POSTGRES_USER=db -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=bank postgres
dockerExec:
	docker exec -it postgres_bank psql -U db bank

migrateUp:
	migrate -path db/migration -database "postgresql://db:secret@localhost:5432/bank?sslmode=disable" -verbose up
migrateDown:
	migrate -path db/migration -database "postgresql://db:secret@localhost:5432/bank?sslmode=disable" -verbose down
migrateForce:
	migrate -path db/migration -database "postgresql://db:secret@localhost:5432/bank?sslmode=disable" force 