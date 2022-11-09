# Step1：打包go成執行檔
FROM golang:1.19-alpine AS builder

LABEL mantainer="POABOB <zxc752166@gmail.com>"

ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 
# GOPATH=/app

RUN mkdir /app
RUN mkdir /app/go
WORKDIR /app/go

COPY . /app/go
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-pano.output *.go
RUN go install github.com/rubenv/sql-migrate/...@latest
COPY ./config-prod.yml /app/go/config.yml

# GIN
WORKDIR /app/go/
EXPOSE 80
# Command to run the executable
CMD ["/app/go/go-pano.output"]