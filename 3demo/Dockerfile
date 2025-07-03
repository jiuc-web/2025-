# Dockerfile
FROM golang:1.21-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
# 复制go.mod和go.sum文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o book-management-system .

# 使用轻量级镜像运行
FROM alpine:latest

WORKDIR /app

# 从builder阶段复制可执行文件
COPY --from=builder /app/book-management-system .

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./book-management-system"]