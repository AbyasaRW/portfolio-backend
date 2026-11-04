ARG GO_VERSION=1.26.5
ARG ALPINE_VERSION=3.22

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio-backend .

FROM alpine:${ALPINE_VERSION}

WORKDIR /app
RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder --chown=app:app /app/portfolio-backend .

USER app
EXPOSE 8080
CMD ["./portfolio-backend"]
