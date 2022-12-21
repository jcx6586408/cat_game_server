package internal

import "github.com/name5566/leaf/log"

type IDManager struct {
	Ids   []int
	Count int
}

func NewIDManager() *IDManager {
	i := &IDManager{}
	i.Ids = []int{}
	return i
}

func (o *IDManager) Get() int {
	if len(o.Ids) > 0 {
		id := o.Ids[len(o.Ids)-1]
		o.Ids = o.delete(o.Ids, id)
		log.Debug("ID列表: %v", o.Ids)
		return id
	}
	o.Count++
	return o.Count
}

func (o *IDManager) Put(id int) {
	for _, v := range o.Ids {
		if v == id {
			return
		}
	}
	o.Ids = append(o.Ids, id)
}

func (o *IDManager) delete(a []int, elem int) []int {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}
