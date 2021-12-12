FROM golang:1.17

WORKDIR /app

ADD go.mod go.sum /app/

RUN go mod download

ADD . /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init

RUN go build -o main main.go

EXPOSE 9000

CMD ["/app/main"]
