-- users 表 CRUD 操作 SQL
-- 生成时间: 2025-12-30 17:50:44

-- ==================== 插入操作 ====================
INSERT INTO users (i_d, created_at, updated_at, deleted_at, username, email, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- ==================== 查询操作 ====================
-- 查询所有记录
SELECT i_d, created_at, updated_at, deleted_at, username, email, status FROM users WHERE deleted_at IS NULL;

-- 根据 ID 查询
SELECT i_d, created_at, updated_at, deleted_at, username, email, status FROM users WHERE id = $1 AND deleted_at IS NULL;

-- 分页查询
SELECT i_d, created_at, updated_at, deleted_at, username, email, status FROM users WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2;

-- ==================== 更新操作 ====================
UPDATE users SET updated_at = $1, username = $2, email = $3, status = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5 AND deleted_at IS NULL;

-- ==================== 删除操作 ====================
-- 软删除
UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL;

-- 硬删除 (谨慎使用)
DELETE FROM users WHERE id = $1;
