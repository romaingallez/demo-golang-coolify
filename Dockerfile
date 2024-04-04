# Build stage
FROM golang:1.16 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o coolify .

# Cache stage
FROM alpine:latest AS cache

WORKDIR /app

COPY --from=build /app/coolify .

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates curl

COPY --from=cache /app/coolify .

CMD ["./coolify"]