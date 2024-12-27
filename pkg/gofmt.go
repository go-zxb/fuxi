package pkg

// FormatGoCode 格式化Go代码
func FormatGoCode(filePath string) error {
	return RunCommand("gofmt", "-w", filePath)
}
