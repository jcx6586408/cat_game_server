package internal

import (
	"errors"
	"excel"
	"sync"

	"github.com/labstack/gommon/log"
)

type BattleRoomManagerer interface {
	Create() BattleRoomer                // 创建房间
	Destroy(roomID int)                  // 通过ID销毁房间
	GetRoomByID(roomID int) BattleRoomer // 通过ID获取已经创建房间
	GetRoomsCount() int                  // 获取当前所有房间数量

	// 运行处理
	Run(roomID int) // 开始运行房间（用于等待成员加入）

	// 游戏
	Play(roomID int)    // 开始游戏
	PlayEnd(roomID int) // 结束游戏

	// 复活
	Relive(roomID int, member Memberer) // 复活

	// 房间处理
	AddRoom(roomID int, room Roomer) (BattleRoomer, error) // 添加成员

	MatchRoom(room Roomer) BattleRoomer // 匹配房间
	MatchRoomzCancel(room Roomer)       // 取消匹配房间
}

type BattleRoomManager struct {
	Pool         sync.Pool
	Rooms        []BattleRoomer
	Done         chan interface{}    // 停用通道
	TableManager *excel.ExcelManager // 表格管理器
}

func (m *BattleRoomManager) Create() BattleRoomer {
	r := manager.Pool.Get().(*BattleRoom)
	log.Debug("创建了战斗房间, ID: %d", r.ID)
	r.OnInit()
	return r
}

func (m *BattleRoomManager) GetRoomsCount() int {
	return len(m.Rooms)
}

func (m *BattleRoomManager) GetRoomByID(roomID int) BattleRoomer {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			return v
		}
	}
	return nil
}

func (m *BattleRoomManager) destroy(room BattleRoomer) bool {
	room.OnClose()
	manager.Pool.Put(room)
	log.Debug("战斗房间回收, ID: %d", room.GetID())
	return true
}

func (m *BattleRoomManager) Destroy(roomID int) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			m.Rooms = m.delete(m.Rooms, v) // 移除
			m.destroy(v)                   // 回收
			return
		}
	}

}

func (m *BattleRoomManager) Run(roomID int) {

}

func (m *BattleRoomManager) Play(roomID int) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			v.Play()
			return
		}
	}
}

func (m *BattleRoomManager) PlayEnd(roomID int) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			v.OnPlayEnd()
			return
		}
	}
}

func (m *BattleRoomManager) Relive(roomID int, member Memberer) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			v.Relive(member)
			return
		}
	}
}

func (m *BattleRoomManager) AddRoom(roomID int, room Roomer) (BattleRoomer, error) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			v.AddRoom(room)
			return v, nil
		}
	}
	return nil, errors.New("")
}

func (m *BattleRoomManager) MatchRoom(room Roomer) BattleRoomer {
	if len(m.Rooms) <= 0 {
		br := m.Create()
		br.AddRoom(room)
		return br
	} else {
		for _, v := range m.Rooms {
			bo := v.AddRoom(room)
			if bo {
				return v
			}
		}
		// 如果找不到合适的战斗房，则创建新的
		br := m.Create()
		br.AddRoom(room)
		return br
	}
}

func (m *BattleRoomManager) MatchRoomzCancel(room Roomer) {

}

func (m *BattleRoomManager) delete(a []BattleRoomer, elem BattleRoomer) []BattleRoomer {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}
