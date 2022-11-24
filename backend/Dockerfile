FROM golang:1.18.4-alpine3.16


RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apk update && apk add --no-cache gcc && apk add build-base
RUN apk add postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o app ./cmd/main/main.go