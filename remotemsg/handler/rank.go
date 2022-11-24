package handler

import (
	"context"
	"fmt"
	"proto/msg"
	"remotemsg"
	"server/client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Rank struct {
	cat         *CatClass
	Conn        *grpc.ClientConn
	innerClient msg.HelloClient
}

func NewRank() *Rank {
	s := Rank{}
	s.cat = &CatClass{}
	s.cat.New()
	AddModel(s.cat)
	return &s
}

var RankInstance *Rank = NewRank()

func (s *Rank) Run(port string) {
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("connect server failed,", err)
	}
	s.Conn = conn
	s.innerClient = msg.NewHelloClient(s.Conn)
	// 注册消息
	s.cat.Register(remotemsg.RANKUPDATE, update)
}

func update(data client.Msg) {
	c := msg.NewHelloClient(RankInstance.Conn)
	r1, err := c.SayHello(context.Background(), &msg.HelloRequest{})
	if err != nil {
		fmt.Println("can't get version,", err)
		return
	}
	fmt.Println("返回消息", r1.Message)
	// msg.Client.Write(constVal.GLOVEINFO, GloveInfoInstance)
}
