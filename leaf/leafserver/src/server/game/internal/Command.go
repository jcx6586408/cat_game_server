package internal

import (
	"excel"
	"fmt"
	"runtime"
)

func init() {
	skeleton.RegisterCommand("echo", "echo user inputs", commandEcho)
	skeleton.RegisterCommand("rooms", "show how many rooms", commandRooms)
	skeleton.RegisterCommand("online", "show how many users", commandOnlines)
	skeleton.RegisterCommand("excelUpdate", "update excel config", commandExcelUpdate)
	skeleton.RegisterCommand("toExcel", "to question excel count", commandToQuestionExcel)
	skeleton.RegisterCommand("msg", "to question excel count", commandMsgExports)
	skeleton.RegisterCommand("go", "count goruntine number", commandGoruntine)
}

func commandEcho(args []interface{}) interface{} {
	return fmt.Sprintf("%v", args)
}

func commandToQuestionExcel(args []interface{}) interface{} {
	Questions.ToExcel()
	return "success"
}

func commandMsgExports(args []interface{}) interface{} {
	var str = "msg_list:\r\n"
	// msg.ProbufProcessor.Range(func(id uint16, t reflect.Type) {
	// 	str += fmt.Sprintf("%v|%v\r\n", id, t)
	// })
	return str
}

func commandGoruntine(args []interface{}) interface{} {
	return fmt.Sprintf("goruntineCount: %v", runtime.NumGoroutine())
}

// 看房间情况
func commandRooms(args []interface{}) interface{} {
	commons := []int{}
	prepare := []int{}
	playing := []int{}
	for _, v := range manager.Rooms {
		commons = append(commons, v.GetID())
	}
	for _, v := range battleManager.Rooms {
		prepare = append(prepare, v.GetID())
	}
	for _, v := range battleManager.PlayingRooms {
		playing = append(playing, v.GetID())
	}
	str := fmt.Sprintln(
		fmt.Sprintf(" commonRooms: %v\r\n", commons),
		fmt.Sprintf("id_Array: %v\r\n", manager.IDManager.Ids),
		"blank--------------------------------------------------------------------: \r\n",
		fmt.Sprintf("battleRooms: %v\r\n", prepare),
		fmt.Sprintf("battlePlayingRooms: %v\r\n", playing),
		fmt.Sprintf("battle_id_Array: %v\r\n", battleManager.IDManager.Ids))
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
