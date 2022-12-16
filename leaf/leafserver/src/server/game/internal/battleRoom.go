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
	Relive(uuid string)                  // 复活
	Answer(a *pmsg.Answer)               // 答题
	SendLeave(member *pmsg.Member)       // 发送离开消息
	OnLeave(Roomer, Memberer)            // 监听成员离开
	Send(msgID int, change *pmsg.Member) // 发送消息
}

type BattleRoom struct {
	ID             int
	Members        []*pmsg.Member
	Rooms          []Roomer
	Max            int
	LibAnswer      *LibAnswer // 当前题库
	QuestionCount  int
	AnswerTime     int // 单次回答问题时间
	Cur            int
	Cancel         context.CancelFunc // 取消
	MatchCancel    context.CancelFunc // 机器人匹配取消
	ReliveDone     chan interface{}
	ReliveWaitTime int
	PrepareTime    int
	isStart        bool
	isRelive       bool
	isAnswerTime   bool // 答题时间
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
	// 挂起机器人
	r.matching()
}

func (r *BattleRoom) OnClose() {
	r.LibAnswer = nil
	r.Members = nil
	r.Rooms = nil
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

func (r *BattleRoom) OnPlayEnd() {
	r.isStart = false
	// 回归所有成员状态
	for _, v := range r.Rooms {
		v.OnEndPlay()
	}
	// 回收房间
	BattleManager.Destroy(r)
}

func (r *BattleRoom) OnLeave(room Roomer, member Memberer) {
	r.SendLeave(member.(*pmsg.Member))
	// 如果房间没人，则清除房间
	if room.GetMemberCount() <= 0 {
		r.Rooms = r.delete(r.Rooms, room)
	}
	// 如果没有真人
	if !r.CheckRealMember() && !r.isStart {
		// 回收房间
		BattleManager.Destroy(r)
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
					v.IsDead = false
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
			if v.State == int32(MEMBERPLAYING) {
				members = append(members, v)
			}
		}
	}
	members = append(members, r.Members...)
	log.Debug("游玩人数: %d", len(members))
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
	log.Debug("当前战斗房间人数: %d", count)
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
		if r.Max == r.GetMemberCount() {
			log.Debug("满员开始游戏")
			r.Play() // 满员开始游戏
		} else {
			log.Debug("当前战斗房人数:%d", r.GetMemberCount())
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
	for _, v := range r.Rooms {
		v.ChangeMemberState(MEMBERPLAYING)
	}

	// 获取题库
	r.LibAnswer = RandAnswerLib(5, GetAnswerLib())
	log.Debug("%s", r.LibAnswer.ToString())
	r.QuestionCount = len(r.LibAnswer.Answers)

	// 转移成员到开始玩
	r.SetDefaultAnswer() // 设置默认答案
	r.Send(remotemsg.ROOMSTARTPLAY, nil)
	log.Debug("开始比赛: %v", r.ID)
	r.PlayRun()
}

func (m *BattleRoom) endjudge() {
	if m.LibAnswer.Progress >= m.QuestionCount-1 {
		log.Debug("答题结束: %v", m.ID)
		m.Send(remotemsg.ROOMENDPLAY, nil)
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
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(m.AnswerTime))
	cur := m.AnswerTime
	m.isAnswerTime = true
	log.Debug("================（单次答题开始%d）===================", m.LibAnswer.Progress+1)
	skeleton.Go(func() {
		for {
			select {
			case <-ctx.Done():
				log.Debug("================单次答题结束==================")
				// 检查所有成员答案
				var allWrong = m.CheckAndHandleDead()
				// 如果没有全错
				m.isAnswerTime = false
				m.Send(remotemsg.ROOMANSWEREND, nil) // 答题结束

				if allWrong {
					if !m.isRelive {
						m.isRelive = true
						m.WaitRelive(func() {
							m.endjudge()
						}, func() {
							log.Debug("全员失败,答题结束: %v", m.ID)
							m.Send(remotemsg.ROOMENDPLAY, nil)
							m.OnPlayEnd()
						})
					} else {
						log.Debug("已经复活等待过, 全员失败,答题结束: %v", m.ID)
						m.Send(remotemsg.ROOMENDPLAY, nil)
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
						m.RandomRobotAnswer(4, 8, 10) // 机器人答题
					} else {
						m.RandomRobotAnswer(3, 5, 7) // 机器人答题
					}
				} else {
					log.Debug("不满足机器人运动条件****************************************")
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
		if !v.IsDead {
			all = false
		}
	})
	return all
}

func (m *BattleRoom) CheckAndHandleDead() bool {
	var allWrong = true
	rightCount := 0
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		if v.IsDead {
			return
		}
		q := m.GetQuestion().RightAnswer
		playerAnswer := v.Answer[m.LibAnswer.Progress]
		tip := ""
		if q == playerAnswer.Result {
			tip = "true---------------"
		}
		if v.IsMaster {
			log.Debug("房主***玩家uuid %v: 第%v题,  正确答案: %v, 玩家答案%v %v", v.Uuid, m.LibAnswer.Progress, q, playerAnswer.Result, tip)
		} else {
			log.Debug("玩家uuid %v: 第%v题, 正确答案: %v, 玩家答案%v %v", v.Uuid, m.LibAnswer.Progress, q, playerAnswer.Result, tip)
		}
		right := (q == playerAnswer.Result)
		if right {
			allWrong = false
			rightCount++
		} else {
			// 标记死亡
			v.IsDead = true
			var allDead = true
			if room != nil {
				// 检查房间是否所有人死亡
				for _, subV := range room.GetMembers() {
					if !subV.IsDead {
						allDead = false
					}
				}
				if allDead {
					log.Debug("发送319消息*****************************************")
					room.Send(remotemsg.ROOMALLFAIL, nil)
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
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		v.Answer = make([]*pmsg.Answer, m.QuestionCount)
		ranAnswer := results[rand.Intn(4)]
		for i, aa := range v.Answer {
			aa = &pmsg.Answer{
				Uuid:       v.Uuid,
				RoomID:     int32(m.ID),
				QuestionID: int32(m.LibAnswer.Progress),
				Result:     ranAnswer,
			}
			v.Answer[i] = aa
		}
	})
}

// 答题
func (m *BattleRoom) Answer(a *pmsg.Answer) {
	if !m.isAnswerTime {
		log.Debug("*********************非答题时间，无法进行答题*********************")
		return
	}
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		if v.Uuid == a.Uuid {
			for i, q := range v.Answer {
				if i >= m.LibAnswer.Progress {
					q.Result = a.Result
				}
			}
			m.SendAnswer(a.Uuid, a.QuestionID, a.Result)
			if !v.IsRobot {
				m.RandomRobotTargetAnswer(a.Result)
			}
		}
	})
}

func (m *BattleRoom) WaitRelive(success, fail func()) {
	subCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.ReliveWaitTime))
	log.Debug("等待复活-----------******************")
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
		case <-time.After(time.Millisecond * time.Duration(20)):
			if !m.CheckAllDead() {
				cancel()
			}
		}
	}
}
