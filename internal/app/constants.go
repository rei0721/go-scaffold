package app

// AppMode 应用启动模式
type AppMode string

const (
	// ModeServer 服务器模式（默认）
	// 完整启动应用，包括 HTTP 服务器
	ModeServer AppMode = "server"

	// ModeInitDB 数据库初始化模式
	// 仅初始化数据库连接，执行初始化脚本后退出
	ModeInitDB AppMode = "initdb"
)

/*
// I18n
*/

const (
	ConstantsI18nMessagesDir     = "./configs/locales"
	ConstantsI18nDefaultLanguage = "zh-CN"
	ConstantsDefaultHost         = "localhost"
)

var ConstantsI18nSupportedLanguages = []string{"zh-CN", "en-US"}

/*
// 数据库初始化相关常量
*/

const (
	// InitDBScriptDir 初始化 SQL 脚本目录
	ConstantsInitDBScriptDir = "./scripts/initdb"

	// InitDBLockFile 初始化锁文件名
	ConstantsInitDBLockFile = ".initialized"

	// InitDBScriptFileName 初始化脚本文件名模板
	ConstantsInitDBScriptFileName = "init" // %s 为数据库类型

	ConstantsInitDBScriptFileSuffix = ".sql"

	ConstantsInitDBScriptFileTemplate = "%s.%s.sql"
)
