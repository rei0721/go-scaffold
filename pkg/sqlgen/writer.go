package sqlgen

import (
	"os"
	"path/filepath"
)

// fileWriter 文件写入器实现
type fileWriter struct{}

// NewFileWriter 创建文件写入器
func NewFileWriter() FileWriter {
	return &fileWriter{}
}

// Write 写入文件
func (w *fileWriter) Write(path string, content []byte) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &GenerateError{
			File:    path,
			Message: "failed to create directory",
			Cause:   err,
		}
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		return &GenerateError{
			File:    path,
			Message: "failed to write file",
			Cause:   err,
		}
	}

	return nil
}

// WriteAtomic 原子性写入文件
// 先写入临时文件，然后重命名，保证写入的原子性
func (w *fileWriter) WriteAtomic(path string, content []byte) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &GenerateError{
			File:    path,
			Message: "failed to create directory",
			Cause:   err,
		}
	}

	// 创建临时文件
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, content, 0644); err != nil {
		return &GenerateError{
			File:    path,
			Message: "failed to write temp file",
			Cause:   err,
		}
	}

	// 原子重命名
	if err := os.Rename(tmpPath, path); err != nil {
		// 清理临时文件
		os.Remove(tmpPath)
		return &GenerateError{
			File:    path,
			Message: "failed to rename temp file",
			Cause:   err,
		}
	}

	return nil
}
