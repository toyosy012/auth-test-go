FROM golang:1.19.5-alpine3.17 AS builder

WORKDIR /root

COPY . /root
RUN go get && go install && GOOS=linux go build main.go

FROM alpine:3.17

EXPOSE 8080
COPY --from=builder /root/main .
COPY --from=builder /root/wait-for-it.sh .
RUN chmod +x wait-for-it.sh