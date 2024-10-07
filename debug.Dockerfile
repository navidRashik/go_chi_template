# Use a base image with Go installed
FROM golang:1.20

ENV CGO_ENABLED=1 GO111MODULE=on GOOS=linux
# Set the working directory inside the container


# Copy the Go module files and download dependencies
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go install github.com/go-delve/delve/cmd/dlv

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN GOPRIVATE=gitlab.upay.dev/golang CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build  -v -o main cmd/main.go

# Set the debug port for Delve
EXPOSE 40000


# Set the debug entrypoint command
ENTRYPOINT ["dlv", "debug", "--headless", "--listen=:40000", "--api-version=2", "exec", "./main"]
# ENTRYPOINT ["sleep", "10000"]
