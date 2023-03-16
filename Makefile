migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/web_scraper?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb migrateup migratedown test server