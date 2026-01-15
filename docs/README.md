# 📚 项目文档中心

> **单一事实来源** — 本目录是 go-scaffold 项目的文档根目录。

---

## 🗂️ 目录结构

```text
docs/
├── ai_prompt.md           # 🤖 AI 协作提示词（最高优先级）
├── architecture/          # 🏗️ 架构文档
│   ├── system_map.md      # 系统架构拓扑图与模块依赖
│   └── variable_index.md  # 全局变量命名索引表
├── guides/                # 📖 操作手册
│   └── getting-started.md # 快速开始指南
├── features/              # 📋 功能规格
│   └── README.md          # PRD、验收标准
├── integrations/          # 🔌 API 与集成契约
│   └── README.md          # 第三方系统集成文档
├── memories/              # 🧠 AI 协作记忆
│   ├── rules.md           # 项目特定的 AI 行为准则
│   ├── context.md         # 当前迭代的背景信息
│   └── universal_prompt.md # 通用记忆提示词
├── incidents/             # 🚨 事故复盘
│   └── ...                # 事故档案
└── archive/               # 📦 已废弃文档
    └── ...                # 历史设计与文档
```

---

## 🚀 快速导航

### 必读文档

| 文档                                                  | 说明            | 阅读时机           |
| ----------------------------------------------------- | --------------- | ------------------ |
| [ai_prompt.md](./ai_prompt.md)                        | AI 协作最高指令 | **每次任务开始前** |
| [system_map.md](./architecture/system_map.md)         | 系统架构全景图  | 理解项目全貌       |
| [variable_index.md](./architecture/variable_index.md) | 变量命名规范    | 创建新常量前       |

### 开发指南

| 文档                                              | 说明               |
| ------------------------------------------------- | ------------------ |
| [getting-started.md](./guides/getting-started.md) | 环境搭建与快速运行 |

### AI 协作记忆

| 文档                                                  | 说明                 |
| ----------------------------------------------------- | -------------------- |
| [rules.md](./memories/rules.md)                       | 项目特定的编码规则   |
| [context.md](./memories/context.md)                   | 当前迭代上下文       |
| [universal_prompt.md](./memories/universal_prompt.md) | 跨任务沉淀的通用规则 |

---

## 📋 记忆层级体系

根据 **IEP v7.2 协议**，文档分为三个记忆层级：

### L1 核心记忆 (Core Memory)

> 定义系统事实与约束，**必须与代码保持同步**

- `architecture/system_map.md` — 架构拓扑
- `architecture/variable_index.md` — 命名规范

### L2 工作记忆 (Working Memory)

> 承载推演过程，**任务结束必须清理**

- `specs/` 目录（位于项目根目录）

### L3 辅助记忆 (Auxiliary Memory)

> 沉淀协作偏好与隐性约束，**持续积累**

- `memories/rules.md` — 项目规则
- `memories/context.md` — 迭代上下文
- `memories/universal_prompt.md` — 通用提示词

---

## ⚡ 标准作业程序 (SOP)

所有任务按以下顺序执行：

1. **认知** — 阅读 `system_map.md` 与 `variable_index.md`
2. **推演** — 在 `specs/` 创建临时推演文档
3. **施工** — 实现代码，确保可运行
4. **测绘** — 同步更新文档，清理 `specs/`

---

## 🔗 相关链接

- **项目仓库**: [github.com/rei0721/rei0721](https://github.com/rei0721/rei0721)
- **协议版本**: IEP v7.2 (工业级工程协作协议)

---

> 💡 **提示**: 如需新增文档，请根据内容类型放入对应子目录，并更新本概述文件。
