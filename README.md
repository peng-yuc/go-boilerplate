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

### Create Docker Compose

[docker-compose.yaml](./docker-compose.yaml) defines two services: `postgres`
and `api`.

- The `postgres` service sets up a PostgreSQL database container with some
  environment variables and a healthcheck.
- The `api` service builds a container from the [Dockerfile](./Dockerfile), sets
  an environment variable with the database connection string (which will be
  used in [start.sh](./start.sh)), and depends on the postgres service to be
  healthy before starting up.

When you run `docker compose up`, Docker Compose starts up the containers for
both services, with the `api` service waiting for the `postgres` service to be
ready before starting.

```bash
# Make the start.sh file executable so that it can be run in Docker container.
chmod +x start.sh

# Start up all the services defined in a Docker Compose file.
docker compose up
```

The screenshot of `docker compose up` shows that `api` service depends on
`postgres` service.

![](https://i.imgur.com/mmHg4NH.png)

### Demo

Let's try the APIs with the following steps:

> TODO: This can be written in integration tests.

**Create User** (bunny)

```bash
curl -X POST \
    -H 'Content-Type: application/json' \
    -d '{"username": "bunny", "full_name": "Bunny", "email": "bunny@email.com", "password": "password"}' \
    http://localhost:8080/users | json_pp

# {
#    "created_at" : "2023-03-26T01:17:27.747643Z",
#    "email" : "bunny@email.com",
#    "full_name" : "Bunny",
#    "password_changed_at" : "0001-01-01T00:00:00Z",
#    "username" : "bunny"
# }
```

**Login User** (bunny)

```bash
curl -X POST \
    -H 'Content-Type: application/json' \
    -d '{"username": "bunny", "password": "password"}' \
    http://localhost:8080/users/login | json_pp

# {
#    "access_token" : "v2.local.mLdrNjZnn-Mvjo8aL3vYXxBmEVZIP4yB56xsdADV-EciMT8j-Ts0Hvpc5amsxsUZFwIyY8ZwRAktFmAF49d6x-z53fVPQP5vQtZ1tUPkq1Ekq_P9dH-V7XmhrHm2gww8d6Cch5mcN8tC0j-TgMEOmNr__dk5vbx9ew8nyOgYNEVELBhDmd9TAKUtuVhoqHPSLiDcYYpb9I1svvdM9Fvj_egmUMpUnF72qSPe7DkhJPaWQ2vtakHf-PxZoxMYnW4iPWRT.bnVsbA",
#    "user" : {
#       "created_at" : "2023-03-26T01:17:27.747643Z",
#       "email" : "bunny@email.com",
#       "full_name" : "Bunny",
#       "password_changed_at" : "0001-01-01T00:00:00Z",
#       "username" : "bunny"
#    }
# }
```

**Create Account** (bunny) without any `access_token`.

```bash
curl -X POST \
    -H 'Content-Type: application/json' \
    -d '{"currency": "EUR"}' \
    http://localhost:8080/accounts | json_pp

# {
#    "error" : "Authorization header is not provided."
# }
```

**Create Account** (bunny, EUR) with a valid and unexpired `access_token`.

```bash
ACCESS_TOKEN="v2.local.mLdrNjZnn-Mvjo8aL3vYXxBmEVZIP4yB56xsdADV-EciMT8j-Ts0Hvpc5amsxsUZFwIyY8ZwRAktFmAF49d6x-z53fVPQP5vQtZ1tUPkq1Ekq_P9dH-V7XmhrHm2gww8d6Cch5mcN8tC0j-TgMEOmNr__dk5vbx9ew8nyOgYNEVELBhDmd9TAKUtuVhoqHPSLiDcYYpb9I1svvdM9Fvj_egmUMpUnF72qSPe7DkhJPaWQ2vtakHf-PxZoxMYnW4iPWRT.bnVsbA"

curl -X POST \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H 'Content-Type: application/json' \
    -d '{"currency": "EUR"}' \
    http://localhost:8080/accounts | json_pp

# {
#    "balance" : 0,
#    "created_at" : "2023-03-26T01:41:27.598016Z",
#    "currency" : "EUR",
#    "id" : 1,
#    "owner" : "bunny"
# }
```

**Create Account** (bunny, USD) with a valid and unexpired `access_token`.

```bash
curl -X POST \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H 'Content-Type: application/json' \
    -d '{"currency": "USD"}' \
    http://localhost:8080/accounts | json_pp

# {
#    "balance" : 0,
#    "created_at" : "2023-03-26T01:48:30.82073Z",
#    "currency" : "USD",
#    "id" : 2,
#    "owner" : "bunny"
# }
```

**Get Account** (bunny) with a valid and unexpired `access_token`.

```bash
curl -X GET \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H 'Content-Type: application/json' \
    http://localhost:8080/accounts/1 | json_pp

# {
#    "balance" : 0,
#    "created_at" : "2023-03-26T01:47:04.48306Z",
#    "currency" : "EUR",
#    "id" : 1,
#    "owner" : "bunny"
# }
```

**Get Accounts** (bunny) with a valid and unexpired `access_token`.

```bash
curl -X GET \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H 'Content-Type: application/json' \
    "http://localhost:8080/accounts?page_id=1&page_size=5" | json_pp

# [
#    {
#       "balance" : 0,
#       "created_at" : "2023-03-26T01:47:04.48306Z",
#       "currency" : "EUR",
#       "id" : 1,
#       "owner" : "bunny"
#    },
#    {
#       "balance" : 0,
#       "created_at" : "2023-03-26T01:48:30.82073Z",
#       "currency" : "USD",
#       "id" : 2,
#       "owner" : "bunny"
#    }
# ]
```

**Create User** (monkey)

```bash
curl -X POST \
    -H 'Content-Type: application/json' \
    -d '{"username": "monkey", "full_name": "Monkey", "email": "monkey@email.com", "password": "password"}' \
    http://localhost:8080/users | json_pp

# {
#    "created_at" : "2023-03-26T01:51:24.055832Z",
#    "email" : "monkey@email.com",
#    "full_name" : "Monkey",
#    "password_changed_at" : "0001-01-01T00:00:00Z",
#    "username" : "monkey"
# }
```

**Login User** (monkey)

```bash
curl -X POST \
    -H 'Content-Type: application/json' \
    -d '{"username": "monkey", "password": "password"}' \
    http://localhost:8080/users/login | json_pp

# {
#    "access_token" : "v2.local.EMsz0VNTGZmi8QYUnohWFf8o08XCoaxsMvZQRBcwHv_OlUk2_gfqgVjSkQf_nq58CdwPFKcv-4wXQdek6kcGTC_KVFJ4pZAV3eXRB5IXQkswKFUj__jkUvpfZZebPWU7t1W3qITj5vn3qfo-x3YhIvPrjqHzaUJYeKeTcEK0oDlJRF06QCwE7_6hvObh-7rmXmac9MLRilXXd0VGd_UAY_1dLVAcXqdXNwJWwezHgHiToy7tayIt85aWgWITWPXNmAXntQ.bnVsbA",
#    "user" : {
#       "created_at" : "2023-03-26T01:51:24.055832Z",
#       "email" : "monkey@email.com",
#       "full_name" : "Monkey",
#       "password_changed_at" : "0001-01-01T00:00:00Z",
#       "username" : "monkey"
#    }
# }
```

**Create Account** (monkey, EUR) with a valid and unexpired `access_token`.

```bash
ACCESS_TOKEN="v2.local.EMsz0VNTGZmi8QYUnohWFf8o08XCoaxsMvZQRBcwHv_OlUk2_gfqgVjSkQf_nq58CdwPFKcv-4wXQdek6kcGTC_KVFJ4pZAV3eXRB5IXQkswKFUj__jkUvpfZZebPWU7t1W3qITj5vn3qfo-x3YhIvPrjqHzaUJYeKeTcEK0oDlJRF06QCwE7_6hvObh-7rmXmac9MLRilXXd0VGd_UAY_1dLVAcXqdXNwJWwezHgHiToy7tayIt85aWgWITWPXNmAXntQ.bnVsbA"

curl -X POST \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H 'Content-Type: application/json' \
    -d '{"currency": "EUR"}' \
    http://localhost:8080/accounts | json_pp

# {
#    "balance" : 0,
#    "created_at" : "2023-03-26T01:58:14.321717Z",
#    "currency" : "EUR",
#    "id" : 3,
#    "owner" : "monkey"
# }
```

**Create Account** (monkey, CAD)

```bash
curl -X POST \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H 'Content-Type: application/json' \
    -d '{"currency": "CAD"}' \
    http://localhost:8080/accounts | json_pp

# {
#    "balance" : 0,
#    "created_at" : "2023-03-26T02:00:09.25404Z",
#    "currency" : "CAD",
#    "id" : 4,
#    "owner" : "monkey"
# }
```

**Create Transfer** ((monkey, CAD) -> (bunny, EUR)) with a valid and unexpired
`access_token`.

```bash
curl -X POST \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H 'Content-Type: application/json' \
    -d '{"from_account_id": 4, "to_account_id": 1, "amount": 10, "currency": "CAD"}' \
    http://localhost:8080/transfers | json_pp

# {
#    "error" : "account [1]'s currency, EUR, mismatches CAD"
# }
```

**Create Transfer** ((monkey, EUR) -> (bunny, EUR)) with a valid and unexpired
`access_token`.

```bash
curl -X POST \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H 'Content-Type: application/json' \
    -d '{"from_account_id": 3, "to_account_id": 1, "amount": 10, "currency": "EUR"}' \
    http://localhost:8080/transfers | json_pp

# {
#    "from_account" : {
#       "balance" : -10,
#       "created_at" : "2023-03-26T01:58:14.321717Z",
#       "currency" : "EUR",
#       "id" : 3,
#       "owner" : "monkey"
#    },
#    "from_entry" : {
#       "account_id" : 3,
#       "amount" : -10,
#       "created_at" : "2023-03-26T02:03:57.27981Z",
#       "id" : 1
#    },
#    "to_account" : {
#       "balance" : 10,
#       "created_at" : "2023-03-26T01:47:04.48306Z",
#       "currency" : "EUR",
#       "id" : 1,
#       "owner" : "bunny"
#    },
#    "to_entry" : {
#       "account_id" : 1,
#       "amount" : 10,
#       "created_at" : "2023-03-26T02:03:57.27981Z",
#       "id" : 2
#    },
#    "transfer" : {
#       "amount" : 10,
#       "created_at" : "2023-03-26T02:03:57.27981Z",
#       "from_account_id" : 3,
#       "id" : 1,
#       "to_account_id" : 1
#    }
# }
```

### Create free-tier AWS account

1. Sign up in [AWS](https://aws.amazon.com).
1. Go to the AWS Management Console.
   - Sign in as Root user using the email you signed up.
1. Option+S to search "EC2".
   - Go to "Instance Types" and filter by "Free-Tier eligible = true".
     ![](https://i.imgur.com/8R2BaIi.png)

### AWS ECR

```bash
REPOSITORY_NAME=bank
```

```bash
aws ecr create-repository --repository-name $REPOSITORY_NAME
```

### AWS Role - Access AWS ECR by GitHub Actions Workflow

```bash
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
REPO_URI_PATTERN="repo:walkccc/go-boilerplate:*"
GITHUB_ACTIONS_ROLE="GitHubActionsRole"
```

```bash
aws iam create-role \
    --role-name GitHubActionsRole \
    --assume-role-policy-document '{
        "Version": "2012-10-17",
        "Statement": [
          {
            "Effect": "Allow",
            "Principal": {
              "Federated": "arn:aws:iam::'"$AWS_ACCOUNT_ID"':oidc-provider/token.actions.githubusercontent.com"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
              "StringEquals": {
                "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
              },
              "StringLike": {
                "token.actions.githubusercontent.com:sub": "'$REPO_URI_PATTERN'"
              }
            }
          }
        ]
      }'
```

```bash
aws iam attach-role-policy \
    --role-name $GITHUB_ACTIONS_ROLE \
    --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser
```

To delete the role, you need to detach the policy first:

```bash
aws iam detach-role-policy \
    --role-name $GITHUB_ACTIONS_ROLE \
    --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser
```

```bash
aws iam delete-role --role-name $GITHUB_ACTIONS_ROLE
```
