package internal

import (
	"context"
	"fmt"
	"math/rand"
	pmsg "proto/msg"
	"remotemsg"
	"time"

	"github.com/google/uuid"
	"github.com/name5566/leaf/log"
)

// 匹配加入机器人
func (m *BattleRoom) matching() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(8))
	m.MatchCancel = cancel
	cur := 0
	skeleton.Go(func() {
		for {
			select {
			case <-ctx.Done():
				// 如果没有真人，直接退出
				if !m.CheckRealMember() {
					return
				}
				var less = m.Max - m.GetMemberCount()
				if less > 0 {
					m.AddRobot(less, NamesLib, IconLib)
					m.Send(remotemsg.ROOMADD, nil)
					// 将房间移入比赛使用房间
				}
				m.Play()
				log.Debug("开始游戏, 退出等待加入============%v", m.ID)
				return
			case <-time.After(time.Duration(1) * time.Second):
				// 如果没有真人，直接退出
				if !m.CheckRealMember() {
					return
				}
				cur++
				if cur <= 3 {
					m.AddRandomCountRobots(4, 7, func() { cancel() })
				}
				if cur >= 5 {
					m.AddRandomCountRobots(1, 3, func() { cancel() })
				}
			}
		}
	}, func() {})
}

func (m *BattleRoom) AddRandomCountRobots(min, max int, callback func()) {
	var less = m.Max - m.GetMemberCount()
	ranNumber := rand.Intn(max+1) + min
	if ranNumber > less {
		log.Debug("补满所有机器人, %d", less)
		m.AddRobot(less, NamesLib, IconLib)
		callback()
	} else {
		log.Debug("补充指定数量机器人, %d", ranNumber)
		m.AddRobot(ranNumber, NamesLib, IconLib)
	}
	m.Send(remotemsg.ROOMADD, nil)
}

func (m *BattleRoom) AddRobot(count int, nameLib []*Names, iconLib []*Icon) {
	subName := RandName(count, nameLib)
	subIcon := RandIcon(count, iconLib)
	for i := 0; i < count; i++ {
		guid := uuid.New().String()
		skinID := 1
		if rand.Intn(10) > 7 {
			skinID = 2 + rand.Intn(len(Skins))
		}
		m.Members = append(m.Members, &pmsg.Member{
			Nickname: fmt.Sprintf("%v", subName[i].ID),
			Uuid:     guid,
			Icon:     fmt.Sprintf("%v", subIcon[i].ID),
			IsMaster: false,
			IsRobot:  true,
			SkinID:   int32(skinID),
		})
	}
}

// 随机机器人答案
func (m *BattleRoom) RandomRobotAnswer(min, max, count int) {
	lenRobot := rand.Intn(max) + min
	startIndex := 0
	if len(m.Members)-lenRobot > 0 {
		startIndex = rand.Intn(len(m.Members) - lenRobot)
	}
	if startIndex >= len(m.Members) {
		return
	}
	if startIndex+lenRobot >= len(m.Members) {
		return
	}
	subArr := m.Members[startIndex : startIndex+lenRobot]
	for _, v := range subArr {
		result := rand.Intn(4)

		var action = rand.Intn(10)
		if action >= count {
			return
		}
		m.Answer(&pmsg.Answer{
			Uuid:       v.Uuid,
			RoomID:     int32(m.ID),
			QuestionID: m.GetQuestion().ID,
			Result:     results[result],
		})
	}
}

func (m *BattleRoom) RandomRobotTargetAnswer(right string) {
	lenRobot := rand.Intn(3) + 1
	startIndex := 0
	if len(m.Members)-lenRobot > 0 {
		startIndex = rand.Intn(len(m.Members) - lenRobot)
	}
	if startIndex >= len(m.Members) {
		return
	}
	if startIndex+lenRobot >= len(m.Members) {
		return
	}
	subArr := m.Members[startIndex : startIndex+lenRobot]
	for _, v := range subArr {
		m.Answer(&pmsg.Answer{
			Uuid:       v.Uuid,
			RoomID:     int32(m.ID),
			QuestionID: m.GetQuestion().ID,
			Result:     right,
		})
	}
}

// 检查真实成员是否还存在
func (m *BattleRoom) CheckRealMember() bool {
	for _, room := range m.Rooms {
		if room.GetMemberCount() > 0 {
			return true
		}
	}
	return false
}
