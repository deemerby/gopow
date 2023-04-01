package communication

import (
	"encoding/json"
	"io"
)

// Type of message
const (
	MsgRequest = iota
	MsgChallenge
	MsgProofOfWork
	MsgQuote
)

// Message
type Message struct {
	Type    int    // type of message
	Payload string // payload
}

func SendMsg(msg *Message, conn io.Writer) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = conn.Write(msgByte)
	return err
}
