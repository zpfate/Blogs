# Go API 容器化实战笔记





## 1. 核心文件

- **main.go**: Go 源代码。
- **Dockerfile**: 镜像构建指令。

## 2. 关键命令流程

### 启用 Go Modules 模式

```shell
 go env -w GO111MODULE=on
```

### 设置国内代理 (推荐七牛云的 goproxy.cn)

```shell
go env -w GOPROXY=https://goproxy.cn,direct
```



### 构建镜像 (确保在项目目录下，注意末尾的点)

```shell
docker build -t go-hello-api .
```

### 运行容器

```
docker run -p 8080:8080 go-hello-api
```

### 常见问题解决

```shell
lsof -i tcp:8080     # 排查端口8080占用
sudo apachectl stop  # 停止 macOS 自带的 httpd 服务
sudo killall httpd    # 强制杀掉进程
```

### 热重载

