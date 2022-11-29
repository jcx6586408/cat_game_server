package room

import (
	"catLog"
	"errors"
	"excel"
	"proto/msg"
	"sync"
)

type RoomManager struct {
	Rooms         []*Room  // 闲置房间列表
	UsingRooms    []*Room  // 正在使用的房间列表
	MatchingRooms []*Room  // 正在匹配的房间列表
	PrepareRooms  []*Room  // 正在准备的房间列表
	GenRoomID     chan int // 房间生成器

	CreateChan            chan *msg.Member       // 创建房间管道
	RecyleChan            chan *Room             // 房间回收管道
	AddFriendChan         chan *ChangeRoom       // 加入好友管道
	AnswerChan            chan *msg.Answer       // 回答请求
	LeaveChan             chan *ChangeRoom       // 离开管道
	OfflineChan           chan string            // 离线管道
	UseChan               chan *Room             // 使用房间管道
	OverRoomChan          chan int               // 关闭房间
	MathRoomChan          chan int               // 匹配房间管道
	MathMemberChan        chan *msg.Member       // 匹配个人管道
	MatchRoomCancelChan   chan int               // 房间匹配取消
	MatchMemberCancelChan chan *msg.LeaveRequest // 匹配个人取消管道
	Done                  chan interface{}       // 停用通道

	TableManager *excel.ExcelManager // 表格管理器
}

type ChangeRoom struct {
	RoomID int
	Member *msg.Member
}

type AnswerRequest struct {
}

var Manager = NewManager()

var lock sync.RWMutex

func NewManager() *RoomManager {
	m := &RoomManager{}
	m.Done = make(chan interface{})
	m.GenRoomID = genRoomID(m.Done)
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
	m.OfflineChan = make(chan string)
	m.MathMemberChan = make(chan *msg.Member)
	m.MathRoomChan = make(chan int)
	m.MatchMemberCancelChan = make(chan *msg.LeaveRequest)
	m.MatchRoomCancelChan = make(chan int)
	m.AnswerChan = make(chan *msg.Answer)
	m.TableManager = excel.Read()
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
			case uuid := <-m.OfflineChan:
				m.OfflineMemeber(uuid)
			case r := <-m.UseChan:
				m.ToUsingRoom(r)
			case roomID := <-m.OverRoomChan:
				m.OverRoom(roomID)
			case a := <-m.AnswerChan:
				m.AnswerQuestion(a)
			case req := <-m.MatchMemberCancelChan:
				m.MatchMemberCancel(req)
			case member := <-m.MathMemberChan:
				m.AddMatchMember(member)
			case roomID := <-m.MatchRoomCancelChan:
				m.MatchRoomCancel(roomID)
			case roomID := <-m.MathRoomChan: // 将房间拉入匹配
				var room *Room
				for _, v := range m.PrepareRooms {
					if v.ID == roomID {
						room = v
						m.Delete(m.PrepareRooms, v)
						break
					}
				}
				m.MatchRoom(room)
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

func (m *RoomManager) AnswerQuestion(a *msg.Answer) {
	for _, v := range m.PrepareRooms {
		if v.ID == int(a.RoomID) {
			v.Answer(a)
			break
		}
	}
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

// 匹配个人，搜索正在匹配的房间
func (m *RoomManager) MatchMember(member *msg.Member) {
	var room *Room
	// 如果不存在匹配房间，则直接创建一个房间
	if len(m.MatchingRooms) <= 0 {
		room = NewRoom(<-m.GenRoomID)
		m.MatchRoom(room)
	} else {
		for _, v := range m.MatchingRooms {
			if !v.IsFull() {
				room = v
				break
			}
		}
	}
	room.AddMember(member) // 加入匹配列表

	// 判断是否满员，满员则开始比赛
	if room.IsFull() {
		// 将房间移入比赛使用房间
		m.MatchingRooms = m.Delete(m.MatchingRooms, room)
		m.UsingRooms = append(m.UsingRooms, room)
		room.StartPrepare()
	}
}

// 取消个人匹配准备
func (m *RoomManager) MatchMemberCancel(req *msg.LeaveRequest) {
	// 如果不存在匹配房间，则直接创建一个房间
	if len(m.MatchingRooms) <= 0 {
		return
	} else {
		for _, v := range m.MatchingRooms {
			if v.ID == int(req.RoomID) {
				v.LeavePrepareMember(req.Member) // 离开房间
				return
			}
		}
	}
}

// 房间匹配取消，所有人回到准备
func (m *RoomManager) MatchRoomCancel(roomID int) error {
	for _, v := range m.MatchingRooms {
		if v.ID == roomID {
			m.MatchingRooms = m.Delete(m.MatchingRooms, v)
			m.Rooms = append(m.Rooms, v)
			return nil
		}
	}
	return errors.New("没有找到准备房间")
}

// 匹配房间, 将准备的房间移入匹配列表 并挂起房间等待
func (m *RoomManager) MatchRoom(room *Room) error {
	// 加入匹配列表
	if room != nil {
		m.MatchingRooms = append(m.MatchingRooms, room)
		return nil
	}
	return errors.New("没有找到准备房间")
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

func (m *RoomManager) OfflineMemeber(uuid string) error {

	// for _, v := range m.PrepareRooms {
	// 	if v.ID == roomID {
	// 		// 通知离开房间
	// 		v.LeavePrepareMember(member)
	// 		return nil
	// 	}
	// }
	// for _, v := range m.UsingRooms {
	// 	if v.ID == roomID {
	// 		// 通知离开房间
	// 		v.LeavePlayingMember(member)
	// 		return nil
	// 	}
	// }
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
			v.AddMember(member)
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
