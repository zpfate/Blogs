FROM golang:1.25-alpine
WORKDIR /app

RUN apk add --no-cache git

ENV GOPROXY=https://goproxy.cn,direct

# 安装 Air
RUN go install github.com/air-verse/air@latest
# 只需要安装好环境，剩下的交给挂载
CMD ["air"]



# # 阶段 1: 编译阶段
# FROM golang:1.21-alpine AS builder
# WORKDIR /app
# COPY . .
# RUN go build -o hello-api main.go

# # 阶段 2: 运行阶段
# FROM alpine:latest
# WORKDIR /root/
# COPY --from=builder /app/hello-api .
# EXPOSE 8080
# CMD ["./hello-api"]



# 完全体

# FROM golang:1.23-alpine

# WORKDIR /app

# # 安装 air (在 2026 年，这是最快捷的安装方式)
# RUN go install github.com/air-verse/air@latest

# # 复制依赖文件
# COPY go.mod go.sum ./
# RUN go mod download

# # 复制源代码
# COPY . .

# # 暴露端口（确保与你代码中的端口一致）
# EXPOSE 8080

# # 使用 air 运行项目
# CMD ["air", "-c", ".air.toml"]