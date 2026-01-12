-- SELECT 操作 SQL
-- 生成时间: 2025-12-30 17:55:22

-- 查询所有记录
SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `username`, `email`, `status` FROM `users` WHERE `deleted_at` IS NULL;

-- 根据 ID 查询
SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `username`, `email`, `status` FROM `users` WHERE `id` = ? AND `deleted_at` IS NULL;

-- 分页查询
SELECT `id`, `created_at`, `updated_at`, `deleted_at`, `username`, `email`, `status` FROM `users` WHERE `deleted_at` IS NULL ORDER BY `id` LIMIT ? OFFSET ?;
