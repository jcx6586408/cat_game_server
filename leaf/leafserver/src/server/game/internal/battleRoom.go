package internal

import (
	"context"
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
	AnswerTime     int                // 单次回答问题时间
	Cur            int                // 答题总时间
	SingleCur      int                // 单次答题时间
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
	r.robotNames = make([]*Names, len(NamesLib))
	r.Done = make(chan interface{})
	copy(r.robotIcons, IconLib)
	copy(r.robotNames, NamesLib)
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
	r.SetDefaultAnswer()
	// 转移成员到开始玩
	r.SendStart()
	log.Debug("开始比赛: ----房间ID: %v;----题库数量: %v;----答题时间: %v", r.ID, levelConf.QuestionNumber, levelConf.QuestionTime)
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
		m.LibAnswer.Progress++ // 进度增长
		log.Debug("%s", m.LibAnswer.SingleToString())
		<-time.After(time.Second * time.Duration(m.PrepareTime))
		m.Send(remotemsg.ROOMANSWERSTART, nil) // 答题开始
		m.singleRun()
	}
}

func (m *BattleRoom) singleRun() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(m.AnswerTime+1))
	// cur := m.AnswerTime + 1
	m.SingleCur = m.AnswerTime + 1
	m.SendTime(m.SingleCur)
	m.isAnswerTime = true
	if m.LibAnswer == nil {
		return
	}
	m.SingleReset() // 状态重置
	skeleton.Go(func() {
		for {
			select {
			case <-m.Done:

				return
			case <-ctx.Done():
				// 检查所有成员答案
				m.CheckAndHandleDead()
				// 如果没有全错
				m.isAnswerTime = false
				m.Send(remotemsg.ROOMANSWEREND, nil) // 答题结束
				m.endjudge()
				return
			case <-time.After(time.Second * time.Duration(1)):
				m.Cur++
				m.SingleCur--
				// log.Debug("房间: %d,当前时间:%d", m.ID, m.SingleCur)
				m.SendTime(m.SingleCur)
				isEnd := m.CheckAllAnswered()
				if isEnd {
					cancel() // 结算
				} else {
					cur := m.AnswerTime + 1 - m.SingleCur
					if cur > 0 && cur <= 3 {
						m.RandomRobotAnswerWithRate(0, 2) // 机器人答题
					} else if cur > 3 && cur <= 10 {
						m.RandomRobotAnswerWithRate(1, 3) // 机器人答题
					} else {
						m.RandomRobotAnswerWithRate(1, 3) // 机器人答题
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

func (m *BattleRoom) SetDefaultAnswer() {
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		v.Answer = []*pmsg.Answer{}
		for i := 0; i < m.QuestionCount; i++ {
			answer := &pmsg.Answer{}
			answer.Uuid = v.Uuid
			answer.RoomID = int32(m.ID)
			answer.QuestionID = ""
			answer.Result = ""
			v.Answer = append(v.Answer, answer)
		}
	})
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
		if v.State == int32(MEMEBERPREPARE) || v.State == int32(MEMEBENONERPREPARE) {
			return
		}
		question := m.GetQuestion()
		q := question.RightAnswer
		if v.Answer == nil || len(v.Answer) <= m.LibAnswer.Progress {
			log.Debug("尚未答题")
			return
		}
		playerAnswer := v.Answer[m.LibAnswer.Progress]

		// 统计
		right := (q == playerAnswer.Result)
		if right {
			allWrong = false
			// 分数增长
			rightCount++
			skeleton.Go(func() {
				if !v.IsRobot {
					Questions.WinChan <- m.GetQuestion().ID
				}
			}, func() {})
		} else {
			skeleton.Go(func() {
				if !v.IsRobot {
					Questions.FailChan <- m.GetQuestion().ID
				}
			}, func() {})
		}
	})
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

// 重置
func (m *BattleRoom) SingleReset() {
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		v.IsAnswered = false
	})
}

// 检查是否所有都已答题
func (m *BattleRoom) CheckAllAnswered() bool {
	all := true
	m.foreachMembers(func(v *pmsg.Member, room Roomer) {
		if v.State == int32(MEMEBERPREPARE) || v.State == int32(MEMEBENONERPREPARE) {
			return
		}
		if !v.IsAnswered {
			all = false
		}
	})
	return all
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
					break
				}
			}
			v.IsAnswered = true
			question := m.GetQuestion()
			q := question.RightAnswer
			right := (q == a.Result) // 是否回答正确
			if right {
				// 增加分数
				index := m.AnswerTime + 1 - m.SingleCur
				if index >= len(Scores) {
					index = len(Scores) - 1
				}
				v.Score += int32(Scores[index].Score)
			}
			m.SendAnswer(a.Uuid, a.QuestionID, a.Result)
		}
	})
}
