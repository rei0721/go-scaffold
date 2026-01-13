# 快速开始

> 5 分钟启动 Go Scaffold 项目

---

## Purpose

帮助新开发者快速搭建本地开发环境并启动项目。

## Scope

本地开发环境搭建。

## Status

**Active**

---

## 先决条件

| 依赖               | 版本要求 | 验证命令                             |
| ------------------ | -------- | ------------------------------------ |
| Go                 | 1.24+    | `go version`                         |
| PostgreSQL / MySQL | 任意版本 | `psql --version` / `mysql --version` |
| Redis（可选）      | 6.0+     | `redis-cli --version`                |

---

## 步骤

### 1. 克隆仓库

```bash
git clone https://github.com/rei0721/go-scaffold.git
cd go-scaffold
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境

```bash
# 复制环境变量模板
cp .env.example .env

# 复制配置文件模板
cp configs/config.example.yaml configs/config.yaml
```

编辑 `.env` 文件，设置数据库连接信息：

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=rei0721
```

### 4. 创建数据库

**PostgreSQL:**

```bash
createdb rei0721
```

**MySQL:**

```bash
mysql -u root -p -e "CREATE DATABASE rei0721;"
```

### 5. 启动项目

```bash
go run ./cmd/server dev
```

服务启动后访问：`http://localhost:9999`

---

## 常见问题

### Q: 如何禁用 Redis？

在 `configs/config.yaml` 中设置：

```yaml
redis:
  enabled: false
```

### Q: 如何切换到 MySQL？

在 `.env` 中设置：

```env
DB_DRIVER=mysql
DB_PORT=3306
```

---

## Evidence

- 入口文件：[cmd/server/main.go](../cmd/server/main.go)
- 环境变量模板：[.env.example](../.env.example)
- 配置模板：[configs/config.example.yaml](../configs/config.example.yaml)

## Related

- [配置说明](configuration.md)
- [开发指南](development.md)

## Changelog

| 日期       | 变更     |
| ---------- | -------- |
| 2026-01-14 | 初始创建 |
