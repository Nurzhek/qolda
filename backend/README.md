# qolda — backend

Go service on gotemplate (Clean Architecture + DDD, gRPC + grpc-gateway, pgx/mobone).
How everything runs together with the frontend: the root `../README.md`. Conventions: `CLAUDE.md`.

### Run locally

```
cp .env.example .env          # Postgres from ../docker-compose.yml on 5435
go run ./cmd                  # API: http://localhost:8090/api, swagger: http://localhost:3004/docs/api.swagger.json
```

### Checks

```
make generate-proto && gofmt -l internal && go vet ./... && go test ./...
```

### Proto generation

```
make generate-proto
```

The proto toolchain lives in `bin/` (git-ignored) and is installed by `make proto-tools`, which
`generate-proto` runs first: `protoc` of the version pinned in the Makefile (`PROTOC_VERSION`) and
plugins of the versions from `go.mod`. Global `protoc` and plugins are not used.

### Migrations

Applied automatically on start (`migrations/`, golang-migrate). To add one:

```
migrate create -ext sql -dir migrations mg_name
```

### DB dump / restore

```
pg_dump --no-owner -Fc -U qolda qolda -f ./qolda.custom
pg_restore --no-owner -d qolda -U qolda ./qolda.custom
```
