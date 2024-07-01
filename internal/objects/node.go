package objects

type NodeId int64

const NodeNodeId NodeId = 0

type Node struct {
	Id       NodeId
	Default  NodeId
	IsSilent bool
	Buttons  []Button
}

func (n *Node) Process(msg string) NodeId {
	if n.IsSilent {
		return n.Default
	}

	for _, btn := range n.Buttons {
		if btn.Match(msg) {
			return btn.NextId
		}
	}

	return n.Default
}

func (n *Node) IsLast() bool {
	return n.Id == NodeNodeId
}
