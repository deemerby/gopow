package main

import (
	"context"
	"log"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var version = "unset"

func main() {
	doneC := make(chan error)
	logger := NewLogger()
	globalCtx, globalCancel := context.WithCancel(context.Background())
	_ = globalCtx
	logger.Infof("Version: %v", version)

	if viper.GetBool("type.server") {
		logger.Info("Server")
		go func() { doneC <- RunServer(globalCtx, logger) }()
	} else {
		logger.Info("Client")
		go func() { doneC <- RunClient(globalCtx, logger) }()
	}

	if err := <-doneC; err != nil {
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
