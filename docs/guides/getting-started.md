# 快速开始 (Getting Started)

> **5 分钟快速上手 go-scaffold**  
> 从零开始启动本项目。

---

## 📋 前置要求

### 必须安装

- **Go**: 1.24.6+ ([下载](https://go.dev/dl/))
- **Git**: 用于克隆仓库

### 可选安装（根据使用的数据库）

- **PostgreSQL**: 5432 端口（生产推荐）
- **MySQL**: 3306 端口
- **SQLite**: 无需 安装（开发推荐）
- **Redis**: 6379 端口（可选，缓存功能）

---

## 🚀 快速启动（SQLite 版）

### 1. 克隆项目

```bash
git clone https://github.com/rei0721/rei0721.git
cd rei0721
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 复制配置文件

```bash
# Windows PowerShell
Copy-Item .env.example .env
Copy-Item configs\config.example.yaml configs\config.yaml

# Linux/macOS
cp .env.example .env
cp configs/config.example.yaml configs/config.yaml
```

### 4. 启动服务

```bash
go run cmd/server/main.go dev
```

### 5. 验证

访问: `http://localhost:8080`

如果看到欢迎页面或 404（未配置路由），说明服务启动成功！

---

## ⚙️ 配置说明

### 配置文件位置

- **主配置**: `configs/config.yaml`
- **环境变量**: `.env`（优先级更高）
- **示例文件**: `configs/config.example.yaml`, `.env.example`

### 三层配置优先级

```
环境变量 (.env) > YAML配置 > 默认值
```

### 最小配置（SQLite）

**configs/config.yaml**:

```yaml
server:
  host: localhost
  port: 8080
  mode: debug # debug | release | test

database:
  driver: sqlite
  dbname: ./data/go-scaffold.db # 自动创建

logger:
  level: debug
  format: console
  output: stdout

i18n:
  default: zh-CN
  supported:
    - zh-CN
    - en-US

redis:
  enabled: false # SQLite版不需要Redis

executor:
  enabled: true
  pools:
    - name: http
      size: 100
      expiry: 10
      nonBlocking: true
```

---

## 🗄️ 使用 PostgreSQL/MySQL

### PostgreSQL 配置

**1. 安装 PostgreSQL**（如已安装跳过）

**2. 创建数据库**:

```sql
CREATE DATABASE go_scaffold;
CREATE USER go_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE go_scaffold TO go_user;
```

**3. 修改配置**:

**configs/config.yaml**:

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: go_user
  password: "" # 留空，使用环境变量
  dbname: go_scaffold
  maxOpenConns: 100
  maxIdleConns: 10
```

**.env**:

```bash
DB_PASSWORD=your_password
```

### MySQL 配置

**1. 创建数据库**:

```sql
CREATE DATABASE go_scaffold CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'go_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON go_scaffold.* TO 'go_user'@'localhost';
FLUSH PRIVILEGES;
```

**2. 修改配置**:

**configs/config.yaml**:

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  user: go_user
  password: ""
  dbname: go_scaffold
  maxOpenConns: 100
  maxIdleConns: 10
```

**.env**:

```bash
DB_PASSWORD=your_password
```

---

## 🔧 使用 Redis（可选）

### 1. 安装 Redis（如已安装跳过）

### 2. 启动 Redis

```bash
redis-server
```

### 3. 修改配置

**configs/config.yaml**:

```yaml
redis:
  enabled: true
  host: localhost
  port: 6379
  password: "" # 如果Redis设置了密码
  db: 0
  poolSize: 20
  minIdleConns: 10
```

**.env** (如果有密码):

```bash
REDIS_PASSWORD=your_redis_password
```

---

## 📦 编译与部署

### 本地编译

```bash
# 编译可执行文件
go build -o bin/go-scaffold cmd/server/main.go

# 运行
./bin/go-scaffold dev
```

### 交叉编译

**Linux**:

```bash
GOOS=linux GOARCH=amd64 go build -o bin/go-scaffold-linux cmd/server/main.go
```

**Windows**:

```bash
GOOS=windows GOARCH=amd64 go build -o bin/go-scaffold.exe cmd/server/main.go
```

### Docker 部署（待补充）

```bash
# TODO: 添加 Dockerfile 和 docker-compose.yml
```

---

## 🧪 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/logger

# 查看覆盖率
go test -cover ./...
```

---

## 📂 项目结构说明

```
go-scaffold/
├── cmd/           # 命令行入口
│   └── server/    # 主服务器程序
├── internal/      # 内部业务逻辑
│   ├── app/       # DI容器
│   ├── config/    # 配置定义
│   ├── handler/   # HTTP处理器
│   ├── service/   # 业务逻辑
│   └── repository/# 数据访问
├── pkg/           # 可复用库
│   ├── cache/     # 缓存
│   ├── cli/       # CLI框架
│   ├── database/  # 数据库
│   ├── executor/  # 协程池
│   ├── httpserver/# HTTP服务器
│   ├── i18n/      # 国际化
│   ├── logger/    # 日志
│   ├── sqlgen/    # SQL生成器
│   └── utils/     # 工具箱
├── types/         # 共享类型
├── configs/       # 配置文件
└── docs/          # 项目文档
```

详见: [架构地图](/docs/architecture/system_map.md)

---

## 🛠️ 常用命令

### CLI 工具

```bash
# 查看帮助
go run ./cmd/server --help

# 启动开发服务器
go run ./cmd/server dev

# SQL 生成器（示例）
go run ./cmd/server sqlgen --help
```

---

## ❓ 常见问题

### Q1: 端口被占用

**错误**: `bind: address already in use`

**解决**:

- 修改 `configs/config.yaml` 中的 `server.port`
- 或设置环境变量: `SERVER_PORT=9090 go run cmd/server/main.go dev`

### Q2: 数据库连接失败

**错误**: `failed to connect to database`

**排查**:

1. 检查数据库服务是否启动
2. 确认 `.env` 中的密码是否正确
3. 确认数据库名称、用户名、端口是否正确

### Q3: SQLite 文件权限错误

**错误**: `unable to open database file`

**解决**:

```bash
mkdir -p ./data
chmod 755 ./data
```

### Q4: 找不到配置文件

**错误**: `config file not found`

**解决**:

- 确保在项目根目录运行
- 确保 `configs/config.yaml` 存在

---

## 📚 进一步学习

- [系统架构地图](/docs/architecture/system_map.md) - 理解整体架构
- [变量命名索引](/docs/architecture/variable_index.md) - 查找常量定义
- [编码规则](/docs/memories/rules.md) - 贡献代码前必读
- [AI 协作提示词](/docs/ai_prompt.md) - AI 协作规范

---

## 🆘 获取帮助

- **文档**: 查看 `docs/` 目录
- **Issues**: GitHub Issues（如有）
- **联系作者**: rei0721 @ GitHub

---

> **提示**: 开发环境推荐使用 SQLite + 控制台日志，生产环境推荐 PostgreSQL + 文件日志（JSON 格式）。
