# Create and run a Docker container with the name `postgres`, using the official
# `postgres:15.2-alpine` Docker image.
postgres:
	docker run --name postgres \
			-p 127.0.0.1:5432:5432/tcp \
			-e POSTGRES_USER=root \
			-e POSTGRES_PASSWORD=password \
			-d postgres:15.2-alpine

# Create a db called "bank".
createdb:
	docker exec -it postgres createdb --username=root --owner=root bank

# Drop a db called "bank".
dropdb:
	docker exec -it postgres dropdb bank

# Migrate up to add tables in "bank".
migrateup:
	migrate -path db/migration \
			-database "postgresql://root:password@localhost:5432/bank?sslmode=disable" \
			-verbose up $(ARGS)

# Migrate down to drop tables in "bank".
migratedown:
	migrate -path db/migration \
			-database "postgresql://root:password@localhost:5432/bank?sslmode=disable" \
			-verbose down $(ARGS)

# Codegen CRUD code from "./db/query/" to "./db/sqlc/".
sqlc:
	sqlc generate

# Run all tests and generate code coverage reports for all packages in the
# current module.
test:
	go test -v -cover ./...

# Start the server.
server:
	go run main.go

# Generate a mock implementation located at `db/mock/store.go` of the `Store`
# interface in the `github.com/walkccc/go-boilerplate/db/sqlc` package.
#
# This mock implementation can then be used for testing purposes.
mockgen:
	mockgen -package mockdb -destination db/mock/store.go \
			github.com/walkccc/go-boilerplate/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock
