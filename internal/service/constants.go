package service

import "time"

const (
	// 缓存键前缀
	CacheKeyPrefixUser = "user:"

	// 缓存过期时间
	CacheTTLUser = 30 * time.Minute
)
