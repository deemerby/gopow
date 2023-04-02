package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/deemerby/gopow/pkg/client"
	"github.com/deemerby/gopow/pkg/options"
	"github.com/deemerby/gopow/pkg/server"
	st "github.com/deemerby/gopow/pkg/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var version = "unset"

func main() {
	doneC := make(chan error)
	logger := NewLogger()
	globalCtx, globalCancel := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	logger.Infof("Version: %v", version)

	opt := &options.AppOptions{
		ZeroCnt:      viper.GetInt("hashcash.zero.cnt"),
		MaxIteration: viper.GetInt("hashcash.max.iteration"),
		MaxNumber:    viper.GetInt64("max.random.number"),
		ConDeadline:  viper.GetDuration("connection.deadline"),
	}

	if viper.GetBool("type.server") {
		logger.Info("Server")
		go func() { doneC <- RunServer(globalCtx, logger, opt) }()
	} else {
		logger.Info("Client")
		go func() { doneC <- RunClient(globalCtx, logger, opt) }()
	}

	select {
	case err := <-doneC:
		globalCancel()
		if err != nil {
			logger.Fatal(err)
		}

		os.Exit(0)
	case err := <-sigc:
		globalCancel()
		logger.Fatal(err)
	}
}

func NewLogger() *logrus.Logger {
	logger := logrus.StandardLogger()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	// Set the log level on the default logger based on command line flag
	if level, err := logrus.ParseLevel(viper.GetString("logging.level")); err != nil {
		logger.Errorf("Invalid %q provided for log level", viper.GetString("logging.level"))
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	return logger
}

func init() {
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatalf("cannot load configuration: %v", err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(viper.GetString("config.source"))
	if viper.GetString("config.file") != "" {
		log.Printf("Serving from configuration file: %s", viper.GetString("config.file"))
		viper.SetConfigName(viper.GetString("config.file"))
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("cannot load configuration: %v", err)
		}
	} else {
		log.Printf("Serving from default values, environment variables, and/or flags")
	}
}

// RunServer ...
func RunServer(ctx context.Context, logger *logrus.Logger, opt *options.AppOptions) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(viper.GetString("server.address"), viper.GetString("server.port")))
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()

		if err := listener.Close(); err != nil {
			logger.Errorf("Failed to close listener: %v", err)
		}
	}()

	logger.Infof("Server is listening... port: %s", viper.GetString("server.port"))
	storage := st.NewMemoryStore(ctx, viper.GetDuration("hashcash.duration"))
	ser := server.NewHandler(logger, storage)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		_ = conn.SetDeadline(time.Now().Add(opt.ConDeadline))
		defer conn.Close()

		go func() {
			ser.HandleRequest(ctx, conn, opt)
		}()
	}
}

// RunClient ...
func RunClient(ctx context.Context, logger *logrus.Logger, opt *options.AppOptions) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(viper.GetString("server.address"), viper.GetString("server.port")))
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()

		if err := conn.Close(); err != nil {
			logger.Errorf("Failed to close connection: %v", err)
		}
	}()

	msg, err := client.HandleResponse(logger, conn, opt)
	if err != nil {
		return err
	}

	logger.Infof("Result of working: %s", msg)
	return nil
}
