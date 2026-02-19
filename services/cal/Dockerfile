FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown

RUN go build \
    -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${BUILD_DATE}" \
    -o /app/server \
    ./services/cal/cmd/server/

FROM alpine:3.19

RUN apk add --no-cache ca-certificates && \
    addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -D appuser

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p /data && chown appuser:appuser /data

USER appuser

ENV CAL_PORT=8085
ENV CAL_DB_PATH=/data/cal.db

EXPOSE 8085

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8085/health || exit 1

ENTRYPOINT ["./server"]
