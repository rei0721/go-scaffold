---
name: test-development
description: 编写符合项目规范的单元测试和集成测试
---

# 测试开发规范

## 概述

本 skill 指导编写符合项目规范的 Go 测试代码，包括单元测试、表驱动测试、Mock 使用和基准测试。

## 文件结构

```
{package}/
├── {source}.go      # 源代码
├── {source}_test.go # 测试文件
└── mock_{interface}.go # Mock 文件（可选）
```

## 表驱动测试模式

### 标准模式

```go
package logger

import (
    "testing"
)

func TestReload_Success(t *testing.T) {
    tests := []struct {
        name    string
        cfg     *Config
        wantErr bool
    }{
        {
            name: "valid config",
            cfg: &Config{
                Level:  "info",
                Format: "console",
                Output: "stdout",
            },
            wantErr: false,
        },
        {
            name: "invalid level",
            cfg: &Config{
                Level:  "invalid",
                Format: "console",
                Output: "stdout",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            log, err := New(tt.cfg)
            if (err != nil) != tt.wantErr {
                t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && log == nil {
                t.Error("New() returned nil logger")
            }
        })
    }
}
```

### 使用 testify/assert

```go
import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUserService_GetByID(t *testing.T) {
    tests := []struct {
        name     string
        userID   int64
        wantUser *User
        wantErr  error
    }{
        {
            name:     "existing user",
            userID:   1,
            wantUser: &User{ID: 1, Username: "test"},
            wantErr:  nil,
        },
        {
            name:     "not found",
            userID:   999,
            wantUser: nil,
            wantErr:  ErrUserNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := setupTestService()
            user, err := svc.GetByID(context.Background(), tt.userID)

            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.wantUser.ID, user.ID)
            assert.Equal(t, tt.wantUser.Username, user.Username)
        })
    }
}
```

## Mock 接口测试

### 使用 testify/mock

```go
import (
    "github.com/stretchr/testify/mock"
)

// MockRepository 是 Repository 的 Mock 实现
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestService_WithMock(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("FindByID", mock.Anything, int64(1)).Return(&User{ID: 1}, nil)

    svc := NewService(mockRepo)
    user, err := svc.GetByID(context.Background(), 1)

    require.NoError(t, err)
    assert.Equal(t, int64(1), user.ID)
    mockRepo.AssertExpectations(t)
}
```

## 并发测试

```go
func TestConcurrent(t *testing.T) {
    log, err := New(&Config{Level: "info", Format: "console", Output: "stdout"})
    require.NoError(t, err)

    var wg sync.WaitGroup
    const numGoroutines = 100

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            log.Info("concurrent log", "goroutine", id)
        }(i)
    }

    wg.Wait()
}
```

## 基准测试

```go
func BenchmarkLogger_Info(b *testing.B) {
    log, _ := New(&Config{Level: "info", Format: "json", Output: "stdout"})

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        log.Info("benchmark message", "iteration", i)
    }
}
```

## 测试命令

```bash
# 运行所有测试
go test ./...

# 详细输出
go test -v ./...

# 运行特定包测试
go test -v ./pkg/logger/...

# 运行特定测试函数
go test -v -run TestReload_Success ./pkg/logger/

# 测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 基准测试
go test -bench=. ./pkg/logger/

# 竞态检测
go test -race ./...
```

## 测试辅助函数

```go
// setupTestService 创建测试用的 Service
func setupTestService() *userService {
    return &userService{
        // 初始化测试依赖
    }
}

// cleanupTestData 清理测试数据
func cleanupTestData(t *testing.T) {
    t.Cleanup(func() {
        // 清理逻辑
    })
}
```

## 检查清单

- [ ] 测试文件命名 `*_test.go`
- [ ] 使用表驱动测试模式
- [ ] 每个测试用例有清晰的 `name`
- [ ] 使用 `t.Run()` 运行子测试
- [ ] 错误断言使用 `assert.ErrorIs`
- [ ] 关键断言使用 `require.NoError`
- [ ] Mock 对象使用 `testify/mock`
- [ ] 并发测试使用 `sync.WaitGroup`
- [ ] 运行 `-race` 检测竞态
