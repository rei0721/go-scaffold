---
name: development-paths
description: 常见开发场景的完整路径指引
updated: 2026-01-19
---

# 开发路径索引

本索引提供常见开发场景的完整步骤指引，帮助快速定位开发方向。

## 场景 1：新增 RESTful API 接口

**复杂度**：中等  
**预计时间**：30-60分钟

### 开发步骤

1. **定义数据模型**
   - 文件：`internal/models/{entity}.go`
   - 操作：使用 GORM 标签定义表结构和关联
   - Skill：[model-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/model-development/SKILL.md)

2. **创建 Repository 接口和实现**
   - 文件：`internal/repository/{entity}_repository.go`
   - 操作：实现 CRUD 和自定义查询方法
   - Skill：[repository-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/repository-development/SKILL.md)

3. **创建 Service 接口**
   - 文件：`internal/service/{module}/interface.go`
   - 操作：定义业务接口方法
   - Skill：[service-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/service-development/SKILL.md)

4. **实现 Service 逻辑**
   - 文件：`internal/service/{module}/{entity}_service.go`
   - 操作：实现业务逻辑，调用 Repository
   - Skill：[service-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/service-development/SKILL.md)

5. **创建 Handler**
   - 文件：`internal/handler/{entity}.go`
   - 操作：实现 HTTP 处理，调用 Service
   - Skill：[handler-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/handler-development/SKILL.md)

6. **注册路由**
   - 文件：`internal/router/router.go`
   - 操作：在路由中注册新的接口

7. **在 App 中初始化**
   - 文件：`internal/app/app_business.go`
   - 操作：初始化 Repository 和 Service

8. **编写测试**
   - 文件：`internal/service/{module}/{entity}_service_test.go`
   - 操作：编写单元测试和集成测试
   - Skill：[test-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/test-development/SKILL.md)

### 注意事项

- 遵循分层架构，从下往上开发（Model → Repository → Service → Handler）
- 先定义接口，再实现逻辑
- 每完成一层就编写测试

---

## 场景 2：新增 Gin 中间件

**复杂度**：简单  
**预计时间**：15-30分钟

### 开发步骤

1. **创建中间件文件**
   - 文件：`internal/middleware/{name}.go`
   - 操作：实现 `gin.HandlerFunc` 函数
   - Skill：[middleware-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/middleware-development/SKILL.md)

2. **注册中间件**
   - 文件：`internal/router/router.go`
   - 操作：在路由中注册中间件（全局或路由组）

3. **编写测试**
   - 文件：`internal/middleware/{name}_test.go`
   - 操作：测试中间件逻辑
   - Skill：[test-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/test-development/SKILL.md)

### 注意事项

- 中间件应该职责单一
- 考虑中间件的执行顺序
- 错误处理要完善

---

## 场景 3：开发可复用工具包

**复杂度**：中等  
**预计时间**：1-3小时

### 开发步骤

1. **创建工具包目录**
   - 路径：`pkg/{tool-name}/`
   - 操作：在 pkg/ 下创建新目录

2. **定义接口**
   - 文件：`pkg/{tool-name}/{tool-name}.go`
   - 操作：定义工具包的接口和主要类型
   - Skill：[pkg-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/pkg-development/SKILL.md)

3. **实现功能**
   - 文件：`pkg/{tool-name}/*.go`
   - 操作：实现具体功能
   - Skill：[pkg-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/pkg-development/SKILL.md)

4. **编写文档**
   - 文件：`pkg/{tool-name}/doc.go`
   - 文件：`pkg/{tool-name}/README.md`
   - 操作：添加包文档和使用说明

5. **编写测试**
   - 文件：`pkg/{tool-name}/*_test.go`
   - 操作：编写单元测试
   - Skill：[test-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/test-development/SKILL.md)

6. **在 App 中集成**
   - 文件：`internal/app/app_infrastructure.go`
   - 操作：初始化并注入到需要的组件

### 注意事项

- pkg/ 下的工具应该独立，不依赖 internal/
- 接口设计要通用和灵活
- 文档要完整清晰

---

## 场景 4：实现用户认证功能

**复杂度**：高  
**预计时间**：3-5小时

### 开发步骤

1. **定义用户模型**
   - 文件：`internal/models/user.go`
   - 操作：用户表结构，包含密码字段

2. **实现密码加密工具**
   - 路径：`pkg/crypto/`
   - 操作：使用 bcrypt 加密密码

3. **实现 JWT 工具**
   - 路径：`pkg/jwt/`
   - 操作：Token 生成和验证

4. **创建 Auth Repository**
   - 文件：`internal/repository/auth_repository.go`
   - 操作：用户查询和验证

5. **创建 Auth Service**
   - 路径：`internal/service/auth/`
   - 操作：登录、注册、Token 刷新逻辑

6. **创建 Auth Handler**
   - 文件：`internal/handler/auth.go`
   - 操作：登录、注册、登出接口

7. **创建认证中间件**
   - 文件：`internal/middleware/auth.go`
   - 操作：JWT 验证中间件

8. **注册路由**
   - 文件：`internal/router/router.go`
   - 操作：注册认证接口和应用中间件

9. **编写测试**
   - 文件：`internal/service/auth/*_test.go`
   - 文件：`internal/handler/auth_test.go`
   - 操作：测试认证流程

### 注意事项

- 密码必须加密存储，不能明文
- Token 应该设置合理的过期时间
- 实现 Token 刷新机制
- 考虑安全性（防暴力破解、CSRF等）

---

## 场景 5：实现 RBAC 权限控制

**复杂度**：高  
**预计时间**：4-6小时

### 开发步骤

1. **定义角色和权限模型**
   - 文件：`internal/models/role.go`
   - 文件：`internal/models/permission.go`
   - 操作：定义 RBAC 相关模型

2. **封装 Casbin RBAC**
   - 路径：`pkg/rbac/`
   - 操作：封装 Casbin，提供 RBAC 接口

3. **创建 RBAC Repository**
   - 文件：`internal/repository/rbac_repository.go`
   - 操作：角色和权限的数据访问

4. **创建 RBAC Service**
   - 路径：`internal/service/rbac/`
   - 操作：角色管理、权限分配逻辑

5. **创建 RBAC 中间件**
   - 文件：`internal/middleware/rbac.go`
   - 操作：权限验证中间件

6. **创建 RBAC Handler**
   - 文件：`internal/handler/rbac.go`
   - 操作：角色和权限管理接口

7. **应用中间件到路由**
   - 文件：`internal/router/router.go`
   - 操作：为需要权限控制的路由添加中间件

### 注意事项

- 设计好角色和权限的数据模型
- 考虑权限的粒度（粗粒度 vs 细粒度）
- 实现权限缓存提高性能

---

## 场景 6：为功能添加缓存层

**复杂度**：中等  
**预计时间**：1-2小时

### 开发步骤

1. **确认 Cache 工具可用**
   - 路径：`pkg/cache/`
   - 操作：确保项目已集成 Cache 工具

2. **在 Service 中注入 Cache**
   - 文件：`internal/service/{module}/`
   - 操作：使用 `SetCache()` 方法注入

3. **实现缓存逻辑**
   - 文件：`internal/service/{module}/*_service.go`
   - 操作：
     - 读取：先查缓存，缓存未命中再查数据库
     - 写入：更新数据库后更新缓存
     - 删除：删除数据库后清除缓存

4. **定义缓存失效策略**
   - 操作：
     - 设置合理的 TTL
     - 更新、删除时清除相关缓存

5. **编写测试**
   - 文件：`internal/service/{module}/*_test.go`
   - 操作：测试缓存逻辑

### 示例代码

```go
func (s *userService) GetByID(id uint) (*models.User, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("user:%d", id)
    if s.cache != nil {
        var user models.User
        err := s.cache.Get(cacheKey, &user)
        if err == nil {
            return &user, nil
        }
    }

    // 2. 缓存未命中，查询数据库
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    if s.cache != nil {
        s.cache.Set(cacheKey, user, 5*time.Minute)
    }

    return user, nil
}
```

### 注意事项

- 缓存 Key 的命名要规范
- 缓存失效策略要合理
- 考虑缓存穿透、击穿、雪崩问题

---

## 场景 7：集成新的配置项

**复杂度**：简单  
**预计时间**：15-30分钟

### 开发步骤

1. **在配置文件中添加配置**
   - 文件：`config/config.yaml`
   - 操作：添加新的配置项

2. **生成配置结构体**
   - 操作：使用 yaml2go 工具生成或手动添加到配置结构

3. **在组件中使用配置**
   - 文件：`internal/app/app_*.go`
   - 操作：从配置中读取并传递给组件

### 注意事项

- 配置项应有合理的默认值
- 环境变量覆盖机制
- 配置验证

---

## 场景 8：添加数据库迁移

**复杂度**：简单  
**预计时间**：10-20分钟

### 开发步骤

1. **定义或修改模型**
   - 文件：`internal/models/*.go`
   - 操作：更新数据模型

2. **使用 AutoMigrate**
   - 文件：`internal/app/app_infrastructure.go`
   - 操作：在数据库初始化时调用 AutoMigrate

### 注意事项

- 生产环境应使用版本化的迁移工具
- 备份数据库再迁移
- 测试迁移是否成功

---

## 场景 9：编写单元测试

**复杂度**：中等  
**预计时间**：30分钟-1小时

### 开发步骤

1. **创建测试文件**
   - 文件：`{module}_test.go`
   - 操作：在同目录下创建测试文件

2. **编写测试用例**
   - 操作：使用 testing 包编写测试
   - Skill：[test-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/test-development/SKILL.md)

3. **Mock 依赖**
   - 操作：使用 mock 替代真实依赖

4. **运行测试**
   - 命令：`go test ./...`

### 注意事项

- 测试覆盖率要足够
- 测试应该独立，不依赖外部资源
- 使用表驱动测试提高代码复用

---

## 开发流程最佳实践

### 实践 1：自下而上开发

从 Model → Repository → Service → Handler 的顺序开发。

**原因**：下层是上层的依赖，先实现依赖。

### 实践 2：接口先行

先定义接口，再实现逻辑。

**原因**：接口定义了契约，便于测试和理解。

### 实践 3：测试驱动

编写功能的同时编写测试。

**原因**：确保代码质量，便于重构。

### 实践 4：小步提交

完成一个步骤就提交一次。

**原因**：便于回滚和代码审查。

### 实践 5：文档同步

代码和文档同步更新。

**原因**：保持文档的准确性。

---

## 快速参考

| 场景       | 复杂度 | 时间      | 主要涉及层级   |
| ---------- | ------ | --------- | -------------- |
| 新增 API   | 中等   | 30-60分钟 | 全部           |
| 新增中间件 | 简单   | 15-30分钟 | Presentation   |
| 开发工具包 | 中等   | 1-3小时   | Infrastructure |
| 用户认证   | 高     | 3-5小时   | 全部 + pkg     |
| RBAC 权限  | 高     | 4-6小时   | 全部 + pkg     |
| 添加缓存   | 中等   | 1-2小时   | Business       |
| 集成配置   | 简单   | 15-30分钟 | Infrastructure |
| 数据库迁移 | 简单   | 10-20分钟 | Data           |
| 编写测试   | 中等   | 30-60分钟 | 对应层级       |
