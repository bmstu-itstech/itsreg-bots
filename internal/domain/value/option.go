package value

type Option struct {
	Text string
	Next State
}

func (o Option) Match(s string) bool {
	return s == o.Text
}
