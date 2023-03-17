
postgres:
	docker run --name postgresdb -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	docker exec -it postgresdb createdb --username=postgres --owner=postgres web_scraper
	
migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

server:
	go run main.go

serverdocker:
	docker-compose build && docker-compose up

.PHONY: createdb migrateup migratedown test server serverdocker postgres