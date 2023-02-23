package internal

import (
	"sort"
)

type Score struct {
	ID    int // 时间
	Score int // 分数
}

func ToScoreLib() []*Score {
	tableConf, ok := manager.TableManager.GetTable("score")
	arr := []*Score{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &Score{}
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					obj.ID = cell[index].(int)
				case "score":
					obj.Score = cell[index].(int)
				}
			}
			arr = append(arr, obj)
		}
	}
	sort.SliceStable(arr, func(i, j int) bool {
		return arr[i].ID < arr[j].ID
	})
	return arr
}
