## 阶段一：破冰与基础语法（第 1 ~ 2 周）

> **给大前端开发的提示**：Go 没有类（Class）和继承，只有结构体（Struct）和组合（Composition）；Go 的接口（Interface）是隐式实现的（Duck Typing），这和 Swift/Dart 的显式实现非常不同。

### 1. 基础环境与工具链

- Go 工具链安装（`go env`）、GOPATH 与 Go Modules (`go mod`) 机制。
- IDE 配置：VS Code (Go 扩展) 或 GoLand (强烈推荐，JetBrains 宇宙，你懂的)。
- Go 的代码组织：Package（包）的导入与可见性规则（首字母大小写决定公私有）。

### 2. 核心语法基础

- **变量与类型**：静态类型、类型推导、指针（Go 的指针不能进行运算，别怕，它更像 Swift 的 `inout`）。
- **控制流**：`if-else`、`for` 循环（Go 只有 for 这一种循环）、`switch-case`（默认不需要 `break`）。
- **基础容器**：Array（数组，定长）与 **Slice（切片，动态数组，核心重点）**、Map（哈希表）。

### 3. 函数与方法

- 多返回值、具名返回值。
- **Defer 延迟调用**：Go 特有的资源释放机制（类似 Swift 的 `defer`）。
- **方法（Methods）**：为结构体（Struct）绑定方法，理解结构体接收者与指针接收者的区别。

### 🎯 阶段一目标：

- **完成一个命令行工具（CLI）**：比如一个“本地 Todo List 管理器”，支持通过终端参数（`flag` 包）增删改查任务，并能将数据持久化保存为 JSON 文件。
- **掌握标准**：熟悉 `go build`、`go run`、`go fmt`，能熟练运用切片和 Map 进行数据处理。

## 阶段二：面向接口与工程化起步（第 3 ~ 4 周）

> **给大前端开发的提示**：忘记 Swift 的 `Protocol` 继承和 Dart 的 `abstract class`。Go 的接口是“只要你长得像鸭子，你就是鸭子”。

### 1. 隐式接口（Interfaces）

- 接口的定义与空接口（`interface{}` 或新版 Go 的 `any`）。
- 非侵入式接口的设计哲学：服务提供方不定义接口，由**使用方**定义接口。
- 类型断言（Type Assertion）与 Type Switch（类似 Swift 的 `is` 和 `as?`）。

### 2. 错误处理哲学（Error Handling）

- Go 没有 `try-catch`。理解 `if err != nil` 的显式错误处理哲学。
- 自定义 Error、`errors.Is` 与 `errors.As`（Go 1.13+ 错误包装）。
- Panic 与 Recover：什么时候该崩溃，什么时候该恢复（对比 iOS 的 Crash 机制）。

### 3. Go 泛型（Generics）

- Go 1.18+ 引入的泛型基础：类型约束（Type Constraints）、`any` 与 `comparable`。

### 🎯 阶段二目标：

- **重构命令行工具**：将“数据存储”抽象为接口（`Storage`），分别实现“文件存储”和“内存存储”两种模式，在启动时动态注入。
- **掌握标准**：理解依赖注入的初级概念，习惯编写带有完整错误处理的代码，不再让程序轻易 Panic。

## 阶段三：并发编程与网络基础（第 5 ~ 7 周）—— **核心跨越**

> **给大前端开发的提示**：Flutter 是单线程 Event Loop（通过 Isolate 实现并发），iOS 依赖 GCD（线程池）。而 Go 拥有恐怖的 **CSP 并发模型**（Goroutine + Channel），几KB的协程就能并发数十万个，这是 Go 的杀手锏。

### 1. 协程与通道（Goroutines & Channels）

- `go` 关键字：如何启动一个轻量级协程。
- **Channel（通道）**：有缓冲 vs 无缓冲通道，单向通道。
- **Select 语句**：多路复用，处理 Channel 的超时与退出。

### 2. 同步原语与并发安全

- `sync.Mutex`（互斥锁）与 `sync.RWMutex`（读写锁）。
- `sync.WaitGroup`（类似 GCD 的 `DispatchGroup`）。
- `sync.Once`（单例模式的终极实现）。
- `sync.Map` 与原子操作（`sync/atomic`）。

### 3. Context 机制（上下文）

- **Context 的核心作用**：超时控制、取消信号传递、跨协程元数据传递（微服务必备）。
- `context.WithTimeout`、`context.WithCancel`。

### 4. 基础网络编程

- 使用标准库 `net/http` 编写基础的 HTTP Server 和 Client。

### 🎯 阶段三目标：

- **编写一个高性能图片/文件下载器**：支持限制并发数（如最多同时下载 5 个文件，使用 Channel 缓冲区实现信道限流），支持计算下载进度，并支持通过 `Context` 在中途取消下载。
- **掌握标准**：深刻理解什么是 Race Condition（竞争条件），学会使用 `go test -race` 检查并发安全。

## 阶段四：Web 框架与数据库实战（第 8 ~ 10 周）

> **给大前端开发的提示**：到了这一步，你将正式接触后端业务的核心——CRUD 与数据持久化。你可以把后端的路由（Router）类比为 Flutter 的 Navigator 或 iOS 的 Router。

### 1. 流行 Web 框架（以 Gin 为例）

- Gin 框架的路由、中间件设计（Middleware，类似服务端的“拦截器”）。
- 参数绑定与校验（JSON、Query、Form 绑定到 Struct）。

### 2. 数据库与 ORM 操作

- SQL 基础与 MySQL / PostgreSQL 连接。
- 使用 **GORM** 或 **ent** 框架进行 CRUD、关联查询（One-to-Many, Many-to-Many）。
- 数据库连接池配置、事务处理（Transaction）。

### 3. 项目分层架构（工程化）

- 经典的后端三层架构：Controller（API层） -> Service（业务逻辑层） -> DAO/Repository（数据访问层）。
- 日志选型（Zap 或 Logrus）与配置管理（Viper）。

### 🎯 阶段四目标：

- **开发一套简易的“App 后端 API 系统”**（如新闻资讯 App 的后端）：
  - 包含用户注册/登录（JWT 认证中间件）。
  - 文章列表分页、点赞、评论功能（MySQL 持久化）。
  - 集成日志记录和统一的错误响应格式。
- **掌握标准**：能够独立搭建一个符合业界规范、分层清晰的 Web 后端项目。

## 阶段五：进阶与微服务架构（第 11 ~ 12 周）—— **迈向资深**

> **给大前端开发的提示**：移动端常使用 Restful API，但后端微服务之间通信更倾向于高性能的 **RPC**。这里你将接触到端到端联调的硬核技术。

### 1. gRPC 与 Protobuf

- Protocol Buffers (Proto3) 语法定义。
- 使用 Go 实现 gRPC 服务端与客户端。
- **大前端联动**：尝试用 Flutter/Swift 直接调用你的 gRPC 服务。

### 2. 缓存与中间件进阶

- Redis 的基本数据类型与 Go 客户端（go-redis）使用。
- **缓存旁路策略（Cache Aside）**：如何利用 Redis 为 MySQL 加速。

### 3. 性能分析与优化（Go 专家必备）

- **pprof 工具链**：如何抓取 CPU、内存、Goroutine 性能 profile。
- 理解 Go 的 **GMP 调度模型**与 **GC（垃圾回收）三色标记法**（理论深度）。

### 🎯 阶段五目标：

- **终极实战：高并发秒杀/抢购系统组件**：
  - 使用 gRPC 作为通信协议。
  - 利用 Redis 进行商品库存预扣减，防止超卖。
  - 使用 pprof 压测并找出系统瓶颈。
- **掌握标准**：理解分布式系统的基本概念，能够分析 Go 程序的内存泄漏和性能死锁。