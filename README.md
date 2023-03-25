# Go Boilerplate

> This boilerplate is based on
> [techschool/simplebank](https://github.com/techschool/simplebank). The goal is
> to provide a well-structured way to start a Go-backed backend infrastructure.

## Schema Generation

Create the db schema in [dbdiagram.io](https://dbdiagram.io/home) to **decouple
the design from a specific database**. Then, click "Export" to export to the
language you want. In this boilerplate, we'll choose Postgres.

## Local Development

### Install Visual Studio Code extensions

```bash
code --install-extension esbenp.prettier-vscode
code --install-extension foxundermoon.shell-format
code --install-extension mtxr.sqltools
code --install-extension golang.go
```

### Install [Docker](https://www.docker.com) and [PostgresSQL image](https://hub.docker.com/_/postgres).

```bash
# Install Docker.
brew install docker

# Run Docker app so that we can access the `docker` command.

# Pull the PostgresSQL image.
docker pull postgres:15.2-alpine

# Check the downloaded image.
docker images
```

### Run a Docker container using the official PostgresSQL image.

Creates and runs a Docker container with the name `postgres`, using the official
`postgres:15-alpine` Docker image. The container is started as a background
process (`-d` flag) and is mapped to port `5432` of the host machine
(`-p 127.0.0.1:5432:5432/tcp` flag), which is the default port for PostgreSQL.

The container is also configured with the environment variables `POSTGRES_USER`
and `POSTGRES_PASSWORD`, which set the default username and password for the
PostgreSQL database. In this case, the username is set to `root` and the
password is set to `password`.

```bash
docker run --name postgres \
  -p 127.0.0.1:5432:5432/tcp \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=password \
  -d postgres:15.2-alpine
```

```bash
# Enter the Postgres shell.
docker exec -it postgres psql -U root

# Try the following query in the shell.
SELECT now();
```

### Install [TablePlus](https://tableplus.com)

```bash
# Install TablePlus.
brew install tableplus
```

Connect to Postgres with the setting

![](https://i.imgur.com/jgHY7h3.png)

### Database Migration

```bash
# Install `migrate` command.
brew install golang-migrate

# Check the installed `migrate` command.
migrate --version

# Create the db migration directory.
mkdir -p db/migration

# Create the first migration script.
migrate create -ext sql -dir db/migration -seq init_schema
```

Now, create a [Makefile](./Makefile) to save time and run the following:

```bash
# Run a PostgresSQL container.
make postgres

# Create a db called "bank" in this tutorial.
make createdb

# Migrate up to create tables in the db.
make migrateup
```

### Codegen via sqlc

```bash
# Install sqlc.
brew install sqlc

# Check the installed sqlc.
sqlc version
```

Initialize [`sqlc.yaml`](./sqlc.yaml) and copy the initial config from
[Getting started with PostgreSQL](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html#getting-started-with-postgresql)
with some modifications.

```bash
sqlc init
```

Add the queries in [account.sql](./db/query/account.sql), then `sqlc generate`
to codegen.

```bash
# Codegen.
make sqlc

# Eliminate red lines inside `db/sqlc/account.sql`.
go mod init github.com/walkccc/go-boilerplate
```

### Unit tests

To write unit tests, we need to connect to the DB driver and import it in
[main_test.go](./db/sqlc/main_test.go).

```bash
# Get required driver.
go get github.com/lib/pq

# Remove "indirect" of lib/pq.
go mod tidy
```

Install [GoMock](https://github.com/golang/mock).

```bash
# Install mockgen.
go install github.com/golang/mock/mockgen@v1.6.0

# Export the Go path.
# Add "export PATH=$PATH:~/go/bin" in your .zshrc or .bashrc

# Check the installed mockgen.
which mockgen
```

### Create Docker image

```bash
# Creates a Docker image with the name "bank" and the tag "latest" using the
# Dockerfile in the current directory.
docker build -t bank:latest .
```
