package internal

import "math/rand"

type Names struct {
	ID   int
	Name string
}

func ToNameLib() []*Names {
	tableConf, ok := manager.TableManager.GetTable("name")
	arr := []*Names{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &Names{}
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					obj.ID = cell[index].(int)

				case "name":
					obj.Name = cell[index].(string)

				}
			}
			arr = append(arr, obj)
		}
	}
	return arr
}

func RandName(count int, arr []*Names) []*Names {
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

	// 返回子数组
	return arr[startIndex : startIndex+count]
}

func RandNameClip(count int, arr []*Names) ([]*Names, []*Names) {
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

	r := append(arr[:startIndex], arr[startIndex+count+1:]...)

	// 返回子数组
	return arr[startIndex : startIndex+count], r
}
