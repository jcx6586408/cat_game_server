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

	CreateChan    chan *msg.Member // 创建房间管道
	RecyleChan    chan *Room       // 房间回收管道
	AddFriendChan chan *ChangeRoom // 加入好友管道
	LeaveChan     chan *ChangeRoom // 离开管道
	OfflineChan   chan *ChangeRoom // 离开管道
	UseChan       chan *Room       // 使用房间管道
	OverRoomChan  chan int         // 关闭房间
	Done          chan interface{} // 停用通道
}

type ChangeRoom struct {
	RoomID int
	Member *msg.Member
}

var Manager = NewManager()

var lock sync.RWMutex

func NewManager() *RoomManager {
	m := &RoomManager{}
	m.GenRoomID = genRoomID()
	m.Rooms = []*Room{}
	// 初始化10个房间
	for i := 0; i < 10; i++ {
		r := NewRoom(<-m.GenRoomID)
		m.Rooms = append(m.Rooms, r)
	}
	m.PrepareRooms = []*Room{}
	m.UsingRooms = []*Room{}
	m.CreateChan = make(chan *msg.Member)
	m.RecyleChan = make(chan *Room)
	m.AddFriendChan = make(chan *ChangeRoom)
	m.LeaveChan = make(chan *ChangeRoom)
	m.UseChan = make(chan *Room)
	m.OverRoomChan = make(chan int)
	m.OfflineChan = make(chan *ChangeRoom)
	m.Done = make(chan interface{})
	return m
}

func (m *RoomManager) Run() {
	go func() {
		for {
			select {
			case <-m.Done:
				return
			case member := <-m.CreateChan:
				m.CreateRoom(member)
			case r := <-m.RecyleChan:
				m.RecycleRoom(r)
			case data := <-m.AddFriendChan:
				m.AddFriendMember(data.RoomID, data.Member)
			case data := <-m.LeaveChan:
				m.LeavePrepareMemeber(data.RoomID, data.Member)
			case data := <-m.OfflineChan:
				m.OfflineMemeber(data.RoomID, data.Member)
			case r := <-m.UseChan:
				m.ToUsingRoom(r)
			case roomID := <-m.OverRoomChan:
				m.OverRoom(roomID)
			}

		}
	}()
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

// 关闭房间
func (m *RoomManager) OverRoom(roomID int) error {
	lock.Lock()
	for _, v := range m.PrepareRooms {
		if v.ID == roomID {
			v.Close()
			return nil
		}
	}
	lock.Unlock()
	return errors.New("找不到房间")
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
		room = NewRoom(<-m.GenRoomID)
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
			// 通知离开房间
			v.LeavePrepareMember(member)
			return nil
		}
	}
	return errors.New("没有找到房间")
}

func (m *RoomManager) OfflineMemeber(roomID int, member *msg.Member) error {
	for _, v := range m.PrepareRooms {
		if v.ID == roomID {
			// 通知离开房间
			v.LeavePrepareMember(member)
			return nil
		}
	}
	for _, v := range m.UsingRooms {
		if v.ID == roomID {
			// 通知离开房间
			v.LeavePlayingMember(member)
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
			room = NewRoom(<-m.GenRoomID)
		}
		m.PrepareRooms = append(m.PrepareRooms, room) // 加入准备列表
		lock.Unlock()
	}
	room.AddMember(member)
}
