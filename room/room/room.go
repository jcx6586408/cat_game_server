package room

import (
	"catLog"
	"config"
	"context"
	"fmt"
	"math/rand"
	"proto/msg"
	"sync"
	"time"
)

type Room struct {
	ID             int           // 房间ID
	PrepareMembers []*msg.Member // 处于准备的成员列表
	PlayingMembers []*msg.Member // 处于比赛的成员列表
	PrepareTime    int           // 比赛准备时间
	QuestionCount  int           // 问题总数量
	AnswerTime     int           // 单次回答问题时间
	MaxMember      int           // 最大成员数量
	Cur            int           // 当前房间时间

	CreateChan        chan int         // 房间创建成功
	AddMemberChan     chan *msg.Member // 加入成员管道
	LeaveMemberChan   chan *msg.Member // 离开成员管道
	StartPlayChan     chan interface{} // 开始游戏管道
	EndPlayChan       chan interface{} // 游戏结束管道
	PrepareChan       chan *msg.Member // 游戏准备通知
	PrepareCancelChan chan *msg.Member // 取消准备通知
	ChangeMasterChan  chan *msg.Member // 转移房主通知
	AnswerChan        chan interface{} // 回答消息通知
	MemberAnswerChan  chan *msg.Member // 成员回答消息通知
	OverChan          chan interface{} // 房间结束解散通知
	TimeChan          chan int         // 计时通知
	OfflineChan       chan *msg.Member // 离线通知

	LibAnswer *LibAnswer // 当前题库
	Lock      sync.RWMutex
}

func NewRoom(id int) *Room {
	r := &Room{}
	r.ID = id
	r.PrepareMembers = []*msg.Member{}
	r.PlayingMembers = []*msg.Member{}
	conf := config.ReadRoom()
	r.PrepareTime = conf.PrepareTime
	r.AnswerTime = conf.AnswerTime
	r.QuestionCount = conf.QuestionCount
	r.MaxMember = conf.MaxMember

	r.CreateChan = make(chan int)
	r.AddMemberChan = make(chan *msg.Member)
	r.EndPlayChan = make(chan interface{})
	r.LeaveMemberChan = make(chan *msg.Member)
	r.StartPlayChan = make(chan interface{})
	r.PrepareCancelChan = make(chan *msg.Member)
	r.PrepareChan = make(chan *msg.Member)
	r.ChangeMasterChan = make(chan *msg.Member)
	r.AnswerChan = make(chan interface{})
	r.MemberAnswerChan = make(chan *msg.Member)
	r.OverChan = make(chan interface{})
	r.TimeChan = make(chan int)
	r.OfflineChan = make(chan *msg.Member)
	return r
}

// 开始准备
func (m *Room) StartPrepare() {
	catLog.Log("开始比赛")
	// 获取题库
	var libID = rand.Intn(12) + 1
	m.LibAnswer = ToAnswerLib(fmt.Sprintf("question%v", libID))
	m.QuestionCount = len(m.LibAnswer.Answers)
	catLog.Log("题库总数量", m.QuestionCount)
	m.StartPlayChan <- 1
	catLog.Log("开始比赛通知")
	m.StartPlay()
}

func (m *Room) GetPlayTime() int {
	return m.QuestionCount*m.AnswerTime + 1
}

// 开始比赛
func (m *Room) StartPlay() {
	// 转移成员到开始玩
	Manager.ToUsingRoom(m) // 移动房间
	m.PrepareToPlaying()   // 移动成员
	total := m.GetPlayTime()
	catLog.Log("等待总时间", total, m.AnswerTime)
	// 创建完成通知
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(total)*time.Second)
	m.LibAnswer.Progress++ // 初始追加第一题进度
	m.SetDefaultAnswer()   // 设置默认答案

	m.Cur = 0
	go func() {
		for {
			select {
			case <-m.OverChan:
				return
			case <-ctx.Done():
				m.PlayingToPrepare() // 转移成员到准备
				m.EndPlayChan <- 1
				return
			case <-time.After(time.Duration(1) * time.Second):
				// 广播时间
				m.Cur++
				m.TimeChan <- m.Cur
				if m.Cur%m.AnswerTime == 0 {
					catLog.Log("房间_", m.ID, "答题结束，下一题")
					if m.LibAnswer != nil {
						// 进度增长
						m.LibAnswer.Progress++

						m.AnswerChan <- 1
					}
				}
			}

		}
	}()
}

func (m *Room) SetDefaultAnswer() {
	for _, v := range m.PlayingMembers {
		a := &msg.Answer{
			Uuid:       v.Uuid,
			RoomID:     int32(m.ID),
			QuestionID: int32(m.LibAnswer.Progress),
			Result:     "A",
		}
		bo := false
		for _, q := range v.Answer {
			if q.QuestionID == a.QuestionID {
				q.Result = a.Result // 改变答案
				bo = true
				break
			}
		}
		if !bo {
			v.Answer = append(v.Answer, a) // 追加答案
		}
	}
}

// 答题
func (m *Room) Answer(a *msg.Answer) {
	for _, v := range m.PlayingMembers {
		if v.Uuid == a.Uuid {
			bo := false
			for _, q := range v.Answer {
				if q.QuestionID == a.QuestionID {
					q.Result = a.Result // 改变答案
					bo = true
					break
				}
			}
			if !bo {
				v.Answer = append(v.Answer, a) // 追加答案
			}
			m.MemberAnswerChan <- v
			break
		}
	}
}

// 关闭房间
func (m *Room) Close(done chan interface{}) {
	// 清理房间
	m.PrepareMembers = m.PrepareMembers[0:0]
	m.PlayingMembers = m.PlayingMembers[0:0]
	m.Cur = 0
	m.LibAnswer = nil
	// 回收房间
	go func() {
		Manager.RecyleChan <- m
	}()
	m.OverChan <- 1 // 通知解散
}

func (m *Room) PrepareToPlaying() {
	// m.Lock.Lock()
	m.PlayingMembers = append(m.PlayingMembers, m.PrepareMembers...)
	m.PrepareMembers = m.PrepareMembers[0:0]
	// m.Lock.Unlock()
}

func (m *Room) PlayingToPrepare() {
	// m.Lock.Lock()
	m.PrepareMembers = append(m.PrepareMembers, m.PlayingMembers...)
	m.PlayingMembers = m.PlayingMembers[0:0]
	// m.Lock.Unlock()

}

// 房主加入
func (m *Room) AddMasterMember(member *msg.Member) {
	if m.IsFull() {
		return
	}
	m.PrepareMembers = append(m.PrepareMembers, member)
}

// 加入成员
func (m *Room) AddMember(member *msg.Member) {
	if m.IsFull() {
		return
	}
	member.IsMaster = false
	m.PrepareMembers = append(m.PrepareMembers, member)
	m.AddMemberChan <- member // 成员加入通知
}

// 离开准备成员
func (m *Room) LeavePrepareMember(member *msg.Member) {
	done := make(chan interface{})
	catLog.Log("离开前成员数量", len(m.PrepareMembers))
	m.PrepareMembers = m.Delete(m.PrepareMembers, member)
	catLog.Log("离开后*****成员数量", len(m.PrepareMembers))
	// 如果是房主，则移交房主
	if member.IsMaster {
		if len(m.PrepareMembers) > 0 {
			otherMember := m.PrepareMembers[0]
			otherMember.IsMaster = true
			go func(ot *msg.Member) {
				m.ChangeMasterChan <- ot
			}(otherMember)
		}
	}

	// 如果房间人数为0,则回收房间
	if len(m.PlayingMembers)+len(m.PrepareMembers) <= 0 {
		go m.Close(done)
	}
	m.LeaveMemberChan <- member // 成员离开通知
}

func (m *Room) LeavePlayingMember(member *msg.Member) {
	done := make(chan interface{})
	m.PlayingMembers = m.Delete(m.PlayingMembers, member)
	// 如果是房主，则移交房主
	if member.IsMaster {
		if len(m.PlayingMembers) > 0 {
			otherMember := m.PlayingMembers[0]
			catLog.Log("移交房主", otherMember.Uuid)
			otherMember.IsMaster = true
			go func(ot *msg.Member) {
				m.ChangeMasterChan <- ot
			}(otherMember)
		}
	}

	// 如果房间人数为0,则回收房间
	if len(m.PlayingMembers)+len(m.PrepareMembers) <= 0 {
		go m.Close(done)
	}

	m.LeaveMemberChan <- member // 成员离开通知
}

func (m *Room) GetProgress() int {
	if m.LibAnswer == nil {
		return 0
	}
	return m.LibAnswer.Progress
}

func (m *Room) GetQuestion() *msg.Question {
	if m.LibAnswer == nil {
		return nil
	}
	// return nil
	if m.LibAnswer.Progress < len(m.LibAnswer.Answers) {
		return m.LibAnswer.Answers[m.LibAnswer.Progress]
	}
	return m.LibAnswer.Answers[len(m.LibAnswer.Answers)-1]
}

func (m *Room) OfflinHanlder(uuid string) bool {
	for _, v := range m.PrepareMembers {
		if v.Uuid == uuid {
			m.LeavePrepareMember(v)
			m.OfflineChan <- v
			return true
		}
	}

	for _, v := range m.PlayingMembers {
		if v.Uuid == uuid {
			m.LeavePlayingMember(v)
			m.OfflineChan <- v
			return true
		}
	}
	return false
}

// 判断房间是否满员
func (m *Room) IsFull() bool {
	return m.MaxMember <= (len(m.PlayingMembers) + len(m.PrepareMembers))
}

func (m *Room) Delete(a []*msg.Member, elem *msg.Member) []*msg.Member {
	for i := 0; i < len(a); i++ {
		if a[i].Uuid == elem.Uuid {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}
