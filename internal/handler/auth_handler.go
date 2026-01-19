package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rei0721/go-scaffold/internal/middleware"
	"github.com/rei0721/go-scaffold/internal/service/auth"
	"github.com/rei0721/go-scaffold/pkg/logger"
	"github.com/rei0721/go-scaffold/types"
	"github.com/rei0721/go-scaffold/types/result"
)

// AuthHandler 认证处理器
// 负责处理用户认证相关的 HTTP 请求
// 包括注册、登录、登出等功能
type AuthHandler struct {
	// authService 认证业务服务
	// 负责实际的认证业务逻辑
	authService auth.AuthService

	// logger 日志记录器
	// 用于记录请求处理过程中的日志
	logger logger.Logger
}

// NewAuthHandler 创建新的认证处理器
// 参数:
//   - authService: 认证服务实例
//   - logger: 日志记录器实例
//
// 返回:
//   - *AuthHandler: 认证处理器实例
func NewAuthHandler(authService auth.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register 处理用户注册请求
// POST /api/v1/auth/register
//
//	Body: {
//	  "username": "user123",
//	  "email": "user@example.com",
//	  "password": "password123"
//	}
//
// 响应:
//
//	200 OK - 注册成功，返回用户信息
//	400 Bad Request - 请求参数错误
//	409 Conflict - 用户名或邮箱已存在
//	500 Internal Server Error - 服务器内部错误
func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest

	// 绑定并验证请求数据
	// ShouldBindJSON 会自动验证 binding tag
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid register request", "error", err)
		result.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 调用服务层处理注册逻辑
	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to register user", "username", req.Username, "error", err)
		// 根据错误类型返回不同的状态码
		// 这里简化处理，实际应该根据具体错误类型判断
		result.InternalError(c, "注册失败: "+err.Error())
		return
	}

	h.logger.Info("user registered successfully", "userId", user.UserID, "username", user.Username)
	c.JSON(http.StatusOK, result.Success(user))
}

// Login 处理用户登录请求
// POST /api/v1/auth/login
//
//	Body: {
//	  "username": "user123",
//	  "password": "password123"
//	}
//
// 响应:
//
//	200 OK - 登录成功，返回 token 和用户信息
//	400 Bad Request - 请求参数错误
//	401 Unauthorized - 用户名或密码错误
//	500 Internal Server Error - 服务器内部错误
func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid login request", "error", err)
		result.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 调用服务层处理登录逻辑
	loginResp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Warn("login failed", "username", req.Username, "error", err)
		// 登录失败返回 401
		result.Unauthorized(c, "用户名或密码错误")
		return
	}

	h.logger.Info("user logged in successfully", "userId", loginResp.User.UserID, "username", loginResp.User.Username)
	c.JSON(http.StatusOK, result.Success(loginResp))
}

// Logout 处理用户登出请求
// POST /api/v1/auth/logout
// Headers: Authorization: Bearer <token>
//
// 响应:
//
//	200 OK - 登出成功
//	401 Unauthorized - 未认证或 token 无效
//	500 Internal Server Error - 服务器内部错误
//
// 注意: 需要认证中间件保护
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从上下文中获取当前用户 ID
	// 这个值由 AuthMiddleware 设置
	userID, ok := middleware.GetUserID(c)
	if !ok {
		h.logger.Warn("user ID not found in context")
		result.Unauthorized(c, "未认证")
		return
	}

	// 调用服务层处理登出逻辑
	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		h.logger.Error("failed to logout", "userId", userID, "error", err)
		result.InternalError(c, "登出失败: "+err.Error())
		return
	}

	h.logger.Info("user logged out successfully", "userId", userID)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"message": "登出成功",
	}))
}

// ChangePassword 处理修改密码请求
// POST /api/v1/auth/change-password
// Headers: Authorization: Bearer <token>
//
//	Body: {
//	  "old_password": "oldpass123",
//	  "new_password": "newpass123"
//	}
//
// 响应：
//
//	200 OK - 密码修改成功
//	400 Bad Request - 请求参数错误
//	401 Unauthorized - 未认证或旧密码错误
//	500 Internal Server Error - 服务器内部错误
//
// 注意: 需要认证中间件保护
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req types.ChangePasswordRequest

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid change password request", "error", err)
		result.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 从上下文中获取当前用户 ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		h.logger.Warn("user ID not found in context")
		result.Unauthorized(c, "未认证")
		return
	}

	// 调用服务层处理密码修改逻辑
	if err := h.authService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		h.logger.Error("failed to change password", "userId", userID, "error", err)
		// 根据错误类型返回不同的响应
		// 这里简化处理，实际应该根据具体错误类型判断
		result.InternalError(c, "修改密码失败: "+err.Error())
		return
	}

	h.logger.Info("password changed successfully", "userId", userID)
	c.JSON(http.StatusOK, result.Success(gin.H{
		"message": "密码修改成功",
	}))
}

// RefreshToken 处理刷新 token 请求
// POST /api/v1/auth/refresh
//
//	Body: {
//	  "refresh_token": "eyJhbGc..."
//	}
//
// 响应：
//
//	200 OK - 刷新成功，返回新的 token
//	400 Bad Request - 请求参数错误
//	401 Unauthorized - refresh token 无效或过期
//	500 Internal Server Error - 服务器内部错误
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req types.RefreshTokenRequest

	// 绑定并验证请求数据
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid refresh token request", "error", err)
		result.BadRequest(c, "无效的请求参数: "+err.Error())
		return
	}

	// 调用服务层处理 token 刷新逻辑
	tokenResp, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		h.logger.Warn("refresh token failed", "error", err)
		// Token 验证失败返回 401
		result.Unauthorized(c, "refresh token 无效或过期")
		return
	}

	h.logger.Info("token refreshed successfully")
	c.JSON(http.StatusOK, result.Success(tokenResp))
}
