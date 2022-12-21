package internal

import (
	"errors"
	"excel"
	"sync"

	"github.com/name5566/leaf/log"
)

type BattleRoomManagerer interface {
	Create() BattleRoomer                // 创建房间
	Destroy(room BattleRoomer)           // 通过ID销毁房间
	GetRoomByID(roomID int) BattleRoomer // 通过ID获取已经创建房间
	GetRoomsCount() int                  // 获取当前所有房间数量

	// 游戏
	Play(BattleRoomer)  // 开始游戏
	PlayEnd(roomID int) // 结束游戏

	// 房间处理
	AddRoom(roomID int, room Roomer) (BattleRoomer, error) // 添加成员

	MatchRoom(room Roomer) BattleRoomer // 匹配房间
	MatchRoomzCancel(room Roomer) bool  // 取消匹配房间
}

type BattleRoomManager struct {
	Pool         sync.Pool
	Rooms        []BattleRoomer
	PlayingRooms []BattleRoomer
	Done         chan interface{}    // 停用通道
	TableManager *excel.ExcelManager // 表格管理器
	IDManager    *IDManager
}

func (m *BattleRoomManager) Create() BattleRoomer {
	r := m.Pool.Get().(*BattleRoom)
	r.ID = m.IDManager.Get()
	log.Debug("创建了战斗房间, ID: %d", r.ID)
	for _, v := range m.Rooms {
		log.Debug("***当前准备间: %d", v.GetID())
	}
	for _, v := range m.PlayingRooms {
		log.Debug("***当前准备间: %d", v.GetID())
	}
	m.Rooms = append(m.Rooms, r)
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
	m.IDManager.Put(room.GetID()) // id回收
	m.Pool.Put(room)
	log.Debug("战斗房间回收, ID: %d", room.GetID())
	for _, v := range m.Rooms {
		log.Debug("当前准备间---: %d", v.GetID())
	}
	for _, v := range m.PlayingRooms {
		log.Debug("当前游玩间***: %d", v.GetID())
	}
	return true
}

func (m *BattleRoomManager) Destroy(room BattleRoomer) {
	for _, v := range m.PlayingRooms {
		if v.GetID() == room.GetID() {
			log.Debug("处于游玩间被回收")
			m.PlayingRooms = m.delete(m.PlayingRooms, room) // 移除
			m.destroy(room)                                 // 回收
			return
		}
	}
	for _, v := range m.Rooms {
		if v.GetID() == room.GetID() {
			log.Debug("处于准备间被回收")
			m.Rooms = m.delete(m.Rooms, room) // 移除
			m.destroy(room)                   // 回收
			return
		}
	}

}

func (m *BattleRoomManager) Play(room BattleRoomer) {
	// 将房间剔除
	m.Rooms = m.delete(m.Rooms, room)
	m.PlayingRooms = append(m.PlayingRooms, room) // 加入战斗房
	log.Debug("战斗房ID%d: 转移到正在比赛房间, 匹配房数量: %d, 游戏房数量: %d", room.GetID(), len(m.Rooms), len(m.PlayingRooms))
}

func (m *BattleRoomManager) PlayEnd(roomID int) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			v.OnPlayEnd()
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
		log.Debug("直接创建战斗****房间: %d", br.GetID())
		return br
	} else {
		for _, v := range m.Rooms {
			bo := v.AddRoom(room)
			if bo {
				log.Debug("加入别人//////战斗房间: %d", v.GetID())
				return v
			}
		}
		// 如果找不到合适的战斗房，则创建新的
		br := m.Create()
		br.AddRoom(room)
		log.Debug("直接创建战斗-----房间: %d", br.GetID())
		return br
	}
}

func (m *BattleRoomManager) MatchRoomzCancel(room Roomer) bool {
	for _, v := range m.Rooms {
		return v.DeleRoom(room)
	}
	log.Debug("管理器找不到匹配取消的房间 %v", room.GetID())
	return false
}

func (m *BattleRoomManager) delete(a []BattleRoomer, elem BattleRoomer) []BattleRoomer {
	for i := 0; i < len(a); i++ {
		if a[i].GetID() == elem.GetID() {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}
