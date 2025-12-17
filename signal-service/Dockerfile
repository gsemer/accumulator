FROM golang:1.24.4

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY ./ /app/

RUN go build -o signal cmd/main.go

CMD [ "/app/signal" ]