package room

import (
	"catLog"
	"errors"
	"proto/msg"
	"sync"
)

type RoomManager struct {
	Rooms        []*Room  // 闲置房间列表
	UsingRooms   []*Room  //正在使用的房间列表
	PrepareRooms []*Room  // 正在准备的房间列表
	GenRoomID    chan int // 房间生成器
	Type         int      // 管理器类型
}

var Manager = NewManager(MATCHROOM)
var lock sync.RWMutex

func NewManager(t int) *RoomManager {
	m := &RoomManager{}
	m.GenRoomID = genRoomID()
	m.Rooms = []*Room{}
	// 初始化10个房间
	for i := 0; i < 10; i++ {
		r := NewRoom(<-m.GenRoomID, t)
		m.Rooms = append(m.Rooms, r)
	}
	m.PrepareRooms = []*Room{}
	m.UsingRooms = []*Room{}
	m.Type = t
	return m
}

func (m *RoomManager) Delete(a []*Room, elem *Room) []*Room {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}

// 回收房间
func (m *RoomManager) RecycleRoom(room *Room) {
	lock.Lock()                                 // 锁住
	m.UsingRooms = m.Delete(m.UsingRooms, room) // 移除正在使用房间
	m.Rooms = append(m.Rooms, room)             // 加入闲置房间
	lock.Unlock()                               // 解锁
	catLog.Log("使用房间数量_", len(m.UsingRooms))
	catLog.Log("闲置房间数量_", len(m.Rooms))
}

func (m *RoomManager) GetPrepareRoom(id int32) (*Room, error) {
	for _, v := range m.PrepareRooms {
		if v.ID == int(id) {
			return v, nil
		}
	}
	return nil, errors.New("找不到准备房间")
}

// 移动到玩的房间
func (m *RoomManager) ToUsingRoom(room *Room) {
	lock.Lock()                                     // 锁住
	m.PrepareRooms = m.Delete(m.PrepareRooms, room) // 移除正在使用房间
	m.UsingRooms = append(m.UsingRooms, room)       // 加入闲置房间
	lock.Unlock()                                   // 解锁
	catLog.Log("准备房间数量_", len(m.PrepareRooms))
	catLog.Log("使用房间数量_", len(m.UsingRooms))
}

// 创建房间
func (m *RoomManager) CreateRoom(member *msg.Member) *Room {
	lock.Lock() // 锁住
	var room *Room
	// 判断是否还有闲置房间
	if len(m.Rooms) > 0 {
		room = m.Rooms[len(m.Rooms)-1]    // 取最末尾房间
		m.Rooms = m.Delete(m.Rooms, room) // 脱离闲置列表
	} else {
		room = NewRoom(<-m.GenRoomID, m.Type)
	}
	m.PrepareRooms = append(m.PrepareRooms, room) // 加入准备列表
	lock.Unlock()                                 // 解锁
	room.AddMember(member)
	// 通知创建房间成功, 返回房间ID
	return room
}

func (m *RoomManager) LeavePrepareMemeber(roomID int, member *msg.Member) error {
	for _, v := range m.PrepareRooms {
		if v.ID == roomID {
			// 通知加入房间
			v.LeavePrepareMember(member)
			return nil
		}
	}
	return errors.New("没有找到房间")
}

// 加入好友房间
func (m *RoomManager) AddFriendMember(roomID int, member *msg.Member) error {
	for _, v := range m.PrepareRooms {
		if v.ID == roomID {
			if v.IsFull() {
				// 通知满员，无法加入
				return errors.New("满员")
			}
			// 通知加入房间
			v.AddPrepareMember(member)
			return nil
		}
	}
	return errors.New("没有找到房间")
	// 通知没有找到该房间
}

// 加入匹配成员
func (m *RoomManager) AddMatchMember(member *msg.Member) {
	var room *Room
	// 判断是否有正在准备的房间
	for _, v := range m.PrepareRooms {
		if !v.IsFull() {
			room = v
		}
	}

	if room == nil {
		// 判断是否还有闲置房间
		lock.Lock()
		if len(m.Rooms) > 0 {
			room = m.Rooms[len(m.Rooms)-1]    // 取最末尾房间
			m.Rooms = m.Delete(m.Rooms, room) // 脱离闲置列表
		} else {
			room = NewRoom(<-m.GenRoomID, m.Type)
		}
		m.PrepareRooms = append(m.PrepareRooms, room) // 加入准备列表
		lock.Unlock()
	}
	room.AddMember(member)
}
