# --- web
FROM node:24-alpine AS web
WORKDIR /web
RUN corepack enable
COPY web/package.json web/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build-only

# --- api
FROM golang:1.26-alpine AS api
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X qolda/internal/constant.Version=${VERSION}" -o /svc ./cmd

# --- app: API + SPA in one binary
FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=api /svc ./svc
COPY backend/docs ./docs
COPY backend/migrations ./migrations
COPY --from=web /web/dist ./web
ENV WEB_DIR=/app/web UPLOAD_DIR=/data/uploads HTTP_PORT=8090
EXPOSE 8090
CMD ["./svc"]
