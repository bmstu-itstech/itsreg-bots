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
) ([]Message, error) {
	messages := make([]Message, 0)

	if !prt.IsProcessing() {
		return messages, nil
	}

	current := b.blocks[prt.State]

	if current.Type != MessageBlock {
		err := prt.AddAnswer(text)
		if err != nil {
			return nil, err
		}
	}

	nextState := current.Process(text)
	prt.SwitchTo(nextState)

	next, ok := b.blocks[nextState]
	if ok {
		message, err := next.Message()
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if next.Type == MessageBlock {
		ms, err := b.Process(prt, "")
		if err != nil {
			return nil, err
		}
		messages = append(messages, ms...)
	}

	return messages, nil
}

func (b *Bot) processStart(
	prt *Participant,
) ([]Message, error) {
	messages := make([]Message, 0, 1)

	start := b.blocks[prt.State]

	msg, err := start.Message()
	if err != nil {
		return nil, err
	}
	messages = append(messages, msg)

	if start.Type == MessageBlock {
		ms, err := b.Process(prt, "")
		if err != nil {
			return nil, err
		}
		messages = append(messages, ms...)
	}

	return messages, nil
}
