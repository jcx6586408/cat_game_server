package internal

import (
	"config"
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
	wxConf = config.Read()
}

func (m *Module) OnDestroy() {

}
