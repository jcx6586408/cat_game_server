package internal

import (
	"sort"

	"github.com/name5566/leaf/log"
)

type Level struct {
	ID    int
	Level int
}

func ToLevelLib() []*Level {
	tableConf, ok := manager.TableManager.GetTable("level")
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

				}
			}
			arr = append(arr, obj)
		}
	}
	sort.SliceStable(arr, func(i, j int) bool {
		return arr[i].Level < arr[j].Level
	})
	for _, v := range arr {
		log.Debug("%v", v.ID)
	}
	return arr
}

func GetIDByLevel(level int) int {
	for _, v := range LevelLib {
		if v.Level > level {
			return v.ID - 1
		}
	}
	return LevelLib[len(LevelLib)-1].ID
}
