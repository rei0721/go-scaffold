# 安装指南

本指南将帮助您快速安装和配置 Go Scaffold 项目。

## 📥 获取项目

### 1. 克隆仓库

```bash
# 使用 HTTPS
git clone https://github.com/rei0721/go-scaffold.git

# 或使用 SSH
git clone git@github.com:rei0721/go-scaffold.git

# 进入项目目录
cd go-scaffold
```

### 2. 检查项目结构

```bash
tree -L 2
```

预期输出：

```
.
├── cmd/                 # 应用程序入口
├── configs/            # 配置文件
├── internal/           # 内部包
├── pkg/               # 公共包
├── types/             # 类型定义
├── docs/              # 项目文档
├── go.mod             # Go 模块文件
├── go.sum             # 依赖校验文件
├── Makefile           # 构建脚本
└── README.md          # 项目说明
```

## 📦 安装依赖

### 1. 下载 Go 模块

```bash
# 下载并安装所有依赖
go mod download

# 验证依赖
go mod verify

# 清理未使用的依赖
go mod tidy
```

### 2. 验证依赖安装

```bash
# 查看依赖列表
go list -m all

# 检查是否有安全漏洞
go list -json -deps ./... | nancy sleuth
```

## ⚙️ 配置项目

### 1. 复制配置文件

```bash
# 复制示例配置文件
cp configs/config.example.yaml configs/config.yaml

# 复制环境变量文件
cp .env.example .env
```

### 2. 编辑配置文件

编辑 `configs/config.yaml`：

```yaml
# 应用配置
app:
  name: "go-scaffold"
  version: "0.1.2"
  mode: "development" # development, production, test
  debug: true

# 服务器配置
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

# 数据库配置
database:
  enabled: true
  driver: "mysql" # mysql, postgres, sqlite
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "scaffold"
  charset: "utf8mb4"
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: "1h"

# Redis 缓存配置
cache:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  database: 0
  pool_size: 10
  min_idle_conns: 5

# JWT 配置
jwt:
  secret: "your-secret-key-change-in-production"
  expires_in: "24h"
  refresh_expires_in: "168h" # 7 days

# 日志配置
logger:
  level: "info" # debug, info, warn, error
  format: "json" # json, console
  output: "stdout" # stdout, file
  file_path: "logs/app.log"
  max_size: 100 # MB
  max_backups: 5
  max_age: 30 # days

# 国际化配置
i18n:
  default_language: "en-US"
  languages:
    - "en-US"
    - "zh-CN"

# RBAC 配置
rbac:
  enabled: true
  model_path: "pkg/rbac/model.conf"
  auto_save: true
```

### 3. 配置环境变量

编辑 `.env` 文件：

```bash
# 应用环境
APP_ENV=development
APP_DEBUG=true

# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_USERNAME=root
DB_PASSWORD=password
DB_DATABASE=scaffold

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key-change-in-production

# 日志
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🗄️ 数据库设置

### 1. 创建数据库

#### MySQL

```sql
CREATE DATABASE scaffold CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### PostgreSQL

```sql
CREATE DATABASE scaffold WITH ENCODING 'UTF8';
```

#### SQLite

SQLite 数据库会自动创建，无需手动操作。

### 2. 初始化数据库

```bash
# 使用项目内置命令初始化数据库
go run ./cmd/server initdb

# 或使用 Makefile
make initdb
```

这将执行以下操作：

- 创建必要的数据表
- 插入初始数据
- 设置默认用户和角色

### 3. 验证数据库连接

```bash
# 测试数据库连接
go run ./cmd/server tests

# 或使用 Makefile
make test-db
```

## 🔧 构建项目

### 1. 使用 Makefile

```bash
# 查看可用命令
make help

# 构建项目
make build

# 构建并运行
make run

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint

# 清理构建文件
make clean
```

### 2. 手动构建

```bash
# 构建可执行文件
go build -o bin/server ./cmd/server

# 交叉编译 (Linux)
GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server

# 交叉编译 (Windows)
GOOS=windows GOARCH=amd64 go build -o bin/server.exe ./cmd/server
```

## 🚀 启动应用

### 1. 开发模式

```bash
# 直接运行
go run ./cmd/server server

# 使用 Air 热重载
air

# 使用 Makefile
make dev
```

### 2. 生产模式

```bash
# 构建并运行
make build
./bin/server server

# 或直接运行
go run ./cmd/server server --config=configs/config.yaml
```

### 3. 验证启动

```bash
# 检查健康状态
curl http://localhost:8080/health

# 预期响应
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "0.1.2"
}
```

## 🐳 Docker 部署

### 1. 构建 Docker 镜像

```bash
# 构建镜像
docker build -t go-scaffold:latest .

# 查看镜像
docker images | grep go-scaffold
```

### 2. 使用 Docker Compose

创建 `docker-compose.yml`：

```yaml
version: "3.8"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    volumes:
      - ./configs:/app/configs
      - ./logs:/app/logs

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: scaffold
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  mysql_data:
  redis_data:
```

启动服务：

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down
```

## 🧪 运行测试

### 1. 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/logger

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 2. 集成测试

```bash
# 运行集成测试
go test -tags=integration ./...

# 使用 Makefile
make test-integration
```

### 3. 基准测试

```bash
# 运行基准测试
go test -bench=. ./...

# 运行特定基准测试
go test -bench=BenchmarkLogger ./pkg/logger
```

## 📋 安装验证清单

完成安装后，请验证以下项目：

- [ ] 项目代码已克隆到本地
- [ ] Go 依赖已成功下载
- [ ] 配置文件已正确设置
- [ ] 数据库连接正常
- [ ] Redis 连接正常（如果启用）
- [ ] 项目可以成功构建
- [ ] 应用可以正常启动
- [ ] 健康检查接口返回正常
- [ ] 单元测试通过
- [ ] Docker 镜像构建成功（如果使用）

## 🚨 常见安装问题

### 依赖下载失败

**问题**: `go mod download` 失败
**解决方案**:

```bash
# 设置 Go 代理
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

# 或使用其他代理
go env -w GOPROXY=https://proxy.golang.org,direct
```

### 数据库连接失败

**问题**: 数据库连接被拒绝
**解决方案**:

1. 检查数据库服务是否运行
2. 验证连接参数（主机、端口、用户名、密码）
3. 检查防火墙设置
4. 确认数据库用户权限

### 端口占用

**问题**: 端口 8080 已被占用
**解决方案**:

```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>

# 或修改配置文件中的端口
```

### 权限问题

**问题**: 文件权限不足
**解决方案**:

```bash
# 修改文件权限
chmod +x bin/server

# 修改目录权限
chmod -R 755 logs/
```

## 🔄 更新项目

### 1. 更新代码

```bash
# 拉取最新代码
git pull origin main

# 更新依赖
go mod tidy
go mod download
```

### 2. 数据库迁移

```bash
# 运行数据库迁移
go run ./cmd/server initdb --migrate
```

### 3. 重新构建

```bash
# 清理并重新构建
make clean
make build
```

---

**下一步**: 查看 [快速开始](./quickstart.md) 开始使用项目功能
