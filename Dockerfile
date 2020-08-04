FROM golang:latest

WORKDIR /var/stormy

COPY . .

RUN go mod download
RUN go build -o ./bin/stormy ./cmd/main.go

CMD ["bash", "entrypoint.sh"]
