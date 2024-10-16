FROM golang:1.23.1
WORKDIR /usr/src/app

ENV CGO_ENABLED=1 GO111MODULE=on GOOS=linux


RUN go install github.com/cosmtrek/air@latest

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest



COPY . .

RUN go mod tidy

ENTRYPOINT [ "air", "./service/main.go" ]
