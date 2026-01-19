---
name: skill-maintenance
description: 识别场景并判断何时需要更新或变更已存在的 skills
---

# Skill 维护与更新判断规范

## 概述

本 skill 指导如何识别需要更新 skill 的场景，提供判断标准和更新决策流程，确保 skills 系统始终保持最新和准确。

### 适用场景

- 项目技术栈升级
- 发现更好的实践方式
- skill 内容与实际代码不一致
- 收到使用反馈
- 定期审查 skills 质量

## 更新触发场景

### 场景分类矩阵

| 触发场景           | 严重程度 | 建议行动   | 优先级 |
| ------------------ | -------- | ---------- | ------ |
| 代码规范已变更     | 🔴 高    | 立即更新   | P0     |
| skill 示例代码失效 | 🔴 高    | 立即更新   | P0     |
| 项目架构重构       | 🟡 中    | 尽快更新   | P1     |
| 依赖版本升级       | 🟡 中    | 评估后更新 | P1     |
| 发现更优实践       | 🟢 低    | 计划更新   | P2     |
| 内容表述优化       | 🟢 低    | 批量优化   | P3     |

### 场景分析详细示例

| 场景描述                                                                                                        | 是否需要更新 | 影响的 Skills                                                          | 建议操作                                                              | 优先级 |
| --------------------------------------------------------------------------------------------------------------- | ------------ | ---------------------------------------------------------------------- | --------------------------------------------------------------------- | ------ |
| 项目从 `errors.New()` 迁移到统一的 `pkg/errors` 错误码系统                                                      | ✅           | `service-development`, `handler-development`, `repository-development` | 立即更新所有涉及错误处理的示例代码,添加错误码使用说明                 | P0     |
| 依赖注入模式从 Options 模式改为 `SetXXX` 延迟注入模式                                                           | ✅           | `service-development`, `middleware-development`, `pkg-development`     | 重写构造函数和依赖注入部分的示例,补充新模式的最佳实践                 | P0     |
| `internal/app/app_business.go` 的初始化流程发生重大改变                                                         | ✅           | 所有涉及依赖注入的 skills                                              | 检查并更新所有 skills 中关于应用初始化的说明和引用                    | P1     |
| 新增了 `pkg/logger` 包,要求所有 service 和 handler 统一使用                                                     | ✅           | `service-development`, `handler-development`, `middleware-development` | 在所有相关 skill 中补充 logger 的集成和使用示例                       | P1     |
| 项目从 GORM v1 升级到 GORM v2,部分 API 已变更                                                                   | ✅           | `repository-development`, `model-development`                          | 评估 Breaking Changes,更新数据访问层的示例代码                        | P1     |
| 发现 `handler-development` 中的示例代码缺少请求参数验证                                                         | ✅           | `handler-development`                                                  | 补充参数验证的完整示例和最佳实践                                      | P1     |
| 代码审查发现 `middleware-development` 中的中间件注册方式已过时                                                  | ✅           | `middleware-development`                                               | 更新中间件注册方式,补充新的路由配置说明                               | P0     |
| 项目新增了缓存层(`pkg/cache`),需要在 service 中集成                                                             | ✅           | `service-development`                                                  | 补充缓存集成章节,包括缓存策略和失效策略                               | P1     |
| `docs/architecture/system_map.md` 更新,目录结构描述与实际代码不一致                                             | ✅           | `code-navigator`                                                       | 同步更新 code-navigator 中的目录结构索引                              | P1     |
| 发现了更简洁的事务处理写法,减少了 30% 的样板代码                                                                | ✅           | `repository-development`, `service-development`                        | 评估新写法的优势,如果显著提升效率则更新示例                           | P2     |
| 新人反馈 `test-development` 中的测试用例不够清晰                                                                | ✅           | `test-development`                                                     | 优化示例代码注释,补充常见测试场景的示例                               | P2     |
| `pkg/jwt` 包新增了 `RefreshToken` 功能                                                                          | ✅           | `handler-development` (涉及认证的部分)                                 | 补充 RefreshToken 的使用示例和最佳实践                                | P2     |
| 代码风格指南 `docs/contributing/code-style.md` 新增了注释规范要求                                               | ✅           | 所有包含代码示例的 skills                                              | 按照新规范为所有示例代码添加或优化注释                                | P2     |
| 持续集成 CI 配置更新,新增了代码覆盖率检查要求                                                                   | ✅           | `test-development`                                                     | 补充测试覆盖率相关的说明和目标要求                                    | P2     |
| 项目从 Go 1.19 升级到 Go 1.21,可以使用泛型简化部分代码                                                          | ✅           | `pkg-development`, `repository-development`                            | 评估泛型的适用场景,在合适的地方补充泛型使用示例                       | P2     |
| 发现 `service-development` 中引用的文件路径 `internal/types/service.go` 已移动到 `internal/types/interfaces.go` | ✅           | `service-development`                                                  | 立即修正文件路径引用                                                  | P0     |
| 代码审查反馈,多个开发者按照 `repository-development` 实现时遗漏了连接池配置                                     | ✅           | `repository-development`                                               | 补充连接池配置的检查项和示例                                          | P1     |
| 新增了 `pkg/rbac` 权限控制包,需要在 handler 和 middleware 中集成                                                | ✅           | `handler-development`, `middleware-development`                        | 创建独立的 `rbac-integration` skill 或在现有 skill 中补充权限控制章节 | P1     |
| 发现 `error-handling` skill 中的错误码定义方式与实际代码库不一致                                                | ✅           | `error-handling`                                                       | 立即同步更新错误码的定义方式和示例                                    | P0     |
| 性能测试发现推荐的数据库查询方式存在 N+1 问题                                                                   | ✅           | `repository-development`                                               | 立即更新为预加载(Preload)方式,补充性能优化说明                        | P0     |
| 项目引入了 Swagger 自动生成 API 文档,handler 需要添加特定注释                                                   | ✅           | `handler-development`                                                  | 补充 Swagger 注释规范和示例                                           | P1     |
| 代码 review 发现按照 `workflow-development` 创建的 workflow 缺少 turbo 标记说明                                 | ✅           | `workflow-development`                                                 | 补充 `// turbo` 和 `// turbo-all` 的使用场景和说明                    | P2     |
| 多个 skills 中都提到了配置管理,但描述方式不一致                                                                 | ✅           | 多个 skills                                                            | 统一配置管理的描述,或创建独立的 `config-integration` skill            | P2     |
| 发现 `pkg-development` 中关于接口设计的建议与新的项目规范不符                                                   | ✅           | `pkg-development`                                                      | 更新接口设计部分,对齐最新的接口规范                                   | P1     |
| 项目新增了 `docs/architecture/dependency_map.md` 依赖关系文档                                                   | ✅           | `code-navigator`                                                       | 在 code-navigator 中补充依赖关系导航                                  | P2     |
| 发现使用 `skill-development` 创建的新 skill 普遍缺少"常见问题"章节                                              | ✅           | `skill-development`                                                    | 在模板中补充"常见问题"章节的说明和示例                                | P2     |
| 代码库中实际的测试文件命名从 `_test.go` 改为统一放在 `test/` 目录                                               | ✅           | `test-development`                                                     | 更新测试文件的组织方式和路径说明                                      | P1     |
| 安全审计发现推荐的密码加密方式已不符合最新安全标准                                                              | ✅           | 涉及认证的 skills                                                      | 立即更新为推荐的加密算法(如 bcrypt/argon2)                            | P0     |
| 仅修改了 `service-development` 中的措辞,使表述更清晰                                                            | ✅           | `service-development`                                                  | 微调优化,可以在低优先级时批量处理                                     | P3     |
| Gin 框架从 v1.7 升级到 v1.9,但 API 完全向后兼容,无新特性                                                        | ❌           | -                                                                      | 无需更新 skills                                                       | -      |
| 开发者询问如何实现某功能,但实际上现有 skill 已充分覆盖                                                          | ❌           | -                                                                      | 无需更新,考虑通过培训或 FAQ 改善                                      | -      |
| 某个冷门的 pkg 工具包更新,但项目中几乎不使用                                                                    | ❌           | -                                                                      | 评估使用频率,如不常用则不更新                                         | -      |
| 项目切换到新的日志库(如从 logrus 到 zap),涉及所有代码层                                                         | ✅           | 所有包含日志使用的 skills                                              | 分批次更新所有相关 skills,优先更新最常用的                            | P0/P1  |
| 项目采用了新的目录结构标准(如 feature-based 替代 layer-based)                                                   | ✅           | 几乎所有开发类 skills                                                  | 制定详细的更新计划,这是重大架构变更,需要全面重构 skills 系统          | P0     |
| 发现 `docs-sync` skill 中的检查脚本在 Windows 环境下无法运行                                                    | ✅           | `docs-sync`                                                            | 补充跨平台兼容的脚本版本或说明                                        | P1     |
| 团队决定统一使用 make 命令简化常用操作,新增了 Makefile                                                          | ✅           | `workflow-development` 以及涉及命令执行的 skills                       | 补充 make 命令的使用说明,更新示例命令                                 | P2     |
| 项目接入了 OpenTelemetry 分布式追踪,service 和 handler 需要添加 trace                                           | ✅           | `service-development`, `handler-development`, `middleware-development` | 创建 `distributed-tracing` skill 或在现有 skill 中补充追踪集成        | P1     |
| 代码审查发现 `model-development` skill 中缺少数据库索引的设计指导                                               | ✅           | `model-development`                                                    | 补充索引设计的最佳实践和性能考量                                      | P2     |
| 新人入职时反馈 skills 之间的引用关系不清晰,难以找到相关内容                                                     | ✅           | `skills-map` 或所有 skills                                             | 优化 skills-map 的导航结构,或在各 skill 中补充"相关 Skills"章节       | P2     |
| 敏感操作(如删除用户、清空数据)引入了二次确认机制                                                                | ✅           | `handler-development`, 涉及危险操作的 skills                           | 补充二次确认的实现示例和安全检查清单                                  | P1     |
| 项目开始使用特性开关(Feature Flag)进行灰度发布                                                                  | ✅           | `service-development`, `handler-development`                           | 创建 `feature-flag` skill 或补充特性开关的集成方式                    | P2     |
| 发现三个不同的 skills 用不同方式描述同一个概念(如"业务逻辑层"、"服务层"、"Service")                             | ✅           | 相关的多个 skills                                                      | 统一术语表述,考虑建立项目术语表                                       | P2     |
| 项目启用了严格的代码检查工具(golangci-lint),新增了多条规则                                                      | ✅           | 所有包含代码示例的 skills                                              | 运行 lint 检查所有示例代码,修复不符合规范的部分                       | P1     |
| 性能优化后,推荐的缓存策略从"缓存所有查询"改为"只缓存热点数据"                                                   | ✅           | `service-development` (缓存相关部分), `repository-development`         | 更新缓存策略说明,补充缓存决策的判断标准                               | P2     |
| 发现按照 `work-log` skill 记录的工作日志格式不统一                                                              | ✅           | `work-log`                                                             | 强化日志格式规范,补充标准模板和示例                                   | P2     |
| 项目文档规范要求所有 pkg 必须有 README.md 和 doc.go                                                             | ✅           | `pkg-development`                                                      | 补充文档要求的检查项和模板                                            | P1     |
| 月度审查发现 5 个 skills 引用的示例文件路径已失效                                                               | ✅           | 相关的 5 个 skills                                                     | 批量修复文件路径,考虑添加自动化路径检查                               | P1     |
| 发现新创建的 skills 普遍缺少代码注释,不符合质量标准                                                             | ✅           | `skill-development`, 相关新 skills                                     | 在 `skill-development` 中强调代码注释要求,为已创建的 skills 补充注释  | P2     |
| 项目从单体架构演进为分层架构,层与层之间的调用规范明确                                                           | ✅           | `handler-development`, `service-development`, `repository-development` | 重大架构变更,需要全面更新这些核心 skills,明确层级调用规范             | P0     |
| 引入了响应统一封装格式(如 `{code, message, data}`),所有 API 返回需要遵循                                        | ✅           | `handler-development`                                                  | 更新所有 API 响应示例,补充统一响应格式的说明                          | P0     |
| 新增了请求上下文(Context)传递用户信息的标准方式                                                                 | ✅           | `middleware-development`, `handler-development`, `service-development` | 补充 Context 使用规范和用户信息提取示例                               | P1     |
| 数据库连接改为读写分离,repository 需要区分主从库                                                                | ✅           | `repository-development`                                               | 补充读写分离的实现方式和最佳实践                                      | P1     |
| 静态代码分析工具指出多个 skill 示例代码存在潜在的空指针风险                                                     | ✅           | 相关 skills                                                            | 立即修复示例代码的安全隐患,补充安全编码检查项                         | P0     |
| 定期审查发现某个两年前创建的 skill 已完全过时,不再适用于当前技术栈                                              | ✅           | 该 skill                                                               | 评估是否需要保留,如不需要则归档或删除                                 | P2     |
| 新增了 AI 辅助代码生成工具,可以基于 skills 自动生成样板代码                                                     | ✅           | `skill-development`, 相关 skills                                       | 补充 AI 工具使用说明,确保 skills 格式适合 AI 理解                     | P2     |
| 项目引入容器化部署,开发环境配置方式发生变化                                                                     | ✅           | 涉及环境配置的 skills                                                  | 更新环境配置说明,补充 Docker 相关的开发流程                           | P1     |

### 1. 代码规范变更

**触发条件：**

- 项目代码风格指南更新（`docs/contributing/code-style.md`）
- 新增编码规范或约束
- 废弃某些旧的实践方式

**判断标准：**

```
如果 skill 中的示例代码或指导与新规范冲突：
  → 必须立即更新
否则：
  → 可延后处理
```

**示例：**

```markdown
旧规范：使用 `errors.New()` 创建错误
新规范：统一使用 `pkg/errors` 包的错误码系统

受影响的 skills：

- service-development
- handler-development
- repository-development
- error-handling

更新操作：

1. 更新所有涉及错误处理的示例代码
2. 添加新的错误码使用说明
3. 标记旧方式为已废弃
```

### 2. 项目架构调整

**触发条件：**

- 目录结构重组
- 模块拆分或合并
- 依赖注入方式改变
- 新增或移除核心组件

**判断标准：**

检查以下文件的变更：

- `internal/app/` - 应用初始化逻辑
- `internal/types/` - 核心接口定义
- `pkg/` - 基础组件库
- `docs/architecture/` - 架构文档

**决策树：**

```
架构变更是否影响 skill 中的开发流程？
├─ 是 → 需要更新
│   ├─ 影响范围 > 3个skill → 优先级 P1
│   └─ 影响范围 ≤ 3个skill → 优先级 P2
└─ 否 → 无需更新
```

**示例：**

```markdown
架构变更：从 Options 模式改为 SetXXX 延迟注入

受影响的 skills：

- service-development（依赖注入方式变更）
- middleware-development（中间件初始化变更）
- pkg-development（接口设计模式变更）

更新检查点：
□ 构造函数签名是否变更
□ 依赖注入方式是否调整
□ 示例代码是否需要重写
□ 最佳实践部分是否需要补充
```

### 3. 依赖版本升级

**触发条件：**

- 主要框架版本升级（如 Gin v1 → v2）
- 核心库 API 变更
- 新增重要第三方库

**评估流程：**

```
1. 检查 CHANGELOG
   ├─ 是否有 Breaking Changes？
   │   ├─ 是 → 必须更新 skill
   │   └─ 否 → 继续评估
   └─ 是否有新增推荐用法？
       ├─ 是 → 建议更新 skill
       └─ 否 → 无需更新
```

**检查清单：**

```markdown
依赖升级影响评估

基础信息：

- 包名: \***\*\_\_\_\*\***
- 旧版本: \***\*\_\_\_\*\***
- 新版本: \***\*\_\_\_\*\***

影响分析：

- [ ] API 签名变更
- [ ] 新增推荐功能
- [ ] 废弃旧方法
- [ ] 配置格式变更
- [ ] 性能优化建议

关联 skills：

- [ ] skill-1: \***\*\_\_\_\*\***
- [ ] skill-2: \***\*\_\_\_\*\***

更新决策：

- [ ] 必须更新（Breaking Changes）
- [ ] 建议更新（新增最佳实践）
- [ ] 无需更新（向后兼容且无新特性）
```

### 4. 实际代码与 Skill 不一致

**识别方法：**

**方法 1：代码审查发现**

```bash
# 审查最近的代码变更
git log --since="1 month ago" --name-only --pretty=format: | sort -u

# 对比 skill 中的示例
```

**方法 2：使用反馈**

- 开发者报告 skill 指导与实际不符
- 新人按 skill 操作遇到问题
- CI/CD 检查发现不一致

**方法 3：自动化检测（可选）**

```bash
# 示例：检查 skill 中引用的文件路径是否存在
#!/bin/bash
for skill in .agent/skills/*/SKILL.md; do
  echo "Checking $skill..."
  grep -oP '(?<=`)[^`]+\.go(?=`)' "$skill" | while read file; do
    if [ ! -f "$file" ]; then
      echo "  ⚠️  File not found: $file"
    fi
  done
done
```

**更新优先级：**

```
不一致严重程度评分（每项1分）：
□ skill 中的代码无法运行（+3分）
□ 文件路径已变更（+2分）
□ 接口签名已改变（+2分）
□ 最佳实践已过时（+1分）
□ 仅措辞不准确（+0.5分）

总分 ≥ 3分 → P0 立即修复
总分 2-2.5分 → P1 本周内修复
总分 1-1.5分 → P2 本月内修复
总分 < 1分 → P3 计划优化
```

### 5. 发现更优实践

**判断标准：**

新实践应满足以下条件才值得更新 skill：

- ✅ **已验证**：在项目中实际使用并证明有效
- ✅ **可复用**：适用于多个类似场景
- ✅ **显著改进**：相比旧方式有明显优势
- ✅ **无副作用**：不引入新问题或技术债

**评估表：**

| 评估维度          | 旧实践得分 | 新实践得分 | 改进幅度 |
| ----------------- | ---------- | ---------- | -------- |
| 代码简洁性（1-5） | \_\_\_     | \_\_\_     | \_\_\_   |
| 性能（1-5）       | \_\_\_     | \_\_\_     | \_\_\_   |
| 可维护性（1-5）   | \_\_\_     | \_\_\_     | \_\_\_   |
| 易理解性（1-5）   | \_\_\_     | \_\_\_     | \_\_\_   |
| **总分**          | \_\_\_     | \_\_\_     | \_\_\_   |

```
如果 改进幅度总分 ≥ 5：
  → 值得更新 skill
否则：
  → 保持现状或作为补充说明
```

### 6. 用户反馈

**反馈来源：**

- 开发者提出的疑问
- Code Review 中的讨论
- 新人入职时的困惑
- 项目复盘会议记录

**处理流程：**

```
收到反馈
  ↓
是否有效反馈（真实问题）？
  ├─ 否 → 归档
  └─ 是 → 记录到反馈日志
      ↓
分析根因
  ├─ skill 内容错误 → P0 立即修复
  ├─ skill 内容不完整 → P1 补充说明
  ├─ skill 示例不够清晰 → P2 优化示例
  └─ skill 无问题（用户误解）→ 考虑增加 FAQ
```

**反馈日志模板：**

```markdown
## 反馈记录

### 反馈 #001

- **日期**: 2026-01-19
- **来源**: 开发者反馈
- **相关 Skill**: service-development
- **问题描述**: 依赖注入部分示例代码缺少错误处理
- **严重程度**: 中
- **处理状态**: 待处理
- **计划更新**: 补充错误处理示例

### 反馈 #002

...
```

## 定期审查机制

### 审查周期

| 审查类型 | 频率   | 范围              | 负责人      |
| -------- | ------ | ----------------- | ----------- |
| 快速检查 | 每周   | 近期变更的 skills | 开发 Leader |
| 全面审查 | 每月   | 所有 skills       | 技术团队    |
| 深度优化 | 每季度 | 核心 skills       | 架构师      |

### 每周快速检查

**检查内容：**

```bash
# 1. 检查本周代码变更影响的 skills
git diff --name-only HEAD~7 HEAD | grep -E "internal|pkg" > changed_files.txt

# 2. 查找可能受影响的 skills
# （手动审查或使用脚本）

# 3. 验证示例代码
# （运行 skill 中的示例，确保可用）
```

### 每月全面审查

**审查清单：**

```markdown
## Skills 月度审查

### 基础检查

- [ ] 所有 skill 的 YAML frontmatter 格式正确
- [ ] 所有引用的文件路径有效
- [ ] 所有示例代码可以运行
- [ ] 所有链接可以访问

### 内容检查

- [ ] 内容与当前代码库一致
- [ ] 最佳实践仍然适用
- [ ] 没有过时的技术栈引用
- [ ] 检查清单完整且准确

### 质量检查

- [ ] 语言表达清晰准确
- [ ] 代码格式规范统一
- [ ] 示例代码有足够注释
- [ ] 结构层次清晰

### 覆盖度检查

- [ ] 是否有新的开发场景未覆盖
- [ ] 是否需要新增 skill
- [ ] 是否有 skill 已不再需要
```

### 季度深度优化

**优化方向：**

1. **结构优化**
   - 合并相似内容的 skills
   - 拆分过于复杂的 skills
   - 调整 skill 之间的引用关系

2. **内容优化**
   - 更新过时的示例
   - 补充新的最佳实践
   - 增加常见问题解答

3. **形式优化**
   - 统一术语和表达
   - 改进代码示例的注释
   - 优化排版和格式

## 更新决策流程

### 决策树

```
发现潜在更新需求
  ↓
评估影响范围
  ├─ 影响核心开发流程？
  │   ├─ 是 → P0/P1 高优先级
  │   └─ 否 → 继续评估
  ↓
评估紧迫性
  ├─ 当前 skill 是否会误导？
  │   ├─ 是 → P0 立即更新
  │   └─ 否 → 继续评估
  ↓
评估改进价值
  ├─ 更新后能显著提升效率？
  │   ├─ 是 → P1/P2 计划更新
  │   └─ 否 → P3 延后或不更新
```

### 更新成本评估

**估算更新工作量：**

| 更新类型          | 预计时间 | 示例         |
| ----------------- | -------- | ------------ |
| 修正错误路径/命令 | 10分钟   | 文件路径变更 |
| 更新单个示例代码  | 30分钟   | API 签名调整 |
| 重写整个章节      | 1-2小时  | 架构模式变更 |
| 全面重构 skill    | 3-4小时  | 技术栈升级   |

**ROI 计算：**

```
更新价值分 = (受影响开发者数量 × 每次使用节省时间 × 预计使用次数)
更新成本分 = 更新所需时间

如果 更新价值分 / 更新成本分 > 5：
  → 高优先级更新
如果 更新价值分 / 更新成本分 > 2：
  → 适时更新
否则：
  → 低优先级或不更新
```

## 更新操作指南

### 更新流程

```
1. 创建更新计划
   ├─ 明确更新范围
   ├─ 列出具体变更点
   └─ 预估工作量

2. 执行更新
   ├─ 备份原 skill（可选）
   ├─ 修改内容
   └─ 验证示例代码

3. 质量检查
   ├─ 自查检查清单
   ├─ 运行示例代码
   └─ 同行审阅（重要更新）

4. 记录变更
   ├─ 在 skill 顶部添加更新记录（可选）
   └─ 在变更日志中记录

5. 通知相关人员
   ├─ 重要更新：通知全部开发者
   └─ 一般更新：在周会中提及
```

### 更新记录格式（可选）

如果 skill 更新频繁，可在文件顶部添加更新记录：

```markdown
---
name: service-development
description: 在 internal/service/ 目录下创建新的业务服务
---

> **更新记录**
>
> - 2026-01-19: 调整依赖注入方式，从 Options 改为 SetXXX 模式
> - 2026-01-15: 补充错误处理示例
> - 2026-01-10: 更新缓存集成说明

# 服务开发规范

...
```

### 批量更新策略

当多个 skills 需要同类更新时：

```
1. 识别更新模式
   示例：所有 skills 都需要更新错误处理方式

2. 创建更新模板
   示例：标准的错误处理代码块

3. 批量应用
   ├─ 使用脚本辅助（如果可行）
   └─ 或逐个手动更新

4. 批量验证
   └─ 确保所有更新一致且正确
```

## 版本控制（可选）

### 是否需要版本号

**适用场景：**

- skills 系统非常成熟稳定
- 有多人协作维护
- 需要追踪重大变更

**版本号格式：**

```
v{Major}.{Minor}.{Patch}

Major: 重大架构变更，不兼容旧版本
Minor: 功能新增或重要更新
Patch: 错误修正或小幅优化
```

**示例：**

```markdown
---
name: service-development
description: 在 internal/service/ 目录下创建新的业务服务
version: 2.1.0
---
```

### 不使用版本号

**大多数情况推荐：**

- 直接更新 skill 内容
- 通过 Git 历史追踪变更
- 保持 skills 简单易维护

## 工具辅助

### 一致性检查脚本

```bash
#!/bin/bash
# check-skills-consistency.sh

echo "🔍 检查 Skills 一致性..."

# 1. 检查引用的文件是否存在
echo "1. 检查文件路径..."
for skill in .agent/skills/*/SKILL.md; do
  grep -oP '(?<=`)[^`]+\.(go|md)(?=`)' "$skill" | while read file; do
    if [ ! -f "$file" ]; then
      echo "  ⚠️  [$skill] 引用的文件不存在: $file"
    fi
  done
done

# 2. 检查 YAML frontmatter
echo "2. 检查 YAML frontmatter..."
for skill in .agent/skills/*/SKILL.md; do
  if ! head -5 "$skill" | grep -q "^name:"; then
    echo "  ⚠️  [$skill] 缺少 name 字段"
  fi
  if ! head -5 "$skill" | grep -q "^description:"; then
    echo "  ⚠️  [$skill] 缺少 description 字段"
  fi
done

# 3. 检查是否有过时的包引用
echo "3. 检查过时的包引用..."
DEPRECATED_PACKAGES=("github.com/pkg/errors")
for skill in .agent/skills/*/SKILL.md; do
  for pkg in "${DEPRECATED_PACKAGES[@]}"; do
    if grep -q "$pkg" "$skill"; then
      echo "  ⚠️  [$skill] 使用了过时的包: $pkg"
    fi
  done
done

echo "✅ 检查完成"
```

### 更新提醒配置

```yaml
# .agent/skills-maintenance.yaml（可选配置文件）

# 审查提醒
reviews:
  weekly: true
  monthly: true
  quarterly: true

# 自动检查
auto_checks:
  - file_existence
  - yaml_format
  - deprecated_patterns

# 关注的代码路径
watched_paths:
  - internal/
  - pkg/
  - docs/architecture/

# 当这些文件变更时，触发 skill 审查提醒
triggers:
  - docs/contributing/code-style.md
  - internal/types/interfaces.go
  - go.mod
```

## 最佳实践

### 1. 预防性维护

**主动识别更新需求：**

- ✅ 代码变更时立即思考是否影响 skills
- ✅ 新技术引入时同步创建/更新 skill
- ✅ 定期审查而不是等问题出现

### 2. 渐进式更新

**避免一次性大规模重构：**

- ✅ 优先更新最常用的 skills
- ✅ 分批次逐步改进
- ✅ 每次更新后收集反馈

### 3. 保持简洁

**避免过度设计：**

- ✅ 不需要的 skill 及时删除
- ✅ 相似内容考虑合并
- ✅ 说明简洁明了，避免冗长

### 4. 文档同步

**确保多处文档一致：**

```
更新 skill 时，同步检查：
  ├─ docs/architecture/ - 架构文档
  ├─ docs/contributing/ - 贡献指南
  ├─ README.md - 项目说明
  └─ pkg/*/README.md - 组件文档
```

### 5. 用户视角

**从使用者角度评估：**

- ✅ 新手能否理解
- ✅ 示例是否足够清晰
- ✅ 步骤是否容易遵循

## 检查清单

### 更新决策检查

- [ ] 已识别更新触发场景
- [ ] 已评估影响范围和优先级
- [ ] 已估算更新成本和价值
- [ ] 已确定是否需要更新

### 更新执行检查

- [ ] 已创建更新计划
- [ ] 已更新相关内容
- [ ] 所有示例代码已验证
- [ ] 已通过检查清单自查
- [ ] 已记录变更（如需要）

### 更新后验证

- [ ] skill 内容与代码库一致
- [ ] 示例代码可运行
- [ ] 文件路径引用正确
- [ ] 相关文档已同步更新
- [ ] 已通知相关人员（如需要）
