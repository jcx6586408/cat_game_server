package internal

import (
	"math/rand"
	"sort"
)

type Icon struct {
	ID      int
	Picture string
}

func ToIconLib() []*Icon {
	tableConf, ok := manager.TableManager.GetTable("picture")
	arr := []*Icon{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &Icon{}
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					obj.ID = cell[index].(int)

				case "picture":
					obj.Picture = cell[index].(string)

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

func RandIcon(count int, arr []*Icon) []*Icon {
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
	return arr[startIndex : startIndex+count]
}

func RandIconClip(count int, arr []*Icon) ([]*Icon, []*Icon) {
	// 获得其实随机数
	startIndex := 0
	length := len(arr)
	if len(arr)-count > 0 {
		startIndex = rand.Intn(length - count)
	}

	// 长度不够则补充
	if count >= length && len(arr) > 0 {
		for i := 0; i < count; i++ {
			arr = append(arr, arr[0])
		}
	}
	endIndex := startIndex + count
	subArr := arr[startIndex:endIndex]
	ex := []*Icon{}
	ex = append(ex, subArr...)
	other := append(arr[:startIndex], arr[endIndex:]...)
	// 返回子数组
	return ex, other
}
