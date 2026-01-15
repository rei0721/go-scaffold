package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rei0721/rei0721/types/errors"
)

// Unauthorized 返回401未授权错误响应
// 用于认证失败的场景
// 参数:
//
//	c: Gin上下文
//	message: 错误消息
//
// HTTP状态码: 401 Unauthorized
// 错误码: errors.ErrUnauthorized
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorWithTrace(
		errors.ErrUnauthorized,
		message,
		GetTraceID(c),
	))
}

// BadRequest 返回400参数错误响应
// 用于请求参数不合法的场景
// 参数:
//
//	c: Gin上下文
//	message: 错误消息
//
// HTTP状态码: 400 Bad Request
// 错误码: errors.ErrInvalidParams
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorWithTrace(
		errors.ErrInvalidParams,
		message,
		GetTraceID(c),
	))
}

// NotFound 返回404资源不found错误响应
// 用于资源不存在的场景
// 参数:
//
//	c: Gin上下文
//	message: 错误消息
//
// HTTP状态码: 404 Not Found
// 错误码: errors.ErrResourceNotFound
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorWithTrace(
		errors.ErrResourceNotFound,
		message,
		GetTraceID(c),
	))
}

// InternalError 返回500内部服务器错误响应
// 用于系统内部错误的场景
// 参数:
//
//	c: Gin上下文
//	message: 错误消息
//
// HTTP状态码: 500 Internal Server Error
// 错误码: errors.ErrInternalServer
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorWithTrace(
		errors.ErrInternalServer,
		message,
		GetTraceID(c),
	))
}

// OK 返回200成功响应（带数据）
// 用于请求成功的场景
// 参数:
//
//	c: Gin上下文
//	data: 响应数据
//
// HTTP状态码: 200 OK
func OK[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Success(data))
}

// GetTraceID 从上下文获取TraceID
// 如果未设置则返回空字符串
// 参数:
//
//	c: Gin上下文
//
// 返回:
//
//	string: TraceID
func GetTraceID(c *gin.Context) string {
	traceID, exists := c.Get("trace_id")
	if !exists {
		return ""
	}
	id, ok := traceID.(string)
	if !ok {
		return ""
	}
	return id
}

// Forbidden 返回403禁止访问错误响应
// 用于权限不足的场景
// 参数:
//
//	c: Gin上下文
//	message: 错误消息
//
// HTTP状态码: 403 Forbidden
// 错误码: errors.ErrPermissionDenied
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorWithTrace(
		errors.ErrPermissionDenied,
		message,
		GetTraceID(c),
	))
}

// Fail 返回指定错误码的错误响应
// 用于通用错误处理
// 参数:
//
//	c: Gin上下文
//	httpStatus: HTTP状态码
//	message: 错误消息
func Fail(c *gin.Context, httpStatus int, message string) {
	code := errors.ErrInternalServer
	if httpStatus == http.StatusBadRequest {
		code = errors.ErrInvalidParams
	} else if httpStatus == http.StatusUnauthorized {
		code = errors.ErrUnauthorized
	} else if httpStatus == http.StatusForbidden {
		code = errors.ErrPermissionDenied
	} else if httpStatus == http.StatusNotFound {
		code = errors.ErrResourceNotFound
	}

	c.JSON(httpStatus, ErrorWithTrace(
		code,
		message,
		GetTraceID(c),
	))
}

// Page 返回分页响应
// 参数:
//
//	c: Gin上下文
//	list: 当前页数据列表
//	total: 总记录数
//	page: 当前页码
//	pageSize: 每页大小
func Page[T any](c *gin.Context, list []T, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Success(NewPageResult(list, page, pageSize, total)))
}
