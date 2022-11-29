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

	AddMemberChan     chan *msg.Member // 加入成员管道
	LeaveMemberChan   chan *msg.Member // 离开成员管道
	StartPlayChan     chan interface{} // 开始游戏管道
	EndPlayChan       chan interface{} // 游戏结束管道
	PrepareChan       chan *msg.Member // 游戏准备通知
	PrepareCancelChan chan *msg.Member // 取消准备通知
	ChangeMasterChan  chan *msg.Member // 转移房主通知
	AnswerChan        chan interface{} // 回答消息通知
	OverChan          chan interface{} // 房间结束解散通知
	TimeChan          chan int         // 计时通知

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

	r.AddMemberChan = make(chan *msg.Member)
	r.EndPlayChan = make(chan interface{})
	r.LeaveMemberChan = make(chan *msg.Member)
	r.StartPlayChan = make(chan interface{})
	r.PrepareCancelChan = make(chan *msg.Member)
	r.PrepareChan = make(chan *msg.Member)
	r.ChangeMasterChan = make(chan *msg.Member)
	r.AnswerChan = make(chan interface{})
	r.OverChan = make(chan interface{})
	r.TimeChan = make(chan int)
	return r
}

// 开始准备
func (m *Room) StartPrepare() {
	go func() {
		<-time.After(time.Duration(m.PrepareTime) * time.Second)
		catLog.Log("准备时间结束")
		// 获取题库
		var libID = rand.Intn(12) + 1
		m.LibAnswer = ToAnswerLib(fmt.Sprintf("question%v", libID))
		m.QuestionCount = len(m.LibAnswer.Answers)
		m.StartPlay()
	}()
}

func (m *Room) GetPlayTime() int {
	return m.QuestionCount*m.AnswerTime + 1
}

// 开始比赛
func (m *Room) StartPlay() {
	// 转移成员到开始玩
	Manager.ToUsingRoom(m) // 移动房间
	m.prepareToPlaying()   // 移动成员
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
			case <-ctx.Done():
				m.playingToPrepare() // 转移成员到准备
				return
			case <-time.After(time.Duration(1) * time.Second):
				// 广播时间
				m.Cur++
				m.TimeChan <- m.Cur
				if m.Cur%m.AnswerTime == 0 {
					catLog.Log("房间_", m.ID, "答题结束，下一题")
					// 进度增长
					m.LibAnswer.Progress++
					m.AnswerChan <- 1
				}
			}

		}
	}()
}

func (m *Room) SetDefaultAnswer() {
	for _, v := range m.PlayingMembers {
		m.Answer(&msg.Answer{
			Member:     v,
			RoomID:     int32(m.ID),
			QuestionID: int32(m.LibAnswer.Progress),
			Result:     "A",
		})
	}
}

// 答题
func (m *Room) Answer(a *msg.Answer) {
	for _, v := range m.PlayingMembers {
		if v.Uuid == a.Member.Uuid {
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
			break
		}
	}
}

// 关闭房间
func (m *Room) Close() {
	catLog.Log("房间关闭", len(m.PlayingMembers), len(m.PrepareMembers))
	m.OverChan <- 1
	// 清理房间
	m.PrepareMembers = m.PrepareMembers[0:0]
	m.PlayingMembers = m.PlayingMembers[0:0]
	m.Cur = 0
	// 回收房间
	Manager.RecyleChan <- m
}

func (m *Room) prepareToPlaying() {
	m.Lock.Lock()
	m.PlayingMembers = append(m.PlayingMembers, m.PrepareMembers...)
	m.PrepareMembers = m.PrepareMembers[0:0]
	m.Lock.Unlock()
}

func (m *Room) playingToPrepare() {
	m.Lock.Lock()
	m.PrepareMembers = append(m.PrepareMembers, m.PlayingMembers...)
	m.PlayingMembers = m.PlayingMembers[0:0]
	m.Lock.Unlock()
}

// 加入成员
func (m *Room) AddMember(member *msg.Member) {
	if m.IsFull() {
		return
	}
	m.Lock.Lock()
	m.PrepareMembers = append(m.PrepareMembers, member)
	m.AddMemberChan <- member // 成员加入通知
	m.Lock.Unlock()
}

// 离开准备成员
func (m *Room) LeavePrepareMember(member *msg.Member) {
	m.Lock.Lock()
	m.PrepareMembers = m.Delete(m.PrepareMembers, member)
	m.LeaveMemberChan <- member // 成员离开通知
	m.Lock.Unlock()
}

func (m *Room) LeavePlayingMember(member *msg.Member) {
	m.Lock.Unlock()
	m.PlayingMembers = m.Delete(m.PlayingMembers, member)
	m.LeaveMemberChan <- member // 成员离开通知
	m.Lock.Unlock()
}

// 广播满员
func (m *Room) BroadcastFull() {

}

// 判断房间是否满员
func (m *Room) IsFull() bool {
	return m.MaxMember <= (len(m.PlayingMembers) + len(m.PrepareMembers))
}

func (m *Room) Delete(a []*msg.Member, elem *msg.Member) []*msg.Member {
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}
