# Build stage
FROM golang:1.19-alpine3.17 AS builder

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/gmc cmd/gmc/main.go

# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/ .

EXPOSE 8080

ENTRYPOINT [ "/app/gmc" ]