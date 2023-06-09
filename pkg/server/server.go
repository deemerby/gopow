package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"

	"github.com/sirupsen/logrus"

	cm "github.com/deemerby/gopow/pkg/communication"
	"github.com/deemerby/gopow/pkg/options"
	"github.com/deemerby/gopow/pkg/pow"
	st "github.com/deemerby/gopow/pkg/storage"
)

// Quotes - https://www.azquotes.com/quotes/topics/wisdom.html
var Quotes = []string{
	"To acquire knowledge, one must study; but to acquire wisdom, one must observe.",
	"Yesterday I was clever, so I wanted to change the world. Today I am wise, so I am changing myself.",
	"A man only becomes wise when he begins to calculate the approximate depth of his ignorance.",
	"Knowing others is intelligence; knowing yourself is true wisdom. Mastering others is strength; mastering yourself is true power.",
	"Don’t depend too much on anyone in this world because even your own shadow leaves you when you are in darkness.",
}

type ServerHandler struct {
	log     *logrus.Logger
	storage *st.MemoryStore
}

// HandleRequest - handle client request
func NewHandler(logger *logrus.Logger, storage *st.MemoryStore) *ServerHandler {
	return &ServerHandler{
		log:     logger,
		storage: storage,
	}
}

// HandleRequest - handle client request
func (h *ServerHandler) HandleRequest(ctx context.Context, conn net.Conn, opt *options.AppOptions) {
	h.log.Infof("New client: %s", conn.RemoteAddr())
	defer conn.Close()

	d := json.NewDecoder(conn)
	msg := &cm.Message{}

	for {
		err := d.Decode(msg)
		if err != nil {
			if err == io.EOF {
				return
			}
			h.log.Infof("Failed to read request: %v", err)
			return
		}

		res, err := h.processRequest(msg, conn, opt)
		if err != nil {
			h.log.Errorf("Failed to process request error: %v", err)
			return
		}
		if res != nil {
			err := cm.SendMsg(res, conn)
			if err != nil {
				h.log.Errorf("Failed to send message: %v", err)
			}
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// processRequest - Processing request from client
func (h *ServerHandler) processRequest(msgReq *cm.Message, conn net.Conn, opt *options.AppOptions) (*cm.Message, error) {
	clName := conn.RemoteAddr().String()

	// check type of message
	switch msgReq.Type {
	case cm.MsgRequest:
		h.log.Infof("Client: %s requests challenge", clName)

		nBig, err := rand.Int(rand.Reader, big.NewInt(opt.MaxNumber))
		if err != nil {
			return nil, fmt.Errorf("failed to get random value: %v", err)
		}
		rand := nBig.Int64()
		h.log.Infof("Add Rand: %d", rand)
		if err := h.storage.Add(int(rand)); err != nil {
			return nil, fmt.Errorf("failed to add client date to storage: %v", err)
		}

		hashcash, err := pow.NewHashcash(clName, int(rand), opt.ZeroCnt)
		if err != nil {
			return nil, fmt.Errorf("failed to get new hashcash: %v", err)
		}
		hashcashMarshaled, err := json.Marshal(hashcash)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal hashcash: %v", err)
		}
		msgRes := &cm.Message{
			Type:    cm.MsgChallenge,
			Payload: string(hashcashMarshaled),
		}

		return msgRes, nil
	case cm.MsgProofOfWork:
		h.log.Infof("Client: %s requests result", clName)

		hashcash := &pow.HashcashData{}
		if err := json.Unmarshal([]byte(msgReq.Payload), hashcash); err != nil {
			return nil, err
		}

		var randV int
		var err error
		if randV, err = hashcash.GetRand(); err != nil {
			return nil, err
		}

		h.log.Debugf("Get rand value from storage: %d", randV)
		if err := h.storage.Get(randV); err != nil {
			h.storage.Delete(randV)
			return nil, fmt.Errorf("failed to get client date from storage: %v", err)
		}

		if _, err := hashcash.CalculateHashcash(hashcash.Counter); err != nil {
			return nil, fmt.Errorf("failed to check hashcash: %v", err)
		}

		// get random quote
		h.log.Infof("Client: %s work confirmed", clName)
		nBig, err := rand.Int(rand.Reader, big.NewInt(4))
		if err != nil {
			panic(err)
		}
		rand := nBig.Int64()
		msgResult := &cm.Message{
			Type:    cm.MsgQuote,
			Payload: Quotes[int(rand)],
		}

		h.storage.Delete(randV)

		return msgResult, nil
	default:
		return nil, fmt.Errorf("unknown type of message")
	}
}
