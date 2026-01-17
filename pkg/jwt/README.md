# JWT - JSON Web Token 认证组件

## 概述

JWT 是一个轻量级的 JWT（JSON Web Token）认证库，提供 token 的生成、验证和刷新功能，适用于基于 token 的用户认证系统。

### 特性

- ✅ **HMAC-SHA256 签名** - 使用安全的签名算法
- ✅ **线程安全** - 所有操作都是并发安全的
- ✅ **配置驱动** - 支持自定义过期时间、签发者等
- ✅ **接口抽象** - 易于测试和扩展
- ✅ **详细错误** - 提供明确的错误类型
- ✅ **标准遵循** - 符合 RFC 7519 标准
- ✅ **密钥安全** - 强制要求最小密钥长度

## 安装

```bash
go get github.com/golang-jwt/jwt/v5
```

## 快速开始

### 1. 创建 JWT 管理器

```go
package main

import (
    "log"

    "github.com/rei0721/go-scaffold/pkg/jwt"
)

func main() {
    // 创建配置
    config := &jwt.Config{
        Secret:    "your-secret-key-at-least-32-characters-long!!!",
        ExpiresIn: 3600,  // 1小时
        Issuer:    "my-app",
    }

    // 创建 JWT 管理器
    jwtManager, err := jwt.New(config)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 2. 生成 Token

```go
// 用户登录成功后生成 token
token, err := jwtManager.GenerateToken(12345, "john_doe")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Token:", token)
// 输出: Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3. 验证 Token

```go
// 验证 token 并获取用户信息
claims, err := jwtManager.ValidateToken(token)
if err != nil {
    if errors.Is(err, jwt.ErrExpiredToken) {
        // token 已过期
    } else if errors.Is(err, jwt.ErrInvalidToken) {
        // token 无效
    }
    log.Fatal(err)
}

fmt.Printf("UserID: %d, Username: %s\n", claims.UserID, claims.Username)
// 输出: UserID: 12345, Username: john_doe
```

### 4. 刷新 Token

```go
// 刷新 token（生成新的 token）
newToken, err := jwtManager.RefreshToken(oldToken)
if err != nil {
    log.Fatal(err)
}

fmt.Println("New Token:", newToken)
```

## API 文档

### Config 配置

```go
type Config struct {
    Secret    string // 签名密钥（至少 32 个字符）
    ExpiresIn int    // 有效期（秒），默认 3600
    Issuer    string // 签发者，默认 "go-scaffold"
}
```

**配置说明**：

| 字段        | 类型     | 必填 | 说明                     | 默认值         |
| ----------- | -------- | ---- | ------------------------ | -------------- |
| `Secret`    | `string` | ✅   | 签名密钥，至少 32 个字符 | -              |
| `ExpiresIn` | `int`    | ❌   | Token 有效期（秒）       | 3600（1 小时） |
| `Issuer`    | `string` | ❌   | Token 签发者标识         | "go-scaffold"  |

### JWT 接口

```go
type JWT interface {
    GenerateToken(userID int64, username string) (string, error)
    ValidateToken(tokenString string) (*Claims, error)
    RefreshToken(tokenString string) (string, error)
}
```

#### GenerateToken

生成访问令牌。

**参数**：

- `userID` (int64) - 用户 ID
- `username` (string) - 用户名

**返回**：

- `string` - JWT token 字符串
- `error` - 生成失败时的错误

**示例**：

```go
token, err := jwtManager.GenerateToken(123, "alice")
```

#### ValidateToken

验证并解析令牌。

**参数**：

- `tokenString` (string) - JWT token 字符串

**返回**：

- `*Claims` - 解析后的载荷信息
- `error` - 验证失败时的错误

**可能的错误**：

- `ErrInvalidToken` - token 格式无效
- `ErrExpiredToken` - token 已过期
- `ErrTokenNotYetValid` - token 尚未生效
- `ErrInvalidSignature` - 签名验证失败

**示例**：

```go
claims, err := jwtManager.ValidateToken(tokenString)
if err != nil {
    // 处理错误
}
```

#### RefreshToken

刷新令牌，生成新的 token。

**参数**：

- `tokenString` (string) - 旧的 JWT token

**返回**：

- `string` - 新的 JWT token
- `error` - 刷新失败时的错误

**示例**：

```go
newToken, err := jwtManager.RefreshToken(oldToken)
```

### Claims 结构

```go
type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

**字段说明**：

- `UserID` - 用户 ID
- `Username` - 用户名
- `RegisteredClaims` - JWT 标准声明
  - `Issuer` - 签发者
  - `IssuedAt` - 签发时间
  - `ExpiresAt` - 过期时间
  - `NotBefore` - 生效时间

## 使用场景

### 场景 1: 用户登录

```go
func login(username, password string) (string, error) {
    // 1. 验证用户名和密码
    user, err := userService.Authenticate(username, password)
    if err != nil {
        return "", err
    }

    // 2. 生成 JWT token
    token, err := jwtManager.GenerateToken(user.ID, user.Username)
    if err != nil {
        return "", err
    }

    return token, nil
}
```

### 场景 2: HTTP 认证中间件

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/rei0721/go-scaffold/pkg/jwt"
)

func AuthMiddleware(jwtManager jwt.JWT) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从 Header 获取 token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
            return
        }

        // 2. 移除 "Bearer " 前缀
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // 3. 验证 token
        claims, err := jwtManager.ValidateToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }

        // 4. 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)

        c.Next()
    }
}

// 使用中间件
func main() {
    router := gin.Default()

    // 受保护的路由
    protected := router.Group("/api")
    protected.Use(AuthMiddleware(jwtManager))
    {
        protected.GET("/profile", getProfile)
        protected.POST("/posts", createPost)
    }

    router.Run(":8080")
}
```

### 场景 3: Token 自动刷新

```go
func refreshTokenHandler(c *gin.Context) {
    // 1. 获取旧 token
    oldToken := c.GetHeader("Authorization")
    oldToken = strings.TrimPrefix(oldToken, "Bearer ")

    // 2. 刷新 token
    newToken, err := jwtManager.RefreshToken(oldToken)
    if err != nil {
        c.JSON(401, gin.H{"error": "failed to refresh token"})
        return
    }

    // 3. 返回新 token
    c.JSON(200, gin.H{
        "token": newToken,
    })
}
```

## 最佳实践

### 1. 密钥管理

```go
// ✅ 从环境变量读取
config := &jwt.Config{
    Secret: os.Getenv("JWT_SECRET"),
}

// ❌ 不要硬编码密钥
config := &jwt.Config{
    Secret: "my-secret-key", // 危险！
}
```

**建议**：

- 使用至少 32 个字符的随机字符串
- 从环境变量或配置文件读取，不要硬编码
- 定期轮换密钥（配合 refresh token 机制）
- 不同环境使用不同的密钥

### 2. 过期时间设置

根据业务敏感度调整过期时间：

| 场景       | 推荐时长   | 配置示例           |
| ---------- | ---------- | ------------------ |
| 高安全场景 | 15-30 分钟 | `ExpiresIn: 1800`  |
| 一般场景   | 1-2 小时   | `ExpiresIn: 3600`  |
| 低敏感场景 | 24 小时    | `ExpiresIn: 86400` |

**配合 Refresh Token**：

```go
// Access Token: 短期（1小时）
accessToken, _ := jwtManager.GenerateToken(userID, username)

// Refresh Token: 长期（7天），用于刷新 access token
// 可以使用单独的 JWT 实例，配置更长的过期时间
```

### 3. Token 传输

**HTTP Header（推荐）**：

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameter（不推荐）**：

```
GET /api/data?token=eyJhbG... # 会被记录在日志中，不安全
```

**Cookie（可选）**：

```go
// 使用 httpOnly cookie，防止 XSS 攻击
c.SetCookie("token", token, 3600, "/", "", true, true)
```

### 4. 错误处理

```go
claims, err := jwtManager.ValidateToken(tokenString)
if err != nil {
    switch {
    case errors.Is(err, jwt.ErrExpiredToken):
        // token 已过期，引导用户重新登录或刷新 token
        return http.StatusUnauthorized, "token expired, please login again"
    case errors.Is(err, jwt.ErrInvalidSignature):
        // 签名无效，可能是伪造的 token
        log.Warn("invalid signature detected", "token", tokenString)
        return http.StatusUnauthorized, "invalid token"
    default:
        // 其他错误统一返回 401
        return http.StatusUnauthorized, "authentication failed"
    }
}
```

**注意**：

- 验证失败时不要泄露具体原因给客户端
- 统一返回 401 Unauthorized
- 在服务端日志中记录详细错误信息

### 5. 安全注意事项

**不要在 JWT 中存储敏感信息**：

```go
// ❌ 危险：JWT 可以被解码
type Claims struct {
    Password   string // 千万不要！
    CreditCard string // 千万不要！
}

// ✅ 安全：只存储用户标识
type Claims struct {
    UserID   int64
    Username string
}
```

**其他安全建议**：

- JWT 是 base64 编码，不是加密，任何人都可以解码查看内容
- Token 一旦签发无法主动撤销，只能等待过期
- 使用 HTTPS 传输 token，防止中间人攻击
- 防止暴力破解：使用足够长的密钥（至少 32 字符）

## 错误参考

### 预定义错误

| 错误                  | 说明           | 场景                     |
| --------------------- | -------------- | ------------------------ |
| `ErrInvalidToken`     | Token 无效     | Token 格式错误           |
| `ErrExpiredToken`     | Token 已过期   | 超过有效期               |
| `ErrTokenNotYetValid` | Token 尚未生效 | 在 NotBefore 之前使用    |
| `ErrInvalidSignature` | 签名无效       | 签名验证失败，可能被篡改 |
| `ErrMissingSecret`    | 缺少密钥       | 配置中未提供 Secret      |

### 错误处理示例

```go
import "errors"

claims, err := jwtManager.ValidateToken(token)
if err != nil {
    if errors.Is(err, jwt.ErrExpiredToken) {
        // 处理过期 token
    } else if errors.Is(err, jwt.ErrInvalidToken) {
        // 处理无效 token
    }
}
```

## 常量定义

### 默认配置

```go
DefaultExpiresIn = 3600        // 1小时
DefaultIssuer    = "go-scaffold"
```

## 与其他包的配合

### 与 pkg/cache 配合

缓存 token 验证结果，减少重复验证：

```go
func validateTokenWithCache(cache cache.Cache, token string) (*jwt.Claims, error) {
    // 1. 尝试从缓存获取
    key := "token:valid:" + token
    if cached, err := cache.Get(ctx, key); err == nil {
        // 缓存命中，token 有效
        var claims jwt.Claims
        json.Unmarshal([]byte(cached), &claims)
        return &claims, nil
    }

    // 2. 验证 token
    claims, err := jwtManager.ValidateToken(token)
    if err != nil {
        return nil, err
    }

    // 3. 缓存验证结果
    data, _ := json.Marshal(claims)
    cache.Set(ctx, key, string(data), 5*time.Minute)

    return claims, nil
}
```

### 与 pkg/rbac 配合

JWT 提供身份认证，RBAC 提供权限控制：

```go
// 1. JWT 验证身份
claims, err := jwtManager.ValidateToken(token)
if err != nil {
    return errors.New("authentication failed")
}

// 2. RBAC 检查权限
hasPermission, err := rbacService.CheckPermission(ctx, claims.UserID, "posts", "write")
if !hasPermission {
    return errors.New("permission denied")
}
```

## 项目结构

```
pkg/jwt/
├── constants.go    # 常量和错误定义
├── jwt.go          # JWT 接口定义和 Claims 结构
├── jwt_impl.go     # JWT 接口实现
├── doc.go          # 包文档
└── README.md       # 本文档
```

## 线程安全

所有公开方法都是线程安全的，可以在并发环境下安全使用。内部使用 `sync.RWMutex` 保护配置字段的读写。

## 依赖

- `github.com/golang-jwt/jwt/v5` - JWT 标准实现库

## 参考链接

- [JWT 官方网站](https://jwt.io/)
- [RFC 7519 - JWT 标准](https://tools.ietf.org/html/rfc7519)
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- [JWT 最佳实践](https://tools.ietf.org/html/rfc8725)
