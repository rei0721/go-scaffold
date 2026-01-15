package rbac

import "time"

const (
	// CacheKeyPrefixRole 角色缓存键前缀
	CacheKeyPrefixRole = "role:"

	// CacheKeyPrefixPermission 权限缓存键前缀
	CacheKeyPrefixPermission = "permission:"

	// CacheKeyPrefixUserPermissions 用户权限集合缓存键前缀
	// 格式: "user:perms:{userID}"
	CacheKeyPrefixUserPermissions = "user:perms:"

	// CacheTTLPermission 权限缓存过期时间
	// 用户权限变更频率较低，缓存时间可以设置较长
	CacheTTLPermission = 60 * time.Minute
)
