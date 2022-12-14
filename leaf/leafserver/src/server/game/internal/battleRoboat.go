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
				var less = m.Max - m.GetMemberCount()
				if less > 0 {
					m.AddRobot(less, NamesLib, IconLib)
					m.send(remotemsg.ROOMADD, nil)
					// 将房间移入比赛使用房间
				}
				// Manager.MatchingToUsingRoom(m)
				m.Play()
				log.Debug("开始游戏, 退出等待加入============%v", m.ID)
				return
			case <-time.After(time.Duration(1) * time.Second):
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
		m.AddRobot(less, NamesLib, IconLib)
		callback()
	} else {
		m.AddRobot(ranNumber, NamesLib, IconLib)
	}
	m.send(remotemsg.ROOMADD, nil)
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
