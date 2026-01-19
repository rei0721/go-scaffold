# Create New Skills System

## 任务概述

创建两个新的 skills：`skills-map`（技能地图）和 `code-navigator`（代码导航系统），用于增强项目的可导航性和可维护性。

## 完成内容

### skills-map（技能地图）

- ✅ 创建 `.agent/skills/skills-map/` 目录结构
- ✅ 编写 `SKILL.md` 主文档（10,657 bytes）
  - Skills 分类体系（开发/管理/工具）
  - Mermaid 全景图
  - 完整的 skills 清单表格
  - 快速索引和查找
  - 变更历史记录
  - 维护指南
- ✅ 创建 `templates/skill-entry.md` 模板文件

**核心功能**：

- 📊 可视化展示所有 13 个 skills
- 🏷️ 按功能分类组织
- 📝 提供多维度快速索引
- 🔄 记录 skills 变更历史

### code-navigator（代码导航系统）

- ✅ 创建 `.agent/skills/code-navigator/` 目录结构
- ✅ 编写 `SKILL.md` 主文档（7,591 bytes）
- ✅ 创建 `README.md` 索引说明文件
- ✅ 创建 `indices/` 目录，包含 5 个 YAML 索引文件：

#### YAML 索引文件

1. **structure.yaml**（5,051 bytes）
   - 完整的项目目录树
   - 每个目录的用途说明
   - 关键文件标注
   - 文件组织规则

2. **modules.yaml**（3,924 bytes）
   - 模块划分清单
   - 模块职责定义
   - 模块依赖关系
   - 入口文件列表

3. **dependencies.yaml**（3,685 bytes）
   - 模块依赖图谱
   - Go 包引用关系
   - 第三方依赖列表
   - 依赖规则说明

4. **layers.yaml**（5,579 bytes）
   - 分层架构定义
   - 各层职责划分
   - 层间通信规则
   - 文件分布详情

5. **development-paths.yaml**（10,537 bytes）
   - 9 个常见开发场景
   - 每个场景的完整步骤
   - 关联的 skills 引用
   - 开发最佳实践

**核心功能**：

- 🗂️ 多文件索引架构
- 🔗 YAML 格式结构化数据
- 🎯 多维度快速定位
- 📋 按需加载索引

## 关键文件

### skills-map

- `.agent/skills/skills-map/SKILL.md` - 主文档
- `.agent/skills/skills-map/templates/skill-entry.md` - 条目模板

### code-navigator

- `.agent/skills/code-navigator/SKILL.md` - 主文档
- `.agent/skills/code-navigator/README.md` - 索引说明
- `.agent/skills/code-navigator/indices/structure.yaml` - 项目结构
- `.agent/skills/code-navigator/indices/modules.yaml` - 模块划分
- `.agent/skills/code-navigator/indices/dependencies.yaml` - 依赖关系
- `.agent/skills/code-navigator/indices/layers.yaml` - 分层架构
- `.agent/skills/code-navigator/indices/development-paths.yaml` - 开发路径

## 设计亮点

### 1. skills-map - 可视化地图

**Mermaid 图表**：

- 使用流程图展示 skills 分类和关系
- 清晰的视觉层次
- 快速了解 skills 体系

**多维度索引**：

- 按开发场景查找
- 按层级查找
- 按功能分类

### 2. code-navigator - 多文件索引

**YAML 元信息**：

```yaml
---
name: index-name
description: 简短描述
updated: YYYY-MM-DD
---
```

**结构化数据**：

- 易于解析和查询
- 支持程序化访问
- 便于维护更新

**按需加载**：

- 5 个独立索引文件
- 只加载需要的信息
- 提高查询效率

## 使用场景

### 新开发者入职

1. 查看 `skills-map` 了解 skills 体系
2. 查看 `code-navigator/indices/structure.yaml` 了解项目结构
3. 查看 `code-navigator/indices/layers.yaml` 理解架构分层

### 日常开发

1. 在 `skills-map` 快速查找所需 skill
2. 在 `development-paths.yaml` 查找开发场景
3. 按场景步骤进行开发

### Skills 维护

1. 新增 skill 后更新 `skills-map`
2. 更新 Mermaid 图表和表格
3. 记录变更历史

## 技术决策

### 为什么使用 YAML?

**优势**：

- 结构化数据，易于解析
- ✅ 人类可读，便于编辑
- ✅ 支持 frontmatter 元信息
- ✅ 可以被程序读取和处理

### 为什么分多个索引文件？

**优势**：

- ✅ 关注点分离，职责单一
- ✅ 按需加载，减少信息过载
- ✅ 独立维护，降低耦合
- ✅ 便于扩展新的索引维度

### 为什么使用 Mermaid 图表？

**优势**：

- ✅ 可视化，直观清晰
- ✅ 代码即文档，易于维护
- ✅ 支持版本控制
- ✅ 在 Markdown 中原生渲染

## 经验总结

### 1. 设计模式

**分离关注点**：

- skills-map 关注"有什么"
- code-navigator 关注"在哪里"

**多维度索引**：

- 不同用户有不同查找习惯
- 提供多种查找方式提高易用性

### 2. 数据组织

**结构化优于自由文本**：

- YAML 比纯文本更易于维护
- 结构化数据支持程序化处理

**元信息的重要性**：

- 每个索引都有 name, description, updated
- 便于追踪和管理

### 3. 文档维护

**强制更新机制**：

- skills-map 明确规定何时必须更新
- 通过 work-log 强制记录变更

**模板化**：

- 提供模板降低添加新内容的成本
- 统一格式提高一致性

---

**Git Commit Message**:

```
feat(skills): add skills-map and code-navigator

Created two new navigation skills:
- skills-map: visual map of all skills with categorization
- code-navigator: multi-file YAML index system for project structure

Features:
- Mermaid diagram for skills visualization
- 5 YAML indices (structure, modules, dependencies, layers, paths)
- Quick reference tables and development scenarios
- Maintenance guidelines

See docs/worklogs/2026/01/19_create_new_skills_system.md
```
