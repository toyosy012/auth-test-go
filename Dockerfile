FROM golang:1.20-alpine3.17 AS builder

WORKDIR /root
COPY go.* ./
RUN go mod download && go install github.com/swaggo/swag/cmd/swag@v1.8.0
COPY . /root
RUN swag init && GOOS=linux go build main.go

FROM alpine:3.17

EXPOSE 8080
COPY --from=builder /root/main .
COPY --from=builder /root/wait-for-it.sh .
RUN chmod +x wait-for-it.sh
