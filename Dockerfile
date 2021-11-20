FROM golang:latest

WORKDIR /go/src/github.com/higordasneves/e-corp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

CMD ["./app"]