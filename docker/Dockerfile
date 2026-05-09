# syntax=docker/dockerfile:1.7
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/swaggo/swag/cmd/swag@v1.16.4

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    swag init -g cmd/api/main.go -o docs --parseInternal

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -trimpath -o /out/gym-pro-api ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata wget && \
    addgroup -S app && adduser -S app -G app

ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /app

COPY --from=builder --chown=app:app /out/gym-pro-api ./gym-pro-api
COPY --from=builder --chown=app:app /app/migrations ./migrations

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=5 \
    CMD wget -qO- http://localhost:8080/health >/dev/null 2>&1 || exit 1

CMD ["./gym-pro-api"]
