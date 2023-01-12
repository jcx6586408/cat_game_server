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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(MATCHINGTIME))
	m.MatchCancel = cancel
	cur := 0
	skeleton.Go(func() {
		for {
			select {
			case <-ctx.Done():
				// 如果没有真人，直接退出
				if !m.CheckRealMember() {
					// 回收房间
					log.Debug("房间ID: %d, 匹配中真人全部退出，回收房间", m.GetID())
					return
				}
				var less = m.Max - m.GetMemberCount()
				if less > 0 {
					m.AddRobot(less)
					m.Send(remotemsg.ROOMADD, nil)
				}
				// 将房间移入比赛使用房间
				m.Play()
				return
			case <-time.After(time.Duration(1) * time.Second):
				// 如果没有真人，直接退出
				if !m.CheckRealMember() {
					log.Debug("房间ID: %d, 匹配中真人全部退出，回收房间", m.GetID())
					return
				}
				var less = m.Max - m.GetMemberCount()
				// log.Debug("机器人匹配: %d----%d----%d", less, m.Max, m.GetMemberCount())
				// 人数已满，直接退出
				if less <= 0 {
					cancel()
				}
				cur++
				if cur <= MATCHINGTIME {
					switch cur {
					case 1:
						if m.GetMemberCount() <= 20 {
							m.AddRandomCountRobots(4, 8, func() { cancel() })
						}
					case 2:
						if m.GetMemberCount() <= 20 {
							m.AddRandomCountRobots(4, 8, func() { cancel() })
						}
					case 3:
						if m.GetMemberCount() <= 30 {
							m.AddRandomCountRobots(4, 6, func() { cancel() })
						}
					case 4:
						if m.GetMemberCount() <= 30 {
							m.AddRandomCountRobots(4, 6, func() { cancel() })
						}
					}

				}
			}
		}
	}, func() {})
}

func (m *BattleRoom) AddRandomCountRobots(min, max int, callback func()) {
	var less = m.Max - m.GetMemberCount()
	ranNumber := rand.Intn(max+1) + min
	if ranNumber > less {
		m.AddRobot(less)
		callback()
	} else {
		m.AddRobot(ranNumber)
	}
	m.Send(remotemsg.ROOMADD, nil)
}

func (m *BattleRoom) AddRobot(count int) {
	subIcon, ir := RandIconClip(count, m.robotIcons)
	m.robotIcons = ir
	for i := 0; i < count; i++ {
		guid := uuid.New().String()
		skinID := 1
		if rand.Intn(10) > 7 {
			skinID = 2 + rand.Intn(len(Skins)-1)
		}
		m.Members = append(m.Members, &pmsg.Member{
			Nickname: "",
			Uuid:     guid,
			Icon:     fmt.Sprintf("%v", subIcon[i].ID),
			IsMaster: false,
			IsRobot:  true,
			SkinID:   int32(skinID),
			State:    int32(MEMBERPLAYING),
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
