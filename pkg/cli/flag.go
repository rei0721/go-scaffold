package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// flagParser 选项解析器
type flagParser struct {
	flags  []Flag
	values map[string]interface{}
	fs     *flag.FlagSet
}

// newFlagParser 创建选项解析器
func newFlagParser(cmdName string, flags []Flag) *flagParser {
	return &flagParser{
		flags:  flags,
		values: make(map[string]interface{}),
		fs:     flag.NewFlagSet(cmdName, flag.ContinueOnError),
	}
}

// parse 解析命令行参数
func (p *flagParser) parse(args []string) ([]string, error) {
	// 禁用默认的错误输出
	p.fs.SetOutput(io.Discard)

	// 注册所有选项
	for _, f := range p.flags {
		p.registerFlag(f)
	}

	// 解析参数
	if err := p.fs.Parse(args); err != nil {
		return nil, &UsageError{Message: err.Error()}
	}

	// 提取解析后的值
	if err := p.extractValues(); err != nil {
		return nil, err
	}

	// 验证必填选项
	if err := p.validate(); err != nil {
		return nil, err
	}

	// 返回剩余的位置参数
	return p.fs.Args(), nil
}

// registerFlag 注册单个选项到 flag.FlagSet
func (p *flagParser) registerFlag(f Flag) {
	// 从环境变量获取默认值
	defaultVal := f.Default
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			defaultVal = envVal
		}
	}

	switch f.Type {
	case FlagTypeString:
		defStr := ""
		if defaultVal != nil {
			defStr = fmt.Sprint(defaultVal)
		}
		p.fs.String(f.Name, defStr, f.Description)
		if f.ShortName != "" {
			p.fs.String(f.ShortName, defStr, f.Description)
		}

	case FlagTypeInt:
		defInt := 0
		if defaultVal != nil {
			if i, ok := defaultVal.(int); ok {
				defInt = i
			} else if s, ok := defaultVal.(string); ok {
				defInt, _ = strconv.Atoi(s)
			}
		}
		p.fs.Int(f.Name, defInt, f.Description)
		if f.ShortName != "" {
			p.fs.Int(f.ShortName, defInt, f.Description)
		}

	case FlagTypeBool:
		defBool := false
		if defaultVal != nil {
			if b, ok := defaultVal.(bool); ok {
				defBool = b
			} else if s, ok := defaultVal.(string); ok {
				defBool, _ = strconv.ParseBool(s)
			}
		}
		p.fs.Bool(f.Name, defBool, f.Description)
		if f.ShortName != "" {
			p.fs.Bool(f.ShortName, defBool, f.Description)
		}

	case FlagTypeStringSlice:
		// 字符串数组暂时不使用标准库 flag，需要自定义实现
		// 这里简化处理，使用逗号分隔的字符串
		defStr := ""
		if defaultVal != nil {
			if slice, ok := defaultVal.([]string); ok {
				defStr = strings.Join(slice, ",")
			}
		}
		p.fs.String(f.Name, defStr, f.Description+" (comma-separated)")
		if f.ShortName != "" {
			p.fs.String(f.ShortName, defStr, f.Description+" (comma-separated)")
		}
	}
}

// extractValues 从 flag.FlagSet 提取解析后的值
func (p *flagParser) extractValues() error {
	for _, f := range p.flags {
		var val interface{}

		// 优先检查短选项（如果用户使用了短选项）
		// 然后再检查长选项
		var flagToUse *flag.Flag
		if f.ShortName != "" {
			if flg := p.fs.Lookup(f.ShortName); flg != nil {
				flagToUse = flg
			}
		}
		// 如果短选项没有被使用，或者没有短选项，使用长选项
		if flagToUse == nil {
			flagToUse = p.fs.Lookup(f.Name)
		}

		if flagToUse == nil {
			continue
		}

		switch f.Type {
		case FlagTypeString:
			val = flagToUse.Value.String()

		case FlagTypeInt:
			v, _ := strconv.Atoi(flagToUse.Value.String())
			val = v

		case FlagTypeBool:
			v, _ := strconv.ParseBool(flagToUse.Value.String())
			val = v

		case FlagTypeStringSlice:
			str := flagToUse.Value.String()
			if str != "" {
				val = strings.Split(str, ",")
			} else {
				val = []string{}
			}
		}

		p.values[f.Name] = val
	}

	return nil
}

// validate 验证必填选项
func (p *flagParser) validate() error {
	for _, f := range p.flags {
		if !f.Required {
			continue
		}

		val, exists := p.values[f.Name]
		if !exists {
			return &UsageError{
				Message: fmt.Sprintf("required flag --%s not provided", f.Name),
			}
		}

		// 检查是否为零值
		switch f.Type {
		case FlagTypeString:
			if val == "" {
				return &UsageError{
					Message: fmt.Sprintf("required flag --%s cannot be empty", f.Name),
				}
			}
		case FlagTypeStringSlice:
			if slice, ok := val.([]string); !ok || len(slice) == 0 {
				return &UsageError{
					Message: fmt.Sprintf("required flag --%s cannot be empty", f.Name),
				}
			}
		}
	}

	return nil
}

// getValues 返回解析后的选项值
func (p *flagParser) getValues() map[string]interface{} {
	return p.values
}
