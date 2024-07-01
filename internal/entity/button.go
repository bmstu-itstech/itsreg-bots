package entity

type Button struct {
	Text   string
	NextId NodeId
}

func (b *Button) Match(txt string) bool {
	return b.Text == txt
}
