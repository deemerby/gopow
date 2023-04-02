package client

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	cm "github.com/deemerby/gopow/pkg/communication"
	"github.com/deemerby/gopow/pkg/options"
	"github.com/deemerby/gopow/pkg/pow"
)

// HandleRequest - handle server response
func HandleResponse(logger *logrus.Logger, conn net.Conn, opt *options.AppOptions) (string, error) {
	logger.Infof("Client is able to connect to server: %s", conn.RemoteAddr())
	defer conn.Close()

	// send request to get challenge
	msg := &cm.Message{Type: cm.MsgRequest}
	if err := cm.SendMsg(msg, conn); err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	d := json.NewDecoder(conn)
	if err := d.Decode(msg); err != nil {
		return "", fmt.Errorf("failed to read response error: %v", err)
	}

	if msg.Type != cm.MsgChallenge {
		return "", fmt.Errorf("got wrong message type: %d", msg.Type)
	}
	logger.Debugln("process challenge response")

	// calculate result
	logger.Debugln("calculate result")
	hashcashRes := &pow.HashcashData{}
	if err := json.Unmarshal([]byte(msg.Payload), hashcashRes); err != nil {
		return "", fmt.Errorf("failed to unmarshal payload error: %v", err)
	}
	logger.Debugf("hashcash response: %+v", hashcashRes)

	hashcashResult, err := hashcashRes.CalculateHashcash(opt.MaxIteration)
	if err != nil {
		return "", err
	}
	logger.Debugf("hashcash calculated: %+v", hashcashResult)

	hashcashMarshaled, err := json.Marshal(hashcashResult)
	if err != nil {
		return "", fmt.Errorf("failed to marshal hashcash: %v", err)
	}
	resultReq := &cm.Message{
		Type:    cm.MsgProofOfWork,
		Payload: string(hashcashMarshaled),
	}

	// send request with result
	if err := cm.SendMsg(resultReq, conn); err != nil {
		return "", fmt.Errorf("failed to send pow result: %v", err)
	}

	// process quote
	if err := d.Decode(msg); err != nil {
		return "", fmt.Errorf("failed to read response error: %v", err)
	}

	if msg.Type != cm.MsgQuote {
		return "", fmt.Errorf("got wrong message type: %d", msg.Type)
	}
	logger.Debugln("process quote response")

	return msg.Payload, nil
}
