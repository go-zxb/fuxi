package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func Info(info string) {
	log.Println(info)
}

func InfoWithPath(path string, info string) {
	_ = log.Output(1, fmt.Sprintf("✅  %s:%d  %s", path, 1, info))
}

func Error(err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok && err != nil {
		_ = log.Output(1, fmt.Sprintf("❌ %s:%d  %s", file, line, err.Error()))
	}
}

func InfoOutPath(info string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		_ = log.Output(1, fmt.Sprintf("ℹ️ %s:%d  %s", file, line, info))
	}
}

func Fatalln(info string, err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok && err != nil {
		_ = log.Output(1, fmt.Sprintf("🔚 %s:%d  %s  %s", file, line, info, err.Error()))
		os.Exit(1)
	}
}
