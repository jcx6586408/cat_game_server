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
	skeleton.RegisterCommand("toExcel", "to question excel count", commandToQuestionExcel)
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}

func commandToQuestionExcel(args []interface{}) interface{} {
	Questions.ToExcel()
	return "success"
}

// 看房间情况
func commandRooms(args []interface{}) interface{} {
	prepare := []int{}
	playing := []int{}
	for _, v := range battleManager.Rooms {
		prepare = append(prepare, v.GetID())
	}
	for _, v := range battleManager.PlayingRooms {
		playing = append(playing, v.GetID())
	}
	str := fmt.Sprintln(
		fmt.Sprintf("commonRooms: %v\n", len(manager.Rooms)),
		fmt.Sprintf("battleRooms: %v\n", prepare),
		fmt.Sprintf("battlePlayingRooms: %v\n", playing),
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
