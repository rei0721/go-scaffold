/*
Package utils 提供常用的工具函数和组件

# 概述

utils 包是一个通用工具库，包含多个独立的工具函数，用于解决常见的技术问题。
每个工具都经过精心设计，线程安全，可以在生产环境中放心使用。

# 包含的工具

1. Snowflake ID 生成器 - 分布式唯一 ID 生成
2. IP 地址验证 - HTTP 监听地址合法性验证
3. 设备 ID 生成 - 基于硬件信息的设备指纹
4. 端口查找 - 自动查找可用 TCP 端口

# 使用示例

## Snowflake ID 生成器

生成分布式唯一 ID，适用于数据库主键、用户 ID、订单号等场景。

	import "github.com/rei0721/go-scaffold/pkg/utils"

	// 单机环境：使用默认生成器
	gen := utils.DefaultSnowflake()
	id := gen.NextID()           // 1747234567890123456
	idStr := gen.NextIDString()  // "1747234567890123456"

	// 分布式环境：指定 nodeID（0-1023）
	gen, err := utils.NewSnowflake(nodeID)
	if err != nil {
		log.Fatal(err)
	}

Snowflake ID 结构（64位）:
  - 1 位: 未使用（始终为0）
  - 41 位: 时间戳（毫秒级，可用约 69 年）
  - 10 位: 节点 ID（支持 1024 个节点）
  - 12 位: 序列号（每毫秒可生成 4096 个 ID）

优点:
  - 时间递增：按生成时间排序
  - 分布式唯一：不同节点不会冲突
  - 高性能：每毫秒可生成 4096 个 ID
  - 无依赖数据库：在应用层生成

注意事项:
  - 相同 nodeID 的实例会生成冲突的 ID
  - 时钟回拨会导致 ID 重复，生产环境需使用 NTP 同步时间
  - 分布式环境中每个实例必须使用不同的 nodeID

## IP 地址验证

验证 HTTP 监听地址是否合法且可被绑定。

	import "github.com/rei0721/go-scaffold/pkg/utils"

	// 方式1：逻辑验证
	err := utils.IsValidListenAddr(":8080")
	if err != nil {
		log.Fatal("invalid listen address:", err)
	}

	// 方式2：实际绑定测试（更严格）
	err = utils.IsValidHTTPListenAddr("127.0.0.1:8080")
	if err != nil {
		log.Fatal("cannot bind to address:", err)
	}

允许的地址格式:
  - :8080              - 监听所有网卡
  - 0.0.0.0:8080       - 监听所有 IPv4
  - 127.0.0.1:8080     - 本地回环
  - localhost:8080     - 本地主机
  - [::]:8080          - 所有 IPv6
  - 本机网卡IP:端口    - 本机真实 IP

禁止的地址:
  - 公网 IP
  - 非本机 IP
  - 非法 host 或端口

## 设备 ID 生成

生成基于硬件信息的设备唯一标识，适用于软件授权、设备绑定等场景。

	import "github.com/rei0721/go-scaffold/pkg/utils"

	// 生成设备 ID
	appSalt := "my-app-v1.0" // 应用盐值，一旦发布不要修改
	deviceID := utils.GenerateDeviceID(appSalt)
	// 返回: "a1b2c3d4e5f6..."（64位十六进制字符串）

生成原理:

	基于以下硬件信息生成指纹：
	1. 操作系统信息（GOOS、GOARCH）
	2. 主机名
	3. 真实网卡 MAC 地址（过滤虚拟网卡）
	4. 应用盐值（防止同一台机器在不同软件中生成相同设备码）

注意事项:
  - appSalt 一旦发布不要轻易修改
  - 虚拟网卡会被自动过滤
  - 返回值为 SHA256 哈希值（64 位十六进制字符串）
  - 建议缓存设备 ID，不要频繁生成

使用场景:
  - 软件授权验证
  - 设备绑定
  - 审计日志
  - 防止多开

## 端口查找

在指定范围内查找可用的 TCP 端口。

	import "github.com/rei0721/go-scaffold/pkg/utils"

	// 查找可用端口（30000-40000）
	port, err := utils.GetAvailablePort(30000, 40000)
	if err != nil {
		log.Fatal("no available port:", err)
	}

	// 排除某些端口
	port, err = utils.GetAvailablePort(30000, 40000, 30001, 30002)

	// 在找到的端口上启动服务
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)

验证方式:
  - 通过实际绑定 net.Listen("tcp", "0.0.0.0:port") 测试端口是否可用
  - 使用互斥锁防止并发抢占同一端口

使用场景:
  - 开发环境自动分配端口
  - 微服务动态端口
  - 测试环境避免端口冲突

推荐端口范围:
  - 开发环境：8000-9000
  - 测试环境：10000-20000
  - 微服务：30000-40000（Kubernetes NodePort 范围）

# 最佳实践

## Snowflake ID 生成器

1. 单例模式：

	var (
		idGenOnce sync.Once
		idGen     utils.IDGenerator
	)

	func GetIDGenerator() utils.IDGenerator {
		idGenOnce.Do(func() {
			gen, err := utils.NewSnowflake(getNodeID())
			if err != nil {
				panic(err)
			}
			idGen = gen
		})
		return idGen
	}

2. 分布式环境配置 nodeID：

	// 方式1: 从配置文件
	nodeID := viper.GetInt64("node_id")

	// 方式2: 从环境变量
	nodeID, _ := strconv.ParseInt(os.Getenv("NODE_ID"), 10, 64)

	// 方式3: K8s StatefulSet（推荐）
	podName := os.Getenv("HOSTNAME") // my-app-1
	index := extractIndex(podName)   // 1
	nodeID := int64(index)

## 设备 ID 生成

1. 固定应用盐值：

	const AppSalt = "my-app-v1.0.0" // 发布后不要修改

2. 缓存设备 ID：

	var cachedDeviceID string

	func GetDeviceID() string {
		if cachedDeviceID == "" {
			cachedDeviceID = utils.GenerateDeviceID(AppSalt)
		}
		return cachedDeviceID
	}

## 端口查找

合理设置端口范围，避免查找时间过长：

	// 开发环境
	port, _ := utils.GetAvailablePort(8000, 9000)

	// 微服务
	port, _ := utils.GetAvailablePort(30000, 40000)

# 性能考虑

- Snowflake ID 生成器: 每毫秒可生成 4096 个 ID，性能优异
- 设备 ID 生成: SHA256 哈希计算，建议缓存结果
- 端口查找: 顺序尝试绑定，性能取决于端口范围大小
- 所有工具都是线程安全的，可以在并发环境下使用

# 依赖项

必须依赖:
  - github.com/bwmarrin/snowflake - Snowflake ID 生成算法实现

标准库依赖:
  - net - 网络操作
  - crypto/sha256 - SHA256 哈希
  - runtime - 运行时信息
  - os - 操作系统信息

# 与其他包的区别

- pkg/logger: 用于记录日志
- pkg/cache: 用于缓存数据
- pkg/jwt: 用于用户认证
- pkg/utils: 通用工具函数集合

utils 包提供的是独立的、通用的工具函数，不涉及业务逻辑。

# 参考链接

- Snowflake ID 算法: https://github.com/twitter-archive/snowflake
- bwmarrin/snowflake: https://github.com/bwmarrin/snowflake
- 设备指纹技术: https://en.wikipedia.org/wiki/Device_fingerprint
*/
package utils
