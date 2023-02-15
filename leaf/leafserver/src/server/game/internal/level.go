package internal

import (
	"sort"
)

type Level struct {
	ID             int   // 段位
	Level          int   // 等级
	QuestionNumber int   // 答题数量
	QuestionTime   int   // 答题时间
	WinRate        []int // 胜率
	AnswerPhase    []int // 对应题库
	RightNumber    int   // 正确数量
	WrongNumber    int   // 错误数量
}

func ToLevelLib() []*Level {
	tableConf, ok := manager.TableManager.GetTable("Rank")
	arr := []*Level{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &Level{}
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					obj.ID = cell[index].(int)
				case "level":
					obj.Level = cell[index].(int)
				case "questionNumber":
					obj.QuestionNumber = cell[index].(int)
				case "questionTime":
					obj.QuestionTime = cell[index].(int)
				case "winRate":
					obj.WinRate = cell[index].([]int)
				case "AnswerPhase":
					obj.AnswerPhase = cell[index].([]int)
				}
			}
			arr = append(arr, obj)
		}
	}
	sort.SliceStable(arr, func(i, j int) bool {
		return arr[i].Level < arr[j].Level
	})
	return arr
}

func GetIDByLevel(level int) int {
	if level == 0 {
		return LevelLib[0].ID
	}
	for _, v := range LevelLib {
		if level <= v.Level {
			if level == v.Level {
				return v.ID + 1
			}
			return v.ID
		}
	}
	return LevelLib[len(LevelLib)-1].ID
}

func GetMaxLevel() int {
	max := 0
	for _, v := range LevelLib {
		if v.Level > max {
			max = v.Level
		}
	}
	return max
}
