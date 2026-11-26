# CLAUDE.md

Guide for Claude Code in this repository. A Go service on Clean Architecture + DDD,
with gRPC/grpc-gateway, Postgres (pgx) and proto generation through the Makefile.
Part of the qolda monolith: the same binary serves the built frontend (`../web/dist`); see the root README.

The stack conventions live in global Claude Code skills (`crud`, `mobone`,
`golang-service`, `golang-samber-lo`), picked up automatically by their descriptions.
They are installed from `github.com/rendau/go-service-skill` (`install.sh`); if they are missing from
the skill list, suggest that the developer install them.

Code comments are in English; explain *why*, not *what*.

---

## Project structure

### Top level
- `cmd/main.go`: entrypoint, starts `internal/app.App`.
- `internal/`: business logic and infrastructure (private packages).
- `api/proto/`: this service's `.proto` sources.
- `pkg/proto/`: generated protobuf/grpc/gateway code (**do not edit by hand**).
- `migrations/`: Postgres SQL migrations.
- `docs/`: swagger JSON and Swagger UI (`docs/index.html`), served at `/docs/*` on the app port ("Try it out" uses
  the web app's session) and on the system port.
- `vendor-proto/`: external `.proto` dependencies (`make refresh-vendor-proto` from `PROTOS_REPO`
  in the `Makefile`; empty means they are updated only together with the template).
- `Makefile`: build and proto generation. The monolith Docker image is `../Dockerfile`.
- `.env.example`: sample environment for running locally.

### Internal packages (`internal/`)
- `internal/app/`: application assembly: servers, migrations, HTTP gateway, DI.
  - `app.go`: dependency graph and component startup.
  - `grpc.go`: gRPC server + session, recovery and error interceptors.
  - `grpc_gateway.go`: HTTP gateway, CORS, error handler.
  - `system_http_server.go`: system HTTP server (`SYSTEM_HTTP_PORT`, 3003): /healthcheck, /docs/*.
  - `web.go`: root HTTP handler: `/api/*` → gateway, `/files/*` → uploaded photos, everything else → SPA.
  - `migration.go`: runs migrations from `migrations/`.
- `internal/config/`: configuration from env (see `config.go`).
- `internal/handler/`: transport layer.
  - `grpc/`: gRPC handlers.
  - `grpc/dto/`: protobuf ↔ domain model conversion; shared helpers in `dto/common.go`
    (`DecodeListParams`, `JsonObjToGrpcStruct`, `GrpcStructToJsonObj`).
  - `http/`: plain HTTP handlers (Google OAuth redirects, file upload).
- `internal/usecase/`: usecase layer (validation, orchestration of services and domain services).
- `internal/domain/`: domain model, services and repositories.
  - `*/model/`: domain structures (entities).
  - `*/service/`: domain services (invariants/logic).
  - `*/repo/`: repositories.
  - `common/`: shared models/utils/PG base repository.
- `internal/service/`: infrastructure/integration services (Google OAuth, file store).
- `internal/errs/` and `internal/constant/`: shared error codes and constants (`Version`, `Commit`,
  `BuiltAt` come from ldflags at build time).

---

## Architecture: layers and dependencies

- **Transport** (`internal/handler/*`):
  - Works only with protobuf DTOs and the usecase layer: a handler takes `*usecase.Usecase` directly,
    without an interface (handlers are thin and not tested separately).
  - Never touches repositories or services directly.
- **Usecase** (`internal/usecase/*`):
  - Entry point from the transport layer (requests from outside).
  - Validates input parameters.
  - Orchestrates domain services and `internal/service/*` services.
  - Input/output are domain models (not protobuf).
  - Takes dependencies (domain services, services) through interfaces declared in its own
    package (`interfaces.go`).
  - Preferably doesn't call neighbouring usecases.
- **Domain** (`internal/domain/*`):
  - `model/`: data structures, entities.
    - **No serialization tags** (`json:"..."`, `yaml:"..."` etc.) on domain model fields:
      the model knows nothing about transport or storage format.
  - `service/`: domain operations and invariants. May use only the repository.
  - `repo/`: storage access. May contain subfolders for different storage types or mocks.
- **Service** (infrastructure/background/external integrations, `internal/service/*`):
  - Background processes, integrations with external systems.
  - Extracted or reusable logic.
  - May use other `internal/service/*` services and `internal/domain/*/service` domain services.
  - Never calls the usecase layer.
- **Composition** (`internal/app/`):
  - Wires dependencies, starts servers, migrations and background services.

### Dependency rule
```
handler        → usecase
usecase        → domain service
usecase        → service
service        → service
service        → domain service
domain service → repo
```
- Reverse dependencies are **forbidden**.
- Only `domain service` accesses the `repo` layer.

### Storage
- Postgres: all domain entities (see `migrations/*`).
- **Table names are always singular**: `usr`, `listing`, `item`
  (not `usrs`, `listings`, `items`). The same value goes into the repository's `TableName`.

### Migrations
- Files in `migrations/` named `NNNNNN_<name>.up.sql` / `.down.sql` (golang-migrate).
- Every `DROP` in a `down` migration must have `CASCADE`
  (e.g. `drop table if exists <table> cascade;`).
- `down` drops objects in the reverse order of `up` (respecting foreign keys).

### API
- gRPC services: `api/proto/qolda_v1/*`.
- HTTP gateway: grpc-gateway + swagger (`docs/api.swagger.json`), mounted under `/api`.
- DTO mapping: `internal/handler/grpc/dto`.
- **REST paths are always singular**:
  `/listing`, `/listing/{id}`, `/usr/{id}` (not `/listings`, `/usrs`).
- **Pagination: `page` starts at `0` in every `list` request.**
- **Updates use HTTP `PUT`** (not `PATCH`): `put: "/<entity>/{id}"`.

### Errors and validation
- Semantic errors go through `internal/errs` (see the gRPC interceptor in `internal/app/grpc.go`).
  A business refusal (not found, bad request, wrong code) is always an `errs.Err` (the service's own
  code or a common one) and is logged at info level: it's a refusal, not a failure. Any other error
  becomes `service_not_available` (`errs.IsFailure`) and is logged at error level.
- Parameter validation happens in the usecase.
- Never return errors without wrapping (`fmt.Errorf("...: %w")`).
- Always use `errors.Is` and `errors.AsType` with errors. Avoid comparing errors directly (`==`)
  and type assertions (`err.(*MyError)`) so wrapped errors are handled correctly.

### Rules for changing code
- gRPC DTOs must not leak into domain services.
- `pkg/proto` and `docs/api.swagger.json` are generated (updated through the Makefile).
- Generation uses only the local toolchain in `bin/` (git-ignored): `make proto-tools` installs
  protoc and plugins at pinned versions (protoc from `PROTOC_VERSION` in the Makefile, protoc-gen-go and
  grpc-gateway from `go.mod`); `generate-proto` runs it itself. Global protoc/plugins
  are not used; bump versions in the Makefile/`go.mod`, for the whole team at once.
- In tests always prefer `testify`: `require` for checks that stop the test,
  `assert` for everything else.
- For worker pools / parallel processing use `errgroup`
  (golang.org/x/sync/errgroup) rather than hand-rolled `sync.WaitGroup` + channels.

---

## Composition root (`internal/app/app.go`)

`app.go` is the single composition point of the application: it builds the dependency graph
(`repo → service → usecase → handler`) and describes the lifecycle. No business logic here,
only wiring components and managing their start/stop.

### The `App` type
- Fields hold **only what must be managed after `Init`**: things to stop explicitly
  (servers), wait for (background services/handlers), close (pgx pool), plus the root `ctx` with its
  `ctxCancel` and `exitCode`.
- Local links of the graph (repo, service, usecase, handler) that are only needed for assembly and
  passed on immediately are **not** fields; they are local variables inside `Init`.

### Composition root imports
- Grouped into blocks separated by a blank line: standard library → external
  dependencies → internal project packages.
- Internal constructor packages are imported with a `P` suffix alias spelling the whole path in camelCase
  (e.g. `internal/service/mdm` → `serviceMdmP`, `internal/handler/grpc` → `handlerGrpcP`,
  `internal/service/2gis/service` → `service2gisServiceP`).

### Lifecycle phase methods
A fixed set of methods, each doing exactly one thing:
- `Init`: create and wire all dependencies.
- `PreStartHook`: actions before start.
- `Start`: start servers and background services.
- `Listen`: block until an OS signal (`SIGINT`/`SIGTERM`).
- `Stop`: cancel the context and gracefully stop servers.
- `WaitJobs`: wait for background jobs to finish.
- `Exit`: close resources and exit with `exitCode`.

### `Init` style
- Assembly goes **top-down in dependency order**: infrastructure (logger,
  pgx pool, migrations) → domain blocks → services → servers.
- Each logical block starts with a short lowercase label comment
  (`// listing`, `// auth`, `// grpc server`).
- Inside a block the chain is uniform:
  `repo := ...New(...)` → `service := ...New(repo)` → `usecase := ...New(service)` →
  `handler := ...New(usecase)`.
- Blocks that don't need outside variables are wrapped in an anonymous block `{ ... }`
  to limit scope and separate them visually.

### Error handling during init
- The `errCheck(err, msg)` helper is used: during assembly any error is fatal
  (log + `os.Exit(1)`). Errors from `Init` are not returned upward.
- A single `var err error` declared up front is reused throughout `Init`.

### Start / WaitJobs / Stop pairing
- Every background component started in `Start()` has a matching call in
  `WaitJobs()` (`.Wait()`) and/or `Stop()`. **Components are listed in the same order in
  all three methods**, which makes them easy to read and compare.

### Configuration and optional components
- All parameters come from the single global config (`config.Conf.*`) right where they are used.
- Optional components are enabled by a config condition (`if config.Conf.X != "" { ... }`),
  falling back to an in-memory/no-op implementation behind a shared interface.

---

## Runtime and configuration

### Startup
- Entry: `cmd/main.go` → `internal/app.App`.
- On start:
  - load env (autoload `.env`),
  - set up the logger,
  - pgx pool,
  - migrations (`internal/app/migration.go`),
  - start gRPC + HTTP gateway.

### Environment variables
- Described in `internal/config/config.go`.
- Sample: `.env.example`.

### System HTTP server
- A separate server on `SYSTEM_HTTP_PORT` (default `3003`).
- Serves service endpoints: `/healthcheck`, `/docs/*`.

### Observability
- pulse, Prometheus metrics and Jaeger tracing from the template were **removed**; don't bring them
  back unless asked. Logs: `log/slog` (JSON when `DEBUG=false`).

### Docs and healthcheck
- `/healthcheck`: HTTP healthcheck (200 OK), system port.
- `/docs/*`: swagger (`docs/api.swagger.json`), system port.

### Build
- `make build` produces the `cmd/build/svc` binary (version, full commit and build time go in via
  ldflags into `internal/constant` and are logged on start).
- `../Dockerfile` (multi-stage) builds the frontend and the binary and puts `docs/`, `migrations/` and
  `web/` into `/app`.

### Session and access
- `internal/domain/session`: JWT (HS256): access (`ACCESS_TTL`) + refresh (`REFRESH_TTL`).
- The `GrpcInterceptorSession` gRPC interceptor and `handler/http/withSession` put the session into ctx;
  an invalid token means anonymous. Only the usecase checks access: `NotAuthorized` → 401, `NoPermission` → 403.
- Sign-in: Google OAuth (code + PKCE, `handler/http/auth.go`) or dev sign-in (`AUTH_DEV_LOGIN=true`).
  The e-mail domain is checked on our side (`ALLOWED_DOMAINS`, subdomains allowed).

### Change verification flow
```
make generate-proto  →  gofmt  →  go test ./...  →  go run ./cmd/.
```
