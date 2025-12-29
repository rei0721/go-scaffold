# 快速开始指南

本指南将帮助您快速搭建 Rei0721 项目的本地开发环境。

## 前置要求

### 系统要求

- **操作系统**: Linux, macOS, Windows
- **Go 版本**: 1.24.6 或更高
- **内存**: 至少 2GB
- **磁盘**: 至少 500MB 可用空间

### 依赖软件

#### 必需

- **Go 1.24.6+** - [下载](https://golang.org/dl/)
- **Git** - [下载](https://git-scm.com/)
- **数据库** (选择一个):
  - PostgreSQL 12+ - [下载](https://www.postgresql.org/download/)
  - MySQL 5.7+ - [下载](https://dev.mysql.com/downloads/mysql/)
  - SQLite 3 - 通常已预装

#### 可选

- **Redis 6.0+** - [下载](https://redis.io/download)
- **Docker** - [下载](https://www.docker.com/)
- **Docker Compose** - [下载](https://docs.docker.com/compose/)

## 安装步骤

### 1. 克隆项目

```bash
git clone https://github.com/your-org/rei0721.git
cd rei0721
```

### 2. 验证 Go 版本

```bash
go version
# 输出应该是 go version go1.24.6 或更高
```

### 3. 下载依赖

```bash
go mod download
go mod tidy
```

### 4. 配置数据库

#### 使用 PostgreSQL

```bash
# 创建数据库
createdb rei0721

# 创建用户 (可选)
createuser -P rei0721_user
```

#### 使用 MySQL

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE rei0721;"

# 创建用户 (可选)
mysql -u root -p -e "CREATE USER 'rei0721_user'@'localhost' IDENTIFIED BY 'password';"
mysql -u root -p -e "GRANT ALL PRIVILEGES ON rei0721.* TO 'rei0721_user'@'localhost';"
```

#### 使用 SQLite

SQLite 无需额外配置，数据库文件会自动创建。

### 5. 配置环境变量

复制环境变量模板：

```bash
cp configs/.env.example .env
```

编辑 `.env` 文件，根据您的数据库选择配置：

**PostgreSQL 配置**:
```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rei0721
```

**MySQL 配置**:
```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=rei0721
```

**SQLite 配置**:
```bash
DB_DRIVER=sqlite
DB_NAME=./rei0721.db
```

**其他配置**:
```bash
SERVER_PORT=8080
SERVER_MODE=debug
LOG_LEVEL=info
LOG_FORMAT=json
```

### 6. 加载环境变量

```bash
# Linux/macOS
source .env

# Windows (PowerShell)
Get-Content .env | ForEach-Object {
    if ($_ -match '^\s*([^=]+)=(.*)$') {
        [Environment]::SetEnvironmentVariable($matches[1], $matches[2])
    }
}
```

或者在运行命令时直接指定：

```bash
DB_DRIVER=postgres DB_HOST=localhost go run ./cmd/server/main.go
```

## 启动服务

### 开发模式

```bash
go run ./cmd/server/main.go
```

预期输出：
```
{"level":"info","ts":1735475400.123,"msg":"database connected successfully"}
{"level":"info","ts":1735475400.124,"msg":"scheduler initialized","poolSize":10000}
{"level":"info","ts":1735475400.125,"msg":"application initialized successfully"}
{"level":"info","ts":1735475400.126,"msg":"starting HTTP server","addr":":8080"}
```

### 生产模式

```bash
SERVER_MODE=release go run ./cmd/server/main.go
```

### 使用 Docker

```bash
# 构建镜像
docker build -t rei0721:latest .

# 运行容器
docker run -p 8080:8080 \
  -e DB_DRIVER=postgres \
  -e DB_HOST=db \
  -e DB_USER=postgres \
  -e DB_PASSWORD=password \
  -e DB_NAME=rei0721 \
  rei0721:latest
```

## 验证安装

### 检查服务健康状态

```bash
curl http://localhost:8080/health
```

预期响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok"
  },
  "serverTime": 1735475400
}
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/router/... -v

# 生成覆盖率报告
go test ./... -cover
```

### 代码检查

```bash
# 检查代码质量
go vet ./...

# 格式化代码
go fmt ./...

# 运行 linter (需要安装 golangci-lint)
golangci-lint run ./...
```

## 常见问题

### Q: 如何更改服务端口？

A: 修改 `configs/config.yaml` 中的 `server.port`，或设置环境变量：

```bash
SERVER_PORT=8081 go run ./cmd/server/main.go
```

### Q: 如何启用 Redis？

A: 在 `configs/config.yaml` 中配置 Redis：

```yaml
redis:
  enabled: true
  host: localhost
  port: 6379
  password: ""
  db: 0
  poolSize: 10
```

### Q: 如何查看详细日志？

A: 设置日志级别为 debug：

```bash
LOG_LEVEL=debug go run ./cmd/server/main.go
```

### Q: 如何使用 SQLite 进行开发？

A: 配置环境变量：

```bash
DB_DRIVER=sqlite DB_NAME=./rei0721.db go run ./cmd/server/main.go
```

### Q: 如何重置数据库？

A: 删除数据库文件并重新启动：

```bash
# SQLite
rm rei0721.db

# PostgreSQL
dropdb rei0721
createdb rei0721

# MySQL
mysql -u root -p -e "DROP DATABASE rei0721; CREATE DATABASE rei0721;"
```

## 下一步

- 阅读 [架构设计](./architecture.md) 了解项目结构
- 查看 [API 规范](./api.md) 了解可用的 API 端点
- 参考 [开发规范](./protocol.md) 了解代码风格和最佳实践
- 查看 [配置管理](./configuration.md) 了解高级配置选项

## 获取帮助

如遇到问题，请：

1. 检查 [常见问题](./faq.md)
2. 查看日志输出
3. 验证环境变量配置
4. 检查数据库连接

---

**最后更新**: 2025-12-30
