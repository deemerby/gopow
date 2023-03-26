package communication

import (
	"encoding/json"
	"net"
)

// Type of message
const (
	MsgRequest = iota
	MsgChallenge
	MsgProofOfWork
	MsgQuote
)

const ByteDelim = byte(0x0A)

// Message
type Message struct {
	Type    int    // type of message
	Payload string // payload
}

func SendMsg(msg *Message, conn net.Conn) error {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	msgByte = append(msgByte, ByteDelim)
	_, err = conn.Write(msgByte)
	return err
}