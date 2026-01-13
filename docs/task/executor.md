### 通用并发调度组件 (General-Purpose Async Executor)

#### 1. 组件定位

为当前系统 `rei0721` 定制，在 `pkg` 基础库层封装基于 `ants` 的**通用并发任务调度器**。该组件旨在作为项目级的底层并发基础设施，提供统一的、资源隔离的异步任务执行环境，适用于 HTTP 服务、RPC 服务、CLI 工具及后台 Job 等多种业务场景。

#### 2. 核心能力 (Core Capabilities)

- **多维资源隔离 (Multi-dimensional Isolation)**：
  支持根据业务特征（如 `IO密集`、`CPU密集`、`高优先级`、`后台批处理`）配置多个独立的协程池，实现“舱壁模式”，防止单一类别的任务积压导致整个应用瘫痪。
- **全链路安全防护 (Safety Net)**：
  内置无死角的 Panic 捕获（Recover）与堆栈快照记录能力，确保任意异步任务的异常都不会导致主进程（Process）崩溃。
- **动态热重载 (Dynamic Hot-Reload)**：
  支持在应用运行时（Runtime）动态更新协程池策略（如扩缩容），通过原子切换机制实现配置变更的“即时生效”，无需重启进程。
- **标准化集成 (Standard Integration)**：
  提供面向接口（Interface-Oriented）的编程范式，适配依赖注入（DI）架构，支持 Mock 测试，确保业务逻辑与并发实现层完全解耦。

#### 3. 设计约束 (Design Constraints)

> 为了保证组件的通用性与稳定性，需遵守以下约束：

**3.1 架构规范**

- **包路径**：`pkg/executor`。
- **零全局依赖**：组件必须设计为**无状态、无全局变量**的实例结构。严禁使用 `init()` 函数初始化全局池，必须通过 `NewManager(config)` 显式实例化。
- **接口定义**：
  对外暴露的核心接口必须包含以下能力，且不依赖任何具体业务逻辑：

```go
type Manager interface {
    // Execute 投递任务到指定名称的池中
    Execute(poolName PoolName, task func()) error
    // Reload 根据新配置热重载所有池，实现零停机配置变更
    Reload(configs []Config) error
    // Shutdown 优雅关闭，等待积压任务处理完毕
    Shutdown()
}

```

**3.2 配置与扩展**

- **配置结构化**：配置项应包含 `Name` (唯一标识), `Size` (容量), `Expiry` (回收时间), `NonBlocking` (是否阻塞)。
- **常量约束**：池的标识符（PoolName）应在业务层定义为常量，组件层仅负责透传，不硬编码任何具体业务名称（如不应在 pkg 里写 "SendEmail" 这种具体业务名）。
- **默认策略**：推荐默认采用 Non-blocking（非阻塞）模式，当资源耗尽时返回标准错误 `ErrPoolOverload`，由调用方决定是重试、同步执行还是丢弃，以适应不同场景（如 HTTP 需丢弃，CLI 可能需阻塞等待）。

**3.3 并发安全与热更机制**

- **原子性保证**：`Reload` 接口必须保证并发安全。在切换配置时，应采用**Copy-On-Write (COW)** 或 **读写锁** 机制，确保在重载瞬间正在提交的任务不会丢失或 Panic。
- **平滑过渡**：热重载发生时，旧的协程池必须执行“优雅退出”流程（停止接收新任务 -> 等待存量任务完成 -> 释放内存），严禁直接 Close 导致正在运行的协程被强制中断。

**3.4 依赖注入 (DI) 规范**

- **接口隔离**：应用层（`internal`）的所有模块（Service/Biz）仅依赖 `pkg/executor` 接口。
- **生命周期托管**：组件的生命周期（Start/Stop）应由应用的**主引导程序** `internal/app` 统一管理，确保在应用退出的最后阶段才关闭并发池。
