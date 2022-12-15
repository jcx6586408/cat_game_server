package internal

import (
	"excel"
	"fmt"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
	skeleton.RegisterCommand("rooms", "show how many rooms", commandRooms)
	skeleton.RegisterCommand("online", "show how many users", commandOnlines)
	skeleton.RegisterCommand("excelUpdate", "update excel config", commandExcelUpdate)
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}

// 看房间情况
func commandRooms(args []interface{}) interface{} {
	str := fmt.Sprintln(
		fmt.Sprintf("commonRooms: %v\n", len(manager.Rooms)),
		fmt.Sprintf("battleRooms: %v\n", len(battleManager.Rooms)),
		fmt.Sprintf("id_Array: %v\n", manager.IDManager.Ids),
		fmt.Sprintf("battle_id_Array: %v\n", battleManager.IDManager.Ids))
	return str
}

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
