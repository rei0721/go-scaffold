-- 数据库初始化脚本
-- 生成时间: 2025-12-30 17:56:50
-- 此文件包含所有表的建表语句

-- ==================== users 表 ====================
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    status INTEGER DEFAULT 1
);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);
CREATE UNIQUE INDEX idx_users_username ON users (username);
CREATE UNIQUE INDEX idx_users_email ON users (email);

