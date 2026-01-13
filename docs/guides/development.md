# 开发指南

> 项目结构与开发工作流

---

## Purpose

帮助开发者理解项目结构和开发规范。

## Scope

日常开发工作。

## Status

**Active**

---

## 项目结构

```
rei0721/
├── cmd/                    # 可执行入口
│   └── server/             # 主服务（CLI 应用）
├── internal/               # 内部模块（不对外暴露）
│   ├── app/                # DI 容器与生命周期管理
│   ├── config/             # 配置管理
│   ├── handler/            # HTTP 处理器
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   ├── repository/         # 数据访问层
│   ├── router/             # 路由定义
│   └── service/            # 业务服务层
├── pkg/                    # 公共库（可对外暴露）
│   ├── cache/              # Redis 缓存封装
│   ├── cli/                # CLI 框架
│   ├── database/           # GORM 数据库封装
│   ├── executor/           # 任务执行器
│   ├── i18n/               # 国际化
│   ├── logger/             # Zap 日志封装
│   ├── sqlgen/             # SQL 生成工具
│   └── utils/              # 通用工具
├── types/                  # 公共类型定义
├── configs/                # 配置文件
├── locales/                # 国际化资源
└── docs/                   # 文档
```

---

## 分层架构

```
┌─────────────────────────────────────┐
│              Handler                │  ← HTTP 入口
├─────────────────────────────────────┤
│              Service                │  ← 业务逻辑
├─────────────────────────────────────┤
│            Repository               │  ← 数据访问
├─────────────────────────────────────┤
│          Database / Cache           │  ← 基础设施
└─────────────────────────────────────┘
```

**依赖规则**：上层可以依赖下层，下层不能依赖上层。

---

## 开发工作流

### 启动开发服务

```bash
go run ./cmd/server dev
```

### 代码格式化

```bash
go fmt ./...
```

### 运行测试

```bash
go test ./...
```

### 构建

```bash
go build -o bin/server ./cmd/server
```

---

## CLI 命令

项目使用 `pkg/cli` 框架，目前支持：

| 命令     | 说明                          |
| -------- | ----------------------------- |
| `dev`    | 启动开发服务器                |
| `sqlgen` | SQL 脚本生成工具（Model→SQL） |

**sqlgen 使用示例**：

```bash
# 从 Go Model 生成 MySQL DDL
go run ./cmd/server sqlgen --type=mysql --mode=script

# 生成 PostgreSQL DDL
go run ./cmd/server sqlgen --type=postgres --mode=script
```

更多 CLI 用法：

```bash
go run ./cmd/server --help
```

---

## 新增功能步骤

1. **Model** - 在 `internal/models/` 定义数据模型
2. **Repository** - 在 `internal/repository/` 实现数据访问
3. **Service** - 在 `internal/service/` 实现业务逻辑
4. **Handler** - 在 `internal/handler/` 实现 HTTP 处理
5. **Router** - 在 `internal/router/` 注册路由
6. **DI** - 在 `internal/app/app.go` 注入依赖

---

## 代码规范

- 遵循 Go 官方代码规范
- 使用 `golangci-lint` 进行静态检查
- 所有导出函数必须有注释
- 错误处理使用 `types/errors` 包

---

## Evidence

- 入口文件：[cmd/server/main.go](../cmd/server/main.go)
- DI 容器：[internal/app/app.go](../internal/app/app.go)
- CLI 框架：[pkg/cli/](../pkg/cli/)

## Related

- [快速开始](quickstart.md)
- [配置说明](configuration.md)
- [架构概述](../architecture/overview.md)

## Changelog

| 日期       | 变更                 |
| ---------- | -------------------- |
| 2026-01-14 | 初始创建             |
| 2026-01-14 | 补充 sqlgen 命令说明 |
