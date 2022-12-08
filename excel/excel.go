package excel

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/google/uuid"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ExcelConfig struct {
	Name           string                   `json:"name"`
	Excel          map[string][]interface{} `json:"excel"`
	AttributeNames []string                 `json:"attributeNames"`
}

type RequestConfig struct {
	Name string `json:"name"`
}

type ExcelManager struct {
	Tables map[string]*ExcelConfig
}

func (m *ExcelManager) ToString() string {
	var str = "\n"
	for _, v := range m.Tables {
		str += v.Name + "\n"
	}
	return str
}

func (m *ExcelManager) GetTable(name string) (*ExcelConfig, bool) {
	t, ok := m.Tables[name]
	return t, ok
}

func (m *ExcelManager) GetCell(tableName string, attributeName string, line int) interface{} {
	return nil
}

func Read() *ExcelManager {
	tableMap := make(map[string]*ExcelConfig)
	err := filepath.Walk("./table",
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			bo, err := regexp.MatchString(".xlsx", f.Name())

			if err != nil {
				panic("匹配出错")
			}

			if bo {
				fileExcel, err := excelize.OpenFile(path)
				if err != nil {
					panic(err)
				}

				for _, sheetName := range fileExcel.GetSheetMap() {
					// 创建字典
					fileMap := make(map[string][]interface{})

					tipsNames := []string{}
					attributeNames := []string{}
					attributeTypes := []string{}
					for i, row := range fileExcel.GetRows(sheetName) {
						if i == 0 {
							tipsNames = row
							continue
						}

						if i == 1 {
							attributeNames = row
							continue
						}

						if i == 2 {
							attributeTypes = row
							continue
						}

						sMap := make(map[string]interface{})
						aMap := []interface{}{}
						nMap := []string{}
						for j, s := range row {
							isTip, err := regexp.MatchString(`^#.*`, tipsNames[j])
							// println("是否是注释列：", isTip)
							if err != nil {
								panic("对比注释字段文字出错")
							}

							if isTip {
								continue
							}
							nMap = append(nMap, attributeNames[j])
							switch attributeTypes[j] {
							case "string":
								sMap[attributeNames[j]] = s
								aMap = append(aMap, s)
							case "int":
								v, err := strconv.Atoi(s)
								if err != nil {
									println("整型", sheetName, v, s, attributeNames[j])
									panic("字段解析错误")
								}
								sMap[attributeNames[j]] = v
								aMap = append(aMap, v)
							case "float":
								v, err := strconv.ParseFloat(s, 64)
								if err != nil {
									println("浮点型", sheetName, v)
									panic("字段解析错误")
								}

								sMap[attributeNames[j]] = v
								aMap = append(aMap, v)
							case "arraystring":
								v := strings.Split(s, "#")
								sMap[attributeNames[j]] = v
								aMap = append(aMap, v)
							case "arrayint":
								v := strings.Split(s, "#")
								intvs := []int{}
								for _, sid := range v {
									intv, err := strconv.Atoi(sid)
									if err == nil {
										intvs = append(intvs, intv)
									}
								}
								sMap[attributeNames[j]] = intvs
								aMap = append(aMap, intvs)
							case "arrayfloat":
								v := strings.Split(s, "#")
								intvs := []float64{}
								for _, sid := range v {
									intv, err := strconv.ParseFloat(sid, 64)
									if err == nil {
										intvs = append(intvs, intv)
									}
								}
								sMap[attributeNames[j]] = intvs
								aMap = append(aMap, intvs)
							}

						}
						fileMap[row[0]] = aMap
						excelData := &ExcelConfig{}
						excelData.Name = sheetName
						excelData.AttributeNames = nMap
						excelData.Excel = fileMap
						tableMap[sheetName] = excelData
						// catLog.Log("读取表格", sheetName)
					}
				}

			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return &ExcelManager{
		Tables: tableMap,
	}
}
