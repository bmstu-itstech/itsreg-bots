package objects

type State int64

const EndState State = 0

type Node struct {
	Id       State
	Default  State
	IsSilent bool
	Buttons  []Button
}

func (n *Node) Process(msg string) State {
	if n.IsSilent {
		return n.Default
	}

	for _, btn := range n.Buttons {
		if btn.Match(msg) {
			return btn.Next
		}
	}

	return n.Default
}

func (n *Node) IsLast() bool {
	return n.Id == EndState
}
