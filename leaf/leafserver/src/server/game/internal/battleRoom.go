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
	Play()                    // 开始游戏
	OnPlayEnd()               // 游戏结束
	AddRoom(room Roomer) bool // 加入房间
	Relive(uuid string)       // 复活
	Answer(a *pmsg.Answer)    // 答题
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
	ReliveChan     chan int           // 复活通知
	ReliveWaitTime int
}

func (r *BattleRoom) GetID() int {
	return r.ID
}

func (r *BattleRoom) OnInit() {
	r.Members = []*pmsg.Member{}
	r.Rooms = []Roomer{}
}

func (r *BattleRoom) OnClose() {
	r.LibAnswer = nil
	r.Members = nil
	r.Rooms = nil
	r.ReliveChan = nil
	r.Cancel = nil
}

func (r *BattleRoom) OnPlayEnd() {
	// 回归所有成员状态
	for _, v := range r.Rooms {
		v.ChangeMemberState(MEMEBERPREPARE)
	}
}

func (r *BattleRoom) Relive(uuid string) {
	for _, room := range r.Rooms {
		for _, v := range room.GetMembers() {
			if v.Uuid == uuid {
				v.IsDead = false
				skeleton.Go(func() {
					if r.ReliveChan != nil {
						log.Debug("发送复活通知***************************************")
						r.ReliveChan <- 1
					} else {
						log.Debug("管道已经关闭, 无需通知")
					}
				}, func() {})
				break
			}
		}
	}
}

func (r *BattleRoom) AddMember(member Memberer) {

}

func (r *BattleRoom) LeaveMember(member Memberer) {

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
	less := r.Max - r.GetMemberCount()
	if less >= room.GetMemberCount() {
		r.Rooms = append(r.Rooms, room)
		return true
	}
	// 填充机器人准备战斗

	return false
}

func (m *BattleRoom) GetPlayTime() int {
	return m.QuestionCount*m.AnswerTime + 1
}

func (r *BattleRoom) Play() {
	// 设置所有成员状态为游玩
	for _, v := range r.Rooms {
		v.ChangeMemberState(MEMBERPLAYING)
	}

	// 获取题库
	r.LibAnswer = RandAnswerLib(5, GetAnswerLib())
	r.QuestionCount = len(r.LibAnswer.Answers)

	// 转移成员到开始玩
	r.SetDefaultAnswer() // 设置默认答案
	// r.Send(remotemsg.ROOMSTARTPLAY)
	log.Debug("开始比赛: %v", r.ID)
	r.PlayRun()
}

func (m *BattleRoom) PlayRun() {
	total := m.GetPlayTime()
	log.Debug("房间ID: %v,等待总时间:%v, 单词答题时间:%v", m.ID, total, m.AnswerTime)
	// 创建完成通知
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(total))
	m.Cancel = cancel
	m.Cur = 0 // 答题总时间计时
	cur := 0  // 单局答题时间计时
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debug("答题结束: %v", m.ID)
				// m.Send(remotemsg.ROOMENDPLAY)
				m.send(remotemsg.ROOMENDPLAY, nil)
				m.OnPlayEnd()
				// m.PlayingToPrepare()          // 转移成员到准备
				// Manager.UsingToPrepareRoom(m) // 转移房间位置
				m.Cancel = nil
				return
			case <-time.After(time.Duration(1) * time.Second):
				// 广播时间
				m.Cur++
				cur++
				log.Debug("房间ID:%v,当前时间:%v, 单次答题时间: %v", m.ID, m.Cur, cur)
				if m.Cur >= 6 && cur >= 3 && cur <= m.AnswerTime-2 {
					if cur <= 10 {
						// m.RandomRobotAnswer(2, 8, 10) // 机器人答题
					} else {
						// m.RandomRobotAnswer(1, 3, 5) // 机器人答题
					}
				}
				// m.SendTime(m.Cur)
				if m.Cur%m.AnswerTime == 0 {
					cur = 0
					log.Debug("房间_: %v,  房间游戏人数: %v, 当前进度: %v", m.ID, m.GetMemberCount(), m.LibAnswer.Progress+1)
					if m.LibAnswer != nil {
						// 检查所有成员答案
						var allWrong = m.CheckAndHandleDead()
						// 如果没有全错
						m.LibAnswer.Progress++               // 进度增长
						m.send(remotemsg.ROOMANSWEREND, nil) // 答题结束
						if allWrong {
							m.WaitRelive(func() {
								cancel()
							}) // 等待复活
						}
					}
				}
			}

		}
	}()
}

func (m *BattleRoom) CheckAndHandleDead() bool {
	var allWrong = true
	rightCount := 0
	for _, room := range m.Rooms {
		for _, v := range room.GetMembers() {
			if v.IsDead {
				continue
			}
			q := m.GetQuestion().RightAnswer
			playerAnswer := v.Answer[m.LibAnswer.Progress]
			tip := ""
			if q == playerAnswer.Result {
				tip = "true---------------"
			}
			if v.IsMaster {
				log.Debug("房主***玩家uuid %v: 正确答案: %v, 玩家答案%v %v", v.Uuid, q, playerAnswer.Result, tip)
			} else {
				log.Debug("玩家uuid %v: 正确答案: %v, 玩家答案%v %v", v.Uuid, q, playerAnswer.Result, tip)
			}
			right := (q == playerAnswer.Result)
			if right {
				allWrong = false
				rightCount++
			} else {
				// 标记死亡
				v.IsDead = true
			}
		}
	}
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
	for _, room := range m.Rooms {
		for _, v := range room.GetMembers() {
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
		}
	}
}

// 答题
func (m *BattleRoom) Answer(a *pmsg.Answer) {
	for _, room := range m.Rooms {
		for _, v := range room.GetMembers() {
			if v.Uuid == a.Uuid {
				for i, q := range v.Answer {
					if i >= m.LibAnswer.Progress {
						q.Result = a.Result
					}
				}
				m.SendAnswer(a.Uuid, a.QuestionID, a.Result)
				if !v.IsRobot {
					// m.RandomRobotTargetAnswer(a.Result)
				}
				break
			}
		}
	}
}

func (m *BattleRoom) WaitRelive(fail func()) {
	subCtx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(m.ReliveWaitTime))
	log.Debug("创建复活管道通知++++++++++++++++++++++++%d", m.ID)
	m.ReliveChan = make(chan int)
	skeleton.Go(
		func() {
			for {
				select {
				case <-m.ReliveChan:
					close(m.ReliveChan)
					m.ReliveChan = nil
					log.Debug("复活成功--------继续游戏%v", m.ReliveChan)
					return
				case <-subCtx.Done():
					close(m.ReliveChan)
					log.Debug("比赛结束-------------------房间ID: %v", m.ID)
					fail()
					return
				}
			}
		},
		func() {},
	)
}
