package internal

import (
	"sync"

	"github.com/name5566/leaf/log"
)

var (
	ROOMFULL    = 1 // 满员
	ROOMNULL    = 2 // 没找到对应得房间
	ROOMSTARTED = 3 // 房间已经开始游戏

	MEMEBERPREPARE = 1 // 成员正在等待
	MEMBERPLAYING  = 2 // 成员正在游玩

	Skins      []*Skin
	NamesLib   []*Names
	IconLib    []*Icon
	AnswerLibs []Answers

	RoomManager   Managerer
	manager       *Manager
	BattleManager BattleRoomManagerer
	battleManager *BattleRoomManager
	results       = []string{"A", "B", "C", "D"}
)

func init() {
	manager = new(Manager)
	manager.Pool = sync.Pool{
		New: func() any {
			r := &Room{}
			return r
		},
	}

	battleManager = new(BattleRoomManager)
	battleManager.Pool = sync.Pool{
		New: func() any {
			r := &BattleRoom{}
			return r
		},
	}
	RoomManager = manager
	BattleManager = battleManager
	ExcelConfigUpdate()
}

func ExcelConfigUpdate() {
	Skins = ToSkinLib()
	NamesLib = ToNameLib()
	IconLib = ToIconLib()
	AnswerLibs = []Answers{}
	AnswerLibs = append(AnswerLibs, ToAnswerLib("question1"))
	log.Debug("皮肤数量: %v", len(Skins))
	log.Debug("名字数量: %v", len(NamesLib))
	log.Debug("Icon数量: %v", len(IconLib))
}
