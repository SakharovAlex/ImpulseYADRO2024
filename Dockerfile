FROM golang:latest

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . /app

CMD ["go", "test", "-v", "./computerclub/..."]
