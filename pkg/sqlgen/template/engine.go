// Package template 提供代码生成模板引擎
package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// TemplateEngine 模板引擎接口 (本地定义避免循环引用)
type TemplateEngine interface {
	Render(tmplName string, data interface{}) (string, error)
	RegisterFunc(name string, fn interface{}) error
	LoadTemplate(name string, content string) error
	LoadTemplateFile(name string, path string) error
}

// engine 模板引擎实现
type engine struct {
	templates *template.Template
	funcMap   template.FuncMap
}

// NewEngine 创建模板引擎
func NewEngine() TemplateEngine {
	e := &engine{
		funcMap: make(template.FuncMap),
	}

	// 注册默认函数
	e.registerDefaultFuncs()

	// 创建模板
	e.templates = template.New("").Funcs(e.funcMap)

	// 加载内置模板
	e.loadBuiltinTemplates()

	return e
}

// registerDefaultFuncs 注册默认模板函数
func (e *engine) registerDefaultFuncs() {
	e.funcMap["toCamel"] = toCamelCase
	e.funcMap["toPascal"] = toPascalCase
	e.funcMap["toSnake"] = toSnakeCase
	e.funcMap["toPlural"] = toPlural
	e.funcMap["toSingular"] = toSingular
	e.funcMap["lower"] = toLower
	e.funcMap["upper"] = toUpper
	e.funcMap["title"] = toTitle
	e.funcMap["join"] = join
	e.funcMap["contains"] = contains
	e.funcMap["hasPrefix"] = hasPrefix
	e.funcMap["hasSuffix"] = hasSuffix
	e.funcMap["add"] = add
	e.funcMap["sub"] = sub
}

// loadBuiltinTemplates 加载内置模板
func (e *engine) loadBuiltinTemplates() {
	// 加载模型模板
	e.templates.New("model").Parse(modelTemplate)

	// 加载 DAO 模板
	e.templates.New("dao").Parse(daoTemplate)
}

// Render 渲染模板
func (e *engine) Render(tmplName string, data interface{}) (string, error) {
	t := e.templates.Lookup(tmplName)
	if t == nil {
		return "", fmt.Errorf("template render error: template '%s' not found", tmplName)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template render error: %v", err)
	}

	return buf.String(), nil
}

// RegisterFunc 注册自定义模板函数
func (e *engine) RegisterFunc(name string, fn interface{}) error {
	e.funcMap[name] = fn
	e.templates = e.templates.Funcs(e.funcMap)
	return nil
}

// LoadTemplate 加载模板
func (e *engine) LoadTemplate(name string, content string) error {
	_, err := e.templates.New(name).Parse(content)
	if err != nil {
		return fmt.Errorf("template render error: failed to parse template '%s': %v", name, err)
	}
	return nil
}

// LoadTemplateFile 从文件加载模板
func (e *engine) LoadTemplateFile(name string, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("template render error: failed to read template file '%s': %v", path, err)
	}
	return e.LoadTemplate(name, string(content))
}
