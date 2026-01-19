---
name: code-navigator
description: 代码导航系统，通过多文件索引快速定位开发方向
---

# 代码导航系统

## 概述

本导航系统通过多个 Markdown + YAML Frontmatter 索引文件提供项目的结构化信息，支持按需查询和快速定位开发方向。

## 索引架构

```
.agent/skills/code-navigator/
├── SKILL.md                    # 本文档（使用说明）
├── README.md                   # 索引文件详细说明
└── indices/                    # Markdown 索引目录
    ├── structure.md            # 项目结构索引
    ├── modules.md              # 模块划分索引
    ├── dependencies.md         # 依赖关系索引
    ├── layers.md               # 分层架构索引
    └── development-paths.md    # 开发路径索引
```

## 索引文件说明

### 📁 structure.md - 项目结构索引

**用途**：完整的目录结构树，包含每个目录的用途说明

**包含信息**：

- 目录层级关系
- 每个目录的职责描述
- 关键文件标注
- 文件组织规则

**查看**：[indices/structure.md](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/code-navigator/indices/structure.md)

---

### 🧩 modules.md - 模块划分索引

**用途**：项目的模块划分，每个模块的职责和入口文件

**包含信息**：

- 模块列表和路径
- 每个模块的职责定义
- 模块的入口文件
- 模块间依赖关系

**查看**：[indices/modules.md](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/code-navigator/indices/modules.md)

---

### 🔗 dependencies.md - 依赖关系索引

**用途**：模块间的依赖关系，包引用关系

**包含信息**：

- 模块依赖图谱
- Go 包导入关系
- 第三方依赖列表
- 内部包引用

**查看**：[indices/dependencies.md](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/code-navigator/indices/dependencies.md)

---

### 📚 layers.yaml - 分层架构索引

**用途**：应用的分层架构，每层的职责和文件分布

**包含信息**：

- 架构层级定义
- 每层的职责说明
- 文件分布情况
- 层间通信规则

**查看**：[indices/layers.md](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/code-navigator/indices/layers.md)

---

### 🛤️ development-paths.md - 开发路径索引

**用途**：常见开发场景的完整路径指引

**包含信息**：

- 常见开发场景列表
- 每个场景的完整步骤
- 需要修改的文件
- 关联的 skills 引用

**查看**：[indices/development-paths.md](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/code-navigator/indices/development-paths.md)

## 使用方法

### 场景 1：我想了解项目整体结构

```bash
# 查看项目结构索引
cat .agent/skills/code-navigator/indices/structure.md
```

→ 获取完整的目录树和每个目录的用途

### 场景 2：我想知道某个模块的职责

```bash
# 查看模块划分索引
cat .agent/skills/code-navigator/indices/modules.md
```

→ 了解各个模块的职责、入口文件和依赖

### 场景 3：我想了解模块间的依赖关系

```bash
# 查看依赖关系索引
cat .agent/skills/code-navigator/indices/dependencies.md
```

→ 查看模块依赖图谱和包引用关系

### 场景 4：我想了解架构分层

```bash
# 查看分层架构索引
cat .agent/skills/code-navigator/indices/layers.md
```

→ 了解应用的分层架构和职责划分

### 场景 5：我想开发新功能，不知道从哪里开始

```bash
# 查看开发路径索引
cat .agent/skills/code-navigator/indices/development-paths.md
```

→ 找到对应场景的完整开发路径

## 快速查询命令

### PowerShell 快速查询

```powershell
# 函数：查询项目结构
function Get-ProjectStructure {
    Get-Content .agent/skills/code-navigator/indices/structure.md
}

# 函数：查询模块信息
function Get-Modules {
    Get-Content .agent/skills/code-navigator/indices/modules.md
}

# 函数：查询依赖关系
function Get-Dependencies {
    Get-Content .agent/skills/code-navigator/indices/dependencies.md
}

# 函数：查询开发路径
function Get-DevPaths {
    param([string]$Scenario)
    $content = Get-Content .agent/skills/code-navigator/indices/development-paths.md -Raw
    if ($Scenario) {
        $content | Select-String -Pattern "name: $Scenario" -Context 0,10
    } else {
        $content
    }
}
```

**使用示例**：

```powershell
Get-ProjectStructure
Get-Modules
Get-DevPaths -Scenario "add-new-api"
```

## 多维度定位

### 按功能定位

**需求**：我要实现用户认证功能

**步骤**：

1. 查看 `development-paths.md` 找到 `add-authentication` 场景
2. 按照步骤依次查看需要修改的文件
3. 参考关联的 skills 进行开发

### 按层级定位

**需求**：我要在 Service 层添加逻辑

**步骤**：

1. 查看 `layers.md` 找到 Business 层定义
2. 查看该层的文件分布（`internal/service/`）
3. 参考 `service-development` skill 进行开发

### 按模块定位

**需求**：我要修改用户模块

**步骤**：

1. 查看 `modules.md` 找到 user 模块
2. 查看模块的入口文件和依赖
3. 定位到具体文件进行修改

## 索引文件格式说明

所有索引文件均使用 **Markdown + YAML Frontmatter 格式**，包含以下结构：

### Markdown + YAML Frontmatter

每个Markdown索引文件开头都有 YAML Frontmatter：

```yaml
---
name: index-name # 索引名称
description: 索引描述 # 简短描述
updated: YYYY-MM-DD # 最后更新日期
---
```

### 数据结构

参见 [README.md](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/code-navigator/README.md) 了解各索引文件的详细数据结构。

## 维护规则

### 何时更新索引

**必须更新**：

- ✅ 新增目录或模块
- ✅ 删除目录或模块
- ✅ 架构调整
- ✅ 依赖关系变更

**建议更新**：

- 📝 目录职责调整
- 📝 模块职责优化
- 📝 开发路径优化
- 📝 定期审查（每月）

### 更新对应关系

| 变更类型      | 需要更新的索引                  |
| ------------- | ------------------------------- |
| 新增/删除目录 | `structure.md`                  |
| 新增/删除模块 | `modules.md`, `structure.md`    |
| 模块依赖变化  | `dependencies.md`, `modules.md` |
| 架构分层调整  | `layers.md`, `modules.md`       |
| 新增开发场景  | `development-paths.md`          |

### 更新步骤

1. **修改对应的 YAML 文件**
2. **更新 `updated` 字段**为当前日期
3. **验证 YAML 格式**是否正确
4. **测试查询功能**是否正常

## 最佳实践

### 新开发者入职

1. 先查看 `structure.md` 了解项目整体结构
2. 然后查看 `layers.md` 理解架构分层
3. 最后查看 `modules.md` 了解模块划分

### 日常开发

1. 根据需求查看 `development-paths.md` 找到开发路径
2. 参考路径中的文件列表和关联 skills
3. 查看 `modules.md` 了解相关模块的依赖

### 架构优化

1. 查看 `dependencies.md` 分析当前依赖关系
2. 识别循环依赖和不合理的依赖
3. 重构后更新所有相关索引

## 相关 Skills

- **[skills-map](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/skills-map/SKILL.md)** - 了解所有可用的 skills
- **[skill-development](file:///d:/coder/go/go-scaffold/main/go-scaffold/.agent/skills/skill-development/SKILL.md)** - 创建新的 development skill

---

**最后更新**：2026-01-19
