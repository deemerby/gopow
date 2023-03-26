package main

import (
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	"github.com/deemerby/gopow/pkg/server"
	st "github.com/deemerby/gopow/pkg/storage"
)

func RunServer(ctx context.Context, logger *logrus.Logger) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(viper.GetString("server.address"), viper.GetString("server.port")))
	if err != nil {
		return err
	}
	defer listener.Close()

	logger.Infof("Server is listening... port: %s", viper.GetString("server.port"))
	storage := st.NewMemoryStore()
	ser := server.NewHandler(logger, storage)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()

		go func() {
			ser.HandleRequest(ctx, conn)
		}()
	}
}
