package log

import (
	"fmt"
)

// Info ...
func Info(value ...interface{}) {
	fmt.Println("\033[1;32m")
	for _, val := range value {
		fmt.Printf("%v ", val)
	}
	fmt.Println("\033[0m")
}

// Error ...
func Error(value ...interface{}) {
	fmt.Println("\033[1;31m")
	for _, val := range value {
		fmt.Printf("%v ", val)
	}
	fmt.Println("\033[0m")
}

// Warning ...
func Warning(value ...interface{}) {
	fmt.Println("\033[1;33m")
	for _, val := range value {
		fmt.Printf("%v ", val)
	}
	fmt.Println("\033[0m")
}
