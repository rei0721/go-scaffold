-- UPDATE 操作 SQL
-- 生成时间: 2025-12-30 17:51:33

UPDATE users SET updated_at = ?, username = ?, email = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND deleted_at IS NULL;
