package bots

import "strconv"

const (
	userIDColumnName = "UserID"
)

type AnswersTable struct {
	Head []string
	Body [][]string
}

func NewAnswersTable(bot *Bot, prts []*Participant) *AnswersTable {
	head, m := thead(bot)
	body := tbody(prts, m)
	return &AnswersTable{
		Head: head,
		Body: body,
	}
}

type mapStateToIndex map[int]int

func thead(bot *Bot) ([]string, mapStateToIndex) {
	blocks := bot.Blocks()

	head := make([]string, 0, len(blocks)+1)
	m := make(mapStateToIndex)

	head = append(head, userIDColumnName)

	i := 1
	for _, block := range blocks {
		if block.Type != MessageBlock {
			head = append(head, block.Title)
			m[block.State] = i
			i++
		}
	}

	return head, m
}

func trow(prt *Participant, m mapStateToIndex) []string {
	row := make([]string, len(m)+1)
	row[0] = strconv.FormatInt(prt.UserID, 10)
	for _, ans := range prt.Answers() {
		i, ok := m[ans.State]
		if ok {
			row[i] = ans.Text
		}
	}
	return row
}

func tbody(prts []*Participant, m mapStateToIndex) [][]string {
	body := make([][]string, len(prts))
	for i, prt := range prts {
		body[i] = trow(prt, m)
	}
	return body
}
