# 部署指南

```
项目: Rei0721 | 版本: v1.0 | 更新: 2025-12-29
```

## 环境要求

| 组件 | 版本 | 必需 |
|------|------|------|
| Go | 1.21+ | ✅ |
| PostgreSQL | 12+ | 可选 |
| MySQL | 8.0+ | 可选 |
| SQLite | 3.35+ | 可选 |
| Redis | 6.0+ | 可选 |
| Docker | 20.10+ | 可选 |

> 数据库三选一，生产环境推荐 PostgreSQL

---

## 配置说明

### 环境变量 (.env)

```bash
cp configs/.env.example configs/.env
```

```bash
# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=rei0721

# Redis (可选)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 应用
APP_ENV=production
APP_SECRET_KEY=your_secret_key
```

### 主配置 (config.yaml)

```yaml
server:
  port: 8080
  mode: release  # debug | release

database:
  driver: postgres  # postgres | mysql | sqlite
  host: ${DB_HOST}
  port: ${DB_PORT}
  maxOpenConns: 100
  maxIdleConns: 10

redis:
  enabled: true
  host: ${REDIS_HOST}
  poolSize: 10

logger:
  level: info  # debug | info | warn | error
  format: json
  output: file

i18n:
  default: zh-CN
  supported: [zh-CN, en-US]
```

---

## 本地部署

```bash
# 1. 克隆
git clone <repository-url> && cd rei0721

# 2. 依赖
go mod download

# 3. 配置
cp configs/.env.example configs/.env
# 编辑 configs/.env 和 configs/config.yaml

# 4. 数据库
# PostgreSQL: createdb rei0721
# MySQL: CREATE DATABASE rei0721 CHARACTER SET utf8mb4;
# SQLite: 自动创建

# 5. 运行
go run cmd/server/main.go

# 6. 验证
curl http://localhost:8080/health
```

---

## Docker 部署

### docker-compose.yml

```yaml
version: "3.8"

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=postgres
      - REDIS_HOST=redis
    volumes:
      - ./configs:/app/configs
      - ./logs:/app/logs
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: rei0721
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres123
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rei0721 cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/rei0721 .
COPY --from=builder /build/configs ./configs
ENV TZ=Asia/Shanghai
EXPOSE 8080
CMD ["./rei0721"]
```

### 启动

```bash
docker-compose up -d
docker-compose logs -f app
docker-compose down
```

---

## 生产环境部署

### 部署前检查

- [ ] 配置文件正确
- [ ] 数据库连接测试
- [ ] Redis 连接测试 (如使用)
- [ ] 日志路径有写入权限
- [ ] 防火墙规则配置
- [ ] SSL/TLS 证书准备

### Systemd 服务

```bash
# 编译
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w' -o /opt/rei0721/rei0721 cmd/server/main.go
```

`/etc/systemd/system/rei0721.service`:

```ini
[Unit]
Description=Rei0721 API Server
After=network.target postgresql.service

[Service]
Type=simple
User=rei0721
WorkingDirectory=/opt/rei0721
ExecStart=/opt/rei0721/rei0721
Restart=on-failure
RestartSec=5s
Environment="APP_ENV=production"

[Install]
WantedBy=multi-user.target
```

```bash
sudo useradd -r -s /bin/false rei0721
sudo chown -R rei0721:rei0721 /opt/rei0721
sudo systemctl daemon-reload
sudo systemctl enable --now rei0721
sudo journalctl -u rei0721 -f
```

### Nginx 反向代理

```nginx
upstream rei0721_backend {
    server 127.0.0.1:8080;
}

server {
    listen 443 ssl http2;
    server_name api.example.com;

    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;

    location / {
        proxy_pass http://rei0721_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 常见问题

### 数据库连接失败

```bash
# 检查服务状态
sudo systemctl status postgresql

# 测试连接
psql -h localhost -U postgres -d rei0721

# 防火墙
sudo ufw allow 5432/tcp
```

### Redis 连接失败

```bash
sudo systemctl status redis
redis-cli ping  # 应返回 PONG

# 不使用 Redis 时禁用
# redis:
#   enabled: false
```

### 端口被占用

```bash
sudo lsof -i :8080
# 终止进程或更改配置端口
```

### 日志权限错误

```bash
mkdir -p logs && chmod 755 logs
sudo chown rei0721:rei0721 /opt/rei0721/logs
```

### 配置热重载不生效

```bash
# 验证 YAML 语法
yamllint configs/config.yaml

# 查看日志
tail -f logs/app.log | grep "config"
```

### Docker 容器无法启动

```bash
docker logs rei0721-app
docker run -it --rm rei0721:latest sh
```

---

## 性能优化

| 组件 | 优化建议 |
|------|----------|
| 数据库 | 添加索引、配置连接池、定期 VACUUM |
| Redis | 合理过期时间、使用连接池、监控内存 |
| 应用 | 启用 Gzip、使用缓存、配置协程池大小 |

---

## 监控与维护

**推荐工具:**
- Prometheus + Grafana (指标监控)
- ELK Stack (日志分析)
- Sentry (错误追踪)

**日常维护:**
- 定期备份数据库
- 监控日志文件大小
- 检查磁盘空间
- 更新依赖包

---

[← README](./README.md) | [api.md](./api.md)
