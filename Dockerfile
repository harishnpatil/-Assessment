
FROM golang:1.21-alpine AS builder

WORKDIR /app


COPY go.mod ./

COPY go.sum* ./


RUN go mod download


COPY *.go ./


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go-service .


FROM alpine:3.19


RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser


COPY --from=builder /go-service .


RUN chown appuser:appgroup go-service

USER appuser


EXPOSE 8080


ENV APP_HOST=0.0.0.0 \
    APP_PORT=8080 \
    LOG_LEVEL=info \
    ALLOWED_ORIGINS="" \
    ENABLE_METRICS=false

ENTRYPOINT ["./go-service"]
