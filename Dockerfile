FROM golang:1.20.1-alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o crudly

ENTRYPOINT ["./crudly"]
