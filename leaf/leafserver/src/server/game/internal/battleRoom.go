package internal

import (
	"context"
	"math/rand"
	pmsg "proto/msg"
	"remotemsg"
	"time"

	"github.com/name5566/leaf/log"
)

type BattleRoomer interface {
	Roombaseer
	Play()                               // 开始游戏
	OnPlayEnd()                          // 游戏结束
	AddRoom(room Roomer) bool            // 加入房间
	DeleRoom(room Roomer) bool           // 退出房间
	Relive(uuid string)                  // 复活
	Answer(a *pmsg.Answer)               // 答题
	SendLeave(member *pmsg.Member)       // 发送离开消息
	OnLeave(Roomer, *pmsg.Member)        // 监听成员离开
	Send(msgID int, change *pmsg.Member) // 发送消息
	SendAdd(member *pmsg.Member)         // 发送加入消息
	Say(uuid, word string)
}

type BattleRoom struct {
	ID             int
	Members        []*pmsg.Member
	Rooms          []Roomer
	Max            int
	LibAnswer      *LibAnswer // 当前题库
	Level          int        // 题库段位
	QuestionCount  int
	AnswerTime     int // 单次回答问题时间
	Cur            int
	Cancel         context.CancelFunc // 取消
	MatchCancel    context.CancelFunc // 机器人匹配取消
	ReliveDone     chan interface{}
	ReliveWaitTime int
	PrepareTime    int
	isStart        bool
	isRelive       bool // 是否使用了复活机会
	isAnswerTime   bool // 答题时间

	robotIcons []*Icon  // 机器人头像
	robotNames []*Names // 机器人名字

	Done chan interface{} // 结束管道
}

func (r *BattleRoom) GetID() int {
	return r.ID
}

func (r *BattleRoom) OnInit() {
	r.Members = []*pmsg.Member{}
	r.Rooms = []Roomer{}
	r.Cur = 0
	r.AnswerTime = RoomConf.AnswerTime
	r.Max = RoomConf.MaxMember
	r.PrepareTime = RoomConf.PrepareTime
	r.ReliveWaitTime = RoomConf.ReliveWaitTime
	log.Debug("当前最大战斗房间人数, %d", r.Max)
	r.isStart = false
	r.isRelive = false
	r.isAnswerTime = false
	// 复制头像名字
	r.robotIcons = make([]*Icon, len(IconLib))
	r.Done = make(chan interface{})
	copy(r.robotIcons, IconLib)
	// 挂起机器人
	r.matching()
}

func (r *BattleRoom) OnClose() {
	defer recover()
	r.LibAnswer = nil
	r.Members = nil
	r.Rooms = nil
	r.isStart = false
	if r.Done != nil {
		close(r.Done)
		r.Done = nil
	}
}

func (m *BattleRoom) delete(a []Roomer, elem Roomer) []Roomer {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}

func (m *BattleRoom) DeleRoom(room Roomer) bool {
	for _, v := range m.Rooms {
		if v.GetID() == room.GetID() {
			log.Debug("房间取消---开始：%v|%v", len(m.Rooms), room.GetID())
			m.Rooms = m.delete(m.Rooms, room)
			room.OnEndPlay()
			log.Debug("房间取消成功%v|%v", len(m.Rooms), room.GetID())
			if len(m.Rooms) <= 0 {
				// 回收房间
				BattleManager.Destroy(m)
			}
			return true
		}
	}
	return false
}

func (r *BattleRoom) OnPlayEnd() {
	r.isStart = false
	// 回收房间
	BattleManager.Destroy(r)
}

func (r *BattleRoom) OnLeave(room Roomer, member *pmsg.Member) {
	r.SendLeave(member)
	// 如果房间没人，则清除房间
	if room.GetMemberCount() <= 0 {
		// r.Rooms = r.delete(r.Rooms, room)
		r.DeleRoom(room)
		return
	}
	if len(room.GetPlayingMembers()) <= 0 {
		r.DeleRoom(room)
		return
	}
}

func (r *BattleRoom) Relive(uuid string) {
	if !r.isStart {
		log.Debug("本次答题已经结束，无法进行复活")
		return
	}
	for _, room := range r.Rooms {
		for _, v := range room.GetMembers() {
			if v.Uuid == uuid {
				if v.IsDead {
					log.Debug("*************************复活成功**************************%v", uuid)
					v.IsDead = false
					v.State = int32(MEMBERPLAYING)
					if r.LibAnswer == nil {
						return
					}
					r.SendRelive(uuid)
				} else {
					log.Debug("该成员没有死亡, 无需复活")
				}
				break
			}
		}
	}
}

func (r *BattleRoom) GetPlayingMembers() []*pmsg.Member {
	var members = []*pmsg.Member{}
	for _, room := range r.Rooms {
		for _, v := range room.GetMembers() {
			if v.State == int32(MEMBERPLAYING) || v.State == int32(MEMBERMATCHING) {
				members = append(members, v)
			}
		}
	}
	members = append(members, r.Members...)
	// log.Debug("战斗房间ID:%d,游玩人数: %d", r.GetID(), len(members))
	return members
}

func (r *BattleRoom) GetPrepareMembers() []*pmsg.Member {
	var members = []*pmsg.Member{}
	for _, room := range r.Rooms {
		for _, v := range room.GetMembers() {
			if v.State == int32(MEMEBERPREPARE) {
				members = append(members, v)
			}
		}
	}
	return members
}

func (r *BattleRoom) Full() bool {
	return r.GetMemberCount() >= r.Max
}

func (r *BattleRoom) GetMemberCount() int {
	count := 0
	for _, v := range r.Rooms {
		count += v.GetMemberCount()
	}
	count += len(r.Members)
	// log.Debug("当前战斗房间人数: %d", count)
	return count
}

func (r *BattleRoom) AddRoom(room Roomer) bool {
	if r.isStart {
		log.Debug("******************房间已经开始, 无法加入******************")
		return false
	}
	less := r.Max - r.GetMemberCount()
	if less >= room.GetMemberCount() {
		r.Rooms = append(r.Rooms, room)
		log.Debug("%d****号战斗房成功加入房间: %d", r.GetID(), room.GetID())
		if r.Max == r.GetMemberCount() {
			log.Debug("满员开始游戏")
			r.Play() // 满员开始游戏
		} else {
			// log.Debug("当前战斗房人数:%d", r.GetMemberCount())
		}
		return true
	}
	return false
}

func (m *BattleRoom) GetPlayTime() int {
	return m.QuestionCount*m.AnswerTime + 1
}

func (r *BattleRoom) Play() {
	if r.isStart {
		return
	}
	r.isStart = true
	battleManager.Play(r)
	// 设置所有成员状态为游玩
	r.foreachMembers(func(v *pmsg.Member, room Roomer) {
		// 未准备成员跳过
		if v.State == int32(MEMEBENONERPREPARE) {
			return
		}
		v.State = int32(MEMBERPLAYING)
	})

	// 获取题库
	var max = r.GetMaxLevel() // 获得最高等级
	l := GetIDByLevel(max)    // 获得段位
	r.Level = l
	log.Debug("段位长度: %v, 索引: %v", len(LevelLib), l-1)
	var levelConf = LevelLib[l-1] // 获得段位配置
	r.LibAnswer = Questions.RandAnswerLib(l, levelConf.QuestionNumber)
	r.AnswerTime = levelConf.QuestionTime
	log.Debug("%s", r.LibAnswer.ToString())
	r.QuestionCount = len(r.LibAnswer.Answers)

	// 转移成员到开始玩
	r.SetDefaultAnswer() // 设置默认答案
	r.SendStart()
	log.Debug("开始比赛: %v|题库数量: %v|答题时间: %v", r.ID, levelConf.QuestionNumber, levelConf.QuestionTime)
	r.PlayRun()
}

func (m *BattleRoom) SendEndPlay() {
	// 回归所有成员状态
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		if v.State == int32(MEMEBENONERPREPARE) {
			return
		}
		v.State = int32(MEMEBERPREPARE)
	})
	m.SendEnd()
	for _, v := range m.Rooms {
		v.OnEndPlay()
	}
}

func (m *BattleRoom) endjudge() {
	if m.LibAnswer == nil || m.LibAnswer.Progress >= m.QuestionCount-1 {
		log.Debug("答题结束: %v", m.ID)
		// 回归所有成员状态
		m.SendEndPlay()
		m.OnPlayEnd()
		m.Cancel = nil
	} else {
		m.LibAnswer.Progress++                 // 进度增长
		m.Send(remotemsg.ROOMANSWERSTART, nil) // 答题开始
		log.Debug("%s", m.LibAnswer.SingleToString())
		<-time.After(time.Second * time.Duration(m.PrepareTime))
		m.singleRun()
	}
}

func (m *BattleRoom) singleRun() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(m.AnswerTime+1))
	cur := m.AnswerTime + 1
	m.SendTime(cur)
	m.isAnswerTime = true
	if m.LibAnswer == nil {
		return
	}
	log.Debug("================（单次答题开始%d）===================", m.LibAnswer.Progress+1)
	num := float32(len(m.Members))
	min := int(num * RoomConf.RobotActionMin)
	max := int(num * RoomConf.RobotActionMax)
	skeleton.Go(func() {
		for {
			select {
			case <-m.Done:
				return
			case <-ctx.Done():
				log.Debug("================单次答题结束==================")
				// 检查所有成员答案
				var allWrong = m.CheckAndHandleDead()
				// 如果没有全错
				m.isAnswerTime = false
				m.Send(remotemsg.ROOMANSWEREND, nil) // 答题结束
				if m.LibAnswer == nil || m.LibAnswer.Progress >= m.QuestionCount-1 {
					m.endjudge()
					return
				}

				if allWrong {
					if !m.isRelive {
						m.isRelive = true
						m.endjudge()
						m.WaitRelive(func() {

						}, func() {
							log.Debug("全员失败,答题结束: %v", m.ID)
							m.SendEndPlay()
							m.OnPlayEnd()
						})
					} else {
						log.Debug("已经复活等待过, 全员失败,答题结束: %v", m.ID)
						m.SendEndPlay()
						m.OnPlayEnd()
					}
				} else {
					m.endjudge()
				}
				return
			case <-time.After(time.Second * time.Duration(1)):
				m.Cur++
				cur--
				log.Debug("房间: %d,当前时间:%d", m.ID, cur)
				m.SendTime(cur)
				if cur > 2 {
					if cur > 5 {
						m.RandomRobotAnswer(min, max, 10) // 机器人答题
					} else {
						// m.RandomRobotAnswer(3, 5, 7) // 机器人答题
						m.robotAnswer(min, max, 7, func() int {
							question := m.GetQuestion()
							q, ok := Questions.QuestionMap[question.ID]
							if !ok {
								log.Debug("找不到题库: %v", question.ID)
								return rand.Intn(4)
							}
							log.Debug("0胜率: win: %v|fail: %v", q.win, q.fail)
							if q.fail+q.win <= 0 {
								if m.Level <= 1 {
									return GetRightNumberAnswer(question.RightAnswer, question)
								} else {
									return rand.Intn(4)
								}
							}
							rate := float32(q.win) * 100 / (float32(q.win) + float32(q.fail))
							log.Debug("当前题目胜率: %v|win: %v|fail: %v", rate, q.win, q.fail)
							if rand.Intn(100) < int(rate) {
								return GetRightNumberAnswer(question.RightAnswer, question)
							}
							return rand.Intn(4)
						})
					}
				}
			}
		}
	}, func() {})
}

func (m *BattleRoom) PlayRun() {
	m.Send(remotemsg.ROOMANSWERSTART, nil) // 答题开始
	<-time.After(time.Second * time.Duration(m.PrepareTime))
	m.singleRun()
}

func (m *BattleRoom) foreachMembers(call func(v *pmsg.Member, room Roomer)) {
	for _, room := range m.Rooms {
		for _, v := range room.GetMembers() {
			call(v, room)
		}
	}
	for _, v := range m.Members {
		call(v, nil)
	}
}

func (m *BattleRoom) CheckAllDead() bool {
	all := true
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		if v.State == int32(MEMEBERPREPARE) || v.State == int32(MEMEBENONERPREPARE) {
			return
		}
		if !v.IsDead {
			all = false
		}
	})
	return all
}

func (m *BattleRoom) CheckRoomAllDead(room Roomer) bool {
	all := true
	for _, v := range room.GetMembers() {
		if v.State == int32(MEMEBERPREPARE) || v.State == int32(MEMEBENONERPREPARE) {
			continue
		}
		if !v.IsDead {
			all = false
		}
	}
	return all
}

func (m *BattleRoom) GetMaxLevel() int {
	var i = 0
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		if v.State == int32(MEMEBERPREPARE) || v.State == int32(MEMEBENONERPREPARE) {
			return
		}
		log.Debug("玩家等级: %v|%v", v.Level, v.Uuid)
		if v.Level > int32(i) {
			i = int(v.Level)
		}
	})
	// if i <= 0 {
	// 	log.Debug("非法段位0**********************")
	// 	i = 1
	// }
	log.Debug("最后采用等级: %v", i)
	return i
}

func (m *BattleRoom) CheckAndHandleDead() bool {
	var allWrong = true
	rightCount := 0
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		if v == nil {
			log.Debug("玩家不存在")
			return
		}
		if v.IsDead {
			return
		}
		if v.State == int32(MEMEBERPREPARE) || v.State == int32(MEMEBENONERPREPARE) {
			return
		}
		question := m.GetQuestion()
		q := question.RightAnswer
		if v.Answer == nil || len(v.Answer) <= m.LibAnswer.Progress {
			log.Debug("答题不存在, 或超出数组")
			return
		}
		playerAnswer := v.Answer[m.LibAnswer.Progress]
		tip := ""
		if q == playerAnswer.Result {
			tip = "true---------------"
		}
		var resultAnswer = GetRightAnswer(question.RightAnswer, question)
		var userAnswer = GetRightAnswer(playerAnswer.Result, question)
		if v.IsMaster {
			log.Debug("房主***玩家uuid %v: 第%v题,  正确答案: %v, 玩家答案: %v %v", v.Uuid, m.LibAnswer.Progress, resultAnswer, userAnswer, tip)
		} else {
			log.Debug("玩家uuid %v: 第%v题, 正确答案: %v, 玩家答案: %v %v", v.Uuid, m.LibAnswer.Progress, resultAnswer, userAnswer, tip)
		}
		right := (q == playerAnswer.Result)
		if right {
			allWrong = false
			rightCount++
			skeleton.Go(func() {
				if !v.IsRobot {
					// redis.AddWinTable(fmt.Sprintf("%v_%v", m.GetQuestion().Table, m.GetQuestion().ID), 1)
					Questions.WinChan <- m.GetQuestion().ID
				}
			}, func() {})

		} else {
			skeleton.Go(func() {
				if !v.IsRobot {
					// redis.AddFailTable(fmt.Sprintf("%v_%v", m.GetQuestion().Table, m.GetQuestion().ID), 1)
					Questions.FailChan <- m.GetQuestion().ID
				}
			}, func() {})
			// 标记死亡
			v.IsDead = true
			var allDead = true

			if room != nil {
				// 检查房间是否所有人死亡
				allDead = m.CheckRoomAllDead(room)
				if allDead {
					m.WaitRoomRelive(room, func() {

					}, func() {
						log.Debug("全员死亡，退出战斗房间, 退出房间ID: %v", room.GetID())
						room.ResetToPrepare()
						room.Send(remotemsg.ROOMALLFAIL, nil)
						m.DeleRoom(room) // 移除房间
					})

				}
			}
		}
	})
	log.Debug("房间: %v, 当前正确人数: %v", m.ID, rightCount)
	return allWrong
}

func (m *BattleRoom) GetQuestion() *pmsg.Question {
	if m.LibAnswer == nil {
		return nil
	}
	if m.LibAnswer.Progress < len(m.LibAnswer.Answers) {
		return m.LibAnswer.Answers[m.LibAnswer.Progress]
	}
	return m.LibAnswer.Answers[len(m.LibAnswer.Answers)-1]
}

func (m *BattleRoom) GetProgress() int {
	if m.LibAnswer == nil {
		return 0
	}
	return m.LibAnswer.Progress
}

func (m *BattleRoom) SetDefaultAnswer() {
	// m.foreachMembers(func(v *pmsg.Member, room Roomer) {
	// 	v.Answer = make([]*pmsg.Answer, m.QuestionCount)
	// 	ranAnswer := results[rand.Intn(4)]
	// 	for i, aa := range v.Answer {
	// 		aa = &pmsg.Answer{
	// 			Uuid:       v.Uuid,
	// 			RoomID:     int32(m.ID),
	// 			QuestionID: "",
	// 			Result:     ranAnswer,
	// 		}
	// 		v.Answer[i] = aa
	// 	}
	// })
}

// 答题
func (m *BattleRoom) Answer(a *pmsg.Answer) {
	if !m.isAnswerTime {
		log.Debug("*********************非答题时间，无法进行答题*********************")
		return
	}
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		if v.Uuid == a.Uuid {
			// 是否已经答题
			if v.IsAnswered {
				return
			}
			for i, q := range v.Answer {
				if i == m.LibAnswer.Progress {
					q.Result = a.Result
				}
			}
			v.IsAnswered = true
			m.SendAnswer(a.Uuid, a.QuestionID, a.Result)
			if !v.IsRobot {
				m.RandomRobotTargetAnswer(a.Result)
			}
		}
	})
}

func (m *BattleRoom) WaitRoomRelive(room Roomer, success, fail func()) {
	subCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.ReliveWaitTime))
	log.Debug("等待复活-----------******************")
	skeleton.Go(func() {
		for {
			select {
			case <-subCtx.Done():
				log.Debug("复活等待结束-------------")
				if m.CheckRoomAllDead(room) {
					fail()
				} else {
					success()
				}
				return
			case <-time.After(time.Millisecond * time.Duration(500)):
				if !m.CheckRoomAllDead(room) {
					cancel()
				}
			}
		}
	}, func() {})
}

func (m *BattleRoom) WaitRelive(success, fail func()) {
	subCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.ReliveWaitTime))
	log.Debug("等待复活-----------******************")
	skeleton.Go(func() {
		for {
			select {
			case <-subCtx.Done():
				log.Debug("复活等待结束-------------")
				if m.CheckAllDead() {
					fail()
				} else {
					success()
				}
				return
			case <-time.After(time.Millisecond * time.Duration(500)):
				if !m.CheckAllDead() {
					cancel()
				}
			}
		}
	}, func() {})
}
