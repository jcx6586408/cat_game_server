package internal

import (
	"config"
	"excel"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/name5566/leaf/log"
)

var (
	ROOMFULL    = 1 // 满员
	ROOMNULL    = 2 // 没找到对应得房间
	ROOMSTARTED = 3 // 房间已经开始游戏

	MEMEBERPREPARE = 1 // 成员正在等待
	MEMBERPLAYING  = 2 // 成员正在游玩

	Skins      []*Skin
	NamesLib   []*Names
	IconLib    []*Icon
	AnswerLibs []Answers

	RoomManager   Managerer
	manager       *Manager
	BattleManager BattleRoomManagerer
	battleManager *BattleRoomManager
	results       = []string{"A", "B", "C", "D"}
	RoomConf      *config.RoomConfig
	Questions     *QuestionLib
)

func ConstInit() {
	RoomConf = config.ReadRoom()
	manager = new(Manager)
	manager.IDManager = NewIDManager()
	manager.Pool = sync.Pool{
		New: func() any {
			r := &Room{}
			return r
		},
	}
	battleManager = new(BattleRoomManager)
	battleManager.IDManager = NewIDManager()
	battleManager.Pool = sync.Pool{
		New: func() any {
			r := &BattleRoom{}
			return r
		},
	}
	RoomManager = manager
	BattleManager = battleManager
	manager.TableManager = excel.Read()
	ExcelConfigUpdate()
}

func ExcelConfigUpdate() {
	Skins = ToSkinLib()
	NamesLib = ToNameLib()
	IconLib = ToIconLib()
	Questions = &QuestionLib{
		QuestionMap:      make(map[int]*Question),
		Question:         make(map[int][]*Question),
		PhaseQuestionLib: make(map[int][]int),
		WinRates:         make(map[int][]float32),
	}
	for i, v := range RoomConf.Question {
		Questions.PhaseQuestionLib[i+1] = v.AnswerPhase
		rates := []float32{}
		for _, ran := range v.WinRate {
			rates = append(rates, float32(ran))
		}
		Questions.WinRates[i+1] = rates
		Questions.Question[i+1] = make([]*Question, 0)
	}
	AnswerLibs = []Answers{}
	AnswerLibs = append(AnswerLibs, ToAnswerLib("question1"))

	for i := 0; i < 80; i++ {
		Questions.WinCount(10)
	}

	for i := 0; i < 20; i++ {
		Questions.FailCount(10)
	}

	re := Questions.GetQuestions(5)
	log.Debug("5段位题库长度%v", len(re))

	log.Debug("皮肤数量: %v", len(Skins))
	log.Debug("名字数量: %v", len(NamesLib))
	log.Debug("Icon数量: %v", len(IconLib))
	OnExit()
}

func OnExit() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGINT)

	skeleton.Go(func() {
		s := <-ch
		switch s {
		case syscall.SIGINT:
			//SIGINT 信号，在程序关闭时会收到这个信号
			fmt.Println("SIGINT", "退出程序，执行退出前逻辑")
			Questions.ToExcel()
			return
		case syscall.SIGKILL:
			fmt.Println("SIGKILL关闭********************")
			Questions.ToExcel()
			return
		default:
			fmt.Println("未知关闭")
			Questions.ToExcel()
		}
	}, func() {})
}
