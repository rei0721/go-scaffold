-- CREATE TABLE 操作 SQL
-- 生成时间: 2025-12-30 17:51:33

CREATE TABLE users (
    i_d INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    status INTEGER DEFAULT 1
);

CREATE INDEX idx_users_deleted_at ON users (deleted_at);
CREATE UNIQUE INDEX idx_users_username ON users (username);
CREATE UNIQUE INDEX idx_users_email ON users (email);
