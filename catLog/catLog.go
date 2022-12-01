package catLog

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func Fatal(msg ...interface{}) {
	var s = color.YellowString("致命错误: ")
	fmt.Printf("\n%v%v-%v", s, time.Now().Format("2006/1/02 15:04"), color.RedString("%v", msg))
	panic(msg)
}

func Err(msg ...interface{}) {
	var s = color.YellowString("错误: ")
	fmt.Printf("\n%v%v-%v", s, time.Now().Format("2006/1/02 15:04"), color.RedString("%v", msg))
}

func Log(msg ...interface{}) {
	// var s = color.YellowString("日志: ")
	// fmt.Printf("\n%v%v-%v", s, time.Now().Format("2006/1/02 15:04"), color.WhiteString("%v", msg))
}

func Warn(msg ...interface{}) {
	// var s = color.YellowString("警告: ")
	// fmt.Printf("\n%v%v-%v", s, time.Now().Format("2006/1/02 15:04"), color.YellowString("%v", msg))
}
