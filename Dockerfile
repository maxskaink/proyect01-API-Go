FROM golang:alpine3.20

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "run", "-c", ".air.toml"]