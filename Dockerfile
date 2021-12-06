FROM golang:1.17

WORKDIR /app

ADD go.mod go.sum /app/

RUN go mod download

ADD . /app

RUN go build -o main .

EXPOSE 9000

CMD ["/app/main"]
