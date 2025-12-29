# Rei0721

> Go 高性能脚手架 | 资源受控 · 配置热重载 · 泛型驱动

```
版本: v1.0.0 | 更新: 2025-12-29 | 状态: Active
```

## 核心特性

| 特性 | 说明 |
|------|------|
| 🚀 资源受控 | 统一协程调度，消除泄漏 |
| 🔄 配置热重载 | 中间件式无感知更新 |
| 🌐 国际化 | 完整 i18n 方案 |
| 🗄️ 多数据库 | PostgreSQL / MySQL / SQLite |
| 🔧 泛型驱动 | 充分解耦 |

## 技术栈

```
Go 1.21+ | Gin | GORM | go-redis/v9 | Viper | Zap/Logrus | ants | go-i18n/v2 | Snowflake
```

## 快速开始

```bash
git clone <repository-url> && cd rei0721
cp configs/.env.example configs/.env
go mod download
go run cmd/server/main.go
```

## 项目结构

```
rei0721/
├── cmd/server/          # 启动入口
├── configs/             # 配置文件 + i18n
├── docs/                # 文档
├── internal/
│   ├── app/             # 依赖注入容器
│   ├── config/          # 配置定义
│   ├── handler/         # HTTP 处理层
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   ├── repository/      # 数据访问层
│   ├── router/          # 路由
│   └── service/         # 业务逻辑层
├── pkg/                 # 通用工具 (严禁引用 internal)
├── types/               # 类型定义 (无依赖)
│   ├── constants/
│   ├── errors/
│   └── result/
└── logs/
```

## 文档索引

| 文档 | 说明 | 必读 |
|------|------|------|
| [protocol.md](./protocol.md) | 工程规范协议 | ⭐⭐⭐ |
| [design.md](./design.md) | 架构设计 | ⭐⭐ |
| [api.md](./api.md) | API 规范 | ⭐⭐ |
| [deployment.md](./deployment.md) | 部署指南 | ⭐ |

## 层级依赖规则

```
cmd → internal/app → handler → service → repository → models
                  ↘ config
                  ↘ pkg (可引用)

types: 无依赖
pkg: 严禁引用 internal
```

## 贡献前必读

1. 阅读 [protocol.md](./protocol.md)
2. `go fmt && go vet`
3. 完成自检清单

---

MIT License | Built with Go
