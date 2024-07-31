package pkg

import (
	"strconv"
)

// StringToUint 将字符串转换为 uint 类型，错误时返回 0
func StringToUint(s string) uint {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(i)
}

// StringToInt 将字符串转换为 int 类型，错误时返回 0
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// StringToFloat32 将字符串转换为 float32 类型，错误时返回 0
func StringToFloat32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}

// StringToFloat64 将字符串转换为 float64 类型，错误时返回 0
func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
