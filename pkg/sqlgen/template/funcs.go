package template

import (
	"strings"
	"unicode"
)

// 模板辅助函数

// toLower 转小写
func toLower(s string) string {
	return strings.ToLower(s)
}

// toUpper 转大写
func toUpper(s string) string {
	return strings.ToUpper(s)
}

// toTitle 首字母大写
func toTitle(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// join 连接字符串
func join(sep string, items []string) string {
	return strings.Join(items, sep)
}

// contains 检查是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// hasPrefix 检查前缀
func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// hasSuffix 检查后缀
func hasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// add 加法
func add(a, b int) int {
	return a + b
}

// sub 减法
func sub(a, b int) int {
	return a - b
}

// toCamelCase 转驼峰命名
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if pascal == "" {
		return pascal
	}
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// toPascalCase 转帕斯卡命名
func toPascalCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	result.Grow(len(s))

	capitalizeNext := true
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			capitalizeNext = true
			continue
		}
		if capitalizeNext {
			result.WriteRune(unicode.ToUpper(r))
			capitalizeNext = false
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}

	return result.String()
}

// toSnakeCase 转蛇形命名
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	result.Grow(len(s) + 5)

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && s[i-1] != '_' {
				if i+1 < len(s) && unicode.IsLower(rune(s[i+1])) ||
					(i > 0 && unicode.IsLower(rune(s[i-1]))) {
					result.WriteByte('_')
				}
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// toPlural 转复数
func toPlural(s string) string {
	if s == "" {
		return s
	}

	switch {
	case strings.HasSuffix(s, "s") ||
		strings.HasSuffix(s, "x") ||
		strings.HasSuffix(s, "z") ||
		strings.HasSuffix(s, "ch") ||
		strings.HasSuffix(s, "sh"):
		return s + "es"
	case strings.HasSuffix(s, "y") && len(s) > 1 && !isVowel(rune(s[len(s)-2])):
		return s[:len(s)-1] + "ies"
	default:
		return s + "s"
	}
}

// toSingular 转单数
func toSingular(s string) string {
	if s == "" {
		return s
	}

	switch {
	case strings.HasSuffix(s, "ies") && len(s) > 3:
		return s[:len(s)-3] + "y"
	case strings.HasSuffix(s, "es") && len(s) > 2:
		if len(s) > 3 {
			beforeEs := s[len(s)-3]
			if beforeEs == 's' || beforeEs == 'x' || beforeEs == 'z' {
				return s[:len(s)-2]
			}
		}
		if strings.HasSuffix(s, "ches") || strings.HasSuffix(s, "shes") {
			return s[:len(s)-2]
		}
		return s[:len(s)-1]
	case strings.HasSuffix(s, "s") && len(s) > 1:
		return s[:len(s)-1]
	default:
		return s
	}
}

// isVowel 判断是否元音
func isVowel(r rune) bool {
	switch unicode.ToLower(r) {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}
