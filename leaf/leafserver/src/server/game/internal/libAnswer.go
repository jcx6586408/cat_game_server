package internal

import (
	"math/rand"
	pmsg "proto/msg"
)

type Answers []*pmsg.Question

type LibAnswer struct {
	Answers  Answers
	Progress int
}

func GetAnswerLib() Answers {
	ran := rand.Intn(len(AnswerLibs))
	return AnswerLibs[ran]
}

func ToAnswerLib(table string) Answers {
	tableConf, ok := Manager.TableManager.GetTable(table)
	arr := []*pmsg.Question{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &pmsg.Question{}
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					obj.ID = int32(cell[index].(int))
				case "question":
					obj.Question = cell[index].(string)
				case "answerA":
					obj.AnswerA = cell[index].(string)
				case "answerB":
					obj.AnswerB = cell[index].(string)
				case "answerC":
					obj.AnswerC = cell[index].(string)
				case "answerD":
					obj.AnswerD = cell[index].(string)
				case "rightAnswer":
					obj.RightAnswer = cell[index].(string)
				}
			}
			arr = append(arr, obj)
		}
	}
	return arr
}

func RandAnswerLib(count int, arr []*pmsg.Question) *LibAnswer {
	// 获得其实随机数
	startIndex := 0
	length := len(arr)
	if len(arr)-count > 0 {
		startIndex = rand.Intn(length - count)
	}

	// 长度不够则补充
	if count >= length {
		for i := 0; i < count; i++ {
			arr = append(arr, arr[0])
		}
	}

	// 返回子数组
	return &LibAnswer{
		Answers:  arr[startIndex : startIndex+count],
		Progress: 0,
	}
}
