# Go语法定义

## 1. Go 语法关键定义
- **package main**: 声明这是一个可执行程序，而不是一个插件库。
- **import**: 导入工具包。`encoding/json` 用于处理 JSON，`net/http` 用于网络通信。
- **struct (结构体)**: Go 没有“类”的概念，用结构体来定义数据的形状。
- **首字母大写**: 在 Go 中，只有首字母大写的变量才能被外部（如 JSON 编码器）看到。
- **`json:"..."` 标签**: 控制 JSON 输出时的键名（Key）。

## 2. JSON 接口实战代码
```go
// 1. 定义返回的数据模型
type Response struct {
    Message string `json:"message"` // 返回给前端叫 message
    Status  int    `json:"status"`  // 返回给前端叫 status
}

// 2. 处理请求
http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
    // 设置响应头，告诉浏览器：我要发 JSON 了
    w.Header().Set("Content-Type", "application/json")
    
    // 实例化对象
    res := Response{ Message: "Success", Status: 200 }
    
    // 【核心句】将对象转为 JSON 并输出
    json.NewEncoder(w).Encode(res)
})
```

## 3. 容器更新步骤

每次修改代码后，请按此顺序操作以节省资源：

```shell
# 清理旧容器：
docker stop $(docker ps -q)  # (停止正在跑的)

docker rm $(docker ps -aq)  # (删除过期的记录)

# 构建新镜像：
docker build -t go-hello-api .

# 后台运行 (释放终端控制权)：
docker run -d -p 8080:8080 go-hello-api
```

## 4. 路由与多路径处理 (Routing)

### 核心语法定义
- **http.HandleFunc(path, handler)**: 
    - `path`: 字符串，如 `/health`。
    - `handler`: 一个匿名函数或具名函数，决定了当用户访问这个路径时，服务器执行什么逻辑。
- **w.WriteHeader(code)**: 发送 HTTP 状态码（如 200 OK, 404 Not Found）。

### 路由分发的意义
- **模块化**: 不同功能拆分到不同的 URL 路径下。
- **职责分离**: 比如 `/` 给浏览器看，`/api` 给移动端 App 解析 JSON。

## 5. 高级开发流：热重载 (Hot Reload)

### 核心工具
- **Air**: Go 语言的实时重载工具。

### 核心命令
- **挂载卷**: `-v [宿主机路径]:[容器路径]`。
- **作用**: 实现本地代码与容器环境的“同步呼吸”，极大减少 `docker build` 产生的 CPU 负载。

### 为什么对低配 Mac 友好？
- 避开了繁重的 Docker 构建过程。
- 容器只在启动时消耗一次大资源，之后仅做增量编译。