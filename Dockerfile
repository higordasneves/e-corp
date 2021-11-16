FROM golang:latest

WORKDIR /go/src/github.com/higordasneves/e-corp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

CMD ["./main"]