# AI 协作提示词 (AI Prompt)

> **最高指令 (Prime Directive)**  
> 任务开始前务必阅读本文件。本文件包含核心逻辑，优先级最高。

---

## 项目元信息

- **项目名称**: go-scaffold
- **代码仓库**: `github.com/rei0721/rei0721`
- **作者**: Rei (`github.com/rei0721`)
- **协议版本**: IEP v7.2 (工业级工程协作协议)
- **当前版本**: v0.1.2

---

## 核心开发准则

### 1. 地图先行 (Map First, Code Follows)

在进行任何代码修改前，必须：

1. **阅读架构地图**: [`docs/architecture/system_map.md`](/docs/architecture/system_map.md)
2. **检索变量索引**: [`docs/architecture/variable_index.md`](/docs/architecture/variable_index.md)
3. **查阅包文档**: 每个 `pkg/*` 模块都有详细的 `doc.go`

### 2. 代码即真相 (Code is Truth)

当引入新库、重构或变更接口时，必须**在同一原子操作中**同步更新：

- 架构地图 (`system_map.md`)
- 变量索引 (`variable_index.md`)
- 相关模块文档

### 3. 变量命名宪法 (Variable Naming Constitution)

- **先检索，后复用**: 在创建新常量前，必须检索 `variable_index.md`
- **禁止同义词**: 如果已存在语义相同的变量，必须复用，不得造新词
- **立即登记**: 新增全局常量后，必须立即更新 `variable_index.md`

### 4. 脚手架兼容优先 (Scaffold Compatibility First)

本项目具有成熟的工程化体系，任何修改必须：

- ✅ 保持现有目录结构 (`cmd/`, `internal/`, `pkg/`, `types/`, `configs/`)
- ✅ 遵循现有设计模式（DI 容器、热重载、接口抽象）
- ✅ 延续现有代码风格（见 `docs/memories/rules.md`）
- ❌ 不强行引入与项目风格冲突的工具或模式

---

## 技术栈约束

### 语言与框架

- **Go**: 1.24.6+
- **HTTP 框架**: Gin
- **ORM**: GORM
- **配置管理**: Viper
- **日志**: Zap
- **缓存**: Redis (go-redis v9)
- **协程池**: ants v2

### 架构模式

- **依赖注入**: 通过 `internal/app.App` 容器管理生命周期
- **接口抽象**: 所有基础设施组件使用接口（`database.Database`, `cache.Cache`, `executor.Manager`, `logger.Logger`）
- **配置热重载**: 支持运行时更新配置（实现 `Reloader` 接口）

---

## 标准作业程序 (SOP)

所有任务必须按以下顺序执行：

### 1. 认知 (Understand)

- 阅读 `system_map.md` 了解全局
- 阅读 `variable_index.md` 了解命名规范
- 阅读相关模块的 `doc.go`

### 2. 推演 (Design)

- 在 `specs/` 目录创建临时推演文档（如 `specs/add_xxx_feature.md`）
- 编写伪代码、选择依赖库、列出风险
- **不允许跳过此步骤直接编码**

### 3. 施工 (Implement)

- 实现代码，确保可编译、可运行
- 遵循 Go 编码规范（`gofmt`, `golint`）
- Context 传递、错误处理必须显式

### 4. 测绘 (Document)

- 同步更新架构地图
- 同步更新变量索引
- 清理 `specs/` 临时文件
- 将可复用经验写入 `docs/memories/universal_prompt.md`

---

## 质量标准

### Go 特别条款

1. **Context 透传**: 所有 I/O 链路必须透传 `context.Context`
2. **错误处理**: 必须显式处理错误，禁止 `_ = err`
3. **依赖注入**: 优先使用接口，便于测试和替换
4. **包文档**: 每个包必须有 `doc.go`，遵循 GoDoc 规范

### 通用规则

1. **IPO 模型**: Input → Process → Output，业务逻辑尽量无副作用
2. **胶水编程**: 能连不造，能抄不写，第三方库只写 Adapter 不改源码
3. **测试优先**: 公开接口必须有单元测试（`*_test.go`）

---

## 记忆体系

### L1 核心记忆 (Core Memory)

- `docs/architecture/` - 系统事实与约束
- 必须始终保持与代码同步

### L2 工作记忆 (Working Memory)

- `specs/` - 临时推演过程
- 任务结束必须清理

### L3 辅助记忆 (Auxiliary Memory)

- `docs/memories/` - 协作偏好与隐性约束
- 持续沉淀可复用规则

---

## 与 IEP v7.2 协议的衔接

本文件是 **IEP v7.2 标准工程模式**的项目级实例化配置。

- **协议地位**: 本文件内容与 IEP v7.2 协议具有同等效力
- **冲突解法**: 当本文件与协议存在冲突时，以本文件为准（适配项目特性）
- **持续演进**: 本文件会随项目发展持续更新

---

## 快速链接

- 📖 [架构地图](/docs/architecture/system_map.md)
- 📋 [变量索引](/docs/architecture/variable_index.md)
- 🧠 [协作规则](/docs/memories/rules.md)
- 🚀 [快速开始](/docs/guides/getting-started.md)

---

> **提醒**: 每次开始任务时，请先阅读本文件和架构地图，确保理解项目上下文。
