# 代码规范

本文档定义了 Go Scaffold 项目的代码风格和编程规范，确保代码的一致性、可读性和可维护性。

## 🎯 总体原则

### 1. 简洁性 (Simplicity)
- 优先选择简单的解决方案
- 避免过度设计和不必要的复杂性
- 代码应该易于理解和维护

### 2. 一致性 (Consistency)
- 遵循统一的命名规范
- 保持代码风格的一致性
- 使用相同的模式和惯例

### 3. 可读性 (Readability)
- 代码应该自解释
- 使用有意义的变量和函数名
- 适当添加注释和文档

### 4. 可维护性 (Maintainability)
- 模块化设计
- 低耦合高内聚
- 易于测试和调试

## 📝 Go 语言规范

### 基础规范

遵循官方的 Go 代码规范：
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Style Guide](https://google.github.io/styleguide/go/)

### 格式化工具

使用以下工具确保代码格式一致：

```bash
# 代码格式化
go fmt ./...

# 导入排序和格式化
goimports -w .

# 代码检查
golangci-lint run
```

## 🏷️ 命名规范

### 包命名 (Package Names)

```go
// ✅ 好的包名
package user
package auth
package logger
package cache

// ❌ 避免的包名
package userManager
package authenticationService
package loggerUtils
```

**规则**：
- 使用小写字母
- 简短且有意义
- 避免下划线和驼峰
- 不要使用复数形式

### 变量命名 (Variable Names)

```go
// ✅ 好的变量名
var userID int64
var userName string
var isActive bool
var maxRetryCount int

// ❌ 避免的变量名
var uid int64           // 太简短
var user_name string    // 使用下划线
var IsActive bool       // 不必要的导出
var MAX_RETRY_COUNT int // 常量应该用 const
```

**规则**：
- 使用驼峰命名法 (camelCase)
- 导出的标识符使用大驼峰 (PascalCase)
- 布尔变量使用 `is`、`has`、`can` 等前缀
- 避免使用缩写，除非是广泛认知的

### 函数命名 (Function Names)

```go
// ✅ 好的函数名
func GetUserByID(id int64) (*User, error)
func CreateUser(user *User) error
func IsValidEmail(email string) bool
func ParseConfig(data []byte) (*Config, error)

// ❌ 避免的函数名
func get_user(id int64) (*User, error)    // 使用下划线
func userCreate(user *User) error         // 不清晰的动词位置
func CheckEmailValid(email string) bool   // 冗余的词汇
```

**规则**：
- 使用动词开头描述行为
- 导出函数使用大驼峰命名
- 返回布尔值的函数使用 `Is`、`Has`、`Can` 等前缀
- 避免冗余的词汇

### 常量命名 (Constants)

```go
// ✅ 好的常量名
const (
    DefaultTimeout = 30 * time.Second
    MaxRetryCount  = 3
    APIVersion     = "v1"
)

// 枚举类型常量
const (
    StatusPending Status = iota
    StatusRunning
    StatusCompleted
    StatusFailed
)

// ❌ 避免的常量名
const (
    DEFAULT_TIMEOUT = 30 * time.Second  // 使用下划线
    maxRetryCount   = 3                 // 应该导出
    apiVersion      = "v1"              // 应该导出
)
```

### 接口命名 (Interface Names)

```go
// ✅ 好的接口名
type Reader interface {
    Read([]byte) (int, error)
}

type UserService interface {
    GetUser(id int64) (*User, error)
    CreateUser(user *User) error
}

type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
}

// ❌ 避免的接口名
type IUserService interface {  // 避免 I 前缀
    GetUser(id int64) (*User, error)
}

type UserServiceInterface interface {  // 避免 Interface 后缀
    GetUser(id int64) (*User, error)
}
```

**规则**：
- 单方法接口通常以 `-er` 结尾
- 避免 `I` 前缀或 `Interface` 后缀
- 使用名词或动词+er 的形式

## 📁 文件和目录结构

### 文件命名

```bash
# ✅ 好的文件名
user.go
user_test.go
auth_handler.go
database_config.go

# ❌ 避免的文件名
User.go              # 首字母大写
userHandler.go       # 驼峰命名
user-handler.go      # 连字符
```

**规则**：
- 使用小写字母和下划线
- 测试文件以 `_test.go` 结尾
- 按功能分组相关文件

### 目录结构

```bash
# ✅ 推荐的目录结构
internal/
├── handler/
│   ├── auth_handler.go
│   ├── user_handler.go
│   └── handler_test.go
├── service/
│   ├── auth/
│   │   ├── auth.go
│   │   ├── auth_impl.go
│   │   └── auth_test.go
│   └── user/
└── repository/
```

## 🔧 代码组织

### 导入顺序

```go
package main

import (
    // 1. 标准库
    "context"
    "fmt"
    "net/http"
    "time"

    // 2. 第三方库
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    "go.uber.org/zap"

    // 3. 本项目包
    "github.com/rei0721/go-scaffold/internal/config"
    "github.com/rei0721/go-scaffold/pkg/logger"
)
```

### 结构体定义

```go
// ✅ 好的结构体定义
type User struct {
    // 导出字段在前
    ID       int64     `json:"id" gorm:"primaryKey"`
    Username string    `json:"username" gorm:"uniqueIndex;not null"`
    Email    string    `json:"email" gorm:"uniqueIndex;not null"`
    
    // 时间字段
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // 私有字段在后
    password string
}

// ❌ 避免的结构体定义
type user struct {  // 应该导出
    id       int64   // 字段顺序混乱
    Username string
    password string
    Email    string
}
```

**规则**：
- 导出字段在前，私有字段在后
- 相关字段分组
- 添加适当的标签
- 使用有意义的字段名

### 函数定义

```go
// ✅ 好的函数定义
func (s *userService) GetUserByID(ctx context.Context, id int64) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

// ❌ 避免的函数定义
func (s *userService) GetUserByID(id int64) (*User, error) {
    // 缺少 context 参数
    // 缺少参数验证
    return s.repo.GetByID(id)
}
```

**规则**：
- 第一个参数通常是 `context.Context`
- 进行参数验证
- 使用错误包装提供上下文
- 保持函数简短和专注

## 🔍 错误处理

### 错误定义

```go
// ✅ 好的错误定义
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidEmail     = errors.New("invalid email format")
    ErrDuplicateUser    = errors.New("user already exists")
)

// 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}
```

### 错误处理

```go
// ✅ 好的错误处理
func (s *userService) CreateUser(ctx context.Context, user *User) error {
    if err := s.validateUser(user); err != nil {
        return fmt.Errorf("user validation failed: %w", err)
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        if errors.Is(err, repository.ErrDuplicateKey) {
            return ErrDuplicateUser
        }
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

// ❌ 避免的错误处理
func (s *userService) CreateUser(user *User) error {
    err := s.repo.Create(user)
    if err != nil {
        return err  // 丢失了错误上下文
    }
    return nil
}
```

**规则**：
- 使用 `fmt.Errorf` 和 `%w` 包装错误
- 定义有意义的错误变量
- 在适当的层级处理错误
- 不要忽略错误

## 📝 注释和文档

### 包注释

```go
// Package auth 提供用户认证和授权功能
//
// 本包实现了基于 JWT 的认证机制，支持用户登录、注册、
// 令牌刷新等功能。同时集成了 RBAC 权限控制系统。
//
// 基本使用方法：
//
//     authService := auth.NewService(config, logger)
//     token, err := authService.Login(username, password)
//
package auth
```

### 函数注释

```go
// GetUserByID 根据用户ID获取用户信息
//
// 参数：
//   - ctx: 请求上下文，用于超时控制和取消操作
//   - id: 用户ID，必须大于0
//
// 返回值：
//   - *User: 用户信息，如果用户不存在则返回 nil
//   - error: 错误信息，可能的错误包括 ErrUserNotFound
//
// 示例：
//
//     user, err := service.GetUserByID(ctx, 123)
//     if err != nil {
//         if errors.Is(err, ErrUserNotFound) {
//             // 处理用户不存在的情况
//         }
//         return err
//     }
//
func (s *userService) GetUserByID(ctx context.Context, id int64) (*User, error) {
    // 实现代码...
}
```

### 结构体注释

```go
// User 表示系统中的用户实体
//
// User 包含了用户的基本信息，包括用户名、邮箱、创建时间等。
// 密码字段不会被序列化到 JSON 中，确保安全性。
type User struct {
    ID       int64     `json:"id"`       // 用户唯一标识
    Username string    `json:"username"` // 用户名，必须唯一
    Email    string    `json:"email"`    // 邮箱地址，必须唯一
    Status   UserStatus `json:"status"`   // 用户状态
    
    CreatedAt time.Time `json:"created_at"` // 创建时间
    UpdatedAt time.Time `json:"updated_at"` // 更新时间
    
    password string // 加密后的密码，不导出
}
```

## 🧪 测试规范

### 测试文件组织

```go
// user_test.go
package user

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestUserService_GetUserByID(t *testing.T) {
    tests := []struct {
        name    string
        userID  int64
        want    *User
        wantErr bool
    }{
        {
            name:   "valid user ID",
            userID: 1,
            want: &User{
                ID:       1,
                Username: "testuser",
                Email:    "test@example.com",
            },
            wantErr: false,
        },
        {
            name:    "invalid user ID",
            userID:  0,
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试实现...
        })
    }
}
```

### 基准测试

```go
func BenchmarkUserService_GetUserByID(b *testing.B) {
    service := setupTestService()
    ctx := context.Background()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := service.GetUserByID(ctx, 1)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## 🔧 工具配置

### .golangci.yml

```yaml
run:
  timeout: 5m
  modules-download-mode: readonly

linters-settings:
  gofmt:
    simplify: true
  goimports:
    local-prefixes: github.com/rei0721/go-scaffold
  golint:
    min-confidence: 0.8
  govet:
    check-shadowing: true
  misspell:
    locale: US

linters:
  enable:
    - gofmt
    - goimports
    - golint
    - govet
    - ineffassign
    - misspell
    - unconvert
    - unused
  disable:
    - errcheck
```

### Makefile 集成

```makefile
.PHONY: fmt lint test

fmt:
	go fmt ./...
	goimports -w .

lint:
	golangci-lint run

test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

check: fmt lint test
```

## 📊 性能考虑

### 内存分配

```go
// ✅ 避免不必要的内存分配
func processUsers(users []User) []string {
    names := make([]string, 0, len(users))  // 预分配容量
    for _, user := range users {
        names = append(names, user.Username)
    }
    return names
}

// ❌ 频繁的内存重分配
func processUsers(users []User) []string {
    var names []string  // 没有预分配容量
    for _, user := range users {
        names = append(names, user.Username)
    }
    return names
}
```

### 字符串拼接

```go
// ✅ 使用 strings.Builder 进行高效拼接
func buildQuery(conditions []string) string {
    var builder strings.Builder
    builder.WriteString("SELECT * FROM users WHERE ")
    
    for i, condition := range conditions {
        if i > 0 {
            builder.WriteString(" AND ")
        }
        builder.WriteString(condition)
    }
    
    return builder.String()
}

// ❌ 使用 + 操作符拼接
func buildQuery(conditions []string) string {
    query := "SELECT * FROM users WHERE "
    for i, condition := range conditions {
        if i > 0 {
            query += " AND "
        }
        query += condition
    }
    return query
}
```

## 🔒 安全考虑

### 输入验证

```go
// ✅ 严格的输入验证
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
    if req == nil {
        return errors.New("request cannot be nil")
    }
    
    if len(req.Username) < 3 || len(req.Username) > 50 {
        return errors.New("username must be between 3 and 50 characters")
    }
    
    if !isValidEmail(req.Email) {
        return errors.New("invalid email format")
    }
    
    if len(req.Password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    
    // 处理逻辑...
}
```

### 敏感信息处理

```go
// ✅ 正确处理敏感信息
type User struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    
    // 密码字段不导出，不会被序列化
    password string
}

// String 方法避免泄露敏感信息
func (u *User) String() string {
    return fmt.Sprintf("User{ID: %d, Username: %s, Email: %s}", 
        u.ID, u.Username, u.Email)
}
```

---

**遵循这些代码规范将帮助我们构建高质量、可维护的 Go 应用程序！** 🚀