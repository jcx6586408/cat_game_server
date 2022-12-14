package internal

import (
	"errors"
	"excel"
	pmsg "proto/msg"
	"sync"

	"github.com/labstack/gommon/log"
)

// 房间管理器
type Managerer interface {
	Create() Roomer                // 创建房间
	Destroy(roomID int)            // 通过ID销毁房间
	GetRoomByID(roomID int) Roomer // 通过ID获取已经创建房间
	GetRoomsCount() int            // 获取当前所有房间数量

	// 运行处理
	Run(roomID int) // 开始运行房间（用于等待成员加入）

	// 匹配
	Matching(roomID int)       // 匹配房间
	MatchingCancel(roomID int) // 取消房间创建

	// 成员处理
	AddMember(roomID int, member *pmsg.Member) (Roomer, int, error) // 添加成员
	LeaveMember(roomID int, member Memberer) (Roomer, error)        // 添加成员
}

type Manager struct {
	Pool         sync.Pool
	Rooms        []Roomer            // 准备房间
	Done         chan interface{}    // 停用通道
	TableManager *excel.ExcelManager // 表格管理器
}

func (m *Manager) Create() Roomer {
	r := manager.Pool.Get().(*Room)
	r.OnInit()
	return r
}

func (m *Manager) GetRoomsCount() int {
	return len(m.Rooms)
}

func (m *Manager) GetRoomByID(roomID int) Roomer {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			return v
		}
	}
	return nil
}

func (m *Manager) destroy(room Roomer) bool {
	room.OnClose()
	manager.Pool.Put(room)
	log.Debug("房间回收, ID: %d", room.GetID())
	return true
}

func (m *Manager) Destroy(roomID int) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			m.Rooms = m.delete(m.Rooms, v) // 移除
			log.Debug("准备房间回收")
			m.destroy(v) // 回收
			return
		}
	}

}

func (m *Manager) Run(roomID int) {

}

func (m *Manager) Matching(roomID int) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			v.Matching()
			return
		}
	}
}

func (m *Manager) MatchingCancel(roomID int) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			v.MatchingCancel()
			return
		}
	}
}

func (m *Manager) AddMember(roomID int, member *pmsg.Member) (Roomer, int, error) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			bo := v.AddMember(member)
			if bo {
				return v, 0, nil
			}
			return nil, ROOMFULL, errors.New("")
		}
	}
	return nil, ROOMNULL, errors.New("")
}

func (m *Manager) LeaveMember(roomID int, member Memberer) (Roomer, error) {
	for _, v := range m.Rooms {
		if v.GetID() == roomID {
			v.LeaveMember(member)
			return v, nil
		}
	}
	return nil, errors.New("")
}

func (m *Manager) delete(a []Roomer, elem Roomer) []Roomer {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}

func (m *Manager) AnswerQuestion(a *pmsg.Answer) (*Room, error) {
	for _, v := range m.Rooms {
		if v.GetID() == int(a.RoomID) {
			v.Answer(a)
			return v, nil
		}
	}
	return nil, errors.New("")
}
