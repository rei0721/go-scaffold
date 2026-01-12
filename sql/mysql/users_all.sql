-- users 表 CRUD 操作 SQL
-- 生成时间: 2025-12-30 17:51:26

-- ==================== 插入操作 ====================
INSERT INTO `users` (`i_d`, `created_at`, `updated_at`, `deleted_at`, `username`, `email`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?);

-- ==================== 查询操作 ====================
-- 查询所有记录
SELECT `i_d`, `created_at`, `updated_at`, `deleted_at`, `username`, `email`, `status` FROM `users` WHERE `deleted_at` IS NULL;

-- 根据 ID 查询
SELECT `i_d`, `created_at`, `updated_at`, `deleted_at`, `username`, `email`, `status` FROM `users` WHERE `id` = ? AND `deleted_at` IS NULL;

-- 分页查询
SELECT `i_d`, `created_at`, `updated_at`, `deleted_at`, `username`, `email`, `status` FROM `users` WHERE `deleted_at` IS NULL ORDER BY `id` LIMIT ? OFFSET ?;

-- ==================== 更新操作 ====================
UPDATE `users` SET `updated_at` = ?, `username` = ?, `email` = ?, `status` = ?, `updated_at` = CURRENT_TIMESTAMP WHERE `id` = ? AND `deleted_at` IS NULL;

-- ==================== 删除操作 ====================
-- 软删除
UPDATE `users` SET `deleted_at` = CURRENT_TIMESTAMP WHERE `id` = ? AND `deleted_at` IS NULL;

-- 硬删除 (谨慎使用)
DELETE FROM `users` WHERE `id` = ?;
