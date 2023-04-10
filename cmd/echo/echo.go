package main

import (
	"github.com/galdor/go-program"
	"github.com/galdor/gossip-glomers/pkg/gg"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	p := program.NewProgram("echo",
		"a maelstrom node sending a simple response to echo messages")

	p.SetMain(mainCmd)

	p.ParseCommandLine()
	p.Run()
}

func handler(node *maelstrom.Node, fn func(*maelstrom.Node, maelstrom.Message) error) maelstrom.HandlerFunc {
	return func(msg maelstrom.Message) error {
		return fn(node, msg)
	}
}

func mainCmd(p *program.Program) {
	node := gg.NewNode()
	node.AddHandler("echo", hEcho)

	if err := node.Run(); err != nil {
		p.Fatal("cannot run node: %v", err)
	}
}

func hEcho(node *gg.Node, msg *gg.Message) error {
	echoMsg := msg.Body.(*gg.EchoMessage)

	echoMsg.Type = "echo_ok"
	echoMsg.InReplyTo = echoMsg.MsgId

	return node.Reply(msg)
}
