package main

import (
	"net"

	"github.com/deemerby/gopow/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

func RunClient(ctx context.Context, logger *logrus.Logger) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(viper.GetString("server.address"), viper.GetString("server.port")))
	if err != nil {
		return err
	}

	msg, err := client.HandleResponse(ctx, logger, conn)
	if err != nil {
		return err
	}

	logger.Infof("Result of working: %s", msg)
	return nil
}
