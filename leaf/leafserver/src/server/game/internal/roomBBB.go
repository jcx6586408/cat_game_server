package internal

// import (
// 	"config"
// 	"context"
// 	"fmt"
// 	"math/rand"
// 	pmsg "proto/msg"
// 	"remotemsg"
// 	"sync"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/name5566/leaf/log"
// )

// type Room struct {
// 	ID             int            // 房间ID
// 	PrepareMembers []*pmsg.Member // 处于准备的成员列表
// 	PlayingMembers []*pmsg.Member // 处于比赛的成员列表
// 	PrepareTime    int            // 比赛准备时间
// 	QuestionCount  int            // 问题总数量
// 	AnswerTime     int            // 单次回答问题时间
// 	MaxMember      int            // 最大成员数量
// 	Cur            int            // 当前房间时间

// 	OverChan     chan interface{}   // 房间结束解散通知
// 	ChangeChan   chan *pmsg.Member  // 成员变更通道
// 	Cancel       context.CancelFunc // 取消
// 	MatchCancel  context.CancelFunc // 匹配取消
// 	RunCancel    context.CancelFunc // 运行取消
// 	RunContext   context.Context    // 运行的上下文
// 	CloseContext context.Context    // 关闭房间上下文
// 	ReliveChan   chan int           // 复活通知

// 	LibAnswer      *LibAnswer // 当前题库
// 	Lock           sync.RWMutex
// 	InviteMaxCount int // 邀请最大成员
// 	ReliveWaitTime int // 复活等待时间

// 	RobotMin int // 每次匹配最小机器人数量
// 	RobotMax int // 每次匹配最大机器人数量
// }

// func NewRoom(id int) *Room {
// 	r := &Room{}
// 	r.ID = id
// 	r.PrepareMembers = []*pmsg.Member{}
// 	r.PlayingMembers = []*pmsg.Member{}
// 	conf := config.ReadRoom()
// 	r.PrepareTime = conf.PrepareTime
// 	r.AnswerTime = conf.AnswerTime
// 	r.MaxMember = conf.MaxMember
// 	r.InviteMaxCount = conf.MaxInvite
// 	r.ReliveWaitTime = conf.ReliveWaitTime
// 	r.OverChan = make(chan interface{})
// 	r.ChangeChan = make(chan *pmsg.Member)
// 	r.RobotMin = conf.RobotMin
// 	r.RobotMax = conf.RobotMax

// 	return r
// }

// func (m *Room) Send(msgID int) {
// 	for _, v := range m.PrepareMembers {
// 		a := Users[v.Uuid]
// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			log.Debug("用户不存在%v", v.Uuid)
// 			continue
// 		}

// 		send(m, a, v, msgID)
// 	}
// 	for _, v := range m.PlayingMembers {
// 		a := Users[v.Uuid]
// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			log.Debug("用户不存在%v", v.Uuid)
// 			continue
// 		}
// 		send(m, a, v, msgID)
// 	}
// }

// func (m *Room) SendLeave(msgID int, member *pmsg.Member) {
// 	for _, v := range m.PrepareMembers {
// 		a := Users[v.Uuid]
// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			log.Debug("用户不存在%v", v.Uuid)
// 			continue
// 		}

// 		send(m, a, member, msgID)
// 	}
// 	for _, v := range m.PlayingMembers {
// 		a := Users[v.Uuid]
// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			log.Debug("用户不存在%v", v.Uuid)
// 			continue
// 		}
// 		send(m, a, member, msgID)
// 	}
// }

// // 发送复活
// func (m *Room) SendRelive(req *pmsg.MemberReliveRequest) {
// 	for _, v := range m.PlayingMembers {
// 		a := Users[v.Uuid]
// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			log.Debug("用户不存在%v", v.Uuid)
// 			continue
// 		}
// 		a.WriteMsg(&pmsg.MemberReliveReply{
// 			Uuid:   req.Uuid,
// 			Answer: v.Answer[m.LibAnswer.Progress],
// 		})
// 	}
// }

// // 发送答题
// func (m *Room) SendAnswer(uuid string, qid int32, result string) {
// 	for _, v := range m.PlayingMembers {
// 		a := Users[v.Uuid]
// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			log.Debug("用户不存在%v", v.Uuid)
// 			continue
// 		}
// 		a.WriteMsg(&pmsg.Answer{
// 			Uuid:       uuid,
// 			RoomID:     int32(m.ID),
// 			QuestionID: qid,
// 			Result:     result,
// 		})
// 	}
// }

// // 发送计时
// func (m *Room) SendTime(cur int) {
// 	for _, v := range m.PlayingMembers {
// 		a := Users[v.Uuid]

// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			log.Debug("用户不存在%v", v.Uuid)
// 			continue
// 		}
// 		a.WriteMsg(&pmsg.RoomTime{
// 			Time: int32(cur),
// 		})
// 	}
// }

// // 成员复活
// func (m *Room) Relive(req *pmsg.MemberReliveRequest) {
// 	for _, v := range m.PlayingMembers {
// 		if v.Uuid == req.Uuid {
// 			v.IsDead = false
// 			skeleton.Go(func() {
// 				if m.ReliveChan != nil {
// 					log.Debug("发送复活通知***************************************")
// 					m.ReliveChan <- 1
// 				} else {
// 					log.Debug("管道已经关闭, 无需通知")
// 				}
// 			}, func() {})
// 			break
// 		}
// 	}
// }

// // 开始准备
// func (m *Room) StartPrepare() {
// 	// 获取题库
// 	m.LibAnswer = RandAnswerLib(5, GetAnswerLib())
// 	m.QuestionCount = len(m.LibAnswer.Answers)

// 	// 转移成员到开始玩
// 	Manager.PrepareToUsingRoom(m) // 移动房间
// 	m.PrepareToPlaying()          // 移动成员
// 	m.SetDefaultAnswer()          // 设置默认答案
// 	m.Send(remotemsg.ROOMSTARTPLAY)
// 	log.Debug("开始比赛: %v", m.ID)
// 	m.StartPlay()
// }

// func (m *Room) GetPlayTime() int {
// 	return m.QuestionCount*m.AnswerTime + 1
// }

// func (m *Room) Run() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	m.RunContext = ctx
// 	m.RunCancel = cancel
// 	skeleton.Go(func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				log.Debug("关闭房间, 退出等待加入")
// 				return
// 			case member := <-m.ChangeChan:
// 				bo := m.AddMember(member)
// 				if bo {
// 					// 如果非机器人，则广播
// 					if !member.IsRobot {
// 						m.Send(remotemsg.ROOMADD)
// 					}
// 				} else {
// 					a := Users[member.Uuid]
// 					a.WriteMsg(&pmsg.RoomAddFail{
// 						Code: int32(ROOMFULL),
// 					})
// 				}
// 			case <-time.After(time.Second * time.Duration(90)): // 每10秒检测是否还有玩家
// 				bo := m.IsAnybody()
// 				if bo {
// 					log.Debug("********检测到没有任何玩家, 进程房间回收********")
// 					m.close()
// 					return
// 				}
// 			}
// 		}
// 	}, func() {})
// }

// // 判断是否还有玩家
// func (m *Room) IsAnybody() bool {
// 	var count = 0
// 	for _, v := range m.PrepareMembers {
// 		a := Users[v.Uuid]
// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			continue
// 		}
// 		count++
// 	}
// 	for _, v := range m.PlayingMembers {
// 		a := Users[v.Uuid]
// 		if v.IsRobot {
// 			continue
// 		}
// 		if a == nil {
// 			continue
// 		}
// 		count++
// 	}
// 	return count == 0
// }

// func (m *Room) close() {
// 	// 清理房间
// 	m.PrepareMembers = m.PrepareMembers[0:0]
// 	m.PlayingMembers = m.PlayingMembers[0:0]
// 	m.Cur = 0
// 	m.LibAnswer = nil
// 	if m.Cancel != nil {
// 		m.Cancel()
// 	}
// 	if m.MatchCancel != nil {
// 		m.MatchCancel()
// 	}
// 	m.RunCancel()
// 	// 回收
// 	Manager.PrepareToRooms(m)
// 	Manager.MathingToRooms(m)
// 	Manager.UsingToRooms(m)
// }

// // 关闭房间
// func (m *Room) Close() {
// 	ctx, _ := context.WithTimeout(context.Background(), time.Minute*time.Duration(1))
// 	m.CloseContext = ctx
// 	skeleton.Go(func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				log.Debug("--------------------房间清理回收%v", m.ID)
// 				m.close()
// 				m.CloseContext = nil
// 				return
// 			case <-time.After(time.Second * time.Duration(1)):
// 				bo := m.IsAnybody() // 如果有玩家，则退出
// 				if !bo {
// 					log.Debug("有玩家加入，退出等待")
// 					m.CloseContext = nil
// 					return
// 				}
// 			}
// 		}
// 	}, func() {})

// }

// // 匹配加入机器人
// func (m *Room) Matching() {
// 	ctx, cancel := context.WithTimeout(m.RunContext, time.Second*time.Duration(8))
// 	m.MatchCancel = cancel
// 	cur := 0
// 	skeleton.Go(func() {
// 		for {
// 			select {
// 			case <-m.RunContext.Done():
// 				log.Debug("退出匹配状态")
// 				return
// 			case <-ctx.Done():
// 				var less = m.MaxMember - m.GetMembersCount()
// 				if less > 0 {
// 					m.AddRobot(less, NamesLib, IconLib)
// 					m.Send(remotemsg.ROOMADD)
// 					// 将房间移入比赛使用房间
// 				}
// 				Manager.MatchingToUsingRoom(m)
// 				m.StartPrepare()
// 				log.Debug("开始游戏, 退出等待加入============%v", m.ID)
// 				return
// 			case <-time.After(time.Duration(1) * time.Second):
// 				cur++
// 				if cur <= 3 {
// 					m.AddRandomCountRobots(4, 7, func() { cancel() })
// 				}
// 				if cur >= 5 {
// 					m.AddRandomCountRobots(1, 3, func() { cancel() })
// 				}
// 			}
// 		}
// 	}, func() {})
// }

// // 开始比赛
// func (m *Room) StartPlay() {
// 	total := m.GetPlayTime()
// 	log.Debug("房间ID: %v,等待总时间:%v, 单词答题时间:%v", m.ID, total, m.AnswerTime)
// 	// 创建完成通知
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(total))
// 	m.Cancel = cancel
// 	m.Cur = 0 // 答题总时间计时
// 	cur := 0  // 单局答题时间计时
// 	skeleton.Go(func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				log.Debug("答题结束: %v", m.ID)
// 				m.Send(remotemsg.ROOMENDPLAY)
// 				m.PlayingToPrepare()          // 转移成员到准备
// 				Manager.UsingToPrepareRoom(m) // 转移房间位置
// 				m.Cancel = nil
// 				return
// 			case <-time.After(time.Duration(1) * time.Second):
// 				// 广播时间
// 				m.Cur++
// 				cur++
// 				log.Debug("房间ID:%v,当前时间:%v, 单次答题时间: %v", m.ID, m.Cur, cur)
// 				if m.Cur >= 6 && cur >= 3 && cur <= m.AnswerTime-2 {
// 					if cur <= 10 {
// 						m.RandomRobotAnswer(2, 8, 10) // 机器人答题
// 					} else {
// 						m.RandomRobotAnswer(1, 3, 5) // 机器人答题
// 					}
// 				}
// 				m.SendTime(m.Cur)
// 				if m.Cur%m.AnswerTime == 0 {
// 					cur = 0
// 					log.Debug("房间_: %v, 房间准备人数: %v, 房间游戏人数: %v, 当前进度: %v", m.ID, m.GetMembersCount(), len(m.PlayingMembers), m.LibAnswer.Progress+1)
// 					if m.LibAnswer != nil {
// 						// 检查所有成员答案
// 						var allWrong = m.CheckAndHandleDead()
// 						// 如果没有全错
// 						m.LibAnswer.Progress++          // 进度增长
// 						m.Send(remotemsg.ROOMANSWEREND) // 答题结束
// 						if allWrong {
// 							m.WaitRelive(func() {
// 								cancel()
// 							}) // 等待复活
// 						}
// 					}
// 				}
// 			}

// 		}
// 	}, func() {})
// }

// func (m *Room) AddRandomCountRobots(min, max int, callback func()) {
// 	var less = m.MaxMember - m.GetMembersCount()
// 	ranNumber := rand.Intn(max+1) + min
// 	if ranNumber > less {
// 		m.AddRobot(less, NamesLib, IconLib)
// 		callback()
// 	} else {
// 		m.AddRobot(ranNumber, NamesLib, IconLib)
// 	}
// 	m.Send(remotemsg.ROOMADD)
// }

// func (m *Room) AddRobot(count int, nameLib []*Names, iconLib []*Icon) {
// 	subName := RandName(count, nameLib)
// 	subIcon := RandIcon(count, iconLib)
// 	for i := 0; i < count; i++ {
// 		guid := uuid.New().String()
// 		skinID := 1
// 		if rand.Intn(10) > 7 {
// 			skinID = 2 + rand.Intn(len(Skins))
// 		}
// 		m.ChangeChan <- &pmsg.Member{
// 			Nickname: fmt.Sprintf("%v", subName[i].ID),
// 			Uuid:     guid,
// 			Icon:     fmt.Sprintf("%v", subIcon[i].ID),
// 			IsMaster: false,
// 			IsRobot:  true,
// 			SkinID:   int32(skinID),
// 		}
// 	}
// }

// var results = []string{"A", "B", "C", "D"}

// // 随机机器人答案
// func (m *Room) RandomRobotAnswer(min, max, count int) {
// 	lenRobot := rand.Intn(max) + min
// 	startIndex := 0
// 	if len(m.PlayingMembers)-lenRobot > 0 {
// 		startIndex = rand.Intn(len(m.PlayingMembers) - lenRobot)
// 	}

// 	subArr := m.PlayingMembers[startIndex : startIndex+lenRobot]
// 	for _, v := range subArr {
// 		result := rand.Intn(4)
// 		var action = rand.Intn(10)
// 		if action >= count {
// 			return
// 		}
// 		if v.IsRobot {
// 			m.Answer(&pmsg.Answer{
// 				Uuid:       v.Uuid,
// 				RoomID:     int32(m.ID),
// 				QuestionID: m.GetQuestion().ID,
// 				Result:     results[result],
// 			})
// 		}
// 	}
// }

// func (m *Room) RandomRobotTargetAnswer(right string) {
// 	lenRobot := rand.Intn(3) + 1
// 	startIndex := 0
// 	if len(m.PlayingMembers)-lenRobot > 0 {
// 		startIndex = rand.Intn(len(m.PlayingMembers) - lenRobot)
// 	}
// 	subArr := m.PlayingMembers[startIndex : startIndex+lenRobot]
// 	for _, v := range subArr {
// 		if v.IsRobot {
// 			m.Answer(&pmsg.Answer{
// 				Uuid:       v.Uuid,
// 				RoomID:     int32(m.ID),
// 				QuestionID: m.GetQuestion().ID,
// 				Result:     right,
// 			})
// 		}
// 	}
// }

// func (m *Room) WaitRelive(fail func()) {
// 	subCtx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(m.ReliveWaitTime))
// 	log.Debug("创建复活管道通知++++++++++++++++++++++++%d", m.ID)
// 	m.ReliveChan = make(chan int)
// 	skeleton.Go(
// 		func() {
// 			for {
// 				select {
// 				case <-m.ReliveChan:
// 					close(m.ReliveChan)
// 					m.ReliveChan = nil
// 					log.Debug("复活成功--------继续游戏%v", m.ReliveChan)
// 					return
// 				case <-subCtx.Done():
// 					close(m.ReliveChan)
// 					log.Debug("比赛结束-------------------房间ID: %v", m.ID)
// 					fail()
// 					return
// 				}
// 			}
// 		},
// 		func() {},
// 	)
// }

// func (m *Room) CheckAndHandleDead() bool {
// 	var allWrong = true
// 	rightCount := 0
// 	for _, v := range m.PlayingMembers {
// 		if v.IsDead {
// 			continue
// 		}
// 		q := m.GetQuestion().RightAnswer
// 		playerAnswer := v.Answer[m.LibAnswer.Progress]
// 		tip := ""
// 		if q == playerAnswer.Result {
// 			tip = "true---------------"
// 		}
// 		if v.IsMaster {
// 			log.Debug("房主***玩家uuid %v: 正确答案: %v, 玩家答案%v %v", v.Uuid, q, playerAnswer.Result, tip)
// 		} else {
// 			log.Debug("玩家uuid %v: 正确答案: %v, 玩家答案%v %v", v.Uuid, q, playerAnswer.Result, tip)
// 		}
// 		right := (q == playerAnswer.Result)
// 		if right {
// 			allWrong = false
// 			rightCount++
// 		} else {
// 			// 标记死亡
// 			v.IsDead = true
// 		}
// 	}
// 	log.Debug("房间: %v, 当前正确人数: %v", m.ID, rightCount)
// 	return allWrong
// }

// func (m *Room) GetMembersCount() int {
// 	return len(m.PrepareMembers)
// }

// func (m *Room) SetDefaultAnswer() {
// 	for _, v := range m.PlayingMembers {
// 		v.Answer = make([]*pmsg.Answer, m.QuestionCount)
// 		ranAnswer := results[rand.Intn(4)]
// 		for i, aa := range v.Answer {
// 			aa = &pmsg.Answer{
// 				Uuid:       v.Uuid,
// 				RoomID:     int32(m.ID),
// 				QuestionID: int32(m.LibAnswer.Progress),
// 				Result:     ranAnswer,
// 			}
// 			v.Answer[i] = aa
// 		}
// 	}
// }

// // 答题
// func (m *Room) Answer(a *pmsg.Answer) {
// 	for _, v := range m.PlayingMembers {
// 		if v.Uuid == a.Uuid {
// 			for i, q := range v.Answer {
// 				if i >= m.LibAnswer.Progress {
// 					q.Result = a.Result
// 				}
// 			}
// 			m.SendAnswer(a.Uuid, a.QuestionID, a.Result)
// 			if !v.IsRobot {
// 				m.RandomRobotTargetAnswer(a.Result)
// 			}
// 			break
// 		}
// 	}
// }

// func (m *Room) PrepareToPlaying() {
// 	m.PlayingMembers = append(m.PlayingMembers, m.PrepareMembers...)
// 	m.PrepareMembers = m.PrepareMembers[0:0]
// }

// func (m *Room) PlayingToPrepare() {
// 	// 移除机器人
// 	for i := 0; i < len(m.PlayingMembers); i++ {
// 		v := m.PlayingMembers[i]
// 		if v.IsRobot {
// 			m.PlayingMembers = m.Delete(m.PlayingMembers, v)
// 			i--
// 		}
// 	}
// 	// 移除非邀请成员
// 	for i := 0; i < len(m.PlayingMembers); i++ {
// 		v := m.PlayingMembers[i]
// 		if !v.IsInvite {
// 			m.PlayingMembers = m.Delete(m.PlayingMembers, v)
// 			i--
// 		}
// 	}
// 	// 重置死亡属性
// 	for _, v := range m.PlayingMembers {
// 		v.IsDead = false
// 	}

// 	m.PrepareMembers = append(m.PrepareMembers, m.PlayingMembers...)
// 	m.PlayingMembers = m.PlayingMembers[0:0]
// }

// // 房主加入
// func (m *Room) AddMasterMember(member *pmsg.Member) {
// 	if m.IsFull() {
// 		return
// 	}
// 	member.IsInvite = true
// 	m.PrepareMembers = append(m.PrepareMembers, member)
// }

// // 加入成员
// func (m *Room) AddMember(member *pmsg.Member) bool {
// 	if m.IsFull() {
// 		return false
// 	}
// 	count := m.GetMembersCount()
// 	member.IsMaster = count == 0
// 	m.PrepareMembers = append(m.PrepareMembers, member)
// 	return true
// }

// // 离开准备成员
// func (m *Room) LeavePrepareMember(member *pmsg.Member) {
// 	log.Debug("离开准备房间, %s", member.Uuid)
// 	m.PrepareMembers = m.Delete(m.PrepareMembers, member)
// 	// 如果是房主，则移交房主
// 	if member.IsMaster {
// 		if len(m.PrepareMembers) > 0 {
// 			otherMember := m.PrepareMembers[0]
// 			otherMember.IsMaster = true
// 		}
// 	}

// 	// 如果房间人数为0,则回收房间
// 	if len(m.PlayingMembers)+len(m.PrepareMembers) <= 0 {
// 		m.Close()

// 	} else {
// 		// 检测是否还有真人
// 		realMember := false
// 		for _, v := range m.PrepareMembers {
// 			if !v.IsRobot {
// 				realMember = true
// 				break
// 			}
// 		}
// 		if !realMember {
// 			m.Close()
// 		}
// 	}

// }

// func (m *Room) LeavePlayingMember(member *pmsg.Member) {
// 	m.PlayingMembers = m.Delete(m.PlayingMembers, member)
// 	// 如果是房主，则移交房主
// 	if member.IsMaster {
// 		if len(m.PlayingMembers) > 0 {
// 			otherMember := m.PlayingMembers[0]
// 			otherMember.IsMaster = true
// 		}
// 	}

// 	// 如果房间人数为0,则回收房间
// 	if len(m.PlayingMembers)+len(m.PrepareMembers) <= 0 {
// 		m.Close()
// 	} else {
// 		// 检测是否还有真人
// 		realMember := false
// 		for _, v := range m.PlayingMembers {
// 			if !v.IsRobot {
// 				realMember = true
// 				break
// 			}
// 		}
// 		if !realMember {
// 			m.Close()
// 		}
// 	}
// }

// func (m *Room) LeaveMatchingMember(member *pmsg.Member) {
// 	m.PrepareMembers = m.Delete(m.PrepareMembers, member)
// 	// 如果是房主，则移交房主
// 	if member.IsMaster {
// 		if len(m.PrepareMembers) > 0 {
// 			otherMember := m.PrepareMembers[0]
// 			otherMember.IsMaster = true
// 		}
// 	}

// 	// 如果房间人数为0,则回收房间
// 	if len(m.PlayingMembers)+len(m.PrepareMembers) <= 0 {
// 		m.Close()
// 	} else {
// 		// 检测是否还有真人
// 		realMember := false
// 		for _, v := range m.PrepareMembers {
// 			if !v.IsRobot {
// 				realMember = true
// 				break
// 			}
// 		}
// 		if !realMember {
// 			m.Close()
// 		}
// 	}
// }

// func (m *Room) GetProgress() int {
// 	if m.LibAnswer == nil {
// 		return 0
// 	}
// 	return m.LibAnswer.Progress
// }

// func (m *Room) GetQuestion() *pmsg.Question {
// 	if m.LibAnswer == nil {
// 		return nil
// 	}
// 	// return nil
// 	if m.LibAnswer.Progress < len(m.LibAnswer.Answers) {
// 		return m.LibAnswer.Answers[m.LibAnswer.Progress]
// 	}
// 	return m.LibAnswer.Answers[len(m.LibAnswer.Answers)-1]
// }

// func (m *Room) OfflinHanlder(uuid string) bool {
// 	log.Debug("-------------------离线玩家处理-------------------")
// 	for _, v := range m.PrepareMembers {
// 		if v.Uuid == uuid {
// 			log.Debug("-------------------离开准备或匹配-------------------%v", uuid)
// 			m.LeavePrepareMember(v)  // 离开准备
// 			m.LeaveMatchingMember(v) // 离开匹配
// 			m.Send(remotemsg.ROOMLEAVE)
// 			return true
// 		}
// 	}

// 	for _, v := range m.PlayingMembers {
// 		if v.Uuid == uuid {
// 			log.Debug("-------------------离开正在游戏-------------------%v", uuid)
// 			m.LeavePlayingMember(v)
// 			m.Send(remotemsg.ROOMLEAVE)
// 			return true
// 		}
// 	}
// 	return false
// }

// // 判断房间是否满员
// func (m *Room) IsFull() bool {
// 	return m.MaxMember <= (len(m.PlayingMembers) + len(m.PrepareMembers))
// }

// // 判断是否邀请满员
// func (m *Room) IsInviteFull() bool {
// 	count := 0
// 	for _, v := range m.PrepareMembers {
// 		if v.IsInvite {
// 			count++
// 		}
// 	}
// 	return count >= m.InviteMaxCount
// }

// func (m *Room) Delete(a []*pmsg.Member, elem *pmsg.Member) []*pmsg.Member {
// 	for i := 0; i < len(a); i++ {
// 		if a[i].Uuid == elem.Uuid {
// 			a = append(a[:i], a[i+1:]...)
// 			i--
// 			break
// 		}
// 	}
// 	return a
// }
