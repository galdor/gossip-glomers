package gg

import maelstrom "github.com/jepsen-io/maelstrom/demo/go"

type MessageHandler func(*Node, *Message) error

type Node struct {
	node *maelstrom.Node
}

func NewNode() *Node {
	n := Node{
		node: maelstrom.NewNode(),
	}

	return &n
}

func (n *Node) Run() error {
	return n.node.Run()
}

func (n *Node) AddHandler(msgType string, h MessageHandler) {
	n.node.Handle(msgType, n.wrapHandler(h))
}

func (n *Node) wrapHandler(h MessageHandler) maelstrom.HandlerFunc {
	return func(mmsg maelstrom.Message) error {
		msg, err := DecodeMessage(mmsg)
		if err != nil {
			return err
		}

		return h(n, msg)
	}

}

func (n *Node) Reply(msg *Message) error {
	return n.node.Send(msg.Source, msg.Body)
}
