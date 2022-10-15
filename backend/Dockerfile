FROM golang:1.18.4-alpine3.16 AS builder

COPY . /backend/
WORKDIR /backend/

RUN go mod download
RUN go build -o ./build/app cmd/main/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /backend/build/app .
COPY --from=0 /backend/etc/ etc/

EXPOSE 8000

CMD ["./app"]