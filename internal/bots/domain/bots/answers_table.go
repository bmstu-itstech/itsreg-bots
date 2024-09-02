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

	head := make([]string, len(blocks)+1)
	m := make(mapStateToIndex)

	head[0] = userIDColumnName

	for i, block := range blocks {
		head[i+1] = block.Title
		m[block.State] = i
	}

	return head, m
}

func trow(prt *Participant, m mapStateToIndex) []string {
	row := make([]string, len(m)+1)
	row[0] = strconv.FormatInt(prt.UserID, 10)
	for _, ans := range prt.Answers() {
		i := m[ans.State]
		row[i+1] = ans.Text
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
