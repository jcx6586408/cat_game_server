package internal

// import (
// 	"errors"
// 	"excel"
// 	"fmt"
// 	pmsg "proto/msg"

// 	"github.com/name5566/leaf/log"
// )

// type RoomManager struct {
// 	Rooms         []*Room             // 闲置房间列表
// 	UsingRooms    []*Room             // 正在使用的房间列表
// 	MatchingRooms []*Room             // 正在匹配的房间列表
// 	PrepareRooms  []*Room             // 正在准备的房间列表
// 	GenRoomID     chan int            // 房间生成器
// 	Done          chan interface{}    // 停用通道
// 	TableManager  *excel.ExcelManager // 表格管理器
// }

// type ChangeRoom struct {
// 	RoomID int
// 	Member *pmsg.Member
// }

// type AnswerRequest struct {
// }

// var (
// 	Manager    *RoomManager
// 	Skins      []*Skin
// 	NamesLib   []*Names
// 	IconLib    []*Icon
// 	AnswerLibs []Answers
// )

// func RoomManagerInit() {
// 	log.Debug("房间管理器初始化")
// 	Manager = NewManager()
// 	ExcelConfigUpdate()
// }

// func ExcelConfigUpdate() {
// 	Skins = ToSkinLib()
// 	NamesLib = ToNameLib()
// 	IconLib = ToIconLib()
// 	AnswerLibs = []Answers{}
// 	AnswerLibs = append(AnswerLibs, ToAnswerLib("question1"))

// 	log.Debug("皮肤数量: %v", len(Skins))
// 	log.Debug("名字数量: %v", len(NamesLib))
// 	log.Debug("Icon数量: %v", len(IconLib))
// }

// // 生成房间ID
// func genRoomID(done chan interface{}) chan int {
// 	c := make(chan int)
// 	n := 0
// 	skeleton.Go(func() {
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case c <- n:
// 				n = n + 1
// 			}

// 		}
// 	}, func() {

// 	})
// 	return c
// }

// func NewManager() *RoomManager {
// 	m := &RoomManager{}
// 	m.Done = make(chan interface{})
// 	m.GenRoomID = genRoomID(m.Done)
// 	m.Rooms = []*Room{}
// 	// 初始化10个房间
// 	for i := 0; i < 10; i++ {
// 		r := NewRoom(<-m.GenRoomID)
// 		m.Rooms = append(m.Rooms, r)
// 	}
// 	m.PrepareRooms = []*Room{}
// 	m.UsingRooms = []*Room{}

// 	m.TableManager = excel.Read()
// 	return m
// }

// func (m *RoomManager) Delete(a []*Room, elem *Room) []*Room {
// 	for i := 0; i < len(a); i++ {
// 		if a[i] == elem {
// 			a = append(a[:i], a[i+1:]...)
// 			i--
// 			break
// 		}
// 	}
// 	return a
// }

// func (m *RoomManager) GetPrepareRoom(id int32) (*Room, error) {
// 	for _, v := range m.PrepareRooms {
// 		if v.ID == int(id) {
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("找不到准备房间")
// }

// func (m *RoomManager) AnswerQuestion(a *pmsg.Answer) (*Room, error) {
// 	for _, v := range m.UsingRooms {
// 		if v.ID == int(a.RoomID) {
// 			v.Answer(a)
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("")
// }

// func (m *RoomManager) toString() {
// 	log.Debug("闲置房间:%v, 准备房间:%v,匹配房间:%v, 使用房间:%v", len(m.Rooms), len(m.PrepareRooms), len(m.MatchingRooms), len(m.UsingRooms))
// 	var roomsStr = "闲置房间: "
// 	for _, v := range m.Rooms {
// 		roomsStr += fmt.Sprintf("%v#", v.ID)
// 	}
// 	log.Debug("%v", roomsStr)
// }

// func (m *RoomManager) isInlist(arr []*Room, r *Room) bool {
// 	for _, v := range arr {
// 		if v.ID == r.ID {
// 			return true
// 		}
// 	}
// 	return false
// }

// // 准备房间回收
// func (m *RoomManager) PrepareToRooms(room *Room) {
// 	if m.isInlist(m.PrepareRooms, room) {
// 		m.PrepareRooms = m.Delete(m.PrepareRooms, room) // 移除正在使用房间
// 		m.Rooms = append(m.Rooms, room)                 // 加入闲置房间
// 		log.Debug("prepare当前闲置房间数量: %v", len(m.Rooms))
// 		m.toString()
// 	}
// }

// // 比赛房间回收
// func (m *RoomManager) UsingToRooms(room *Room) {
// 	if m.isInlist(m.UsingRooms, room) {
// 		m.UsingRooms = m.Delete(m.UsingRooms, room) // 移除正在使用房间
// 		m.Rooms = append(m.Rooms, room)             // 加入闲置房间
// 		log.Debug("using当前闲置房间数量: %v", len(m.Rooms))
// 		m.toString()
// 	}
// }

// // 匹配房间回收
// func (m *RoomManager) MathingToRooms(room *Room) {
// 	if m.isInlist(m.MatchingRooms, room) {
// 		m.MatchingRooms = m.Delete(m.MatchingRooms, room) // 移除正在使用房间
// 		m.Rooms = append(m.Rooms, room)                   // 加入闲置房间
// 		log.Debug("matching当前闲置房间数量: %v", len(m.Rooms))
// 		m.toString()
// 	}
// }

// // 准备移动到玩的房间
// func (m *RoomManager) PrepareToUsingRoom(room *Room) {
// 	if m.isInlist(m.PrepareRooms, room) {
// 		m.PrepareRooms = m.Delete(m.PrepareRooms, room) // 移除正在使用房间
// 		m.UsingRooms = append(m.UsingRooms, room)       // 加入闲置房间
// 		log.Debug("准备%v移动到使用%v", len(m.PrepareRooms), len(m.UsingRooms))
// 		m.toString()
// 	}
// }

// // 匹配移动到玩的房间
// func (m *RoomManager) MatchingToUsingRoom(room *Room) {
// 	if m.isInlist(m.MatchingRooms, room) {
// 		m.MatchingRooms = m.Delete(m.MatchingRooms, room) // 移除正在使用房间
// 		m.UsingRooms = append(m.UsingRooms, room)         // 加入闲置房间
// 		log.Debug("匹配%v移动到使用%v", len(m.MatchingRooms), len(m.UsingRooms))
// 		m.toString()
// 	}
// }

// // 移动到准备的房间
// func (m *RoomManager) UsingToPrepareRoom(room *Room) {
// 	if m.isInlist(m.UsingRooms, room) {
// 		m.UsingRooms = m.Delete(m.UsingRooms, room)   // 移除正在使用房间
// 		m.PrepareRooms = append(m.PrepareRooms, room) // 加入准备房间
// 		log.Debug("使用%v移动到准备%v", len(m.UsingRooms), len(m.UsingRooms))
// 		m.toString()
// 	}
// }

// // 匹配移动到准备的房间
// func (m *RoomManager) MatchingToPrepareRoom(room *Room) {
// 	if m.isInlist(m.MatchingRooms, room) {
// 		m.MatchingRooms = m.Delete(m.MatchingRooms, room) // 移除正在使用房间
// 		m.PrepareRooms = append(m.PrepareRooms, room)     // 加入准备房间
// 		log.Debug("匹配%v移动到准备%v", len(m.UsingRooms), len(m.UsingRooms))
// 		m.toString()
// 	}
// }

// // 移动到匹配的房间
// func (m *RoomManager) PrepareToMatchingRoom(room *Room) {
// 	if m.isInlist(m.PrepareRooms, room) {
// 		m.PrepareRooms = m.Delete(m.PrepareRooms, room) // 移除正在使用房间
// 		m.MatchingRooms = append(m.MatchingRooms, room) // 加入匹配房间
// 		log.Debug("准备%v移动到匹配%v", len(m.UsingRooms), len(m.UsingRooms))
// 		m.toString()
// 	}
// }

// // 关闭房间
// func (m *RoomManager) OverRoom(roomID int) (*Room, error) {
// 	for _, v := range m.PrepareRooms {
// 		if v.ID == roomID {
// 			v.Close()           // 关闭房间
// 			m.PrepareToRooms(v) // 回收房间
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("找不到房间")
// }

// func (m *RoomManager) GetRoom(roomID int) (*Room, error) {
// 	for _, v := range m.PrepareRooms {
// 		if v.ID == roomID {
// 			return v, nil
// 		}
// 	}
// 	for _, v := range m.MatchingRooms {
// 		if v.ID == roomID {
// 			return v, nil
// 		}
// 	}
// 	for _, v := range m.UsingRooms {
// 		if v.ID == roomID {
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("找不到房间")
// }

// // 创建房间
// func (m *RoomManager) CreateRoom(member *pmsg.Member) *Room {
// 	var room *Room
// 	// 判断是否还有闲置房间
// 	if len(m.Rooms) > 0 {
// 		room = m.Rooms[len(m.Rooms)-1]    // 取最末尾房间
// 		m.Rooms = m.Delete(m.Rooms, room) // 脱离闲置列表
// 	} else {
// 		room = NewRoom(<-m.GenRoomID)
// 	}
// 	m.PrepareRooms = append(m.PrepareRooms, room) // 加入准备列表

// 	room.AddMasterMember(member)
// 	room.Run() // 运行成员等待加入
// 	return room
// }

// // 匹配个人，搜索正在匹配的房间
// func (m *RoomManager) MatchMember(member *pmsg.Member) (*Room, error) {
// 	var room *Room
// 	// 如果不存在匹配房间，则直接创建一个房间
// 	if len(m.MatchingRooms) <= 0 {
// 		room = m.CreateRoom(member)
// 		m.MatchRoom(room)
// 		log.Debug("匹配时创建了新房间: %v, 当前房间人数: %v", room.ID, room.GetMembersCount())
// 	} else {
// 		for _, v := range m.MatchingRooms {
// 			if !v.IsFull() {
// 				room = v
// 				log.Debug("匹配时找到了房间: %v, 当前房间人数: %v", room.ID, room.GetMembersCount())
// 				break
// 			}
// 		}
// 		// 通知加入房间
// 		skeleton.Go(func() {
// 			log.Debug("搜索加入房间")
// 			member.IsInvite = false
// 			room.ChangeChan <- member
// 		}, func() {})
// 	}

// 	// 判断是否满员，满员则开始比赛
// 	if room.IsFull() {
// 		// 将房间移入比赛使用房间
// 		m.MatchingToUsingRoom(room)
// 		room.StartPrepare()
// 	}
// 	return room, nil
// }

// // 取消个人匹配准备
// func (m *RoomManager) MatchMemberCancel(req *pmsg.LeaveRequest) (*Room, error) {
// 	// 如果不存在匹配房间，则直接创建一个房间
// 	if len(m.MatchingRooms) <= 0 {
// 		return nil, errors.New("")
// 	} else {
// 		for _, v := range m.MatchingRooms {
// 			if v.ID == int(req.RoomID) {
// 				v.LeavePrepareMember(req.Member) // 离开房间
// 				return v, nil
// 			}
// 		}
// 	}
// 	return nil, errors.New("")
// }

// // 房间匹配取消，所有人回到准备
// func (m *RoomManager) MatchRoomCancel(roomID int) (*Room, error) {
// 	for _, v := range m.MatchingRooms {
// 		if v.ID == roomID {
// 			m.MatchingRooms = m.Delete(m.MatchingRooms, v)
// 			m.Rooms = append(m.Rooms, v)
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("没有找到准备房间")
// }

// func (m *RoomManager) StartPlay(roomID int) {
// 	for _, v := range m.PrepareRooms {
// 		log.Debug("开始游戏得房间ID: %v, 查找房间ID: %v", roomID, v.ID)
// 		if v.ID == roomID {
// 			v.StartPrepare()
// 			return
// 		}
// 	}
// 	for _, v := range m.UsingRooms {
// 		log.Debug("开始游戏得房间ID: %v, 查找房间ID: %v", roomID, v.ID)
// 		if v.ID == roomID {
// 			v.StartPrepare()
// 			return
// 		}
// 	}
// }

// // 成员复活
// func (m *RoomManager) Relive(req *pmsg.MemberReliveRequest) (*Room, error) {
// 	for _, v := range m.UsingRooms {
// 		if v.ID == int(req.RoomID) {
// 			v.Relive(req)
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("找不到房间")
// }

// func (m *RoomManager) GetRoomByID(roomID int) (*Room, error) {
// 	for _, v := range m.PrepareRooms {
// 		if v.ID == roomID {
// 			return v, nil
// 		}
// 	}
// 	for _, v := range m.UsingRooms {
// 		if v.ID == roomID {
// 			return v, nil
// 		}
// 	}
// 	for _, v := range m.MatchingRooms {
// 		if v.ID == roomID {
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("找不到房间")
// }

// // 匹配房间, 将准备的房间移入匹配列表 并挂起房间等待
// func (m *RoomManager) MatchRoom(room *Room) (*Room, error) {
// 	// 加入匹配列表
// 	if room != nil {
// 		// 检测是否满员
// 		if room.IsFull() {
// 			log.Debug("=============================满员直接开始")
// 			m.UsingRooms = append(m.UsingRooms, room)
// 			room.StartPrepare()
// 		} else {
// 			log.Debug("=============================等待机器人加入")
// 			m.PrepareToMatchingRoom(room)
// 			room.Matching()
// 		}
// 		return room, nil
// 	}
// 	return nil, errors.New("没有找到准备房间")
// }

// func (m *RoomManager) LeavePrepareMemeber(roomID int, member *pmsg.Member) (*Room, error) {
// 	for _, v := range m.PrepareRooms {
// 		if v.ID == roomID {
// 			// 通知离开房间
// 			v.LeavePrepareMember(member)
// 			return v, nil
// 		}
// 	}
// 	for _, v := range m.MatchingRooms {
// 		if v.ID == roomID {
// 			// 通知离开房间
// 			v.LeaveMatchingMember(member)
// 			return v, nil
// 		}
// 	}
// 	for _, v := range m.UsingRooms {
// 		if v.ID == roomID {
// 			// 通知离开房间
// 			v.LeavePlayingMember(member)
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("没有找到房间")
// }

// func (m *RoomManager) OfflineMemeber(uuid string) (*Room, error) {
// 	for _, v := range m.UsingRooms {
// 		bo := v.OfflinHanlder(uuid)
// 		if bo {
// 			return v, nil
// 		}
// 	}
// 	for _, v := range m.MatchingRooms {
// 		bo := v.OfflinHanlder(uuid)
// 		if bo {
// 			return v, nil
// 		}
// 	}
// 	for _, v := range m.PrepareRooms {
// 		bo := v.OfflinHanlder(uuid)
// 		if bo {
// 			return v, nil
// 		}
// 	}
// 	return nil, errors.New("没有找到房间")
// }

// // 加入好友房间
// func (m *RoomManager) AddFriendMember(roomID int, member *pmsg.Member) (*Room, int, error) {
// 	log.Debug("有人加入房间")
// 	m.toString()
// 	member.IsInvite = true
// 	for _, v := range m.PrepareRooms {
// 		log.Debug("roomID: %v", v.ID)
// 		if v.ID == roomID {
// 			if v.IsInviteFull() {
// 				// 通知满员，无法加入
// 				return nil, ROOMFULL, errors.New("满员")
// 			}
// 			log.Debug("房间加入成功", roomID)
// 			// 通知加入房间
// 			skeleton.Go(func() {
// 				log.Debug("加入房间通知", roomID)
// 				v.ChangeChan <- member
// 			}, func() {})
// 			return v, 0, nil
// 		}
// 	}
// 	for _, v := range m.UsingRooms {
// 		if v.ID == roomID {
// 			return v, ROOMSTARTED, errors.New("游戏已经开始")
// 		}
// 	}
// 	return nil, ROOMNULL, errors.New("没有找到房间")
// 	// 通知没有找到该房间
// }
