MIGRATIONS_DIR=./migrations
DATABASE_URL="mysql://user:password@tcp(localhost:3308)/aws-intern"

.PHONY: build run db test create-migration
build:
	docker-compose up --build

create-migration:
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database $(DATABASE_URL) up

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database $(DATABASE_URL) down

test:
	go test -v

