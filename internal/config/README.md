# Config - 配置管理系统

## 概述

配置管理系统提供统一的配置加载、热更新和环境变量支持。支持多种配置来源，并按优先级自动合并。

## 配置优先级

**优先级规则**（从高到低）：

```
1. 系统环境变量 (最高)
2. .env 文件
3. config.yaml (最低)
```

### 优先级示例

假设有以下配置：

**config.yaml**:

```yaml
redis:
  host: localhost
  port: 6379
```

**.env 文件**:

```env
REDIS_HOST=redis-dev
```

**系统环境变量**:

```bash
export REDIS_PORT=6380
```

**最终生效的配置**:

- `redis.host` = `redis-dev` (来自 .env)
- `redis.port` = `6380` (来自系统环境变量)

## 使用方式

### 1. 本地开发

创建 `.env` 文件：

```bash
# 复制示例文件
cp .env.example .env

# 编辑 .env 文件
vim .env
```

**.env 文件示例**:

```env
# 数据库配置
DB_PASSWORD=local_password
DB_HOST=localhost

# Redis 配置
REDIS_PASSWORD=local_redis_pass
REDIS_ENABLED=true
```

### 2. 生产环境

**不使用 .env 文件**，直接设置系统环境变量：

```bash
# Docker
docker run -e DB_PASSWORD=prod_pass -e REDIS_PASSWORD=redis_pass myapp

# Kubernetes
kubectl create secret generic app-secrets \
  --from-literal=DB_PASSWORD=prod_pass \
  --from-literal=REDIS_PASSWORD=redis_pass

# 直接导出
export DB_PASSWORD=prod_pass
export REDIS_PASSWORD=redis_pass
```

### 3. 配置文件中使用环境变量

`config.yaml` 中可以使用 `${VAR:default}` 语法：

```yaml
server:
  port: ${SERVER_PORT:8080} # 默认 8080
  mode: ${SERVER_MODE:debug}

database:
  host: ${DB_HOST:localhost}
  password: ${DB_PASSWORD:} # 必须从环境变量读取
```

## 支持的环境变量

### 数据库配置

| 环境变量            | 说明         | 示例        |
| ------------------- | ------------ | ----------- |
| `DB_DRIVER`         | 数据库驱动   | `postgres`  |
| `DB_HOST`           | 主机地址     | `localhost` |
| `DB_PORT`           | 端口         | `5432`      |
| `DB_USER`           | 用户名       | `postgres`  |
| `DB_PASSWORD`       | 密码         | `secret`    |
| `DB_NAME`           | 数据库名     | `myapp`     |
| `DB_MAX_OPEN_CONNS` | 最大连接数   | `100`       |
| `DB_MAX_IDLE_CONNS` | 最大空闲连接 | `10`        |

### Redis 配置

| 环境变量               | 说明         | 示例        |
| ---------------------- | ------------ | ----------- |
| `REDIS_ENABLED`        | 是否启用     | `true`      |
| `REDIS_HOST`           | 主机地址     | `localhost` |
| `REDIS_PORT`           | 端口         | `6379`      |
| `REDIS_PASSWORD`       | 密码         | `secret`    |
| `REDIS_DB`             | 数据库索引   | `0`         |
| `REDIS_POOL_SIZE`      | 连接池大小   | `20`        |
| `REDIS_MIN_IDLE_CONNS` | 最小空闲连接 | `10`        |
| `REDIS_MAX_RETRIES`    | 最大重试次数 | `3`         |
| `REDIS_DIAL_TIMEOUT`   | 连接超时(秒) | `5`         |
| `REDIS_READ_TIMEOUT`   | 读取超时(秒) | `3`         |
| `REDIS_WRITE_TIMEOUT`  | 写入超时(秒) | `3`         |

### 服务器配置

| 环境变量               | 说明         | 示例      |
| ---------------------- | ------------ | --------- |
| `SERVER_PORT`          | HTTP 端口    | `8080`    |
| `SERVER_MODE`          | 运行模式     | `release` |
| `SERVER_READ_TIMEOUT`  | 读取超时(秒) | `30`      |
| `SERVER_WRITE_TIMEOUT` | 写入超时(秒) | `30`      |

### 日志配置

| 环境变量     | 说明     | 示例     |
| ------------ | -------- | -------- |
| `LOG_LEVEL`  | 日志级别 | `info`   |
| `LOG_FORMAT` | 日志格式 | `json`   |
| `LOG_OUTPUT` | 日志输出 | `stdout` |

### 国际化配置

| 环境变量         | 说明       | 示例          |
| ---------------- | ---------- | ------------- |
| `I18N_DEFAULT`   | 默认语言   | `zh-CN`       |
| `I18N_SUPPORTED` | 支持的语言 | `zh-CN,en-US` |

## 代码示例

### 加载配置

```go
package main

import (
    "github.com/rei0721/go-scaffold/internal/config"
)

func main() {
    // 创建配置管理器
    manager := config.NewManager()

    // 加载配置
    // 会自动:
    // 1. 加载 .env 文件(如果存在)
    // 2. 读取 config.yaml
    // 3. 用环境变量覆盖配置
    err := manager.Load("configs/config.yaml")
    if err != nil {
        panic(err)
    }

    // 获取配置
    cfg := manager.Get()
    fmt.Println("Port:", cfg.Server.Port)
}
```

### 监听配置变化

```go
// 注册配置变更钩子
manager.RegisterHook(func(old, new *config.Config) {
    log.Info("config changed",
        "oldPort", old.Server.Port,
        "newPort", new.Server.Port)

    // 重新初始化需要更新的组件
    reloadComponents(new)
})

// 开始监听配置文件变化
err = manager.Watch()
```

## 最佳实践

### 1. 敏感信息使用环境变量

❌ **不要**在 config.yaml 中硬编码敏感信息：

```yaml
database:
  password: "my-secret-password" # 不安全！
```

✅ **应该**使用环境变量：

```yaml
database:
  password: ${DB_PASSWORD:} # 从环境变量读取
```

### 2. .env 文件不提交到 Git

确保 `.env` 在 `.gitignore` 中：

```gitignore
.env
.env.local
.env.*.local
```

只提交 `.env.example` 作为模板。

### 3. 生产环境使用系统环境变量

- Kubernetes: 使用 Secrets 和 ConfigMaps
- Docker: 使用 `--env-file` 或 `-e` 参数
- 云平台: 使用各平台的环境变量管理

### 4. 提供合理的默认值

```yaml
server:
  port: ${SERVER_PORT:8080} # 提供默认值
  mode: ${SERVER_MODE:release}
```

### 5. 环境变量命名规范

- 全大写
- 单词间使用下划线
- 使用模块前缀
- 示例: `DB_PASSWORD`, `REDIS_HOST`

## Docker 集成

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/configs /configs
# 不复制 .env 文件!
CMD ["./server"]
```

### docker-compose.yml

```yaml
version: "3.8"
services:
  app:
    build: .
    env_file:
      - .env.production # 使用生产环境配置
    environment:
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_PASSWORD=${REDIS_PASSWORD}
    ports:
      - "8080:8080"
```

## Kubernetes 集成

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  db-password: <base64-encoded-password>
  redis-password: <base64-encoded-password>
```

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
        - name: app
          image: myapp:latest
          env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: app-secrets
                  key: db-password
            - name: REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: app-secrets
                  key: redis-password
```

## 故障排查

### 配置未生效

1. 检查优先级顺序
2. 查看日志确认环境变量是否被加载
3. 确认环境变量名称是否正确（区分大小写）

### .env 文件未加载

1. 确保 .env 文件在项目根目录
2. 检查文件权限
3. 查看是否有语法错误

### 环境变量格式错误

```bash
# ✅ 正确
export REDIS_ENABLED=true
export SERVER_PORT=8080

# ❌ 错误
export REDIS_ENABLED = true  # 不要有空格
export SERVER_PORT="8080"    # 数字不需要引号(但也可以)
```

## 参考文档

- [12-Factor App - Config](https://12factor.net/config)
- [godotenv](https://github.com/joho/godotenv)
- [viper](https://github.com/spf13/viper)
