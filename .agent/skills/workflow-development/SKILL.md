---
name: workflow-development
description: 在 .agent/workflows/ 目录下创建可复用的工作流
---

# 工作流开发规范

## 概述

本 skill 指导在 `.agent/workflows/` 目录下创建可复用的工作流文件，用于自动化常见开发任务。

## 文件结构

```
.agent/
└── workflows/
    ├── {workflow-name}.md    # 工作流定义
    └── ...
```

## 工作流格式

### 基础模板

````markdown
---
description: 简短描述工作流的用途
---

# 工作流名称

## 前置条件

- 条件1
- 条件2

## 步骤

1. 第一步描述
   ```bash
   具体命令
   ```
````

2. 第二步描述
   ```bash
   另一个命令
   ```

````

## Turbo 注解

### 单步自动运行

使用 `// turbo` 注解标记可以自动运行的单个步骤：

```markdown
2. 创建目录
// turbo
3. 安装依赖
   ```bash
   go mod tidy
````

````

只有第 3 步会自动运行，第 2 步仍需确认。

### 全局自动运行

使用 `// turbo-all` 注解标记所有步骤都可自动运行：

```markdown
---
description: 构建和测试工作流
---

// turbo-all

## 步骤

1. 格式化代码
   ```bash
   go fmt ./...
````

2. 运行测试
   ```bash
   go test ./...
   ```

````

所有步骤都会自动运行。

## 常用工作流示例

### 构建部署

```markdown
---
description: 构建并部署应用
---

## 步骤

1. 运行测试确保代码正确
   ```bash
   go test ./...
````

// turbo 2. 构建生产版本

```bash
go build -o bin/server ./cmd/server
```

3. 部署（需要确认）
   ```bash
   docker-compose up -d
   ```

````

### 数据库迁移

```markdown
---
description: 执行数据库迁移
---

## 前置条件
- 数据库服务已启动
- 已配置正确的连接信息

## 步骤

1. 备份当前数据库
   ```bash
   pg_dump -h localhost -U postgres mydb > backup.sql
````

// turbo 2. 运行迁移

```bash
go run ./cmd/server initdb
```

````

### 代码质量检查

```markdown
---
description: 运行完整的代码质量检查
---

// turbo-all

## 步骤

1. 格式化代码
   ```bash
   go fmt ./...
   goimports -w .
````

2. 静态检查

   ```bash
   golangci-lint run
   ```

3. 运行测试

   ```bash
   go test -v -race ./...
   ```

4. 检查覆盖率
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   ```

```

## 调用工作流

在聊天中使用斜杠命令调用工作流：

```

/workflow-name

```

例如：
- `/build-deploy` - 调用构建部署工作流
- `/code-quality` - 调用代码质量检查工作流

## 命名规范

| 规则 | 示例 |
|-----|------|
| 使用小写字母 | `build-deploy.md` |
| 使用连字符分隔单词 | `code-quality.md` |
| 描述性名称 | `database-migration.md` |

## 检查清单

- [ ] 文件位于 `.agent/workflows/` 目录
- [ ] 文件名使用小写字母和连字符
- [ ] 包含 YAML frontmatter（`description`）
- [ ] 步骤编号清晰
- [ ] 命令使用代码块包裹
- [ ] 危险操作不使用 `// turbo`
- [ ] 安全操作可使用 `// turbo` 或 `// turbo-all`
```
