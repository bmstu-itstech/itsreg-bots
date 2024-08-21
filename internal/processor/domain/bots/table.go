package bots

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
)

type Table struct {
	Head []string
	Body [][]string
}

func NewTable(bot *Bot, answers []*Answer) (*Table, error) {
	mappedBlockIDs := createMapBlockIDs(bot)
	prtIDs := participantIDs(answers)
	mappedPrtIDs := createMapParticipantsIDs(prtIDs)

	tbody := make([][]string, len(prtIDs))
	for i, id := range prtIDs {
		tbody[i] = make([]string, len(mappedBlockIDs)+1)
		tbody[i][0] = strconv.FormatInt(id, 10)
	}

	for _, ans := range answers {
		i := mappedPrtIDs[ans.UserID]
		j, ok := mappedBlockIDs[ans.State]
		if !ok {
			return nil, fmt.Errorf("answer has non existent state %d", ans.State)
		}
		tbody[i][j] = ans.Text
	}

	thead := make([]string, len(mappedBlockIDs)+1)
	thead[0] = "UserID"
	for state, block := range bot.Blocks {
		i := mappedBlockIDs[state]
		thead[i] = block.Title
	}

	return &Table{
		Head: thead,
		Body: tbody,
	}, nil
}

func createMapBlockIDs(bot *Bot) map[int]int {
	states := make([]int, 0, len(bot.Blocks))
	for _, block := range bot.Blocks {
		states = append(states, block.State)
	}

	sort.Ints(states)

	i := 1
	mapped := make(map[int]int)
	for _, state := range states {
		mapped[state] = i
		i++
	}

	return mapped
}

func participantIDs(answers []*Answer) []int64 {
	arr := make([]int64, 0)
	for _, ans := range answers {
		if !slices.Contains(arr, ans.UserID) {
			arr = append(arr, ans.UserID)
		}
	}
	return arr
}

func createMapParticipantsIDs(prtIDs []int64) map[int64]int {
	mapped := make(map[int64]int)
	for i, id := range prtIDs {
		mapped[id] = i
	}
	return mapped
}
