package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rei0721/go-scaffold/internal/service/auth"
	"github.com/rei0721/go-scaffold/types"
	"github.com/rei0721/go-scaffold/types/errors"
)

// AuthHandler 认证处理器
// 处理用户注册、登录、登出、密码管理等HTTP请求
type AuthHandler struct {
	authService auth.AuthService
}

// NewAuthHandler 创建认证处理器实例
// 参数:
//
//	authService: 认证服务实例
//
// 返回:
//
//	*AuthHandler: 认证处理器
func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 处理用户注册请求
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body types.RegisterRequest true "注册信息"
// @Success 200 {object} types.Response{data=types.UserResponse} "注册成功"
// @Failure 400 {object} types.Response "请求参数错误"
// @Failure 409 {object} types.Response "用户名或邮箱已存在"
// @Failure 500 {object} types.Response "服务器内部错误"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		types.ErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidRequest, "invalid request body", err)
		return
	}

	// 2. 调用服务层处理注册逻辑
	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		// 服务层会返回详细的业务错误
		if bizErr, ok := err.(*errors.BizError); ok {
			switch bizErr.Code {
			case errors.ErrDuplicateUsername, errors.ErrDuplicateEmail:
				types.ErrorResponse(c, http.StatusConflict, bizErr.Code, bizErr.Message, bizErr.Cause)
			case errors.ErrInvalidRequest:
				types.ErrorResponse(c, http.StatusBadRequest, bizErr.Code, bizErr.Message, bizErr.Cause)
			default:
				types.ErrorResponse(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message, bizErr.Cause)
			}
			return
		}
		types.ErrorResponse(c, http.StatusInternalServerError, errors.ErrInternalServer, "registration failed", err)
		return
	}

	// 3. 返回成功响应
	types.SuccessResponse(c, user)
}

// Login 处理用户登录请求
// @Summary 用户登录
// @Description 用户名密码登录，返回JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body types.LoginRequest true "登录信息"
// @Success 200 {object} types.Response{data=types.LoginResponse} "登录成功"
// @Failure 400 {object} types.Response "请求参数错误"
// @Failure 401 {object} types.Response "用户名或密码错误"
// @Failure 500 {object} types.Response "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		types.ErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidRequest, "invalid request body", err)
		return
	}

	// 2. 调用服务层处理登录逻辑
	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if bizErr, ok := err.(*errors.BizError); ok {
			switch bizErr.Code {
			case errors.ErrInvalidCredentials, errors.ErrUserNotFound:
				types.ErrorResponse(c, http.StatusUnauthorized, bizErr.Code, bizErr.Message, bizErr.Cause)
			case errors.ErrUserDisabled:
				types.ErrorResponse(c, http.StatusForbidden, bizErr.Code, bizErr.Message, bizErr.Cause)
			default:
				types.ErrorResponse(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message, bizErr.Cause)
			}
			return
		}
		types.ErrorResponse(c, http.StatusInternalServerError, errors.ErrInternalServer, "login failed", err)
		return
	}

	// 3. 返回成功响应
	types.SuccessResponse(c, resp)
}

// Logout 处理用户登出请求
// @Summary 用户登出
// @Description 清除用户会话信息
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} types.Response "登出成功"
// @Failure 401 {object} types.Response "未授权"
// @Failure 500 {object} types.Response "服务器内部错误"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 1. 从上下文获取用户ID (由认证中间件注入)
	userID, exists := c.Get("user_id")
	if !exists {
		types.ErrorResponse(c, http.StatusUnauthorized, errors.ErrUnauthorized, "user not authenticated", nil)
		return
	}

	// 2. 类型断言
	uid, ok := userID.(int64)
	if !ok {
		types.ErrorResponse(c, http.StatusInternalServerError, errors.ErrInternalServer, "invalid user ID format", nil)
		return
	}

	// 3. 调用服务层处理登出逻辑
	if err := h.authService.Logout(c.Request.Context(), uid); err != nil {
		if bizErr, ok := err.(*errors.BizError); ok {
			types.ErrorResponse(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message, bizErr.Cause)
			return
		}
		types.ErrorResponse(c, http.StatusInternalServerError, errors.ErrInternalServer, "logout failed", err)
		return
	}

	// 4. 返回成功响应
	types.SuccessResponse(c, gin.H{"message": "logout successful"})
}

// ChangePassword 处理修改密码请求
// @Summary 修改密码
// @Description 用户修改自己的密码
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body types.ChangePasswordRequest true "密码信息"
// @Success 200 {object} types.Response "修改成功"
// @Failure 400 {object} types.Response "请求参数错误"
// @Failure 401 {object} types.Response "未授权或旧密码错误"
// @Failure 500 {object} types.Response "服务器内部错误"
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// 1. 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		types.ErrorResponse(c, http.StatusUnauthorized, errors.ErrUnauthorized, "user not authenticated", nil)
		return
	}

	uid, ok := userID.(int64)
	if !ok {
		types.ErrorResponse(c, http.StatusInternalServerError, errors.ErrInternalServer, "invalid user ID format", nil)
		return
	}

	// 2. 绑定并验证请求参数
	var req types.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		types.ErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidRequest, "invalid request body", err)
		return
	}

	// 3. 调用服务层处理修改密码逻辑
	if err := h.authService.ChangePassword(c.Request.Context(), uid, &req); err != nil {
		if bizErr, ok := err.(*errors.BizError); ok {
			switch bizErr.Code {
			case errors.ErrInvalidCredentials:
				types.ErrorResponse(c, http.StatusUnauthorized, bizErr.Code, bizErr.Message, bizErr.Cause)
			case errors.ErrInvalidRequest:
				types.ErrorResponse(c, http.StatusBadRequest, bizErr.Code, bizErr.Message, bizErr.Cause)
			default:
				types.ErrorResponse(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message, bizErr.Cause)
			}
			return
		}
		types.ErrorResponse(c, http.StatusInternalServerError, errors.ErrInternalServer, "change password failed", err)
		return
	}

	// 4. 返回成功响应
	types.SuccessResponse(c, gin.H{"message": "password changed successfully"})
}

// RefreshToken 处理刷新令牌请求
// @Summary 刷新访问令牌
// @Description 使用refresh token获取新的access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body types.RefreshTokenRequest true "刷新令牌"
// @Success 200 {object} types.Response{data=types.TokenResponse} "刷新成功"
// @Failure 400 {object} types.Response "请求参数错误"
// @Failure 401 {object} types.Response "令牌无效或已过期"
// @Failure 500 {object} types.Response "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 1. 绑定并验证请求参数
	var req types.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		types.ErrorResponse(c, http.StatusBadRequest, errors.ErrInvalidRequest, "invalid request body", err)
		return
	}

	// 2. 调用服务层处理刷新令牌逻辑
	resp, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		if bizErr, ok := err.(*errors.BizError); ok {
			switch bizErr.Code {
			case errors.ErrInvalidToken, errors.ErrTokenExpired:
				types.ErrorResponse(c, http.StatusUnauthorized, bizErr.Code, bizErr.Message, bizErr.Cause)
			default:
				types.ErrorResponse(c, http.StatusInternalServerError, bizErr.Code, bizErr.Message, bizErr.Cause)
			}
			return
		}
		types.ErrorResponse(c, http.StatusInternalServerError, errors.ErrInternalServer, "refresh token failed", err)
		return
	}

	// 3. 返回成功响应
	types.SuccessResponse(c, resp)
}
