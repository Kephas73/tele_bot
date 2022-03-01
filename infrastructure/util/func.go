package util

import (
	"fmt"
	"runtime"
)

func FuncName() string {
	pc, _, line, _ := runtime.Caller(1)
	result := fmt.Sprintf("%s:%v", runtime.FuncForPC(pc).Name(), line)
	return result
}
