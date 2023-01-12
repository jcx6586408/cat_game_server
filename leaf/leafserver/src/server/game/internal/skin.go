package internal

type Skin struct {
	ID      int
	Name    string
	Picture string
	Group   string
	Method  int
}

func ToSkinLib() []*Skin {
	tableConf, ok := manager.TableManager.GetTable("skin")
	arr := []*Skin{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &Skin{}
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					obj.ID = cell[index].(int)
				case "name":
					obj.Name = cell[index].(string)
				case "picture":
					obj.Picture = cell[index].(string)
				case "group":
					obj.Group = cell[index].(string)
				}
			}
			arr = append(arr, obj)
		}
	}
	return arr
}
