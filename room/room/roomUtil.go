package room

// 生成房间ID
func genRoomID() chan int {
	c := make(chan int)
	n := 0
	go func() {
		for {
			c <- n
			n = n + 1
		}
	}()
	return c
}
