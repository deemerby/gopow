package main

import (
	"github.com/spf13/pflag"
)

const (
	// configuration defaults support local development (i.e. "go run ...")

	// Server
	defaultServerAddress = "0.0.0.0"
	defaultServerPort    = "8084"

	defaulTypeServer = true

	defaultLoggingLevel = "info"

	defaultHashcashZeroCnt       = 6
	defaultHashcashDuration      = 20
	defaultHashcashMaxIterations = 1000000000

	defaultRand = 100000
)

var (
	// define flag overrides
	_ = pflag.String("server.address", defaultServerAddress, "address of tcp server")
	_ = pflag.String("server.port", defaultServerPort, "port of tcp server")

	_ = pflag.String("logging.level", defaultLoggingLevel, "log level of application")

	_ = pflag.Bool("type.server", defaulTypeServer, "type of application")

	_ = pflag.Int("hashcash.zero.cnt", defaultHashcashZeroCnt, "hashcash zero count")
	_ = pflag.Int("hashcash.duration", defaultHashcashDuration, "hashcash duration")
	_ = pflag.Int("hashcash.max.iteration", defaultHashcashMaxIterations, "hashcash maximum iteration")

	_ = pflag.Int("max.random.number", defaultRand, "maximum random number")
)
