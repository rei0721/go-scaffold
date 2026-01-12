-- DELETE 操作 SQL
-- 生成时间: 2025-12-30 17:55:22

-- 软删除
UPDATE `users` SET `deleted_at` = CURRENT_TIMESTAMP WHERE `id` = ? AND `deleted_at` IS NULL;

-- 硬删除 (谨慎使用)
DELETE FROM `users` WHERE `id` = ?;
