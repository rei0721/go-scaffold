CREATE TABLE `users` (
  `i_d` BIGINT NOT NULL,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` TEXT,
  `username` VARCHAR(50) NOT NULL,
  `email` VARCHAR(100) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `status` INT DEFAULT 1,
  `roles` TEXT,
  PRIMARY KEY (`i_d`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `roles` (
  `i_d` BIGINT NOT NULL,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` TEXT,
  `name` VARCHAR(50) NOT NULL,
  `description` VARCHAR(255),
  `status` INT DEFAULT 1,
  `permissions` TEXT,
  PRIMARY KEY (`i_d`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `permissions` (
  `i_d` BIGINT NOT NULL,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` TEXT,
  `name` VARCHAR(100) NOT NULL,
  `resource` VARCHAR(100) NOT NULL,
  `action` VARCHAR(50) NOT NULL,
  `description` VARCHAR(255),
  `status` INT DEFAULT 1,
  PRIMARY KEY (`i_d`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_roles` (
  `user_i_d` BIGINT NOT NULL,
  `role_i_d` BIGINT NOT NULL,
  PRIMARY KEY (`user_i_d`, `role_i_d`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `role_permissions` (
  `role_i_d` BIGINT NOT NULL,
  `permission_i_d` BIGINT NOT NULL,
  PRIMARY KEY (`role_i_d`, `permission_i_d`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;