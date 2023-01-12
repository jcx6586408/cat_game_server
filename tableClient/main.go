package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leafserver/src/server/msg"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

func GetCityInfo() {
	host := "http://api01.aliyun.venuscn.com"
	path := "/ip"
	url := host + path
	querys := "ip=218.18.228.178"
	appcode := "c7adf888186e4ceba75f3841a009b17f"
	url = url + "?" + querys
	println("请求路径:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println("============================")
		panic(err)
	}
	req.Header.Add("Authorization", "APPCODE "+appcode)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		println("============================")
		panic(err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	println("请求返回：", response.Status, len(body), string(body))
}

func main() {
	Connect()
}

type BackNew struct {
	TableGet *msg.TableGet
}

func Connect() {
	url := "ws://localhost:3653/" //服务器地址
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("%v", err)
	}
	if len(os.Args) > 1 {
		url = os.Args[1]
	}
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
	count := 0
	go func() {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				log.Fatal("%v", err)
			}
			tableData := &BackNew{}
			err = json.Unmarshal(data, tableData)
			if err != nil {
				panic("解析json文件出错")
			}
			log.Debug("返回数据: %v", tableData)
			if tableData.TableGet != nil {
				log.Debug("数组长度: %d", len(tableData.TableGet.Questions))
				excelData := [][]string{}
				for _, v := range tableData.TableGet.Questions {
					// str := strings.Split(fmt.Sprintf("%v", v.Member), "_")
					subData := []string{fmt.Sprintf("%v", v.ID), fmt.Sprintf("%v", v.Win), fmt.Sprintf("%v", v.Fail)}
					excelData = append(excelData, subData)
				}
				if len(tableData.TableGet.Questions) > 0 {
					log.Debug("写入表格==============")
					CreateXlS(excelData, os.Args[4], []string{"ID", "rightNumber", "wrongNumber"})
					log.Debug("写入表格==============完成")
				}
				count++
				ws.Close()
				os.Exit(2)
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
