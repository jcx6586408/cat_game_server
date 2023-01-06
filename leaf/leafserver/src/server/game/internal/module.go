package internal

import (
	"leafserver/src/server/base"

	"github.com/name5566/leaf/module"
	// "server/base"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	// 模块初始化
	RankInit()
	ConstInit()
}

func (m *Module) OnDestroy() {

}
