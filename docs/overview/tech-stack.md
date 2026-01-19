# 技术栈

## 🛠️ 核心技术栈

Go Scaffold 基于现代化的 Go 生态系统构建，选择了经过生产验证的成熟技术和框架。

## 🔧 开发语言与运行时

### Go 1.24.6

- **选择理由**: 高性能、强类型、并发友好的现代编程语言
- **特性**:
  - 原生协程支持，高并发处理能力
  - 快速编译，单一可执行文件部署
  - 丰富的标准库和第三方生态
  - 内存安全和垃圾回收机制

## 🌐 Web 框架

### Gin v1.11.0

```go
github.com/gin-gonic/gin v1.11.0
```

- **选择理由**: 高性能的 HTTP Web 框架，简单易用
- **核心特性**:
  - 快速路由，基于 Radix 树实现
  - 中间件支持，可扩展的请求处理链
  - JSON 验证和绑定
  - 错误管理和恢复机制
  - 渲染支持（JSON、XML、HTML等）

**使用示例**:

```go
router := gin.New()
router.Use(gin.Logger(), gin.Recovery())
router.GET("/api/v1/users", userHandler.GetUsers)
```

## 🗄️ 数据存储

### GORM v1.31.1

```go
gorm.io/gorm v1.31.1
```

- **选择理由**: 功能丰富的 Go ORM 库，支持多种数据库
- **核心特性**:
  - 自动迁移和模型定义
  - 关联关系处理（一对一、一对多、多对多）
  - 事务支持和连接池管理
  - 钩子函数和插件系统
  - 软删除和批量操作

**支持的数据库驱动**:

```go
gorm.io/driver/mysql v1.6.0      // MySQL
gorm.io/driver/postgres v1.6.0   // PostgreSQL
gorm.io/driver/sqlite v1.6.0     // SQLite
```

### Redis v9.17.2

```go
github.com/redis/go-redis/v9 v9.17.2
```

- **选择理由**: 高性能的内存数据库，用于缓存和会话存储
- **核心特性**:
  - 连接池管理
  - 集群支持
  - 管道和事务
  - 发布订阅模式
  - 多种数据结构支持

## 🔐 安全认证

### JWT v5.3.0

```go
github.com/golang-jwt/jwt/v5 v5.3.0
```

- **选择理由**: 标准的 JSON Web Token 实现
- **核心特性**:
  - 多种签名算法支持（HMAC、RSA、ECDSA）
  - 声明验证和自定义声明
  - 令牌解析和验证
  - 过期时间管理

### Bcrypt 加密

```go
golang.org/x/crypto v0.46.0
```

- **选择理由**: 安全的密码哈希算法
- **核心特性**:
  - 自适应哈希函数
  - 盐值自动生成
  - 可配置的计算成本
  - 时间攻击防护

### Casbin v3.8.1

```go
github.com/casbin/casbin/v3 v3.8.1
github.com/casbin/gorm-adapter/v3 v3.40.0
```

- **选择理由**: 强大的访问控制库，支持多种访问控制模型
- **核心特性**:
  - RBAC、ABAC、ACL 等多种模型
  - 策略持久化存储
  - 动态策略管理
  - 高性能策略匹配

## 📝 日志系统

### Zap v1.27.1

```go
go.uber.org/zap v1.27.1
```

- **选择理由**: 高性能的结构化日志库
- **核心特性**:
  - 零内存分配的日志记录
  - 结构化日志输出（JSON/Console）
  - 多级别日志支持
  - 采样和缓冲机制
  - 自定义编码器和输出器

### Lumberjack v2.2.1

```go
gopkg.in/natefinch/lumberjack.v2 v2.2.1
```

- **选择理由**: 日志文件轮转和管理
- **核心特性**:
  - 自动日志轮转
  - 文件大小和时间控制
  - 压缩和清理机制
  - 线程安全操作

## ⚙️ 配置管理

### Viper v1.21.0

```go
github.com/spf13/viper v1.21.0
```

- **选择理由**: 功能完整的配置管理库
- **核心特性**:
  - 多种配置格式支持（YAML、JSON、TOML等）
  - 环境变量和命令行参数
  - 配置文件热重载
  - 配置值类型转换
  - 远程配置支持

### GoDotEnv v1.5.1

```go
github.com/joho/godotenv v1.5.1
```

- **选择理由**: 环境变量文件加载
- **核心特性**:
  - .env 文件解析
  - 环境变量覆盖
  - 多文件支持
  - 变量展开

## 🌍 国际化

### Go-i18n v2.6.0

```go
github.com/nicksnyder/go-i18n/v2 v2.6.0
golang.org/x/text v0.32.0
```

- **选择理由**: 完整的国际化解决方案
- **核心特性**:
  - 多语言消息管理
  - 复数形式处理
  - 消息模板和变量替换
  - 语言环境检测
  - 翻译文件管理

## ⚡ 性能优化

### Ants v2.11.4

```go
github.com/panjf2000/ants/v2 v2.11.4
```

- **选择理由**: 高性能的协程池库
- **核心特性**:
  - 协程复用和管理
  - 动态扩缩容
  - 任务队列和调度
  - 内存优化
  - 监控和统计

### Snowflake v0.3.0

```go
github.com/bwmarrin/snowflake v0.3.0
```

- **选择理由**: 分布式唯一ID生成器
- **核心特性**:
  - 高性能ID生成
  - 时间有序性
  - 分布式友好
  - 无依赖实现

## 🔧 开发工具

### Air (热重载)

```bash
go install github.com/cosmtrek/air@latest
```

- **选择理由**: 开发时自动重载工具
- **核心特性**:
  - 文件变化监听
  - 自动编译和重启
  - 可配置的监听规则
  - 彩色日志输出

### FSNotify v1.9.0

```go
github.com/fsnotify/fsnotify v1.9.0
```

- **选择理由**: 跨平台文件系统事件监听
- **核心特性**:
  - 文件和目录监听
  - 跨平台支持
  - 事件过滤和处理
  - 低延迟通知

## 🧪 测试框架

### 内置测试

```go
import "testing"
```

- **选择理由**: Go 原生测试框架，简单高效
- **核心特性**:
  - 单元测试和基准测试
  - 测试覆盖率统计
  - 并行测试支持
  - 子测试和表格驱动测试

### Testify (推荐)

```go
github.com/stretchr/testify
```

- **选择理由**: 丰富的测试断言和Mock工具
- **核心特性**:
  - 丰富的断言函数
  - Mock 对象生成
  - 测试套件支持
  - HTTP 测试工具

## 📦 构建和部署

### Docker

```dockerfile
FROM golang:1.24-alpine AS builder
# 构建阶段
FROM alpine:latest
# 运行阶段
```

- **选择理由**: 容器化部署，环境一致性
- **核心特性**:
  - 多阶段构建优化
  - 轻量级运行镜像
  - 环境隔离
  - 可移植性

### Make

```makefile
.PHONY: build test clean
build:
	go build -o bin/server ./cmd/server
```

- **选择理由**: 简单的构建自动化工具
- **核心特性**:
  - 构建任务管理
  - 依赖关系处理
  - 跨平台支持
  - 简单易用

## 📊 技术选型原则

### 1. 性能优先

- 选择高性能的库和框架
- 避免过度抽象和复杂度
- 关注内存使用和GC压力

### 2. 生产就绪

- 选择成熟稳定的技术
- 有活跃的社区支持
- 良好的文档和示例

### 3. 可维护性

- 代码简洁易懂
- 良好的错误处理
- 完善的日志和监控

### 4. 可扩展性

- 模块化设计
- 接口驱动开发
- 支持水平扩展

## 🔄 版本管理策略

### 依赖版本锁定

- 使用 `go.mod` 锁定依赖版本
- 定期更新安全补丁
- 测试兼容性后升级主版本

### 向后兼容

- 遵循语义化版本规范
- 保持API稳定性
- 渐进式升级策略

---

**下一步**: 查看 [快速开始](../getting-started/quickstart.md) 开始使用项目
