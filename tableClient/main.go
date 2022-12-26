package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gorilla/websocket"
	"github.com/name5566/leaf/log"
)

type Back struct {
	BackTable *BackTable `json:"BackTable"`
}

type BackTable struct {
	Arr  []Z `json:"Arr"`
	Name string
}

type Z struct {
	Score  float64     `json:"Score"`
	Member interface{} `json:"Member"`
}

func main() {
	url := "ws://localhost:3653/" //服务器地址
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("%v", err)
	}
	url = os.Args[1]
	min, _ := strconv.Atoi(os.Args[2])
	max, _ := strconv.Atoi(os.Args[3])

	data := []byte(fmt.Sprintf(`{
		"TableCount": {"min": %d, "Max": %d}
	}`, min, max))

	e := ws.WriteMessage(websocket.BinaryMessage, data)
	if e != nil {
		log.Debug("报错：发送消息失败%v", e)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGINT)

	go func() {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				log.Fatal("%v", err)
			}
			tableData := &Back{}
			err = json.Unmarshal(data, tableData)
			if err != nil {
				panic("解析json文件出错")
			}
			log.Debug("receive: %v", string(data))
			if tableData.BackTable != nil {
				log.Debug("数组长度: %d", len(tableData.BackTable.Arr))
				excelData := [][]string{}
				for _, v := range tableData.BackTable.Arr {
					log.Debug("%v|%v", v.Member, v.Score)
					str := strings.Split(fmt.Sprintf("%v", v.Member), "_")
					subData := []string{str[0], str[1], fmt.Sprintf("%v", v.Score)}
					excelData = append(excelData, subData)
				}
				CreateXlS(excelData, tableData.BackTable.Name, []string{"tableName", "questionID", "count"})
				if tableData.BackTable.Name == "fail" {
					os.Exit(2)
					return
				}
			}
		}
	}()

	s := <-ch
	switch s {
	case syscall.SIGINT:
		//SIGINT 信号，在程序关闭时会收到这个信号
		fmt.Println("SIGINT", "退出程序，执行退出前逻辑")
		return
	case syscall.SIGKILL:
		fmt.Println("SIGKILL")
		return
	default:
		fmt.Println("default")
	}
}

// CreateXlS data为要写入的数据,fileName 文件名称, headerNameArray 表头数组
func CreateXlS(data [][]string, fileName string, headerNameArray []string) {
	f := excelize.NewFile()
	sheetName := "sheet1"
	sheetWords := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U",
		"V", "W", "X", "Y", "Z",
	}

	for k, v := range headerNameArray {
		f.SetCellValue(sheetName, sheetWords[k]+"1", v)
	}

	//设置列的宽度
	f.SetColWidth("Sheet1", "A", sheetWords[len(headerNameArray)-1], 18)
	headStyleID, _ := f.NewStyle(`{
   "font":{
      "color":"#333333",
      "bold":false,
      "family":"arial"
   },
   "alignment":{
      "vertical":"center",
      "horizontal":"center"
   }
}`)
	//设置表头的样式
	f.SetCellStyle(sheetName, "A1", sheetWords[len(headerNameArray)-1]+"1", headStyleID)
	line := 1
	// 循环写入数据
	for _, v := range data {
		line++
		for kk, _ := range headerNameArray {
			f.SetCellValue(sheetName, sheetWords[kk]+strconv.Itoa(line), v[kk])
		}
	}
	// 保存文件
	if err := f.SaveAs(fileName + ".xlsx"); err != nil {
		fmt.Println(err)
	}
}
