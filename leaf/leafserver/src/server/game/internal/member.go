package internal

import pmsg "proto/msg"

type Memberer interface {
	GetUuid() string
	GetUid() string
	GetNickname() string
	GetIcon() string
	GetIsMaster() bool
	GetIsRobot() bool
	GetIsInvite() bool
	GetIsDead() bool
	GetAnswer() []*pmsg.Answer
	GetSkinID() int32
	GetRoomID() int32
	GetBattleRoomID() int32
	GetState() int32
}

var (
	PREPARE  = 1 // 准备
	MATCHING = 2 // 匹配
	PLAYINT  = 3 // 游玩
)
