// Package result 提供了 API 响应的通用结果类型
// 这个包除了 Go 标准库外没有外部依赖
// 设计目标:
// - 统一 API 响应格式,保持前后端接口一致性
// - 使用泛型支持任意数据类型,提高代码复用性
// - 包含必要的元数据(错误码、追踪ID、服务器时间)
package result

import "time"

// Result 表示通用的 API 响应结构
// 使用 Go 1.18+ 泛型,T 可以是任意类型
// 这个结构体定义了所有 API 响应的标准格式
// 设计考虑:
// - 前端可以统一处理响应格式
// - 错误码和消息便于错误处理
// - TraceID 用于问题追踪
// - ServerTime 用于客户端时间同步和请求耗时计算
type Result[T any] struct {
	// Code 响应码
	// 0 表示成功,其他值表示各种错误
	// 前端应该首先检查此字段判断请求是否成功
	// 对应 types/errors/codes.go 中定义的错误码
	Code int `json:"code"`

	// Message 响应消息
	// 成功时通常是 "success"
	// 失败时是人类可读的错误描述
	// 可以直接展示给用户(已经是友好的文本)
	Message string `json:"message"`

	// Data 响应数据
	// 泛型字段,可以是任意类型:
	// - 单个对象: Result[UserResponse]
	// - 列表: Result[[]UserResponse]
	// - 分页: Result[PageResult[UserResponse]]
	// - 简单值: Result[bool], Result[string]
	// omitempty: 错误时此字段为空,不会在 JSON 中出现
	Data T `json:"data,omitempty"`

	// TraceID 请求追踪 ID
	// 从中间件中传递过来,用于日志关联
	// 当出现错误时,用户可以提供此 ID 给技术支持人员
	// omitempty: 如果没有设置 TraceID,不会在 JSON 中出现
	TraceID string `json:"traceId,omitempty"`

	// ServerTime 服务器时间戳(Unix 秒)
	// 用途:
	// - 客户端可以计算请求往返时间(RTT)
	// - 可以用于客户端和服务器时间同步
	// - 在日志中方便判断请求处理的时间点
	ServerTime int64 `json:"serverTime"`
}

// Success 创建一个成功的 Result
// 这是一个便捷函数,封装了成功响应的创建逻辑
// 参数:
//   data: 要返回的数据,可以是任意类型
// 返回:
//   *Result[T]: 包含数据的成功响应
// 使用示例:
//   return c.JSON(200, result.Success(user))
//   return c.JSON(200, result.Success(users))
//   return c.JSON(200, result.Success(pageResult))
// 响应格式:
//   {
//     "code": 0,
//     "message": "success",
//     "data": {...},
//     "serverTime": 1640000000
//   }
func Success[T any](data T) *Result[T] {
	return &Result[T]{
		Code:       0,                 // 0 表示成功
		Message:    "success",         // 固定的成功消息
		Data:       data,              // 实际的返回数据
		ServerTime: time.Now().Unix(), // 当前服务器时间(Unix秒)
	}
}

// Error 创建一个错误 Result
// 用于返回不包含 TraceID 的错误响应
// 参数:
//   code: 错误码,应使用 errors 包中定义的常量
//   message: 错误消息,应该是用户友好的描述
// 返回:
//   *Result[any]: 错误响应(不包含 data)
// 使用场景:
//   一般不直接使用,推荐使用 ErrorWithTrace
//   因为大多数场景都需要 TraceID 进行问题追踪
func Error(code int, message string) *Result[any] {
	return &Result[any]{
		Code:       code,              // 错误码
		Message:    message,           // 错误消息
		ServerTime: time.Now().Unix(), // 服务器时间
		// 注意:没有设置 Data 和 TraceID
	}
}

// ErrorWithTrace 创建一个包含 TraceID 的错误 Result
// 这是推荐的创建错误响应的方式
// 参数:
//   code: 错误码,应使用 errors 包中定义的常量
//   message: 错误消息,用户友好的描述
//   traceID: 请求追踪 ID,从中间件获取
// 返回:
//   *Result[any]: 包含 TraceID 的错误响应
// 使用示例:
//   traceID := getTraceID(c)
//   c.JSON(400, result.ErrorWithTrace(errors.ErrInvalidParams, "参数错误", traceID))
// 响应格式:
//   {
//     "code": 1000,
//     "message": "参数错误",
//     "traceId": "123456789",
//     "serverTime": 1640000000
//   }
// 为什么需要 TraceID:
//   - 用户报告问题时可以提供 TraceID
//   - 技术支持可以通过 TraceID 查找相关日志
//   - 在分布式系统中追踪请求链路
func ErrorWithTrace(code int, message string, traceID string) *Result[any] {
	return &Result[any]{
		Code:       code,              // 错误码
		Message:    message,           // 错误消息
		TraceID:    traceID,           // 请求追踪 ID,这是与 Error 的区别
		ServerTime: time.Now().Unix(), // 服务器时间
		// 注意:错误响应不包含 Data 字段
	}
}
