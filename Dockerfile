# Use a minimal Alpine Linux image as the base image
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc g++ git openssh-client

# Set the working directory inside the container
WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=1 go build -a -installsuffix cgo -o monitor ./cmd/monitor

FROM alpine:latest as RUNNER
WORKDIR /app
COPY --from=builder /app/monitor .

ENTRYPOINT ["./monitor"]
