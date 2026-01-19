# 常见问题

本文档收集了 Go Scaffold 项目使用过程中的常见问题和解决方案。

## 🚀 安装和启动

### Q: 项目启动失败，提示 "Config file not found"

**A:** 这通常是因为配置文件不存在导致的。

**解决方案：**

```bash
# 复制示例配置文件
cp configs/config.example.yaml configs/config.yaml
cp .env.example .env

# 或者指定配置文件路径
go run ./cmd/server server --config=configs/config.yaml
```

### Q: 数据库连接失败，提示 "connection refused"

**A:** 数据库服务未启动或连接参数错误。

**解决方案：**

```bash
# 检查数据库服务状态
sudo systemctl status mysql
# 或
sudo systemctl status postgresql

# 启动数据库服务
sudo systemctl start mysql

# 检查连接参数
mysql -h localhost -P 3306 -u root -p

# 修改配置文件中的数据库连接信息
vim configs/config.yaml
```

### Q: Redis 连接失败，但我不想使用 Redis

**A:** 可以在配置文件中禁用 Redis。

**解决方案：**

```yaml
# configs/config.yaml
cache:
  enabled: false # 禁用 Redis 缓存
```

### Q: 端口 8080 被占用

**A:** 修改配置文件中的端口号。

**解决方案：**

```yaml
# configs/config.yaml
server:
  port: 9000 # 修改为其他端口
```

或者找到占用端口的进程：

```bash
# 查找占用端口的进程
lsof -i :8080

# 杀死进程
kill -9 <PID>
```

## 🗄️ 数据库相关

### Q: 数据库初始化失败

**A:** 检查数据库权限和配置。

**解决方案：**

```bash
# 确保数据库存在
mysql -u root -p -e "CREATE DATABASE scaffold CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 重新初始化数据库
go run ./cmd/server initdb --force

# 检查数据库配置
go run ./cmd/server tests --test=database
```

### Q: 如何切换数据库类型？

**A:** 修改配置文件中的数据库驱动。

**解决方案：**

```yaml
# 使用 MySQL
database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "scaffold"

# 使用 PostgreSQL
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  database: "scaffold"

# 使用 SQLite
database:
  driver: "sqlite"
  database: "scaffold.db"
```

### Q: 数据库迁移如何处理？

**A:** 项目使用 GORM 的自动迁移功能。

**解决方案：**

```bash
# 重新运行数据库初始化
go run ./cmd/server initdb

# 或者在代码中手动迁移
db.AutoMigrate(&models.User{}, &models.Role{})
```

## 🔐 认证和权限

### Q: JWT Token 过期时间如何修改？

**A:** 在配置文件中修改 JWT 相关设置。

**解决方案：**

```yaml
# configs/config.yaml
jwt:
  expires_in: "72h" # 访问令牌 3 天过期
  refresh_expires_in: "720h" # 刷新令牌 30 天过期
```

### Q: 如何添加新的用户角色？

**A:** 通过 RBAC 系统添加角色和权限。

**解决方案：**

```go
// 在初始化时添加角色
rbac.AddRole("editor")
rbac.AddPermissionForRole("editor", "articles", "read")
rbac.AddPermissionForRole("editor", "articles", "write")

// 为用户分配角色
rbac.AddRoleForUser("user123", "editor")
```

### Q: 忘记管理员密码怎么办？

**A:** 通过数据库直接重置密码。

**解决方案：**

```bash
# 重新初始化数据库（会重置所有数据）
go run ./cmd/server initdb --force

# 或者直接修改数据库中的密码
# 首先生成新密码的哈希值
go run -c "
package main
import (
    \"fmt\"
    \"golang.org/x/crypto/bcrypt\"
)
func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte(\"newpassword\"), 12)
    fmt.Println(string(hash))
}
"

# 然后更新数据库
mysql -u root -p scaffold -e "UPDATE users SET password='$2a$12$...' WHERE username='admin';"
```

## 🔧 配置相关

### Q: 环境变量不生效

**A:** 检查环境变量名称和格式。

**解决方案：**

```bash
# 确保使用正确的前缀
export REI_DATABASE_HOST=localhost
export REI_DATABASE_PORT=3306

# 检查环境变量是否设置
env | grep REI_

# 或者使用 .env 文件
echo "DB_HOST=localhost" >> .env
echo "DB_PORT=3306" >> .env
```

### Q: 配置热重载不工作

**A:** 确保启用了热重载功能。

**解决方案：**

```yaml
# configs/config.yaml
app:
  hot_reload: true
  reload_interval: "10s"
```

### Q: 如何在不同环境使用不同配置？

**A:** 使用环境特定的配置文件。

**解决方案：**

```bash
# 创建环境特定配置
cp configs/config.yaml configs/config.production.yaml

# 使用特定配置启动
go run ./cmd/server server --config=configs/config.production.yaml

# 或使用环境变量
export REI_CONFIG_PATH=configs/config.production.yaml
go run ./cmd/server server
```

## 📝 日志相关

### Q: 日志文件太大，如何处理？

**A:** 配置日志轮转。

**解决方案：**

```yaml
# configs/config.yaml
logger:
  output: "file"
  file:
    path: "logs/app.log"
    max_size: 100 # 100MB
    max_backups: 5 # 保留 5 个备份
    max_age: 30 # 保留 30 天
    compress: true # 压缩备份文件
```

### Q: 如何调整日志级别？

**A:** 修改配置文件或使用环境变量。

**解决方案：**

```yaml
# configs/config.yaml
logger:
  level: "debug" # debug, info, warn, error
```

或者：

```bash
export REI_LOGGER_LEVEL=debug
```

### Q: 生产环境日志格式建议

**A:** 使用 JSON 格式便于日志分析。

**解决方案：**

```yaml
# configs/config.yaml
logger:
  level: "info"
  format: "json"
  output: "file"
  file:
    path: "logs/app.log"
    max_size: 100
    max_backups: 10
    max_age: 30
    compress: true
```

## 🚀 性能相关

### Q: 应用启动很慢

**A:** 检查数据库连接和依赖初始化。

**解决方案：**

```bash
# 启用调试模式查看启动过程
export REI_LOGGER_LEVEL=debug
go run ./cmd/server server

# 检查数据库连接时间
go run ./cmd/server tests --test=database

# 优化数据库连接池配置
```

```yaml
database:
  max_open_conns: 25 # 减少连接数
  max_idle_conns: 5 # 减少空闲连接
  conn_max_lifetime: "5m"
```

### Q: 内存使用过高

**A:** 检查协程池和连接池配置。

**解决方案：**

```yaml
# 调整协程池大小
executor:
  pool_size: 50 # 减少协程池大小
  max_blocking_tasks: 100

# 调整数据库连接池
database:
  max_open_conns: 10
  max_idle_conns: 5

# 调整 Redis 连接池
cache:
  pool_size: 5
  min_idle_conns: 2
```

### Q: API 响应慢

**A:** 启用数据库查询日志和性能分析。

**解决方案：**

```yaml
# 启用慢查询日志
database:
  log_level: "info"
  slow_threshold: "100ms" # 记录超过 100ms 的查询
```

```bash
# 使用 pprof 分析性能
go tool pprof http://localhost:8080/debug/pprof/profile
```

## 🐳 Docker 相关

### Q: Docker 构建失败

**A:** 检查 Dockerfile 和网络连接。

**解决方案：**

```bash
# 使用国内镜像加速
docker build --build-arg GOPROXY=https://goproxy.cn,direct -t go-scaffold .

# 检查 Docker 版本
docker --version

# 清理 Docker 缓存
docker system prune -a
```

### Q: Docker Compose 启动失败

**A:** 检查端口冲突和依赖关系。

**解决方案：**

```bash
# 检查端口占用
netstat -tlnp | grep 8080

# 修改 docker-compose.yml 中的端口
ports:
  - "9000:8080"  # 修改外部端口

# 查看详细错误信息
docker-compose up --no-deps app
```

### Q: 容器内无法连接数据库

**A:** 检查网络配置和主机名。

**解决方案：**

```yaml
# docker-compose.yml
services:
  app:
    environment:
      - REI_DATABASE_HOST=mysql # 使用服务名作为主机名
    depends_on:
      - mysql

  mysql:
    # 确保 MySQL 服务配置正确
```

## 🧪 测试相关

### Q: 测试运行失败

**A:** 检查测试环境配置。

**解决方案：**

```bash
# 使用测试配置
export REI_APP_MODE=test
go test ./...

# 运行特定测试
go test -v ./internal/service/auth

# 跳过集成测试
go test -short ./...
```

### Q: 测试数据库如何配置？

**A:** 使用独立的测试数据库。

**解决方案：**

```yaml
# configs/config.test.yaml
database:
  database: "scaffold_test" # 使用测试数据库
```

```go
// 在测试中使用内存数据库
func setupTestDB() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    return db
}
```

## 🔄 部署相关

### Q: 生产环境部署注意事项

**A:** 检查以下配置项：

**解决方案：**

```yaml
# configs/config.production.yaml
app:
  mode: "production"
  debug: false

jwt:
  secret: "change-this-in-production" # 必须修改

logger:
  level: "info"
  format: "json"
  output: "file"

server:
  read_timeout: "30s"
  write_timeout: "30s"
```

### Q: 如何进行健康检查？

**A:** 使用内置的健康检查接口。

**解决方案：**

```bash
# 基础健康检查
curl http://localhost:8080/health

# 详细健康检查
curl http://localhost:8080/health/detailed

# 在 Docker 中配置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
```

### Q: 如何优雅关闭应用？

**A:** 应用已内置优雅关闭机制。

**解决方案：**

```bash
# 发送 SIGTERM 信号
kill -TERM <PID>

# 或使用 SIGINT
kill -INT <PID>

# 应用会等待现有请求完成后关闭
```

## 🔍 调试相关

### Q: 如何启用调试模式？

**A:** 修改配置文件或环境变量。

**解决方案：**

```yaml
# configs/config.yaml
app:
  debug: true

logger:
  level: "debug"
  format: "console" # 更易读的格式
```

### Q: 如何查看详细的错误信息？

**A:** 启用调试日志和错误堆栈。

**解决方案：**

```bash
# 启用调试模式
export REI_APP_DEBUG=true
export REI_LOGGER_LEVEL=debug

# 查看应用日志
tail -f logs/app.log

# 或实时查看控制台输出
go run ./cmd/server server
```

## 📞 获取更多帮助

如果以上解决方案无法解决您的问题，请：

1. **查看项目文档**：[docs/README.md](./README.md)
2. **搜索已有问题**：[GitHub Issues](https://github.com/rei0721/go-scaffold/issues)
3. **创建新问题**：详细描述问题和环境信息
4. **参与讨论**：[GitHub Discussions](https://github.com/rei0721/go-scaffold/discussions)

### 问题报告模板

创建问题时，请提供以下信息：

```markdown
## 环境信息

- OS: [例如 macOS 12.0]
- Go 版本: [例如 1.24.6]
- 项目版本: [例如 v0.1.2]

## 问题描述

[详细描述遇到的问题]

## 复现步骤

1. [步骤1]
2. [步骤2]
3. [步骤3]

## 预期行为

[描述期望的行为]

## 实际行为

[描述实际发生的行为]

## 错误日志
```

[粘贴相关的错误日志]

````

## 配置文件
```yaml
[粘贴相关的配置内容]
````

```

---

**希望这些解答能帮助您顺利使用 Go Scaffold！** 🚀
```
