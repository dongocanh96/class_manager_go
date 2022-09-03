setuppostgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.5-alpine3.16

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root class_manager

dropdb:
	docker exec -it postgres14 dropdb class_manager

createmigration:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/class_manager?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/class_manager?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/class_manager?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/class_manager?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -cover ./...

server:
	go run main.go

mock:
	mockgen --build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/dongocanh96/class_manager_go/db/sqlc Store

.PHONY: setuppostgres createdb dropdb createmigration migrateup migratedown sqlc test server mock