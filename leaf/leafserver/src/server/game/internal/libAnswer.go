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
	Q    *pmsg.Question
	win  int
	fail int
}

type QuestionLib struct {
	QuestionMap      map[int]*Question
	Question         map[int][]*Question // 段位题库
	PhaseQuestionLib map[int][]int       // 各个段位对于的题库
	WinRates         map[int][]float32   // 各个段位对应胜率
	WinChan          chan int            // 统计管道
	FailChan         chan int            // 统计管道
	Done             chan interface{}
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

func (l *LibAnswer) ToString() string {
	var str = []string{}
	for i, v := range l.Answers {
		str = append(str, fmt.Sprintf("\n问题%d: %s------答案：%s", i+1, v.Question, v.RightAnswer))
	}
	return fmt.Sprintln(str)
}

func (l *LibAnswer) SingleToString() string {
	i := l.Progress
	v := l.Answers[l.Progress]
	return fmt.Sprintf("\n问题%d: %s------答案：%s", i+1, v.Question, v.RightAnswer)
}

func GetAnswerLib() Answers {
	ran := rand.Intn(len(AnswerLibs))
	return AnswerLibs[ran]
}

func ToAnswerLib(table string) Answers {
	tableConf, ok := manager.TableManager.GetTable(table)
	arr := []*pmsg.Question{}
	if ok {
		for _, cell := range tableConf.Excel {
			obj := &pmsg.Question{}
			obj.Table = table
			for index, v := range tableConf.AttributeNames {
				switch v {
				case "ID":
					obj.ID = int32(cell[index].(int))
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
				}
			}
			arr = append(arr, obj)
			q := &Question{
				Q: obj,
			}
			Questions.QuestionMap[int(obj.ID)] = q
			subA, ok := Questions.Question[1]
			if !ok {
				subA = make([]*Question, 0)
				subA = append(subA, q)
				Questions.Question[1] = subA
			} else {
				Questions.Question[1] = append(Questions.Question[1], q)
			}
		}
	}
	return arr
}

func RandAnswerLib(count int, arr []*pmsg.Question) *LibAnswer {
	// 获得其实随机数
	startIndex := 0
	length := len(arr)
	if len(arr)-count > 0 {
		startIndex = rand.Intn(length - count)
	}

	// 长度不够则补充
	if count >= length {
		for i := 0; i < count; i++ {
			arr = append(arr, arr[0])
		}
	}

	// 返回子数组
	return &LibAnswer{
		Answers:  arr[startIndex : startIndex+count],
		Progress: 0,
	}
}

func (m *QuestionLib) RandAnswerLib(id, count int) *LibAnswer {
	l := GetIDByLevel(id)
	log.Debug("************************输入段位:%v|%v", id, l)
	var lib = m.GetQuestions(l, count)
	return RandAnswerLib(count, lib)
}

func (m *QuestionLib) _getQuestions(id int) []*pmsg.Question {
	arr := []*pmsg.Question{}
	ids, ok := m.PhaseQuestionLib[id]
	if ok {
		for _, phase := range ids {
			for _, v := range m.Question[phase] {
				arr = append(arr, v.Q)
			}
		}
	}

	return arr
}

// 根据段位获取题库
func (m *QuestionLib) GetQuestions(id, total int) []*pmsg.Question {
	ARR := []*pmsg.Question{}
	count := id
	c := 0
	for {
		if len(ARR) <= total {
			ARR = append(ARR, m._getQuestions(count)...)
			count--
			if count < 0 {
				count = 12
			}
			c++
			if c >= 12 {
				break
			}
		} else {
			break
		}
	}
	return ARR
}

func (m *QuestionLib) WinCount(id int) {
	q, ok := m.QuestionMap[id]
	if ok {
		q.win++
		m.updateLib(q)
	}
}

func (m *QuestionLib) FailCount(id int) {
	q, ok := m.QuestionMap[id]
	if ok {
		q.fail++
		m.updateLib(q)
	}
}

func (m *QuestionLib) updateLib(q *Question) {
	total := q.win + q.fail
	// 统计大于100且每格100次更新一次题库
	if total >= RoomConf.QuestionCountMinLimit && total%RoomConf.QuestionCountDur == 0 {
		// 计算题库正确率
		rate := float32(q.win) * 100 / (float32(q.win) + float32(q.fail))
		for i, v := range m.WinRates {
			if rate >= v[0] && rate < v[1] {
				// 找到原有题库
				for k, questions := range m.Question {
					for _, subQ := range questions {
						if subQ.Q.ID == int32(q.Q.ID) {
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
