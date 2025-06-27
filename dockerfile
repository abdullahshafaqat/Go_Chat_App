# GO_CHAT_APP/Dockerfile
FROM golang:1.24.2-alpine


WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8003

CMD ["./main"]
