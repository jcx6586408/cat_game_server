package room

// 生成房间ID
func genRoomID(done chan interface{}) chan int {
	c := make(chan int)
	n := 0
	go func() {
		for {
			select {
			case <-done:
				return
			case c <- n:
				n = n + 1
			}

		}
	}()
	return c
}
