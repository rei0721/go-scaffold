---
name: project-structure
description: Go-Scaffold 项目完整结构索引
updated: 2026-01-19
---

# 项目结构索引

## 根目录：go-scaffold

Go Web 应用脚手架

## 目录树

### cmd/ - 应用程序入口

```
cmd/
└── server/
    └── main.go          # 主入口文件，初始化应用
```

**职责**：应用程序的启动入口

**关键文件**：

- `main.go` - 应用启动，读取配置，初始化 App

---

### config/ - 配置文件

```
config/
├── config.yaml          # 主配置文件
└── config.example.yaml  # 配置示例
```

**职责**：存放应用配置

**说明**：使用 Viper 读取配置

---

### internal/ - 私有应用代码

```
internal/
├── app/                 # 应用核心
│   ├── app.go
│   ├── app_business.go
│   └── app_infrastructure.go
├── handler/             # HTTP 处理器
│   ├── user.go
│   ├── auth.go
│   └── rbac.go
├── service/             # 业务逻辑层
│   ├── user/
│   └── auth/
├── repository/          # 数据访问层
│   ├── user_repository.go
│   ├── auth_repository.go
│   └── rbac_repository.go
├── models/              # 数据模型
│   ├── user.go
│   └── role.go
├── middleware/          # Gin 中间件
│   ├── auth.go
│   ├── rbac.go
│   └── logger.go
├── router/              # 路由配置
│   └── router.go
└── types/               # 类型定义
    ├── interfaces.go
    └── constants.go
```

#### internal/app - 应用核心

**职责**：应用初始化和组装

**关键文件**：

- `app.go` - 应用主结构体
- `app_business.go` - 业务组件初始化（Repository, Service, Handler）
- `app_infrastructure.go` - 基础设施初始化（Logger, Cache, DB）

**层级**：Infrastructure

---

#### internal/handler - HTTP 处理器

**职责**：HTTP 请求处理，参数验证，响应封装

**层级**：Presentation

**文件组织**：

- 按业务实体组织：`user.go`, `auth.go`, `rbac.go`
- 每个文件包含一组相关接口

**关键文件**：

- `user.go` - 用户相关接口（CRUD）
- `auth.go` - 认证接口（登录、注册、登出）
- `rbac.go` - 权限管理接口

---

#### internal/service - 业务逻辑层

**职责**：业务逻辑处理，事务管理

**层级**：Business

**文件组织**：

- 按模块组织子目录
- 每个模块有 `interface.go` 定义接口
- 服务实现在 `*_service.go`

**示例结构**：

```
service/
├── user/
│   ├── interface.go      # UserService 接口
│   └── user_service.go   # 实现
└── auth/
    ├── interface.go      # AuthService 接口
    └── auth_service.go   # 实现
```

---

#### internal/repository - 数据访问层

**职责**：数据持久化，数据库操作

**层级**：Data

**文件命名**：`{entity}_repository.go`

**关键文件**：

- `user_repository.go` - 用户数据访问
- `auth_repository.go` - 认证数据访问
- `rbac_repository.go` - 权限数据访问

---

#### internal/models - 数据模型

**职责**：数据模型定义，表结构映射

**层级**：Data

**文件命名**：`{entity}.go`

**关键文件**：

- `user.go` - 用户模型（GORM）
- `role.go` - 角色模型

---

#### internal/middleware - Gin 中间件

**职责**：请求拦截，认证授权，日志记录

**层级**：Presentation

**关键文件**：

- `auth.go` - JWT 认证中间件
- `rbac.go` - RBAC 权限控制
- `logger.go` - 请求日志

---

#### internal/router - 路由配置

**职责**：路由注册，中间件配置

**层级**：Presentation

**关键文件**：

- `router.go` - 路由注册和配置

---

### pkg/ - 可复用工具包

```
pkg/
├── logger/              # 日志工具
├── cache/               # 缓存工具
├── jwt/                 # JWT 认证
├── rbac/                # RBAC 权限
├── httpserver/          # HTTP 服务器封装
├── executor/            # 协程池
├── sqlgen/              # SQL 生成
└── yaml2go/             # YAML 转 Go
```

**职责**：可复用的工具组件

**层级**：Infrastructure

**文件组织**：

- 每个工具独立目录
- 必须包含 `doc.go` 和 `README.md`

**关键包**：

- `logger` - 基于 Zap 的日志封装
- `cache` - Redis 缓存封装
- `jwt` - JWT 令牌生成和验证
- `rbac` - Casbin RBAC 封装

---

### docs/ - 项目文档

```
docs/
└── worklogs/            # 工作日志
    └── YYYY/MM/         # 按年月组织
```

**职责**：项目文档和工作日志

---

### .agent/ - AI 助手配置

```
.agent/
├── skills/              # Skills 定义
└── workflows/           # 工作流定义
```

**职责**：AI 助手的配置和规范

---

## 文件组织规则

### 规则 1：按层级组织

`internal/` 下按架构分层组织：

- Presentation: handler, middleware, router
- Business: service
- Data: repository, models

### 规则 2：按功能模块组织

每个模块有独立的文件：

- `user.go` - 用户相关
- `auth.go` - 认证相关

### 规则 3：接口分离

服务接口定义在 `interface.go`，实现在独立文件。

### 规则 4：可复用代码放 pkg

通用工具和可复用组件放在 `pkg/` 目录，不依赖 `internal/`。

---

## 快速定位

| 我想找...  | 查看目录                    |
| ---------- | --------------------------- |
| 应用入口   | `cmd/server/main.go`        |
| 应用初始化 | `internal/app/`             |
| HTTP 接口  | `internal/handler/`         |
| 业务逻辑   | `internal/service/`         |
| 数据访问   | `internal/repository/`      |
| 数据模型   | `internal/models/`          |
| 中间件     | `internal/middleware/`      |
| 路由配置   | `internal/router/router.go` |
| 工具包     | `pkg/`                      |
