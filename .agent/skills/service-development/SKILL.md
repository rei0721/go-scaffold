---
name: service-development
description: 在 internal/service/ 目录下创建新的业务服务
---

# 业务服务开发规范

## 概述

本 skill 指导在 `internal/service/` 目录下创建符合项目规范的业务服务。

## 目录结构

```
internal/service/{service-name}/
├── {service}.go        # 接口定义（必需）
├── {service}_impl.go   # 具体实现（必需）
├── {service}_helpers.go # 辅助函数（可选）
└── constants.go        # 常量定义（可选）
```

## 开发步骤

### 1. 创建接口文件 `{service}.go`

```go
// Package {service} 提供 {功能} 服务的实现
// 职责：
// - 职责1
// - 职责2
//
// 设计原则：
// - 与其他 Service 职责分离
// - 支持事务
// - 集成现有组件：JWT、RBAC、Cache、Logger、Executor
package {service}

import (
    "context"

    "github.com/rei0721/go-scaffold/pkg/cache"
    "github.com/rei0721/go-scaffold/pkg/database"
    "github.com/rei0721/go-scaffold/pkg/executor"
    "github.com/rei0721/go-scaffold/pkg/logger"
    "github.com/rei0721/go-scaffold/types"
)

// {Service}Service 定义 {功能} 服务的接口
type {Service}Service interface {
    // 业务方法
    Create(ctx context.Context, req *types.CreateRequest) (*types.Response, error)
    GetByID(ctx context.Context, id int64) (*types.Response, error)
    Update(ctx context.Context, id int64, req *types.UpdateRequest) error
    Delete(ctx context.Context, id int64) error

    // 延迟注入方法
    SetDB(db database.Database)
    SetExecutor(exec executor.Manager)
    SetCache(c cache.Cache)
    SetLogger(l logger.Logger)
}
```

### 2. 创建实现文件 `{service}_impl.go`

```go
package {service}

import (
    "context"
    "sync/atomic"

    "github.com/rei0721/go-scaffold/internal/repository"
    "github.com/rei0721/go-scaffold/pkg/cache"
    "github.com/rei0721/go-scaffold/pkg/database"
    "github.com/rei0721/go-scaffold/pkg/executor"
    "github.com/rei0721/go-scaffold/pkg/logger"
    "github.com/rei0721/go-scaffold/types"
)

// serviceImpl 是 {Service}Service 的具体实现
type serviceImpl struct {
    // 延迟注入的依赖（使用 atomic.Value）
    db       atomic.Value // database.Database
    executor atomic.Value // executor.Manager
    cache    atomic.Value // cache.Cache
    logger   atomic.Value // logger.Logger

    // Repository（如需要）
    repo repository.{Service}Repository
}

// New{Service}Service 创建新的服务实例
func New{Service}Service() {Service}Service {
    return &serviceImpl{}
}

// SetDB 设置数据库依赖（延迟注入）
func (s *serviceImpl) SetDB(db database.Database) {
    s.db.Store(db)
    // 初始化 repository
    s.repo = repository.New{Service}Repository(db.DB())
}

// SetExecutor 设置协程池管理器（延迟注入）
func (s *serviceImpl) SetExecutor(exec executor.Manager) {
    s.executor.Store(exec)
}

// SetCache 设置缓存实例（延迟注入）
func (s *serviceImpl) SetCache(c cache.Cache) {
    s.cache.Store(c)
}

// SetLogger 设置日志记录器（延迟注入）
func (s *serviceImpl) SetLogger(l logger.Logger) {
    s.logger.Store(l)
}

// getDB 获取数据库实例
func (s *serviceImpl) getDB() database.Database {
    if db := s.db.Load(); db != nil {
        return db.(database.Database)
    }
    return nil
}

// getLogger 获取日志实例
func (s *serviceImpl) getLogger() logger.Logger {
    if l := s.logger.Load(); l != nil {
        return l.(logger.Logger)
    }
    return nil
}

// getCache 获取缓存实例
func (s *serviceImpl) getCache() cache.Cache {
    if c := s.cache.Load(); c != nil {
        return c.(cache.Cache)
    }
    return nil
}

// Create 创建实体
func (s *serviceImpl) Create(ctx context.Context, req *types.CreateRequest) (*types.Response, error) {
    log := s.getLogger()

    // 1. 验证请求
    if err := req.Validate(); err != nil {
        return nil, err
    }

    // 2. 业务逻辑
    // ...

    // 3. 调用 Repository
    if err := s.repo.Create(ctx, entity); err != nil {
        if log != nil {
            log.Error("failed to create entity", "error", err)
        }
        return nil, err
    }

    // 4. 清除缓存（如需要）
    if cache := s.getCache(); cache != nil {
        _ = cache.Delete(ctx, cacheKey)
    }

    return response, nil
}
```

### 3. 集成到 App 容器

在 `internal/app/app_business.go` 中添加：

```go
func (a *App) initServices() error {
    // 创建服务
    {service}Svc := {service}.New{Service}Service()

    // 注入依赖
    {service}Svc.SetDB(a.DB)
    {service}Svc.SetLogger(a.Logger)
    {service}Svc.SetCache(a.Cache)
    {service}Svc.SetExecutor(a.Executor)

    a.{Service}Service = {service}Svc
    return nil
}
```

## 事务处理模式

```go
func (s *serviceImpl) CreateWithTransaction(ctx context.Context, req *types.Request) error {
    db := s.getDB()
    if db == nil {
        return errors.New("database not initialized")
    }

    return db.DB().Transaction(func(tx *gorm.DB) error {
        // 在事务中执行多个操作
        if err := s.repo.WithTx(tx).Create(ctx, entity1); err != nil {
            return err
        }
        if err := s.repo.WithTx(tx).Create(ctx, entity2); err != nil {
            return err
        }
        return nil
    })
}
```

## 检查清单

- [ ] 接口定义在 `{service}.go`
- [ ] 实现在 `{service}_impl.go`
- [ ] 所有依赖使用 `Set*` 方法延迟注入
- [ ] 使用 `atomic.Value` 存储延迟注入的依赖
- [ ] 有 `get*` 私有方法安全获取依赖
- [ ] 业务方法第一个参数是 `context.Context`
- [ ] 调用 Repository 处理数据访问
- [ ] 在 `internal/app/app_business.go` 中集成
- [ ] 错误使用 `fmt.Errorf` 包装上下文
