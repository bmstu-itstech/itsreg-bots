package bots

func (b Block) Process(text string) int {
	for _, opt := range b.Options {
		if opt.Match(text) {
			return opt.Next
		}
	}

	return b.NextState
}

func (b *Bot) Process(
	prt *Participant,
	text string,
) ([]Block, *Answer, error) {
	if prt.IsFinished() {
		return make([]Block, 0), nil, nil
	}

	processed := make([]Block, 0, 1)

	current := b.Blocks[prt.State]

	nextState := current.Process(text)
	prt.SwitchTo(nextState)

	next := b.Blocks[nextState]
	if !next.IsZero() {
		processed = append(processed, next)
	}

	var ans *Answer
	if current.Type != MessageBlock {
		var err error
		ans, err = NewAnswer(prt.UserID, current.State, text)
		if err != nil {
			return nil, nil, err
		}
	}

	if next.Type == MessageBlock {
		bs, _, err := b.Process(prt, "")
		if err != nil {
			return nil, nil, err
		}
		processed = append(processed, bs...)
		return processed, ans, nil
	}

	return processed, ans, nil
}
