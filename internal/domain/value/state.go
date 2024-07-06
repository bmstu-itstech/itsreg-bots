package value

type State uint64

const StateNone State = 0

func (s State) IsNone() bool {
	return s == StateNone
}
