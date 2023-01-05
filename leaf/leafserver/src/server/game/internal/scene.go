package internal

import (
	"math/rand"
	"sort"
)

type Scene struct {
	ID    int
	Scene string
}

func ToSceneLib() []*Scene {
	tableConf, ok := manager.TableManager.GetTable("scene")
	arr := []*Scene{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &Scene{}
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					obj.ID = cell[index].(int)

				case "scene":
					obj.Scene = cell[index].(string)

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

func RandScene(arr []*Scene) *Scene {
	index := rand.Intn(len(arr))
	return arr[index]
}
