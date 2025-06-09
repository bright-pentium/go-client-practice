# Stage 1: Build the application
FROM golang:1.23.8 AS builder

WORKDIR /opt/

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY docs ./docs
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server 

# # Stage 2: Minimal runtime image
FROM alpine:3.21

WORKDIR /opt

COPY --from=builder /opt/app ./

CMD ["./app"]