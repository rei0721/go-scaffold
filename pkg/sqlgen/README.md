# SQL 生成器 (SQLGen)

SQLGen 是一个基于 GORM 模型自动生成 SQL 语句的工具，支持 PostgreSQL、MySQL 和 SQLite 三种数据库。

## 功能特性

- ✅ **多数据库支持**: PostgreSQL、MySQL、SQLite
- ✅ **完整的 CRUD**: 自动生成建表、插入、查询、更新、删除 SQL
- ✅ **GORM 标签解析**: 支持 primaryKey、uniqueIndex、index、size、default 等标签
- ✅ **软删除支持**: 自动处理 GORM 的软删除机制
- ✅ **文件生成**: 支持生成单独或合并的 SQL 文件
- ✅ **命令行工具**: 提供便捷的命令行接口

## 快速开始

### 1. 使用命令行工具

```bash
# 生成 PostgreSQL SQL 文件
go run ./cmd/sqlgen/main.go -dialect postgres -output ./sql/postgres

# 生成 MySQL SQL 文件
go run ./cmd/sqlgen/main.go -dialect mysql -output ./sql/mysql

# 生成分离的 SQLite SQL 文件
go run ./cmd/sqlgen/main.go -dialect sqlite -output ./sql/sqlite -separate

# 查看帮助
go run ./cmd/sqlgen/main.go -help
```

### 2. 在代码中使用

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/rei0721/rei0721/internal/models"
    "github.com/rei0721/rei0721/pkg/sqlgen"
)

func main() {
    // 创建 PostgreSQL 方言的生成器
    dialect := sqlgen.NewPostgresDialect()
    generator := sqlgen.New(dialect)

    // 生成 User 模型的 SQL
    result, err := generator.GenerateSQL(models.User{})
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("建表 SQL:")
    fmt.Println(result.CreateTable)
    fmt.Println("\n插入 SQL:")
    fmt.Println(result.Insert)
}
```

### 3. 生成文件

```go
package main

import (
    "log"
    
    "github.com/rei0721/rei0721/internal/models"
    "github.com/rei0721/rei0721/pkg/sqlgen"
)

func main() {
    // 创建生成器
    dialect := sqlgen.NewMySQLDialect()
    generator := sqlgen.New(dialect)
    fileGenerator := sqlgen.NewFileGenerator(generator, "./sql")

    // 生成选项
    options := &sqlgen.GenerateOptions{
        OutputDir:       "./sql",
        SeparateFiles:   true,  // 生成分离的文件
        GenerateSummary: true,  // 生成汇总文件
        IncludeComments: true,  // 包含注释
    }

    // 生成文件
    models := []interface{}{
        models.User{},
        // 添加更多模型...
    }

    if err := fileGenerator.GenerateWithOptions(options, models...); err != nil {
        log.Fatal(err)
    }
}
```

## 支持的 GORM 标签

### 字段标签

| 标签 | 说明 | 示例 |
|------|------|------|
| `primaryKey` | 主键 | `gorm:"primaryKey"` |
| `uniqueIndex` | 唯一索引 | `gorm:"uniqueIndex"` |
| `index` | 普通索引 | `gorm:"index"` |
| `not null` | 非空约束 | `gorm:"not null"` |
| `size:n` | 字段大小 | `gorm:"size:50"` |
| `default:value` | 默认值 | `gorm:"default:1"` |
| `column:name` | 列名 | `gorm:"column:user_name"` |
| `type:type` | 数据类型 | `gorm:"type:jsonb"` |

### 示例模型

```go
type User struct {
    ID       int64  `gorm:"primaryKey" json:"id"`
    Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email    string `gorm:"uniqueIndex;size:100;not null" json:"email"`
    Password string `gorm:"size:255;not null" json:"-"`
    Status   int    `gorm:"default:1" json:"status"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (User) TableName() string {
    return "users"
}
```

## 生成的 SQL 示例

### PostgreSQL

```sql
-- 建表 SQL
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

-- 插入 SQL
INSERT INTO users (id, created_at, updated_at, deleted_at, username, email, status) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- 查询 SQL
SELECT id, created_at, updated_at, deleted_at, username, email, status 
FROM users WHERE deleted_at IS NULL;

-- 更新 SQL
UPDATE users SET updated_at = $1, username = $2, email = $3, status = $4, updated_at = CURRENT_TIMESTAMP 
WHERE id = $5 AND deleted_at IS NULL;

-- 删除 SQL (软删除)
UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL;
```

### MySQL

```sql
-- 建表 SQL
CREATE TABLE `users` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `created_at` DATETIME,
    `updated_at` DATETIME,
    `deleted_at` DATETIME,
    `username` VARCHAR(50) NOT NULL,
    `email` VARCHAR(100) NOT NULL,
    `status` INT DEFAULT 1,
    UNIQUE KEY `idx_users_username` (`username`),
    UNIQUE KEY `idx_users_email` (`email`),
    KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入 SQL
INSERT INTO `users` (`created_at`, `updated_at`, `deleted_at`, `username`, `email`, `status`) 
VALUES (?, ?, ?, ?, ?, ?);
```

## 命令行选项

```bash
go run ./cmd/sqlgen/main.go [选项]

选项:
  -dialect string
        数据库方言 (postgres, mysql, sqlite) (默认: postgres)
  -output string
        输出目录 (默认: ./sql)
  -separate
        是否生成分离的 SQL 文件 (默认: false)
  -summary
        是否生成汇总文件 (默认: true)
  -comments
        是否包含注释 (默认: true)
  -help
        显示帮助信息
```

## 生成的文件结构

### 合并模式 (默认)

```
sql/
├── init_database.sql    # 所有表的建表语句
└── users_all.sql        # users 表的所有 CRUD 操作
```

### 分离模式 (-separate)

```
sql/
├── init_database.sql    # 所有表的建表语句
├── users_create.sql     # 建表语句
├── users_insert.sql     # 插入语句
├── users_select.sql     # 查询语句
├── users_update.sql     # 更新语句
└── users_delete.sql     # 删除语句
```

## 数据类型映射

### Go 类型 → PostgreSQL

| Go 类型 | PostgreSQL 类型 |
|---------|-----------------|
| bool | BOOLEAN |
| int, int32 | INTEGER |
| int64 | BIGINT |
| float32 | REAL |
| float64 | DOUBLE PRECISION |
| string | VARCHAR(n) / TEXT |
| time.Time | TIMESTAMP WITH TIME ZONE |
| gorm.DeletedAt | TIMESTAMP WITH TIME ZONE |

### Go 类型 → MySQL

| Go 类型 | MySQL 类型 |
|---------|-------------|
| bool | TINYINT(1) |
| int, int32 | INT |
| int64 | BIGINT |
| float32 | FLOAT |
| float64 | DOUBLE |
| string | VARCHAR(n) / TEXT |
| time.Time | DATETIME |
| gorm.DeletedAt | DATETIME |

### Go 类型 → SQLite

| Go 类型 | SQLite 类型 |
|---------|-------------|
| bool | BOOLEAN |
| int, int32, int64 | INTEGER |
| float32, float64 | REAL |
| string | TEXT |
| time.Time | DATETIME |
| gorm.DeletedAt | DATETIME |

## 扩展新的数据库方言

要支持新的数据库，只需实现 `Dialect` 接口：

```go
type CustomDialect struct{}

func (d *CustomDialect) GetDataType(fieldType reflect.Type, gormTag string) string {
    // 实现类型映射逻辑
}

func (d *CustomDialect) GetCreateTableSQL(tableName string, fields []Field) string {
    // 实现建表 SQL 生成逻辑
}

// 实现其他方法...
```

## 最佳实践

1. **模型定义**: 使用清晰的 GORM 标签定义模型
2. **表名规范**: 实现 `TableName()` 方法明确指定表名
3. **字段命名**: 使用驼峰命名，工具会自动转换为蛇形命名
4. **索引设计**: 合理使用 `uniqueIndex` 和 `index` 标签
5. **软删除**: 使用 `gorm.DeletedAt` 实现软删除功能

## 注意事项

- 生成的 SQL 仅作为参考，实际使用前请根据具体需求调整
- 复杂的业务逻辑和约束可能需要手动添加
- 建议在开发环境中测试生成的 SQL 后再用于生产环境
- 对于大型项目，建议使用数据库迁移工具管理表结构变更

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个工具！

---

**最后更新**: 2025-12-30