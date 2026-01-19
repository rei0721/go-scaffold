package storage

import "errors"

var (
	// ErrInvalidConfig 无效配置错误
	ErrInvalidConfig = errors.New("fileservice: invalid configuration")

	// ErrPathNotFound 路径不存在错误
	ErrPathNotFound = errors.New("fileservice: path not found")

	// ErrNotDirectory 非目录错误
	ErrNotDirectory = errors.New("fileservice: not a directory")

	// ErrNotFile 非文件错误
	ErrNotFile = errors.New("fileservice: not a file")

	// ErrWatcherNotFound 监听器不存在错误
	ErrWatcherNotFound = errors.New("fileservice: watcher not found")

	// ErrInvalidFSType 无效的文件系统类型
	ErrInvalidFSType = errors.New("fileservice: invalid filesystem type")

	// ErrWatcherAlreadyExists 监听器已存在错误
	ErrWatcherAlreadyExists = errors.New("fileservice: watcher already exists for this path")
)
