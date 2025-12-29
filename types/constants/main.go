package constants

import "time"

const (
	// DefaultConfigPath 是配置文件的默认路径
	// 当环境变量 CONFIG_PATH 未设置时使用此路径
	DefaultConfigPath = "configs/config.yaml"

	// ShutdownTimeout 是优雅关闭的最大等待时间
	// 设置为 30 秒以确保所有正在处理的请求能够完成
	// 超过此时间后将强制关闭,避免进程无限期等待
	ShutdownTimeout = 30 * time.Second

	// EnvConfigPathName 是配置文件路径的环境变量名称
	// 首先尝试从环境变量 CONFIG_PATH 读取配置文件路径
	// 这样做是为了支持不同环境(开发、测试、生产)使用不同的配置文件
	EnvConfigPathName = "CONFIG_PATH"
)
