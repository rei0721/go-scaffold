---
name: skill-development
description: 识别场景并创建所需的 skill 文件
---

# Skill 开发规范

## 概述

本 skill 指导如何识别开发场景中缺失的 skill，并创建符合项目规范的新 skill 文件。

## 识别缺失的 Skill

### 何时需要创建新 Skill

当遇到以下情况时，考虑创建新 skill：

1. **重复性任务**：某个开发任务需要多次执行，但没有现成的 skill 指导
2. **新技术栈**：引入新的技术或工具，需要标准化使用方式
3. **复杂流程**：某个流程涉及多个步骤，需要明确的指导文档
4. **最佳实践**：发现了更好的实现方式，需要形成规范

### 场景分析示例

| 场景                  | 是否需要新 Skill | 建议                           |
| --------------------- | ---------------- | ------------------------------ |
| 添加新的 API 接口     | ❌               | 使用现有 `handler-development` |
| 创建 GraphQL resolver | ✅               | 创建 `graphql-development`     |
| 使用 gRPC 服务        | ✅               | 创建 `grpc-development`        |
| 部署到 Kubernetes     | ✅               | 创建 `k8s-deployment`          |

## Skill 设计流程

### 1. 确定 Skill 范围

回答以下问题：

- **目的**：这个 skill 要解决什么问题？
- **受众**：谁会使用这个 skill？（开发者、运维、测试？）
- **频率**：这个任务多久执行一次？
- **复杂度**：涉及多少个步骤？

### 2. 定义内容结构

标准 skill 包含的部分：

```markdown
---
name: { skill-name }
description: 简短描述（一句话）
---

# {标题}

## 概述

简要说明这个 skill 的用途和适用场景

## 文件结构（如果涉及代码）

目录和文件组织方式

## 开发步骤

1. 步骤一
2. 步骤二
   ...

## 代码示例

实际可用的代码示例

## 最佳实践

注意事项和建议

## 检查清单

- [ ] 检查项1
- [ ] 检查项2
```

### 3. 收集参考信息

创建 skill 前需要收集：

- **现有代码**：查看项目中类似的实现
- **官方文档**：相关技术的官方文档
- **项目规范**：`docs/contributing/code-style.md` 等
- **团队约定**：团队已有的惯例

## 创建 Skill 步骤

### 1. 确定 Skill 名称

命名规范：

- 使用小写字母和连字符
- 动词-名词结构（如 `test-development`）
- 或领域-类型结构（如 `grpc-service`）

```bash
# 好的命名
skill-development
api-documentation
database-migration

# 避免的命名
SkillDevelopment  # 不要使用大驼峰
skill_dev         # 不要使用下划线
create-skill      # 太宽泛
```

### 2. 创建目录和文件

```bash
# 创建 skill 目录
mkdir -p .agent/skills/{skill-name}

# 创建 SKILL.md 文件
touch .agent/skills/{skill-name}/SKILL.md
```

### 3. 编写 YAML Frontmatter

```yaml
---
name: { skill-name } # 必需：与目录名一致
description: 一句话描述 # 必需：简洁明了
---
```

### 4. 编写内容大纲

遵循以下结构：

```markdown
# {标题}

## 概述

- 用途说明
- 适用场景
- 不适用场景（可选）

## 前置条件（可选）

所需的环境、工具、权限等

## 核心概念（可选）

关键术语解释

## 文件结构（如适用）

目录组织方式

## 开发步骤

详细步骤说明

## 代码示例

可直接使用的代码

## 常见问题（可选）

FAQ

## 最佳实践

经验总结

## 参考资源（可选）

相关文档链接

## 检查清单

验证项
```

### 5. 编写代码示例

代码示例原则：

- **可运行**：示例代码应该能直接使用
- **完整**：包含必要的导入和上下文
- **注释**：关键部分添加注释
- **规范**：遵循项目代码规范

### 6. 添加检查清单

每个 skill 都应该有检查清单：

```markdown
## 检查清单

- [ ] 必需的文件已创建
- [ ] 代码遵循项目规范
- [ ] 已添加必要的注释
- [ ] 已通过测试
```

## 批量创建 Skills

当一个场景需要多个相关 skill 时：

### 示例：添加 gRPC 支持

需要创建的 skills：

1. **grpc-service** - gRPC 服务端开发
   - Protocol Buffers 定义
   - 服务实现
   - 拦截器使用

2. **grpc-client** - gRPC 客户端开发
   - 客户端初始化
   - 连接池管理
   - 错误重试

3. **proto-development** - Protobuf 文件规范
   - 消息定义规范
   - 服务定义规范
   - 版本管理

### 创建步骤

```bash
# 1. 创建目录结构
mkdir -p .agent/skills/{grpc-service,grpc-client,proto-development}

# 2. 分别编写各个 skill
# 每个 skill 专注于一个方面

# 3. skill 之间可以互相引用
# 例如在 grpc-service 中提到：
# "参见 proto-development skill 了解 .proto 文件规范"
```

## 质量标准

好的 skill 应该满足：

| 标准         | 描述             | 示例                                  |
| ------------ | ---------------- | ------------------------------------- |
| **明确性**   | 目标清晰，不含糊 | ✅ "创建 gRPC 服务" vs ❌ "使用 gRPC" |
| **完整性**   | 包含所有必要步骤 | 从创建到部署的完整流程                |
| **实用性**   | 可直接应用       | 代码可复制粘贴使用                    |
| **规范性**   | 符合项目规范     | 遵循 `code-style.md`                  |
| **可维护性** | 易于更新         | 结构清晰，章节分明                    |

## 验证新 Skill

创建后需要验证：

```bash
# 1. 检查文件存在
ls -la .agent/skills/{skill-name}/SKILL.md

# 2. 验证 YAML frontmatter
head -5 .agent/skills/{skill-name}/SKILL.md

# 3. 检查内容完整性
# 确保包含：概述、步骤、示例、检查清单
```

## Skill 更新和维护

### 何时更新

- 发现更好的实现方式
- 项目规范变更
- 新版本工具/库发布
- 收到使用反馈

### 更新流程

1. 标记待更新的部分
2. 收集新的实践和示例
3. 更新内容
4. 通知相关开发者

## 检查清单

- [ ] Skill 名称使用小写字母和连字符
- [ ] YAML frontmatter 格式正确
- [ ] 包含清晰的概述
- [ ] 步骤详细且有序
- [ ] 代码示例可运行
- [ ] 有实用的检查清单
- [ ] 符合项目代码规范
- [ ] 与现有 skills 不重复
