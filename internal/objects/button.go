package objects

type Button struct {
	Text string
	Next State
}

func (b *Button) Match(txt string) bool {
	return b.Text == txt
}
