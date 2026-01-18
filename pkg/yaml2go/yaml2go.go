package yaml2go

// Converter YAML 转 Go 结构体转换器接口
// 提供 YAML 字符串到 Go 结构体代码的转换功能
//
// 为什么使用接口:
//   - 定义契约: 明确转换器提供的能力
//   - 依赖倒置: 使用方依赖接口而非实现
//   - 便于测试: 可以创建 mock 实现进行单元测试
//   - 解耦: 可以轻松替换不同的实现方式
type Converter interface {
	// Convert 转换 YAML 字符串为 Go 结构体代码
	// 参数:
	//   yamlStr: YAML 格式的字符串
	// 返回:
	//   string: 生成的 Go 结构体代码
	//   error: 转换失败时的错误
	// 业务流程:
	//   1. 解析 YAML 字符串
	//   2. 推断字段类型
	//   3. 构建结构体树
	//   4. 生成 Go 代码
	Convert(yamlStr string) (string, error)

	// ConvertToFile 转换 YAML 并写入文件
	// 参数:
	//   yamlStr: YAML 格式的字符串
	//   outputPath: 输出文件路径
	// 返回:
	//   error: 转换或写入失败时的错误
	// 注意:
	//   - 如果文件已存在会被覆盖
	//   - 会自动创建父目录
	ConvertToFile(yamlStr string, outputPath string) error

	// SetConfig 更新配置（支持热更新）
	// 参数:
	//   config: 新的配置对象
	// 返回:
	//   error: 配置无效时返回错误
	// 注意:
	//   - 配置会立即生效
	//   - nil 配置会使用默认值
	SetConfig(config *Config) error
}

// Config 转换器配置
// 用于自定义代码生成行为
type Config struct {
	// PackageName 包名
	// 默认: "main"
	// 示例: "config", "model"
	PackageName string

	// StructName 根结构体名称
	// 默认: "Config"
	// 示例: "AppConfig", "Settings"
	StructName string

	// Tags 需要生成的标签列表
	// 默认: ["json", "yaml", "mapstructure", "toml"]
	// 可选值: "json", "yaml", "xml", "mapstructure", "toml", "validate"
	// 注意: mapstructure 是 viper 使用的标签
	Tags []string

	// UsePointer 字段是否使用指针类型
	// true: 字段类型为 *string, *int64 等
	// false: 字段类型为 string, int64 等
	// 默认: false
	// 指针的优势:
	//   - 可以区分零值和未设置
	//   - 节省内存（对于大型结构体）
	// 指针的劣势:
	//   - 使用时需要判空
	//   - 代码略显繁琐
	UsePointer bool

	// OmitEmpty 是否在标签中添加 omitempty 选项
	// true: `json:"field,omitempty"`
	// false: `json:"field"`
	// 默认: false
	// 作用: 序列化时忽略零值字段
	OmitEmpty bool

	// IndentStyle 缩进风格
	// "tab": 使用 tab 缩进（推荐）
	// "space": 使用空格缩进
	// 默认: "tab"
	IndentStyle string

	// AddComments 是否添加字段注释
	// true: 为每个字段生成注释
	// false: 不生成注释
	// 默认: false
	AddComments bool
}

// New 创建一个新的 Converter 实例
// 参数:
//
//	config: 配置对象，nil 时使用默认配置
//
// 返回:
//
//	Converter: 转换器实例
//
// 使用示例:
//
//	// 使用默认配置
//	converter := togo.New(nil)
//
//	// 自定义配置
//	converter := togo.New(&togo.Config{
//	    PackageName: "config",
//	    StructName:  "AppConfig",
//	})
func New(config *Config) Converter {
	return newConverter(config)
}
