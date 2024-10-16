# Start from golang base image
FROM golang:1.20 as builder



# Working directory
WORKDIR /

# Copy go mod and sum files
COPY go.mod go.sum ./


# Download all dependencies
RUN go mod download
# Copy everythings
COPY . .

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build  -v -o main cmd/main.go


# Start a new stage from scratch
FROM debian:stable
# RUN apk --no-cache add tzdata
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Also copy config yml file
# COPY --from=builder .env .
COPY --from=builder main .
COPY --from=builder Makefile .
COPY --from=builder assets/ assets/

RUN apt update -y && apt install ca-certificates make curl -y && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xz -- migrate
# COPY --from=builder /app/VERSION .

EXPOSE 8000
ENTRYPOINT [ "./main" ]
