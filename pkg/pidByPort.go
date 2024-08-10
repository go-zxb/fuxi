package pkg

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
)

func PIDByPort(port int) (string, error) {
	// 执行 netstat -ano | findstr :8080 命令
	cmd := exec.Command("cmd", "/C", "netstat -ano | findstr :"+fmt.Sprintf("%d", port))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "", nil
	}

	// 解析命令输出
	output := out.String()
	re := regexp.MustCompile(`\s*\S+\s+\S+\s+\S+\s+\S+\s+(\d+)`)
	matches := re.FindAllStringSubmatch(output, -1)

	for _, match := range matches {
		if len(match) > 1 {
			return match[1], nil
		}
	}
	return "", nil
}
