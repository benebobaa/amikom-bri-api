postgres:
	sudo docker run --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:16-alpine

createdb:
	sudo docker exec -it postgres16 createdb --username=root --owner=root amikombri

dropdb:
	sudo docker exec -it postgres16 dropdb amikombri

migratecreate:
	migrate create -ext sql -dir db/postgres/migrations -seq $(name)

migrateup:
	migrate -path db/postgres/migrations -database "postgresql://root:root@localhost:5432/amikombri?sslmode=disable" -verbose up

migratedown:
	migrate -path db/postgres/migrations -database "postgresql://root:root@localhost:5432/amikombri?sslmode=disable" -verbose down 1

migratefix:
	migrate -path db/postgres/migrations -database "postgresql://root:root@localhost:5432/amikombri?sslmode=disable" -verbose force $(version)

test:
	go test -v -cover ./...

PHONY: postgres createdb dropdb migrate_create migrate_up migrate_down test migrate_fix