package catLog

import (
	"fmt"

	"github.com/fatih/color"
)

func Fatal(msg ...interface{}) {
	println(color.RedString("致命错误: %v", msg))
	panic(msg)
}

func Log(msg ...interface{}) {
	fmt.Printf("\n日志: %v", msg)
}

func Warn(msg ...interface{}) {
	var s = color.YellowString("警告: ")
	fmt.Printf("\n%v%v", s, color.YellowString("%v", msg))
}
