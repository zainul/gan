package log

import (
	"fmt"
)

// Info is colored log with green color
func Info(value ...interface{}) {
	fmt.Println("\033[1;32m")
	for _, val := range value {
		fmt.Printf("%v ", val)
	}
	fmt.Println("\033[0m")
}

// Error is colored log with red color
func Error(value ...interface{}) {
	fmt.Println("\033[1;31m")
	for _, val := range value {
		fmt.Printf("%v ", val)
	}
	fmt.Println("\033[0m")
}

// Warning is colored log with yellow color
func Warning(value ...interface{}) {
	fmt.Println("\033[1;33m")
	for _, val := range value {
		fmt.Printf("%v ", val)
	}
	fmt.Println("\033[0m")
}
