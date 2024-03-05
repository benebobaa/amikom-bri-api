postgres:
	sudo docker run --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:16-alpine

createdb:
	sudo docker exec -it postgres16 createdb --username=root --owner=root amikombri

dropdb:
	sudo docker exec -it postgres16 dropdb amikombri

migratecreate:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrateup:
	migrate -path db/migrations -database "postgresql://root:root@localhost:5432/amikombri?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:root@localhost:5432/amikombri?sslmode=disable" -verbose down 1

migratefix:
	migrate -path db/migrations -database "postgresql://root:root@localhost:5432/amikombri?sslmode=disable" -verbose force $(version)

migratecreatess:
	migrate create -ext sql -dir db/sqlserver -seq $(name)

migrateupss:
	migrate -path db/sqlserver -database "sqlserver://sa:Secret1!@localhost:1433?database=amikombri" -verbose up

migratedownss:
	migrate -path db/sqlserver -database "sqlserver://sa:Secret1!@localhost:1433?database=amikombri" -verbose down 1

migratefixss:
	migrate -path db/sqlserver -database "sqlserver://sa:Secret1!@localhost:1433?database=amikombri" -verbose force $(version)


test:
	go test -v -cover ./...

PHONY: postgres createdb dropdb migrate_create migrate_up migrate_down test migrate_fix