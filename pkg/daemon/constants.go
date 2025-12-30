package daemon

// 守护进程状态常量
// 用于标识守护进程的运行状态
const (
	// StatusIdle 空闲状态
	// 守护进程已创建但尚未启动
	StatusIdle = "idle"

	// StatusStarting 启动中状态
	// 守护进程正在启动过程中
	StatusStarting = "starting"

	// StatusRunning 运行中状态
	// 守护进程已成功启动并正常运行
	StatusRunning = "running"

	// StatusStopping 停止中状态
	// 守护进程正在执行优雅关闭流程
	StatusStopping = "stopping"

	// StatusStopped 已停止状态
	// 守护进程已完全停止
	StatusStopped = "stopped"

	// StatusFailed 失败状态
	// 守护进程启动或运行过程中发生错误
	StatusFailed = "failed"
)

// 默认超时时间常量
// 用于控制守护进程启动和停止的最大等待时间
const (
	// DefaultStartTimeout 默认启动超时时间
	// 守护进程启动时的最大等待时间
	// 如果超过此时间还未启动成功,将返回超时错误
	DefaultStartTimeout = 30 // 秒

	// DefaultStopTimeout 默认停止超时时间
	// 守护进程停止时的最大等待时间
	// 如果超过此时间还未完全停止,将强制退出
	DefaultStopTimeout = 30 // 秒
)

// 日志消息常量
// 避免在代码中使用魔法字符串,便于统一管理和修改
const (
	// MsgDaemonRegistered 守护进程注册成功消息
	MsgDaemonRegistered = "daemon registered"

	// MsgDaemonStarting 守护进程开始启动消息
	MsgDaemonStarting = "daemon starting"

	// MsgDaemonStarted 守护进程启动成功消息
	MsgDaemonStarted = "daemon started successfully"

	// MsgDaemonStartFailed 守护进程启动失败消息
	MsgDaemonStartFailed = "daemon start failed"

	// MsgDaemonStopping 守护进程开始停止消息
	MsgDaemonStopping = "daemon stopping"

	// MsgDaemonStopped 守护进程停止成功消息
	MsgDaemonStopped = "daemon stopped successfully"

	// MsgDaemonStopFailed 守护进程停止失败消息
	MsgDaemonStopFailed = "daemon stop failed"

	// MsgAllDaemonsStarting 所有守护进程开始启动消息
	MsgAllDaemonsStarting = "starting all daemons"

	// MsgAllDaemonsStarted 所有守护进程启动完成消息
	MsgAllDaemonsStarted = "all daemons started successfully"

	// MsgAllDaemonsStopping 所有守护进程开始停止消息
	MsgAllDaemonsStopping = "stopping all daemons"

	// MsgAllDaemonsStopped 所有守护进程停止完成消息
	MsgAllDaemonsStopped = "all daemons stopped successfully"

	// MsgStopTimeout 停止超时消息
	MsgStopTimeout = "timeout waiting for daemons to stop"
)

// 错误消息常量
// 用于创建错误时的统一消息格式
const (
	// ErrMsgDaemonStartFailed 守护进程启动失败的错误消息格式
	// 使用 fmt.Sprintf(ErrMsgDaemonStartFailed, daemonName, err)
	ErrMsgDaemonStartFailed = "daemon '%s' failed to start: %w"

	// ErrMsgDaemonStopFailed 守护进程停止失败的错误消息格式
	// 使用 fmt.Sprintf(ErrMsgDaemonStopFailed, daemonName, err)
	ErrMsgDaemonStopFailed = "daemon '%s' failed to stop: %w"
)
