# 常见问题 (FAQ)

本文档收集了 Rei0721 项目开发过程中的常见问题和解决方案。

## 安装和配置

### Q: 如何安装 Go？

A: 访问 [golang.org](https://golang.org/dl/) 下载适合您操作系统的 Go 版本。

```bash
# 验证安装
go version
# 输出应该是 go version go1.24.6 或更高
```

### Q: 如何配置 GOPATH？

A: 在 Go 1.11+ 中，GOPATH 不再是必需的。使用 Go Modules 管理依赖：

```bash
# 初始化模块
go mod init rei0721

# 下载依赖
go mod download
go mod tidy
```

### Q: 如何更改服务端口？

A: 有三种方式：

**方式 1: 修改配置文件**
```yaml
# configs/config.yaml
server:
  port: 8081
```

**方式 2: 环境变量**
```bash
SERVER_PORT=8081 go run ./cmd/server/main.go
```

**方式 3: 运行时更新**
```go
configManager.Update(func(cfg *config.Config) {
    cfg.Server.Port = 8081
})
```

### Q: 如何使用不同的数据库？

A: 修改 `DB_DRIVER` 环境变量：

```bash
# PostgreSQL
DB_DRIVER=postgres go run ./cmd/server/main.go

# MySQL
DB_DRIVER=mysql go run ./cmd/server/main.go

# SQLite
DB_DRIVER=sqlite go run ./cmd/server/main.go
```

### Q: SQLite 需要特殊配置吗？

A: SQLite 需要 CGO 支持。如果遇到问题：

```bash
# 启用 CGO
CGO_ENABLED=1 go run ./cmd/server/main.go

# 或使用 PostgreSQL/MySQL
DB_DRIVER=postgres go run ./cmd/server/main.go
```

## 开发和调试

### Q: 如何启用调试模式？

A: 设置 `SERVER_MODE=debug`：

```bash
SERVER_MODE=debug LOG_LEVEL=debug go run ./cmd/server/main.go
```

### Q: 如何查看详细日志？

A: 设置 `LOG_LEVEL=debug`：

```bash
LOG_LEVEL=debug go run ./cmd/server/main.go
```

### Q: 如何运行测试？

A: 使用 `go test` 命令：

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/router/... -v

# 生成覆盖率报告
go test ./... -cover

# 生成覆盖率 HTML 报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Q: 如何格式化代码？

A: 使用 `go fmt` 或 `gofmt`：

```bash
# 格式化当前目录
go fmt ./...

# 或使用 gofmt
gofmt -w .
```

### Q: 如何检查代码质量？

A: 使用 `go vet`：

```bash
# 检查代码
go vet ./...

# 使用 golangci-lint (需要安装)
golangci-lint run ./...
```

### Q: 如何调试 HTTP 请求？

A: 使用 `curl` 或 Postman：

```bash
# 检查健康状态
curl http://localhost:8080/health

# 用户注册
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secure_password_123"
  }'

# 查看详细信息
curl -v http://localhost:8080/health
```

## 数据库问题

### Q: 如何创建数据库？

A: 根据您使用的数据库类型：

**PostgreSQL**:
```bash
createdb rei0721
```

**MySQL**:
```bash
mysql -u root -p -e "CREATE DATABASE rei0721;"
```

**SQLite**:
```bash
# 自动创建
DB_DRIVER=sqlite DB_NAME=./rei0721.db go run ./cmd/server/main.go
```

### Q: 如何重置数据库？

A: 删除数据库并重新创建：

**PostgreSQL**:
```bash
dropdb rei0721
createdb rei0721
```

**MySQL**:
```bash
mysql -u root -p -e "DROP DATABASE rei0721; CREATE DATABASE rei0721;"
```

**SQLite**:
```bash
rm rei0721.db
```

### Q: 如何查看数据库连接状态？

A: 查看应用日志：

```bash
# 应该看到
# {"level":"info","msg":"database connected successfully"}
```

### Q: 如何优化数据库连接？

A: 调整连接池参数：

```yaml
database:
  maxOpenConns: 100    # 最大连接数
  maxIdleConns: 10     # 最大空闲连接数
```

### Q: 如何处理数据库超时？

A: 增加超时时间：

```yaml
server:
  readTimeout: 30      # 读超时 (秒)
  writeTimeout: 30     # 写超时 (秒)
```

## 性能和优化

### Q: 如何提高应用性能？

A: 几个建议：

1. **使用生产模式**:
```bash
SERVER_MODE=release go run ./cmd/server/main.go
```

2. **调整日志级别**:
```bash
LOG_LEVEL=warn go run ./cmd/server/main.go
```

3. **优化数据库连接**:
```yaml
database:
  maxOpenConns: 200
  maxIdleConns: 20
```

4. **启用 Redis 缓存**:
```yaml
redis:
  enabled: true
```

### Q: 如何监控应用性能？

A: 使用 pprof：

```bash
# 导入 pprof
import _ "net/http/pprof"

# 访问性能数据
curl http://localhost:6060/debug/pprof/
```

### Q: 如何处理高并发？

A: 调整协程池大小：

```go
scheduler.Config{
    PoolSize:       10000,
    ExpiryDuration: time.Second,
}
```

## 部署问题

### Q: 如何构建二进制文件？

A: 使用 `go build`：

```bash
# 构建
go build -o bin/server ./cmd/server

# 运行
./bin/server
```

### Q: 如何使用 Docker？

A: 创建 Dockerfile：

```dockerfile
FROM golang:1.24.6 AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY configs/ configs/
EXPOSE 8080
CMD ["./server"]
```

构建和运行：

```bash
# 构建镜像
docker build -t rei0721:latest .

# 运行容器
docker run -p 8080:8080 rei0721:latest
```

### Q: 如何使用 Docker Compose？

A: 创建 `docker-compose.yml`：

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DB_DRIVER: postgres
      DB_HOST: db
      DB_USER: postgres
      DB_PASSWORD: password
      DB_NAME: rei0721
    depends_on:
      - db
  
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: rei0721
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

运行：

```bash
docker-compose up
```

### Q: 如何处理优雅关闭？

A: 应用会自动处理 SIGINT 和 SIGTERM 信号：

```bash
# 启动应用
go run ./cmd/server/main.go

# 按 Ctrl+C 或发送信号
kill -SIGTERM <pid>

# 应用会优雅关闭
```

## 错误处理

### Q: 如何处理 "port already in use" 错误？

A: 更改端口或杀死占用端口的进程：

```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <pid>

# 或使用不同的端口
SERVER_PORT=8081 go run ./cmd/server/main.go
```

### Q: 如何处理 "connection refused" 错误？

A: 检查数据库连接：

```bash
# 验证数据库是否运行
psql -h localhost -U postgres -d rei0721

# 检查配置
cat configs/config.yaml

# 查看日志
LOG_LEVEL=debug go run ./cmd/server/main.go
```

### Q: 如何处理 "permission denied" 错误？

A: 检查文件权限：

```bash
# 查看权限
ls -la configs/config.yaml

# 修改权限
chmod 644 configs/config.yaml
```

### Q: 如何处理 "context deadline exceeded" 错误？

A: 增加超时时间：

```yaml
server:
  readTimeout: 30
  writeTimeout: 30
```

## 最佳实践

### Q: 如何组织代码？

A: 遵循分层架构：

```
cmd/          - 入口
internal/     - 业务代码
  ├── app/    - IoC 容器
  ├── handler/- HTTP 处理
  ├── service/- 业务逻辑
  └── repository/ - 数据访问
pkg/          - 工具库
types/        - 类型定义
```

### Q: 如何处理错误？

A: 使用统一的错误类型：

```go
return nil, &errors.BizError{
    Code:    errors.ErrUserNotFound,
    Message: "user not found",
}
```

### Q: 如何进行异步处理？

A: 使用 Scheduler：

```go
scheduler.Submit(ctx, func(ctx context.Context) {
    // 异步任务
})
```

### Q: 如何进行日志记录？

A: 使用 Logger 接口：

```go
logger.Info("user registered",
    "userId", user.ID,
    "email", user.Email,
)
```

### Q: 如何进行配置管理？

A: 使用 ConfigManager：

```go
cfg := configManager.Get()
port := cfg.Server.Port
```

## 获取帮助

如果您的问题未在此列出，请：

1. 查看 [项目文档](./README.md)
2. 查看 [架构设计](./architecture.md)
3. 查看 [开发规范](./protocol.md)
4. 查看应用日志
5. 提交 Issue 或 Pull Request

---

**最后更新**: 2025-12-30
