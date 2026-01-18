package yaml2go

const (
	// DefaultPackageName 默认包名
	DefaultPackageName = "main"

	// DefaultStructName 默认根结构体名
	DefaultStructName = "Config"

	// DefaultIndentStyle 默认缩进风格（使用 tab）
	DefaultIndentStyle = "tab"

	// IndentStyleTab tab 缩进
	IndentStyleTab = "tab"

	// IndentStyleSpace 空格缩进
	IndentStyleSpace = "space"
)

var (
	// DefaultTags 默认生成的标签列表
	// - json: JSON 序列化
	// - yaml: YAML 序列化
	// - mapstructure: Viper 配置读取
	// - toml: TOML 序列化
	DefaultTags = []string{"json", "yaml", "mapstructure", "toml"}
)
