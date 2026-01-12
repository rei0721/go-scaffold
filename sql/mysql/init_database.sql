-- 数据库初始化脚本
-- 生成时间: 2025-12-30 17:51:26
-- 此文件包含所有表的建表语句

-- ==================== users 表 ====================
CREATE TABLE `users` (
    `i_d` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `created_at` DATETIME,
    `updated_at` DATETIME,
    `deleted_at` DATETIME,
    `username` VARCHAR(50) NOT NULL,
    `email` VARCHAR(100) NOT NULL,
    `status` INT DEFAULT 1,
    KEY `idx_users_deleted_at` (`deleted_at`),
    UNIQUE KEY `idx_users_username` (`username`),
    UNIQUE KEY `idx_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

