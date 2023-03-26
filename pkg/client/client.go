package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"

	cm "github.com/deemerby/gopow/pkg/communication"
	"github.com/deemerby/gopow/pkg/pow"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// HandleRequest - handle server response
func HandleResponse(ctx context.Context, logger *logrus.Logger, conn net.Conn) (string, error) {
	logger.Infof("Client is able to connect to server: %s", conn.RemoteAddr())
	defer conn.Close()

	// send request to get challenge
	response := &cm.Message{Type: cm.MsgRequest}
	if err := cm.SendMsg(response, conn); err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	reader := bufio.NewReader(conn)
	res, err := reader.ReadBytes(cm.ByteDelim)
	if err != nil {
		return "", fmt.Errorf("failed to read response error: %v", err)
	}

	// unmarshal challenge response
	err = json.Unmarshal(res, response)
	if err != nil {
		return "", err
	}
	if response.Type != cm.MsgChallenge {
		return "", fmt.Errorf("got wrong message type: %v", err)
	}
	logger.Debugln("process challenge response")

	// calculate result
	logger.Debugln("calculate result")
	hashcashRes := &pow.HashcashData{}
	if err = json.Unmarshal([]byte(response.Payload), hashcashRes); err != nil {
		return "", err
	}
	logger.Debugf("hashcash response: %+v", hashcashRes)

	hashcashResult, err := hashcashRes.CalculateHashcash(viper.GetInt("hashcash.max.iteration"))
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
	resQuote, err := reader.ReadBytes(cm.ByteDelim)
	if err != nil {
		return "", fmt.Errorf("failed to read response error: %v", err)
	}

	err = json.Unmarshal(resQuote, response)
	if err != nil {
		return "", err
	}
	if response.Type != cm.MsgQuote {
		return "", fmt.Errorf("got wrong message type: %v", err)
	}
	logger.Debugln("process quote response")

	return response.Payload, nil
}
