package internal

import (
	"config"
	"excel"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"

	"github.com/name5566/leaf/db/mongodb"
	"github.com/name5566/leaf/log"
)

var (
	ROOMFULL    = 1 // 满员
	ROOMNULL    = 2 // 没找到对应得房间
	ROOMSTARTED = 3 // 房间已经开始游戏

	MEMEBERPREPARE     = 1 // 成员正在等待
	MEMBERPLAYING      = 2 // 成员正在游玩
	MEMEBENONERPREPARE = 3 // 成员未准备
	MEMBERMATCHING     = 4 // 成员匹配

	MATCHINGTIME = 5 // 匹配最长时间

	Skins            []*Skin
	NamesLib         []*Names
	IconLib          []*Icon
	LevelLib         []*Level
	AnswerLibs       []Answers // 标准题库
	LowestAnswerLibs Answers   // 文盲题库
	Scenes           []*Scene
	Scores           []*Score

	RoomManager   Managerer
	manager       *Manager
	BattleManager BattleRoomManagerer
	battleManager *BattleRoomManager
	results       = []string{"A", "B", "C", "D"}
	RoomConf      *config.RoomConfig
	Questions     *QuestionLib
	ServerConf    *config.Config
	MD            *mongodb.DialContext
	MAX           int // 最大段位等级
)

func ConstInit() {
	RoomConf = config.ReadRoom()
	ServerConf = config.Read()
	// 数据库连接
	// MongoConnect()
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
	LevelLib = ToLevelLib()
	Scenes = ToSceneLib()
	Scores = ToScoreLib()
	log.Release("分数与时间数量: %v|%v|%v", len(Scores), Scores[0].ID, Scores[1].ID)
	Questions = &QuestionLib{
		QuestionMap:      make(map[string]*Question),
		Question:         make(map[int][]*Question),
		PhaseQuestionLib: make(map[int][]int),
		WinRates:         make(map[int][]float32),
		WinChan:          make(chan string),
		FailChan:         make(chan string),
		Done:             make(chan interface{}),
	}

	sort.SliceStable(LevelLib, func(i, j int) bool {
		return LevelLib[i].ID < LevelLib[j].ID
	})

	for i, v := range LevelLib {
		Questions.PhaseQuestionLib[i+1] = v.AnswerPhase
		rates := []float32{}
		for _, ran := range v.WinRate {
			rates = append(rates, float32(ran))
		}
		Questions.WinRates[i+1] = rates
		Questions.Question[i+1] = make([]*Question, 0)
	}
	MAX = len(LevelLib)
	AnswerLibs = []Answers{}
	tables := ServerConf.QuestionTables
	for _, v := range tables {
		AnswerLibs = append(AnswerLibs, ToAnswerLib(v))
	}
	LowestAnswerLibs = ToBaseAnswerLib(ServerConf.NewerTable, nil)
	log.Release("新手题库数量: %v", len(LowestAnswerLibs))
	log.Release("标准题库数量: %v", len(Questions.QuestionMap))
	for i, v := range Questions.Question {
		log.Release("段位: %v --> %v", i, len(v))
	}
	MongoConnect()  // 数据库连接
	Questions.Run() // 题库监听
	log.Release("段位长度: %d", len(LevelLib))
	log.Release("皮肤数量: %v", len(Skins))
	log.Release("名字数量: %v", len(NamesLib))
	log.Release("Icon数量: %v", len(IconLib))
	// OnExit()
}

func OnExit() {
	// util.DeepClone(Users)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGINT)

	skeleton.Go(func() {

		s := <-ch
		close(Questions.Done)
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
