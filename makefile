postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root zhasa_news

dropdb:
	docker exec -it postgres14 dropdb zhasa_news

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/zhasa_news?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/zhasa_news?sslmode=disable" -verbose down

migrateforce:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/zhasa_news?sslmode=disable" force 1 

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc migrateforce
