# 快速开始

本指南将帮助您在 5 分钟内快速体验 Go Scaffold 的核心功能。

## 🚀 5分钟快速体验

### 1. 克隆并启动项目

```bash
# 克隆项目
git clone https://github.com/rei0721/go-scaffold.git
cd go-scaffold

# 安装依赖
go mod download

# 复制配置文件
cp configs/config.example.yaml configs/config.yaml
cp .env.example .env

# 初始化数据库（使用 SQLite，无需额外配置）
go run cmd/server/main.go initdb

# 启动服务
go run cmd/server/main.go server
```

### 2. 验证服务启动

```bash
# 检查健康状态
curl http://localhost:8080/health
```

预期响应：
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "0.1.2"
}
```

🎉 **恭喜！** 您的 Go Scaffold 服务已成功启动！

## 📚 核心功能演示

### 1. 用户认证系统

#### 注册新用户

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

预期响应：
```json
{
  "code": 200,
  "message": "User registered successfully",
  "data": {
    "user_id": "1234567890",
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

预期响应：
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "user": {
      "user_id": "1234567890",
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

#### 访问受保护的接口

```bash
# 使用获得的 access_token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 2. 国际化支持

#### 中文响应

```bash
curl -X GET http://localhost:8080/api/v1/health \
  -H "Accept-Language: zh-CN"
```

预期响应：
```json
{
  "status": "正常",
  "message": "服务运行正常",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### 英文响应

```bash
curl -X GET http://localhost:8080/api/v1/health \
  -H "Accept-Language: en-US"
```

预期响应：
```json
{
  "status": "ok",
  "message": "Service is running normally",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 3. 错误处理演示

#### 参数验证错误

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "email": "invalid-email",
    "password": "123"
  }'
```

预期响应：
```json
{
  "code": 400,
  "message": "Validation failed",
  "errors": [
    {
      "field": "username",
      "message": "Username is required"
    },
    {
      "field": "email",
      "message": "Invalid email format"
    },
    {
      "field": "password",
      "message": "Password must be at least 6 characters"
    }
  ]
}
```

#### 认证错误

```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer invalid-token"
```

预期响应：
```json
{
  "code": 401,
  "message": "Invalid or expired token",
  "error": "UNAUTHORIZED"
}
```

## 🔧 配置自定义

### 1. 切换数据库

编辑 `configs/config.yaml`：

```yaml
# 使用 MySQL
database:
  enabled: true
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "scaffold"

# 或使用 PostgreSQL
database:
  enabled: true
  driver: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  database: "scaffold"
```

### 2. 启用 Redis 缓存

```yaml
cache:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  database: 0
```

### 3. 修改服务端口

```yaml
server:
  host: "0.0.0.0"
  port: 9000  # 修改为其他端口
```

重启服务后，访问 `http://localhost:9000`

## 🧪 开发模式

### 1. 启用热重载

```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 启动热重载
air
```

现在修改代码后，服务会自动重启。

### 2. 启用调试模式

编辑 `configs/config.yaml`：

```yaml
app:
  mode: "development"
  debug: true

logger:
  level: "debug"
  format: "console"  # 更易读的日志格式
```

### 3. 查看详细日志

```bash
# 启动服务并查看日志
go run cmd/server/main.go server | jq
```

## 📊 监控和指标

### 1. 健康检查

```bash
# 基础健康检查
curl http://localhost:8080/health

# 详细健康检查
curl http://localhost:8080/health/detailed
```

### 2. 应用指标

```bash
# 应用统计信息
curl http://localhost:8080/metrics

# 数据库连接状态
curl http://localhost:8080/metrics/database

# 缓存状态
curl http://localhost:8080/metrics/cache
```

## 🧪 测试功能

### 1. 运行单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./pkg/logger -v

# 查看测试覆盖率
go test -cover ./...
```

### 2. API 测试

```bash
# 运行 API 测试
go test ./internal/handler -v

# 运行集成测试
go test -tags=integration ./...
```

## 🐳 Docker 快速启动

### 1. 使用 Docker Compose

```bash
# 启动完整环境（包括 MySQL 和 Redis）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看应用日志
docker-compose logs -f app
```

### 2. 仅启动应用

```bash
# 构建镜像
docker build -t go-scaffold .

# 运行容器
docker run -p 8080:8080 go-scaffold
```

## 📝 API 文档

### 1. 查看 API 文档

访问 `http://localhost:8080/docs` 查看自动生成的 API 文档。

### 2. 导出 OpenAPI 规范

```bash
curl http://localhost:8080/api/openapi.json > api-spec.json
```

## 🔍 故障排除

### 常见问题快速解决

#### 服务启动失败

```bash
# 检查端口占用
lsof -i :8080

# 检查配置文件
go run cmd/server/main.go server --config=configs/config.yaml --dry-run
```

#### 数据库连接失败

```bash
# 测试数据库连接
go run cmd/server/main.go tests --test=database

# 重新初始化数据库
go run cmd/server/main.go initdb --force
```

#### 依赖问题

```bash
# 清理并重新下载依赖
go clean -modcache
go mod download
go mod tidy
```

## 🎯 下一步

现在您已经成功运行了 Go Scaffold，可以继续探索：

1. **[项目结构](../development/project-structure.md)** - 了解代码组织
2. **[API 开发](../development/api-development.md)** - 开发新的 API 接口
3. **[数据库操作](../development/database.md)** - 数据库设计和操作
4. **[配置管理](../development/configuration.md)** - 深入了解配置系统
5. **[部署指南](../deployment/deployment.md)** - 部署到生产环境

## 💡 实用技巧

### 1. 使用 Makefile 简化操作

```bash
# 查看所有可用命令
make help

# 常用命令
make dev      # 开发模式启动
make build    # 构建项目
make test     # 运行测试
make lint     # 代码检查
make clean    # 清理文件
```

### 2. 环境变量覆盖配置

```bash
# 使用环境变量覆盖配置
export DB_HOST=192.168.1.100
export REDIS_HOST=192.168.1.101
go run cmd/server/main.go server
```

### 3. 生产环境配置

```bash
# 使用生产配置启动
go run cmd/server/main.go server --config=configs/production.yaml
```

---

**恭喜！** 您已经完成了 Go Scaffold 的快速入门。现在可以开始构建您的应用了！