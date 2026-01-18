package constants

import "time"

const (
	// CacheKeyPrefixUser 是用户缓存的前缀
	CacheKeyPrefixUser = "user:"

	// CacheTTLUser 是用户缓存的过期时间
	CacheTTLUser = 30 * time.Minute
)
