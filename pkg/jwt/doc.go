/*
Package jwt 提供JWT（JSON Web Token）令牌的生成和验证功能

# 设计目标

- 安全性: 使用HMAC-SHA256算法签名,确保token不可伪造
- 简单易用: 提供清晰的接口,隐藏复杂的JWT实现细节
- 配置驱动: 支持自定义过期时间、签发者等配置
- 线程安全: 所有方法都是并发安全的

# 核心概念

JWT是一种开放标准(RFC 7519),用于在各方之间安全地传输信息。
一个JWT由三部分组成,用点(.)分隔:
  - Header: 包含token类型和签名算法
  - Payload: 包含声明(claims)
  - Signature: 对header和payload的签名

本包使用HMAC-SHA256算法进行签名,适合单体应用场景。

# 使用示例

基本用法:

	import (
		"github.com/rei0721/go-scaffold/pkg/jwt"
	)

	// 1. 创建JWT管理器
	jwtManager, err := jwt.New(&jwt.Config{
		Secret:    "your-secret-key-at-least-32-characters-long",
		ExpiresIn: 3600,  // 1小时
		Issuer:    "my-app",
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. 生成token
	token, err := jwtManager.GenerateToken(12345, "john_doe")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token:", token)

	// 3. 验证token
	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("UserID: %d, Username: %s\n", claims.UserID, claims.Username)

与HTTP中间件配合使用:

	import (
		"github.com/gin-gonic/gin"
		"github.com/rei0721/go-scaffold/pkg/jwt"
	)

	func AuthMiddleware(jwtManager jwt.JWT) gin.HandlerFunc {
		return func(c *gin.Context) {
			token := c.GetHeader("Authorization")
			if token == "" {
				c.AbortWithStatus(401)
				return
			}

			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				c.AbortWithStatus(401)
				return
			}

			c.Set("user_id", claims.UserID)
			c.Next()
		}
	}

# 最佳实践

1. 密钥管理
  - 使用至少32个字符的随机字符串作为密钥
  - 从环境变量读取密钥,不要硬编码在代码中
  - 定期轮换密钥（需要配合token刷新机制）

2. 过期时间设置
  - 根据业务敏感度调整过期时间
  - 高安全场景: 15-30分钟
  - 一般场景: 1-2小时
  - 低敏感场景: 24小时
  - 配合RefreshToken实现长时间会话

3. Token传输
  - 使用HTTPS传输token
  - 在HTTP请求头中使用Authorization: Bearer <token>
  - 前端存储: 使用httpOnly cookie或内存存储,避免XSS攻击

4. 错误处理
  - 验证失败时不要泄露具体原因给客户端
  - 统一返回401 Unauthorized
  - 在服务端日志中记录详细错误信息用于调试

# 线程安全

所有公开方法都是线程安全的,可以在并发环境下安全使用。
内部使用sync.RWMutex保护配置字段的读写。

# 与其他包的区别

- pkg/logger: 用于记录日志
- pkg/cache: 用于缓存数据
- pkg/jwt: 用于用户认证和授权

JWT包专注于token的生成和验证,不处理用户管理、权限控制等业务逻辑。
这些应该在业务层（如internal/service）实现。

# 性能考虑

- JWT生成和验证是CPU密集型操作
- 建议在高并发场景下使用协程池（pkg/executor）异步处理
- 考虑缓存验证结果（使用pkg/cache）减少重复验证

# 安全注意事项

1. 不要在JWT中存储敏感信息（如密码、信用卡号）
2. JWT是base64编码,不是加密,任何人都可以解码查看内容
3. Token一旦签发无法主动撤销,只能等待过期
4. 防止暴力破解: 使用足够长的密钥
5. 防止时序攻击: jwt库已内置防护

# 依赖

- github.com/golang-jwt/jwt/v5: JWT标准实现库
*/
package jwt
