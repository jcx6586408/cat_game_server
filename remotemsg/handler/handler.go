package handler

import (
	"catLog"
	"server/client"
)

var CatModels = []CatModule{}

type CatModule interface {
	Clear() // 服务器关闭时进行清理
	SetClient(*client.Client)
	GetDone() chan interface{}
	GetOfflineChan() chan string
}

type CatClass struct {
	Done        chan interface{} // 关闭chan
	OffLineChan chan string      // 玩家下线通知
	Client      *client.Client   // 每个模块都存有客户端连接引用
}

func AddModel(m CatModule) {
	CatModels = append(CatModels, m)
	catLog.Log("加入模块", len(CatModels))
}

func (cat *CatClass) New() {
	cat.Done = make(chan interface{})
	cat.OffLineChan = make(chan string)
}

func (s *CatClass) GetDone() chan interface{} {
	return s.Done
}

func (s *CatClass) SetClient(c *client.Client) {
	s.Client = c
}

func (s *CatClass) GetOfflineChan() chan string {
	return s.OffLineChan
}

func RegisterOffline(c *client.Client) {
	for _, v := range CatModels {
		v.SetClient(c)
		go func(m CatModule) {
			for {
				select {
				case <-m.GetDone():
					return
				case <-c.C:
					m.GetOfflineChan() <- c.Uuid // 通知下线
					return
				}
			}
		}(v)
	}
}

func (cat *CatClass) Clear() {
	close(cat.Done)
}

func (cat *CatClass) Register(id int, msgHandler func(msg client.Msg)) {
	go func() {
		// 注册消息
		c := make(chan client.Msg)
		handler := &client.MsgHandler{}
		handler.Chan = c
		handler.MsgID = id
		client.RegisterHandler(handler)

		for {
			select {
			case <-cat.Done:
				return
			case msg := <-c:
				msgHandler(msg)
			}
		}
	}()
}
