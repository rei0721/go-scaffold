package httpserver

import "github.com/rei0721/go-scaffold/pkg/executor"

// SetExecutor 设置协程池管理器
// 实现延迟注入模式，支持在服务器初始化后设置
// 使用 atomic.Value 实现原子替换，无需加锁
//
// 参数:
//
//	exec: 协程池管理器实例，为nil时禁用executor功能
//
// 线程安全:
//
//	使用原子操作保证并发安全
func (s *httpServer) SetExecutor(exec executor.Manager) {
	s.executor.Store(exec)
}

// getExecutor 获取当前executor（内部辅助方法）
// 返回nil表示未设置executor
func (s *httpServer) getExecutor() executor.Manager {
	if exec := s.executor.Load(); exec != nil {
		return exec.(executor.Manager)
	}
	return nil
}
