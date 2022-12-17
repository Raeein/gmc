# Build stage
FROM golang:latest AS builder

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/gmc cmd/gmc/main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gmc .

EXPOSE 8080

ENTRYPOINT [ "/app/gmc" ]