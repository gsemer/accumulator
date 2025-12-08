FROM golang:1.24.4

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY ./ /app/

RUN go build -o accumulator cmd/main.go

CMD [ "/app/accumulator" ]