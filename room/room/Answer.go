package room

import "proto/msg"

type LibAnswer struct {
	Answers  []*msg.Question
	Progress int
}

func ToAnswerLib(table string) *LibAnswer {
	tableConf, ok := Manager.TableManager.GetTable(table)
	arr := []*msg.Question{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &msg.Question{}
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
	return &LibAnswer{
		Answers:  arr,
		Progress: 0,
	}
}
