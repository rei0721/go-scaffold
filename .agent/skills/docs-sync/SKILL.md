---
name: docs-sync
description: 检查项目代码变更是否影响文档一致性，并指导同步更新文档
---

# 文档同步检查与更新

## 概述

本 skill 用于检查项目代码变更是否影响文档的一致性，确保 `docs/` 目录下的文档与项目代码保持同步。当进行代码修改、重构或新增功能时，应使用此 skill 判断是否需要更新相关文档。

### 适用场景

- ✅ 修改了项目架构或设计模式
- ✅ 新增/修改/删除了 API 接口
- ✅ 修改了配置文件结构或参数
- ✅ 新增了公共组件（`pkg/` 下的包）
- ✅ 修改了业务模块（`internal/` 下的层次结构）
- ✅ 修改了数据模型或数据库结构
- ✅ 修改了认证/授权机制
- ✅ 修改了部署方式或环境要求

### 不适用场景

- ❌ 仅修改了代码注释
- ❌ 仅修复了 bug 但不影响接口
- ❌ 仅优化了内部实现逻辑但不改变外部行为

## 文档-代码映射关系

### 1. 项目概览文档

#### docs/overview/introduction.md

- **关联代码**: 整体项目特性、核心功能
- **触发同步条件**:
  - 新增核心功能模块（如新增 GraphQL、gRPC 支持）
  - 修改项目定位或目标
  - 技术栈发生重大变化
- **检查点**:
  - [ ] 核心特性列表是否完整
  - [ ] 项目优势描述是否准确
  - [ ] 适用场景是否涵盖新功能

#### docs/overview/architecture.md

- **关联代码**:
  - `internal/app/` - 依赖注入容器实现
  - 整体分层架构（Handler → Service → Repository）
  - 中间件机制
- **触发同步条件**:
  - 修改了分层架构设计
  - 新增/删除了架构组件
  - 修改了依赖注入模式
  - 修改了生命周期管理方式
- **检查点**:
  - [ ] 架构图是否反映当前设计
  - [ ] 分层说明是否准确
  - [ ] 依赖注入流程是否最新
  - [ ] 组件生命周期说明是否完整

#### docs/overview/tech-stack.md

- **关联代码**: `go.mod` 中的依赖项
- **触发同步条件**:
  - 新增主要第三方库
  - 升级核心依赖版本
  - 替换技术选型（如 ORM、日志库）
- **检查点**:
  - [ ] 核心依赖列表是否完整
  - [ ] 版本号是否最新
  - [ ] 技术选型理由是否准确

### 2. 快速开始文档

#### docs/getting-started/prerequisites.md

- **关联代码**:
  - `go.mod` - Go 版本要求
  - `docker-compose.yml` - 依赖服务
  - `configs/config.yaml` - 必需的外部服务
- **触发同步条件**:
  - 修改了 Go 版本要求
  - 新增了必需的外部服务（如 Redis、PostgreSQL）
  - 修改了开发工具要求
- **检查点**:
  - [ ] Go 版本要求是否正确
  - [ ] 外部服务列表是否完整（数据库、缓存等）
  - [ ] 开发工具列表是否准确

#### docs/getting-started/installation.md

- **关联代码**:
  - `Makefile` - 构建命令
  - `.env.example` - 环境变量示例
  - `configs/config.yaml` - 配置文件
- **触发同步条件**:
  - 修改了安装步骤或构建命令
  - 新增了配置项
  - 修改了初始化流程
- **检查点**:
  - [ ] 安装步骤是否完整
  - [ ] 配置示例是否最新
  - [ ] Make 命令是否准确

#### docs/getting-started/quickstart.md

- **关联代码**:
  - `cmd/server/` - 启动命令
  - 主要 API 接口
- **触发同步条件**:
  - 修改了启动方式
  - 核心 API 接口变更
  - 默认端口或配置变更
- **检查点**:
  - [ ] 启动命令是否正确
  - [ ] 示例 API 请求是否有效
  - [ ] 响应格式是否最新

### 3. 开发指南文档

#### docs/development/project-structure.md

- **关联代码**: 整个项目的目录结构
- **触发同步条件**:
  - 新增/删除了目录
  - 新增了 `pkg/` 下的公共包
  - 新增了 `internal/` 下的业务模块
  - 修改了文件组织方式
- **检查点**:
  - [ ] 目录树是否反映当前结构
  - [ ] 新增包的说明是否完整
  - [ ] 文件命名规范是否准确
  - [ ] 目录职责说明是否清晰

#### docs/development/configuration.md

- **关联代码**:
  - `configs/config.yaml` - 配置文件结构
  - `internal/config/` - 配置结构定义
  - `.env.example` - 环境变量
- **触发同步条件**:
  - 新增/修改/删除配置项
  - 修改了配置验证规则
  - 新增了环境变量
  - 修改了配置热重载机制
- **检查点**:
  - [ ] 配置文件示例是否完整
  - [ ] 配置项说明是否准确
  - [ ] 环境变量列表是否最新
  - [ ] 默认值是否正确
  - [ ] 配置结构定义代码是否同步

#### docs/development/database.md

- **关联代码**:
  - `internal/models/` - 数据模型
  - `pkg/database/` - 数据库抽象
  - 数据库迁移文件
- **触发同步条件**:
  - 新增/修改/删除数据模型
  - 修改了数据库连接方式
  - 新增了数据库迁移
- **检查点**:
  - [ ] 数据模型定义是否最新
  - [ ] ER 图是否反映当前设计
  - [ ] 迁移说明是否完整

#### docs/development/api-development.md

- **关联代码**:
  - `internal/handler/` - HTTP 处理器
  - `internal/router/` - 路由定义
  - `types/result/` - 响应格式
- **触发同步条件**:
  - 新增/修改了 API 接口
  - 修改了 API 规范或约定
  - 修改了响应格式
- **检查点**:
  - [ ] API 开发流程是否准确
  - [ ] 代码示例是否有效
  - [ ] 错误处理说明是否完整

#### docs/development/testing.md

- **关联代码**: 测试文件（`*_test.go`）
- **触发同步条件**:
  - 修改了测试框架或工具
  - 新增了测试类型
  - 修改了测试运行方式
- **检查点**:
  - [ ] 测试命令是否正确
  - [ ] 测试框架说明是否准确

#### docs/development/logging.md

- **关联代码**:
  - `pkg/logger/` - 日志系统
  - `internal/app/app_logger.go` - 日志配置
- **触发同步条件**:
  - 修改了日志库
  - 修改了日志格式或级别
  - 新增了日志功能
- **检查点**:
  - [ ] 日志使用示例是否正确
  - [ ] 日志级别说明是否准确

### 4. API 文档

#### docs/api/overview.md

- **关联代码**:
  - `internal/handler/` - 所有 HTTP 处理器
  - `types/result/` - 响应格式定义
- **触发同步条件**:
  - 修改了 API 设计原则
  - 修改了统一响应格式
  - 修改了 API 版本策略
- **检查点**:
  - [ ] API 规范是否最新
  - [ ] 响应格式示例是否准确

#### docs/api/authentication.md

- **关联代码**:
  - `pkg/jwt/` - JWT 实现
  - `internal/middleware/auth.go` - 认证中间件
  - `internal/handler/auth_handler.go` - 认证处理器
- **触发同步条件**:
  - 修改了认证机制
  - 修改了 JWT 配置
  - 新增认证方式
- **检查点**:
  - [ ] 认证流程说明是否准确
  - [ ] JWT 配置说明是否最新
  - [ ] 认证示例是否有效

#### docs/api/error-handling.md

- **关联代码**:
  - `types/errors/` - 错误定义
  - `types/result/` - 错误响应格式
- **触发同步条件**:
  - 新增/修改错误码
  - 修改了错误响应格式
  - 修改了错误处理机制
- **检查点**:
  - [ ] 错误码列表是否完整
  - [ ] 错误响应示例是否准确

#### docs/api/endpoints.md

- **关联代码**:
  - `internal/handler/` - 所有处理器
  - `internal/router/` - 路由定义
- **触发同步条件**:
  - 新增/修改/删除 API 接口
  - 修改了请求/响应格式
  - 修改了接口路径
- **检查点**:
  - [ ] 接口列表是否完整
  - [ ] 请求参数是否准确
  - [ ] 响应示例是否正确
  - [ ] 接口路径是否最新

### 5. 运维部署文档

#### docs/deployment/deployment.md

- **关联代码**:
  - `Dockerfile` - 容器构建
  - `Makefile` - 构建命令
  - `configs/` - 生产配置
- **触发同步条件**:
  - 修改了部署方式
  - 新增了部署步骤
  - 修改了环境要求
- **检查点**:
  - [ ] 部署步骤是否完整
  - [ ] 配置说明是否准确

#### docs/deployment/docker.md

- **关联代码**:
  - `Dockerfile` - Docker 镜像构建
  - `docker-compose.yml` - Docker Compose 配置
- **触发同步条件**:
  - 修改了 Dockerfile
  - 修改了 docker-compose 配置
  - 新增了容器服务
- **检查点**:
  - [ ] Docker 镜像构建说明是否准确
  - [ ] docker-compose 配置是否最新

### 6. 贡献指南文档

#### docs/contributing/code-style.md

- **关联代码**: 整个项目的代码实现
- **触发同步条件**:
  - 修改了代码规范
  - 新增了编程约定
  - 修改了工具配置（如 golangci-lint）
- **检查点**:
  - [ ] 代码规范是否反映实际实践
  - [ ] 示例代码是否符合规范

#### docs/contributing/commit-convention.md

- **关联代码**: Git 提交历史
- **触发同步条件**:
  - 修改了提交规范
- **检查点**:
  - [ ] 提交格式说明是否准确

### 7. 其他文档

#### docs/faq.md

- **关联代码**: 常见问题相关的代码实现
- **触发同步条件**:
  - 发现新的常见问题
  - 修改了问题的解决方案
  - 修改了相关代码导致答案失效
- **检查点**:
  - [ ] 问题解答是否仍然有效
  - [ ] 代码示例是否可运行

#### docs/changelog.md

- **关联代码**: 所有代码变更
- **触发同步条件**:
  - 每次版本发布
  - 重大功能变更
  - 破坏性更新
- **检查点**:
  - [ ] 是否记录了新版本的变更
  - [ ] 变更说明是否完整

## 检查流程

### 步骤 1: 识别代码变更范围

根据代码变更确定影响范围：

```bash
# 查看最近的代码变更
git diff HEAD~1 --name-only

# 或查看具体分支的变更
git diff main..feature-branch --name-only
```

### 步骤 2: 确定受影响的文档

根据上面的"文档-代码映射关系"表，确定需要检查的文档：

| 变更类型                    | 需检查的文档                                                                |
| --------------------------- | --------------------------------------------------------------------------- |
| 修改 `pkg/` 下的包          | `docs/overview/tech-stack.md`, `docs/development/project-structure.md`      |
| 修改 `internal/config/`     | `docs/development/configuration.md`, `docs/getting-started/installation.md` |
| 修改 `internal/models/`     | `docs/development/database.md`, `docs/development/project-structure.md`     |
| 修改 `internal/handler/`    | `docs/api/endpoints.md`, `docs/development/api-development.md`              |
| 修改 `internal/middleware/` | `docs/overview/architecture.md`, `docs/api/authentication.md`               |
| 修改 `types/errors/`        | `docs/api/error-handling.md`                                                |
| 修改 `Dockerfile`           | `docs/deployment/docker.md`                                                 |
| 修改 `go.mod`               | `docs/overview/tech-stack.md`, `docs/getting-started/prerequisites.md`      |
| 新增功能模块                | `docs/overview/introduction.md`, `docs/development/project-structure.md`    |
| 修改架构设计                | `docs/overview/architecture.md`                                             |

### 步骤 3: 逐项检查文档

针对每个受影响的文档，使用对应的"检查点"进行核对：

1. 打开文档文件
2. 对照代码实现，逐项检查"检查点"
3. 标记需要更新的内容

### 步骤 4: 更新文档

根据检查结果更新文档：

#### 更新原则

- **准确性**: 确保描述与代码实现一致
- **完整性**: 新增功能必须有对应的文档说明
- **清晰性**: 使用清晰的语言和示例
- **时效性**: 及时更新版本号和日期

#### 更新模板

每个文档通常包含：

````markdown
# 文档标题

## 概述

简要说明本文档的内容

## 核心内容

详细说明

## 代码示例

```go
// 可运行的代码示例
```
````

## 最佳实践

注意事项和建议

---

**最后更新**: YYYY年MM月  
**版本**: vX.Y.Z

````

### 步骤 5: 验证更新

- [ ] 检查所有链接是否有效
- [ ] 确认代码示例可以运行
- [ ] 验证配置示例是否正确
- [ ] 确认命令示例可以执行

## 同步更新规范

### 同步时机

1. **即时同步**:
   - 新增/修改 API 接口
   - 修改配置结构
   - 修改环境要求

2. **版本发布时同步**:
   - 更新 `docs/changelog.md`
   - 更新版本号
   - 整体检查文档一致性

3. **定期同步**:
   - 每月检查一次 FAQ
   - 季度检查技术栈文档

### 同步检查清单

- [ ] 已识别代码变更范围
- [ ] 已确定受影响的文档列表
- [ ] 已逐项检查所有检查点
- [ ] 已更新所有需要修改的文档
- [ ] 已验证代码示例可运行
- [ ] 已验证配置示例正确
- [ ] 已更新文档底部的版本号和日期
- [ ] 已检查文档内部链接有效性
- [ ] 已提交文档变更到 Git

## 常见文档更新场景

### 场景 1: 新增 pkg 包

**影响文档**:
- `docs/development/project-structure.md`
- `docs/overview/tech-stack.md`（如果是新技术）

**更新步骤**:
1. 在 `project-structure.md` 的 `pkg/` 部分添加新包说明
2. 添加包的职责描述
3. 添加使用示例
4. 更新目录结构树

### 场景 2: 修改配置项

**影响文档**:
- `docs/development/configuration.md`

**更新步骤**:
1. 更新配置文件 YAML 示例
2. 更新配置结构定义代码
3. 更新环境变量列表
4. 更新验证规则说明
5. 如果是破坏性变更，添加到 FAQ

### 场景 3: 新增 API 接口

**影响文档**:
- `docs/api/endpoints.md`
- `docs/development/api-development.md`（可能）

**更新步骤**:
1. 在 `endpoints.md` 中添加新接口
2. 说明请求方法、路径、参数
3. 提供请求和响应示例
4. 说明错误码
5. 添加使用注意事项

### 场景 4: 重构架构

**影响文档**:
- `docs/overview/architecture.md`
- `docs/development/project-structure.md`

**更新步骤**:
1. 更新架构图（Mermaid 图表）
2. 更新分层说明
3. 更新数据流向图
4. 更新依赖注入流程
5. 更新最佳实践

## 文档质量标准

### 必需元素

每个文档都应包含：
- [ ] 清晰的标题和概述
- [ ] 结构化的章节
- [ ] 可运行的代码示例
- [ ] 实用的检查清单或最佳实践
- [ ] 版本号和更新日期

### 代码示例规范

```go
// ✅ 好的示例：完整可运行
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

// ❌ 不好的示例：不完整，缺少 import
func main() {
    fmt.Println("Hello, World!")
}
````

### 配置示例规范

```yaml
# ✅ 好的示例：包含注释说明
database:
  host: "localhost"        # 数据库地址
  port: 3306              # 数据库端口

# ❌ 不好的示例：没有注释
database:
  host: "localhost"
  port: 3306
```

## 工具和自动化

### 文档链接检查

```bash
# 使用 markdown-link-check 检查链接
npm install -g markdown-link-check
find docs -name "*.md" -exec markdown-link-check {} \;
```

### 代码示例验证

````bash
# 提取文档中的 Go 代码并验证编译
# 可以编写自定义脚本提取 ```go 块并验证
````

### 文档生成

```bash
# 使用 make 命令生成文档站点
make docs-serve
```

## 最佳实践

### 1. 及时更新

- 代码变更后立即更新相关文档
- 不要等到版本发布时才批量更新

### 2. 保持一致

- 文档中的术语与代码保持一致
- 配置示例与实际默认值一致
- API 示例与实际接口一致

### 3. 用户视角

- 从用户角度编写文档
- 提供完整的使用场景
- 包含故障排除指南

### 4. 版本管理

- 为重大变更创建迁移指南
- 在 changelog 中记录所有变更
- 标注破坏性变更

### 5. 协作规范

- 代码 PR 应包含文档更新
- 文档变更需要 Review
- 定期进行文档质量审查

## 参考资源

- [项目文档主页](../../docs/README.md)
- [代码规范](../../docs/contributing/code-style.md)
- [提交规范](../../docs/contributing/commit-convention.md)
- [work-log skill](../work-log/SKILL.md) - 完成变更后记录工作日志

## 检查清单

使用此 skill 时的完整检查清单：

- [ ] 已识别本次代码变更的文件和范围
- [ ] 已查阅"文档-代码映射关系"确定受影响文档
- [ ] 已对每个受影响文档使用对应检查点进行核对
- [ ] 已更新所有不一致的文档内容
- [ ] 代码示例已验证可运行
- [ ] 配置示例已验证正确
- [ ] 命令示例已验证可执行
- [ ] 文档内部链接已检查有效
- [ ] 已更新文档版本号和更新日期
- [ ] 已将文档变更提交到 Git
- [ ] 如有破坏性变更，已更新 FAQ 和迁移指南
