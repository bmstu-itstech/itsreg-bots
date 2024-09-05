package bots

import "fmt"

type EntryNotFoundError struct {
	Key string
}

func (e EntryNotFoundError) Error() string {
	return fmt.Sprintf("entry '%s' not found", e.Key)
}

func (b *Bot) Entry(prt *Participant, key string) ([]Message, error) {
	e, ok := b.entryPoints[key]
	if !ok {
		return nil, EntryNotFoundError{Key: key}
	}

	b.cleanAllAnswersFrom(e.State, prt)
	prt.SwitchTo(e.State)

	response := make([]Message, 0, 1)
	response = append(response)

	ms, err := b.processStart(prt)
	if err != nil {
		return nil, err
	}

	response = append(response, ms...)

	return response, nil
}

func (b *Bot) cleanAllAnswersFrom(start int, prt *Participant) {
	blocks := b.Traverse(start)
	for _, block := range blocks {
		prt.CleanAnswerIfExists(block.State)
	}
}
