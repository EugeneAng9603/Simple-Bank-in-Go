postgres:
	docker run --name postgres15_useThis -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres15_useThis createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15_useThis dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

#use ./... to run all unit tests 
test: 
	go test -v -cover ./... 

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test