package sqlgen

import (
	"strings"
	"unicode"
)

// ToSnakeCase 将字符串转换为蛇形命名法
// 例如: UserName -> user_name, HTTPRequest -> http_request
func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	result.Grow(len(s) + 5) // 预分配一些额外空间

	for i, r := range s {
		if unicode.IsUpper(r) {
			// 不是第一个字符，且前一个字符不是下划线
			if i > 0 && s[i-1] != '_' {
				// 如果前一个是小写，或者后一个是小写（处理连续大写如 HTTPRequest）
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

// ToCamelCase 将字符串转换为驼峰命名法 (首字母小写)
// 例如: user_name -> userName, USER_NAME -> userName
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if pascal == "" {
		return pascal
	}
	// 首字母小写
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// ToPascalCase 将字符串转换为帕斯卡命名法 (首字母大写)
// 例如: user_name -> UserName, USER_NAME -> UserName
func ToPascalCase(s string) string {
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

// ToPlural 将英文单词转换为复数形式 (简化版)
// 支持基本规则，不处理不规则变化
func ToPlural(s string) string {
	if s == "" {
		return s
	}

	// 基本规则
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

// ToSingular 将英文单词转换为单数形式 (简化版)
func ToSingular(s string) string {
	if s == "" {
		return s
	}

	switch {
	case strings.HasSuffix(s, "ies") && len(s) > 3:
		return s[:len(s)-3] + "y"
	case strings.HasSuffix(s, "es") && len(s) > 2:
		// 检查是否是 -ses, -xes, -zes, -ches, -shes
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

// isVowel 判断字符是否为元音字母
func isVowel(r rune) bool {
	switch unicode.ToLower(r) {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}

// ConvertName 根据命名规则转换名称
func ConvertName(name string, rule NamingRule) string {
	switch rule {
	case NamingSnakeCase:
		return ToSnakeCase(name)
	case NamingCamelCase:
		return ToCamelCase(name)
	case NamingPascalCase:
		return ToPascalCase(name)
	default:
		return name
	}
}
