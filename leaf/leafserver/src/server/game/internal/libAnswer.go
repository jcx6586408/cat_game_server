package internal

import (
	"excel"
	"fmt"
	"math/rand"
	pmsg "proto/msg"

	"github.com/name5566/leaf/log"
)

type Answers []*pmsg.Question

type LibAnswer struct {
	Answers  Answers
	Progress int
	Name     string
}

type Question struct {
	Q           *pmsg.Question
	win         int
	fail        int
	updateCount int
}

type QuestionLib struct {
	QuestionMap      map[string]*Question
	Question         map[int][]*Question // 段位题库
	PhaseQuestionLib map[int][]int       // 各个段位对于的题库
	WinRates         map[int][]float32   // 各个段位对应胜率
	WinChan          chan string         // 统计管道
	FailChan         chan string         // 统计管道
	Done             chan interface{}    // 完成管道
}

func (q *QuestionLib) Run() {
	skeleton.Go(func() {
		for {
			select {
			case <-q.Done:
				return
			case id := <-q.WinChan:
				q.WinCount(id)
			case id := <-q.FailChan:
				q.FailCount(id)
			}
		}
	}, func() {})
}

func GetRightAnswer(r string, obj *pmsg.Question) string {
	switch r {
	case "A":
		return obj.AnswerA
	case "B":
		return obj.AnswerB
	case "C":
		return obj.AnswerC
	case "D":
		return obj.AnswerD
	default:
		return ""
	}
}

func GetTestQuestion(label string, level int) *pmsg.Question {
	if level == 1 {
		lib := LowestAnswerLibs
		var ran = rand.Intn(len(lib))
		return lib[ran]
	}
	var lib = Questions.Question[level]
	if len(lib) > 0 {
		var ran = rand.Intn(len(lib))
		return lib[ran].Q
	}
	return nil
}

func GetRightNumberAnswer(r string, obj *pmsg.Question) int {
	switch r {
	case "A":
		return 0
	case "B":
		return 1
	case "C":
		return 2
	case "D":
		return 3
	default:
		return 0
	}
}

func (l *LibAnswer) ToString() string {
	var str = []string{}
	for i, v := range l.Answers {
		str = append(str, fmt.Sprintf("\n问题%d: %s------答案：%s", i+1, v.Question, GetRightAnswer(v.RightAnswer, v)))
	}
	return fmt.Sprintln(str)
}

func (l *LibAnswer) SingleToString() string {
	i := l.Progress
	v := l.Answers[l.Progress]
	return fmt.Sprintf("\n问题%d: %s------答案：%s", i+1, v.Question, GetRightAnswer(v.RightAnswer, v))
}

func GetAnswerLib() Answers {
	ran := rand.Intn(len(AnswerLibs))
	return AnswerLibs[ran]
}

func ToBaseAnswerLib(table string, callback func(obj *pmsg.Question)) Answers {
	tableConf, ok := manager.TableManager.GetTable(table)
	arr := []*pmsg.Question{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &pmsg.Question{}
			obj.Table = table
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					// obj.ID = int32(cell[index].(int))
					obj.ID = fmt.Sprintf("%v_%v", table, int32(cell[index].(int)))
				case "question":
					obj.Question = cell[index].(string)
				case "answerA":
					obj.AnswerA = cell[index].(string)
				case "answerB":
					obj.AnswerB = cell[index].(string)
				case "answerC":
					obj.AnswerC = cell[index].(string)
				case "answerD":
					obj.AnswerD = cell[index].(string)
				case "rightAnswer":
					obj.RightAnswer = cell[index].(string)
					switch obj.RightAnswer {
					case "A":
					case "B":
					case "C":
					case "D":
					default:
						log.Error("%v|%v|%v非法答案, %v", table, obj.ID, obj.Question, obj.RightAnswer)
					}
				case "rightNumber":
					obj.RightNumber = int32(cell[index].(int))
				case "wrongNumber":
					obj.WrongNumber = int32(cell[index].(int))
				case "label":
					obj.Label = cell[index].(string)
				}
			}
			arr = append(arr, obj)
			if callback != nil {
				callback(obj)
			}
		}
	}
	return arr
}

func ToAnswerLib(table string) Answers {
	var as = ToBaseAnswerLib(table, func(obj *pmsg.Question) {
		q := &Question{
			Q:    obj,
			win:  int(obj.RightNumber),
			fail: int(obj.WrongNumber),
		}
		Questions.QuestionMap[obj.ID] = q
		subA, ok := Questions.Question[1]
		if !ok {
			subA = make([]*Question, 0)
			subA = append(subA, q)
			Questions.Question[1] = subA // 初始默认设置为段位1的题库
		} else {
			Questions.Question[1] = append(Questions.Question[1], q)
		}
	})
	// 更新题库
	for _, v := range Questions.QuestionMap {
		Questions.updateLib(v)
	}
	return as
}

func RandAnswerLib(count int, arr []*pmsg.Question) *LibAnswer {
	log.Debug("-------------------随机题库长度: %v", len(arr))
	// 获得其实随机数
	startIndex := 0
	length := len(arr)
	if len(arr)-count > 0 {
		startIndex = rand.Intn(length - count)
	}

	// 长度不够则补充
	if count >= length && len(arr) > 0 {
		for i := 0; i < count; i++ {
			arr = append(arr, arr[0])
		}
	}

	ex := []*pmsg.Question{}
	ex = append(ex, arr[startIndex:startIndex+count]...)
	// 返回子数组
	return &LibAnswer{
		Answers:  ex,
		Progress: 0,
	}
}

func (m *QuestionLib) RandAnswerLib(l, count int) *LibAnswer {
	log.Debug("传入进来的等级: %v", l)
	if l == 1 {
		log.Debug("新手题库----------------------------")
		return RandAnswerLib(count, LowestAnswerLibs)
	}
	log.Debug("正规题库----------------------------")
	var lib = m.GetQuestions(l, count)
	return RandAnswerLib(count, lib)
}

func (m *QuestionLib) _getQuestions(id int, usedids []int) ([]*pmsg.Question, []int) {
	arr := []*pmsg.Question{}
	ids, ok := m.PhaseQuestionLib[id]
	if ok {
		for _, phase := range ids {
			var used = false
			for _, v := range usedids {
				if v == phase {
					used = true
					log.Debug("段位题库: %d, 已经被搜索过，无需再次搜索", used)
					break
				}
			}
			if !used {
				for _, v := range m.Question[phase] {
					arr = append(arr, v.Q)
				}
			}
		}
	}
	return arr, ids
}

// 根据段位获取题库
func (m *QuestionLib) GetQuestions(id, total int) []*pmsg.Question {
	ARR := []*pmsg.Question{}
	count := id
	c := 0
	var used = []int{}
	for {
		if len(ARR) <= total {
			sub, ids := m._getQuestions(count, used)
			used = append(used, ids...)
			ARR = append(ARR, sub...)
			count--
			if count < 0 {
				count = MAX
			}
			c++
			if c >= MAX {
				break
			}
		} else {
			break
		}
	}
	return ARR
}

func (m *QuestionLib) WinCount(id string) {
	q, ok := m.QuestionMap[id]
	if ok {
		q.win++
		q.Q.WrongNumber++
		m.updateLib(q)
	}
}

func (m *QuestionLib) FailCount(id string) {
	q, ok := m.QuestionMap[id]
	if ok {
		q.fail++
		q.Q.WrongNumber++
		m.updateLib(q)
	}
}

func (m *QuestionLib) updateQuestion(q *Question) {
	rate := float32(q.win) * 100 / (float32(q.win) + float32(q.fail))
	for i, v := range m.WinRates {
		if rate >= v[0] && rate < v[1] {
			// 找到原有题库
			for k, questions := range m.Question {
				if k == i {
					continue
				}
				for _, subQ := range questions {
					if subQ.Q.ID == q.Q.ID {
						m.Question[k] = m.delete(questions, subQ)   // 移除该题库
						m.Question[i] = append(m.Question[i], subQ) // 加入新题库
						break
					}
				}
			}
			return
		}
	}
}

func (m *QuestionLib) updateLib(q *Question) {
	total := q.win + q.fail
	if q.updateCount > 1 {
		// 统计大于100且每格100次更新一次题库
		if total >= RoomConf.QuestionCountMinLimit && total%RoomConf.QuestionCountDur == 0 {
			q.updateCount++
			// 计算题库正确率

			m.updateQuestion(q)
		}
	} else {
		if total >= 5 {
			q.updateCount = 2
			// 计算题库正确率
			m.updateQuestion(q)
		}
	}
}

func (m *QuestionLib) delete(a []*Question, elem *Question) []*Question {
	for i := 0; i < len(a); i++ {
		if a[i].Q.ID == elem.Q.ID {
			a = append(a[:i], a[i+1:]...)
			i--
			break
		}
	}
	return a
}

func (m *QuestionLib) ToExcel() {
	var str = [][]string{}
	title1 := []string{"ID", "question", "answerA", "answerB", "answerC", "answerD", "rightAnswer", "win", "fail"}
	str = append(str, title1)
	title2 := []string{"int", "string", "string", "string", "string", "string", "string", "int", "int"}
	str = append(str, title2)
	for _, v := range m.QuestionMap {
		subS := []string{fmt.Sprintf("%v", v.Q.ID), v.Q.Question, v.Q.AnswerA, v.Q.AnswerB, v.Q.AnswerC, v.Q.AnswerD, v.Q.RightAnswer, fmt.Sprintf("%v", v.win), fmt.Sprintf("%v", v.fail)}
		str = append(str, subS)
	}
	excel.CreateXlS(str, "question1", []string{"ID", "问题", "选项A", "选项B", "选项C", "选项D", "正确答案", "正确次数", "错误次数"})
}
