FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /auth-service ./src/cmd/main.go

FROM alpine:latest

WORKDIR /

COPY --from=builder /auth-service /auth-service

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

EXPOSE 8081

CMD ["/auth-service"]