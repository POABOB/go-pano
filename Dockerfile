# Dev Stage 
FROM golang:1.19-alpine AS dev

LABEL mantainer="POABOB <zxc752166@gmail.com>"

# 建立環境變數
ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPATH=/app \
    PATH="$PATH:/app/bin/linux_amd64/"

# 指定工作目錄
WORKDIR /app/go/

# 把當前專案複製到/app/go裡
COPY . /app/go

# 安裝環境依賴函式庫
RUN go mod download \
    && go install github.com/rubenv/sql-migrate/...@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest\
    && go install github.com/google/wire/cmd/wire@latest \
    && go install github.com/codegangsta/gin \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# 暴露端口
EXPOSE 80

# 執行
CMD ["gin", "--appPort", "3000", "--port", "80", "run", "main.go"]

# Build Stage
FROM golang:1.19-alpine AS builder

LABEL mantainer="POABOB <zxc752166@gmail.com>"

# 建立環境變數
ENV GO111MODULE=auto \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 指定工作目錄
WORKDIR /app/go/

# 把當前專案複製到/app/go裡
COPY . /app/go

# 安裝環境依賴函式庫
RUN go mod tidy \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-pano.output *.go

# Deploy Stage
FROM alpine:3.16.2 AS prod

# 指定工作目錄
WORKDIR /app/go/

COPY --from=builder /app/go/go-pano.output /app/go/go-pano.output
COPY --from=builder /app/go/config-prod.yml /app/go/config.yml

# 暴露端口
EXPOSE 80

# 執行
CMD ["/app/go/go-pano.output"]