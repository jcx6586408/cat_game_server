package internal

type IDManager struct {
	Ids   []int
	Count int
}

func (o *IDManager) Get() int {
	if len(o.Ids) > 0 {
		id := o.Ids[len(o.Ids)-1]
		o.Ids = o.delete(o.Ids, id)
		return id
	}
	o.Count++
	return o.Count
}

func (o *IDManager) Put(id int) {
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
