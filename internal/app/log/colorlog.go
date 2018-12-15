package log

import (
	"fmt"
)

// Info ...
func Info(value ...interface{}) {
	fmt.Println("\033[1;32m", value, "\033[0m")
}

// Error ...
func Error(value ...interface{}) {
	fmt.Println("\033[1;31m", value, "\033[0m")
}

// Warning ...
func Warning(value ...interface{}) {
	fmt.Println("\033[1;33m", value, "\033[0m")
}
