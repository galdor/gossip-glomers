package gg

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type MessageHeader struct {
	Type      string `json:"type,omitempty"`
	MsgId     int    `json:"msg_id,omitempty"`
	InReplyTo int    `json:"in_reply_to,omitempty"`
	Code      int    `json:"code,omitempty"`
	Text      string `json:"text,omitempty"`
}

type Message struct {
	Source    string
	Recipient string
	Body      MessageBody
}

type MessageBody interface {
	fmt.Stringer
}

type EchoMessage struct {
	MessageHeader
	Echo string `json:"echo"`
}

func (msg *EchoMessage) String() string {
	return fmt.Sprintf("EchoMessage{Echo: %q}", msg.Echo)
}

func DecodeMessage(mmsg maelstrom.Message) (*Message, error) {
	// We are forced to decode the body twice because message-specific members
	// are at the same level as generic members. Nothing we can do about it
	// without changing the protocol itself.

	var header MessageHeader
	if err := json.Unmarshal(mmsg.Body, &header); err != nil {
		return nil, fmt.Errorf("cannot decode message header: %w", err)
	}

	msg := Message{
		Source:    mmsg.Src,
		Recipient: mmsg.Dest,
	}

	switch header.Type {
	case "echo":
		msg.Body = new(EchoMessage)

	default:
		return nil, fmt.Errorf("unknown message type %q", header.Type)
	}

	if err := json.Unmarshal(mmsg.Body, &msg.Body); err != nil {
		return nil, fmt.Errorf("cannot decode message body: %w", err)
	}

	return &msg, nil
}
