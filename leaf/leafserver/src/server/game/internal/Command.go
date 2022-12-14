package internal

import (
	"excel"
	"fmt"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
	// skeleton.RegisterCommand("rooms", "show how many rooms", commandRooms)
	skeleton.RegisterCommand("online", "show how many users", commandOnlines)
	skeleton.RegisterCommand("excelUpdate", "update excel config", commandExcelUpdate)
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}

// // 看房间情况
// func commandRooms(args []interface{}) interface{} {
// 	m := Manager
// 	return fmt.Sprintf("idle:%v, prepare:%v,matching:%v, using:%v", len(m.Rooms), len(m.PrepareRooms), len(m.MatchingRooms), len(m.UsingRooms))
// }

// 查看在线人数
func commandOnlines(args []interface{}) interface{} {
	return fmt.Sprintf("%v", len(Users))
}

// 更新配表
func commandExcelUpdate(args []interface{}) interface{} {
	manager.TableManager = excel.Read()
	state := "\nexcel update complete"
	ExcelConfigUpdate()
	return manager.TableManager.ToString() + state
}
