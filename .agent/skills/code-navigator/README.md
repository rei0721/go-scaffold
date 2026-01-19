# Code Navigator 索引文件说明

本目录包含项目的结构化索引文件，所有索引均使用 YAML 格式。

## 索引文件列表

### structure.yaml

**项目结构索引** - 完整的目录树和文件组织

### modules.yaml

**模块划分索引** - 模块职责和依赖关系

### dependencies.yaml

**依赖关系索引** - 模块和包的依赖图谱

### layers.yaml

**分层架构索引** - 应用分层和职责划分

### development-paths.yaml

**开发路径索引** - 常见开发场景指引

## YAML 格式说明

所有索引文件遵循统一格式：

```yaml
---
#  Frontmatter（元信息）
name: index-name
description: 简短描述
updated: YYYY-MM-DD
---
# 数据内容（各文件格式不同）
...
```

## 维护指南

1. **更新索引**：项目结构变化时及时更新
2. **更新日期**：修改 `updated` 字段
3. **验证格式**：确保 YAML 语法正确
4. **同步更新**：相关索引需要一起更新

详见主文档 [SKILL.md](../SKILL.md)
