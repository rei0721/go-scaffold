// Package handler 提供应用程序的 HTTP 处理器
// 处理器负责参数解析、验证和响应构造
// 所有响应使用 Result 类型,并从请求上下文中提取 TraceID
// 设计原则:
// - 薄处理层:只处理 HTTP 相关逻辑,业务逻辑在 service 层
// - 统一响应格式:使用 Result 类型
// - 错误处理:将业务错误转换为合适的 HTTP 状态码
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/go-scaffold/internal/service"
	"github.com/rei0721/go-scaffold/types"
	bizErr "github.com/rei0721/go-scaffold/types/errors"
	"github.com/rei0721/go-scaffold/types/result"
)

// TraceIDKey 上下文中存储 TraceID 的键
// 用于从 Gin 上下文中获取请求追踪 ID
const TraceIDKey = "X-Request-ID"

// UserHandler 处理用户相关的 HTTP 请求
// 这是 HTTP 层的用户处理器,依赖 service 层
// 职责:
// - 解析 HTTP 请求参数
// - 验证输入格式
// - 调用 service 层处理业务
// - 构造 HTTP 响应
type UserHandler struct {
	// service 用户业务服务
	// 通过接口依赖,便于测试(可以 mock)
	service service.UserService
}

// NewUserHandler 创建一个新的 UserHandler
// 这是工厂函数,遵循依赖注入模式
// 参数:
//
//	svc: 用户服务接口
//
// 返回:
//
//	*UserHandler: 处理器实例
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		service: svc,
	}
}

// getTraceID 从 Gin 上下文中提取 TraceID
// 这是一个辅助函数,简化 TraceID 的获取
// 参数:
//
//	c: Gin 上下文
//
// 返回:
//
//	string: TraceID,如果不存在返回空字符串
//
// 查找顺序:
//  1. 从上下文变量中获取(由 TraceID 中间件设置)
//  2. 从 HTTP header 中获取(客户端传递)
func getTraceID(c *gin.Context) string {
	// 优先从上下文获取
	if traceID, exists := c.Get(TraceIDKey); exists {
		// 类型断言,确保是字符串
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	// 降级:从 header 获取
	return c.GetHeader(TraceIDKey)
}

// Register 处理用户注册请求
// POST /api/v1/users/register
// 请求体:
//
//	{
//	  "username": "alice",
//	  "email": "alice@example.com",
//	  "password": "password123"
//	}
//
// 响应:
//
//	成功: 200 OK with UserResponse
//	失败: 400/422/500 with error
func (h *UserHandler) Register(c *gin.Context) {
	// 1. 获取 TraceID
	// 用于错误响应和日志关联
	traceID := getTraceID(c)

	// 2. 解析请求体
	// ShouldBindJSON 会:
	// - 解析 JSON 请求体
	// - 根据 binding tag 验证参数
	// 如果验证失败,err 包含详细的验证错误
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 参数验证失败,返回 400 Bad Request
		// err.Error() 包含 Gin 的验证错误详情
		// 例如: "Key: 'RegisterRequest.Username' Error:Field validation for 'Username' failed on the 'min' tag"
		c.JSON(http.StatusBadRequest, result.ErrorWithTrace(
			bizErr.ErrInvalidParams,
			err.Error(),
			traceID,
		))
		return
	}

	// 3. 调用 service 层处理业务逻辑
	// c.Request.Context() 传递请求上下文,支持超时和取消
	resp, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		// 业务错误,由 handleServiceError 处理
		// 会根据错误码返回合适的 HTTP 状态码
		handleServiceError(c, err, traceID)
		return
	}

	// 4. 返回成功响应
	// 使用 Result 泛型类型构造响应
	// 状态码 200 表示成功
	c.JSON(http.StatusOK, &result.Result[*types.UserResponse]{
		Code:       bizErr.CodeSuccess,
		Message:    "success",
		Data:       resp,
		TraceID:    traceID,
		ServerTime: result.Success(resp).ServerTime,
	})
}

// Login 处理用户登录请求
// POST /api/v1/users/login
// 请求体:
//
//	{
//	  "username": "alice",
//	  "password": "password123"
//	}
//
// 响应:
//
//	成功: 200 OK with LoginResponse (包含 token)
//	失败: 400/401/500 with error
func (h *UserHandler) Login(c *gin.Context) {
	// 流程与 Register 类似
	traceID := getTraceID(c)

	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorWithTrace(
			bizErr.ErrInvalidParams,
			err.Error(),
			traceID,
		))
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		handleServiceError(c, err, traceID)
		return
	}

	c.JSON(http.StatusOK, &result.Result[*types.LoginResponse]{
		Code:       bizErr.CodeSuccess,
		Message:    "success",
		Data:       resp,
		TraceID:    traceID,
		ServerTime: result.Success(resp).ServerTime,
	})
}

// GetUser 处理根据 ID 获取用户的请求
// GET /api/v1/users/:id
// 路径参数:
//
//	id: 用户 ID (数字)
//
// 响应:
//
//	成功: 200 OK with UserResponse
//	失败: 400/404/500 with error
func (h *UserHandler) GetUser(c *gin.Context) {
	traceID := getTraceID(c)

	// 1. 获取路径参数
	// c.Param("id") 获取 URL 中的 :id 部分
	// 例如: /api/v1/users/123 => idStr = "123"
	idStr := c.Param("id")

	// 2. 解析为整数
	// strconv.ParseInt(s, base, bitSize)
	// - s: 要解析的字符串
	// - base: 10 表示十进制
	// - bitSize: 64 表示 int64
	// 如果解析失败(如 "abc"),err != nil
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// ID 格式错误(不是数字),返回 400
		c.JSON(http.StatusBadRequest, result.ErrorWithTrace(
			bizErr.ErrInvalidParams,
			"invalid user id",
			traceID,
		))
		return
	}

	// 3. 调用 service 查询用户
	resp, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		// 错误处理(可能是 404 用户不存在)
		handleServiceError(c, err, traceID)
		return
	}

	// 4. 返回用户信息
	c.JSON(http.StatusOK, &result.Result[*types.UserResponse]{
		Code:       bizErr.CodeSuccess,
		Message:    "success",
		Data:       resp,
		TraceID:    traceID,
		ServerTime: result.Success(resp).ServerTime,
	})
}

// ListUsers 处理分页查询用户列表的请求
// GET /api/v1/users?page=1&pageSize=10
// 查询参数:
//
//	page: 页码,默认 1,必须 >= 1
//	pageSize: 每页大小,默认 10,范围 1-100
//
// 响应:
//
//	成功: 200 OK with PageResult
//	失败: 400/500 with error
func (h *UserHandler) ListUsers(c *gin.Context) {
	traceID := getTraceID(c)

	// 1. 解析 page 参数
	// DefaultQuery 会:
	// - 获取查询参数 "page"
	// - 如果不存在,使用默认值 "1"
	// Atoi 将字符串转换为整数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		// page 不是数字或 < 1
		c.JSON(http.StatusBadRequest, result.ErrorWithTrace(
			bizErr.ErrInvalidParams,
			"invalid page parameter",
			traceID,
		))
		return
	}

	// 2. 解析 pageSize 参数
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	// 验证 pageSize 范围: 1-100
	// 限制最大值防止:
	// - 查询过多数据导致性能问题
	// - 响应体过大
	// - 数据库负载过高
	if err != nil || pageSize < 1 || pageSize > 100 {
		c.JSON(http.StatusBadRequest, result.ErrorWithTrace(
			bizErr.ErrInvalidParams,
			"invalid pageSize parameter (must be 1-100)",
			traceID,
		))
		return
	}

	// 3. 调用 service 查询列表
	resp, err := h.service.List(c.Request.Context(), page, pageSize)
	if err != nil {
		handleServiceError(c, err, traceID)
		return
	}

	// 4. 返回分页结果
	// 使用嵌套泛型: Result[PageResult[UserResponse]]
	c.JSON(http.StatusOK, &result.Result[*result.PageResult[types.UserResponse]]{
		Code:       bizErr.CodeSuccess,
		Message:    "success",
		Data:       resp,
		TraceID:    traceID,
		ServerTime: result.Success(resp).ServerTime,
	})
}

// UpdateUser 处理用户更新请求
// PUT /api/v1/users/:id
// 路径参数:
//
//	id: 用户 ID (数字)
//
// 请求体:
//
//	{
//	  "username": "newname",  // 可选
//	  "email": "new@example.com",  // 可选
//	  "status": 1  // 可选
//	}
//
// 响应:
//
//	成功: 200 OK with UserResponse
//	失败: 400/404/422/500 with error
func (h *UserHandler) UpdateUser(c *gin.Context) {
	traceID := getTraceID(c)

	// 1. 获取路径参数中的用户 ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorWithTrace(
			bizErr.ErrInvalidParams,
			"invalid user id",
			traceID,
		))
		return
	}

	// 2. 解析请求体
	var req types.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorWithTrace(
			bizErr.ErrInvalidParams,
			err.Error(),
			traceID,
		))
		return
	}

	// 3. 调用 service 层更新用户
	resp, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		handleServiceError(c, err, traceID)
		return
	}

	// 4. 返回更新后的用户信息
	c.JSON(http.StatusOK, &result.Result[*types.UserResponse]{
		Code:       bizErr.CodeSuccess,
		Message:    "success",
		Data:       resp,
		TraceID:    traceID,
		ServerTime: result.Success(resp).ServerTime,
	})
}

// DeleteUser 处理用户删除请求
// DELETE /api/v1/users/:id
// 路径参数:
//
//	id: 用户 ID (数字)
//
// 响应:
//
//	成功: 200 OK
//	失败: 400/404/500 with error
func (h *UserHandler) DeleteUser(c *gin.Context) {
	traceID := getTraceID(c)

	// 1. 获取路径参数中的用户 ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorWithTrace(
			bizErr.ErrInvalidParams,
			"invalid user id",
			traceID,
		))
		return
	}

	// 2. 调用 service 层删除用户
	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err, traceID)
		return
	}

	// 3. 返回成功响应（无数据）
	c.JSON(http.StatusOK, &result.Result[interface{}]{
		Code:       bizErr.CodeSuccess,
		Message:    "user deleted successfully",
		Data:       nil,
		TraceID:    traceID,
		ServerTime: result.Success[interface{}](nil).ServerTime,
	})
}

// handleServiceError 将业务错误转换为合适的 HTTP 响应
// 这是一个辅助函数,统一处理业务层返回的错误
// 参数:
//
//	c: Gin 上下文,用于发送响应
//	err: 业务错误
//	traceID: 请求追踪 ID
//
// 处理逻辑:
//  1. 如果是 BizError,根据错误码映射 HTTP 状态码
//  2. 如果是未知错误,返回 500 Internal Server Error
func handleServiceError(c *gin.Context, err error, traceID string) {
	// 尝试将 error 类型断言为 *BizError
	if bizError, ok := err.(*bizErr.BizError); ok {
		// 是业务错误,根据错误码获取对应的 HTTP 状态码
		statusCode := getHTTPStatusCode(bizError.Code)
		c.JSON(statusCode, result.ErrorWithTrace(
			bizError.Code,
			bizError.Message,
			traceID,
		))
		return
	}

	// 未知错误(不是 BizError)
	// 为了安全,不向客户端暴露错误详情
	// 返回通用的 500 错误
	c.JSON(http.StatusInternalServerError, result.ErrorWithTrace(
		bizErr.ErrInternalServer,
		"internal server error",
		traceID,
	))
}

// getHTTPStatusCode 将业务错误码映射到 HTTP 状态码
// 这是错误码体系和 HTTP 协议之间的桥梁
// 参数:
//
//	code: 业务错误码
//
// 返回:
//
//	int: HTTP 状态码
//
// 映射规则(与 errors/codes.go 中的分段对应):
//
//	0           -> 200 OK
//	1000-1999   -> 400 Bad Request (参数错误)
//	2000-2999   -> 422 Unprocessable Entity (业务错误)
//	3000-3999   -> 401 Unauthorized (认证错误)
//	4000-4999   -> 404 Not Found (资源不存在)
//	5000-5999   -> 500 Internal Server Error (系统错误)
//	其他        -> 500 Internal Server Error (默认)
func getHTTPStatusCode(code int) int {
	switch {
	case code == bizErr.CodeSuccess:
		// 成功
		return http.StatusOK

	case code >= 1000 && code < 2000:
		// 参数错误
		// 客户端输入有误,应该修正参数后重试
		return http.StatusBadRequest

	case code >= 2000 && code < 3000:
		// 业务逻辑错误
		// 例如:用户名已存在、邮箱重复等
		// 422 表示请求格式正确,但因业务规则无法处理
		return http.StatusUnprocessableEntity

	case code >= 3000 && code < 4000:
		// 认证/授权错误
		// 例如:密码错误、token 过期、权限不足
		// 401 要求客户端提供有效的身份凭证
		return http.StatusUnauthorized

	case code >= 4000 && code < 5000:
		// 资源不存在
		// 例如:用户 ID 不存在
		// 404 表示请求的资源未找到
		return http.StatusNotFound

	case code >= 5000 && code < 6000:
		// 系统错误
		// 例如:数据库错误、缓存错误等
		// 500 表示服务器内部错误
		return http.StatusInternalServerError

	default:
		// 未知错误码,返回 500
		return http.StatusInternalServerError
	}
}
