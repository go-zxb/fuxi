package pkg

import "os"

// PathExists 检查路径是否存在。
// `path` 是要检查的路径。
// 返回值：如果路径存在返回true，否则返回false。
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil // 路径存在
	}
	if os.IsNotExist(err) {
		return false, nil // 路径不存在
	}
	return false, err // 其他错误
}
