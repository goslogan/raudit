package main

// Application to receive Redis Enterprise Software audit messages and write them to a log

import (
	"fmt"
	"io"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"

	"github.com/goslogan/raudit/server"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

var port uint16
var address string
var useSyslog bool
var facility string
var exclude bool
var useSameLog bool
var logFile string
var appendLogs bool
var useConsole bool

func main() {
	var err error
	var logger, internalLogger zerolog.Logger
	pflag.Parse()
	logger, err = initialiseMainLogger()

	if useSameLog {
		internalLogger = logger
	} else {
		internalLogger, _ = initialiseInternalLogger()
	}

	if err != nil {
		internalLogger.Fatal().Err(err).Msg("unable to initialise logger")
	}

	server, err := server.NewServer(internalLogger, logger)
	if err != nil {
		internalLogger.Fatal().Err(err).Msg("unable to build server")
	}
	err = server.Listen(fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		internalLogger.Fatal().Uint16("port", port).Str("address", address).Err(err).Msg("unable to listen")
	}

	server.Start()

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	internalLogger.Info().Msg("Shutting down server...")
	server.Stop()
	internalLogger.Info().Msg("Server stopped.")

}

// Build the logger. This may be writing to a  file or to syslog. Additionally, it can write the log to the console if
// configured to. If it fails, a logger writing to stderr is return with the error.
func initialiseMainLogger() (zerolog.Logger, error) {

	writers := []io.Writer{}

	if useSyslog {
		facility, err := getSyslogFacility(facility)
		if err != nil {
			return zerolog.New(os.Stderr), err
		}
		writer, err := syslog.New(facility|syslog.LOG_INFO, "redis-audit")
		if err != nil {
			return zerolog.New(os.Stderr), err
		}
		writers = append(writers, zerolog.SyslogLevelWriter(writer))

	} else {
		flags := os.O_WRONLY | os.O_CREATE
		var err error
		if appendLogs {
			flags |= os.O_APPEND
		}

		var file = os.Stderr

		if logFile != "" {
			file, err = os.OpenFile(logFile, flags, 0644)
			if err != nil {
				return zerolog.New(os.Stderr), err
			}
		}

		writers = append(writers, file)

		if useConsole {
			writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
		}

	}

	return zerolog.New(zerolog.MultiLevelWriter(writers...)), nil

}

// If useSameLog is  not set, the internal logger (used for errors within this package) will be created writing to STDERR
func initialiseInternalLogger() (zerolog.Logger, error) {
	return zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger(), nil
}

func getSyslogFacility(facility string) (syslog.Priority, error) {
	// Convert the facility string to syslog.Priority
	if priority, ok := Facilities[facility]; ok {
		return priority, nil
	}
	return syslog.LOG_LOCAL0, fmt.Errorf("invalid syslog facility: %s", facility)
}

func init() {
	pflag.Uint16VarP(&port, "port", "p", 29001, "port to listen on")
	pflag.StringVarP(&address, "address", "a", "127.0.0.1", "address to listen on")
	pflag.BoolVarP(&useSyslog, "syslog", "s", false, "send logs to syslog")
	pflag.StringVarP(&facility, "facility", "F", "LOCAL0", "syslog facility to use")
	pflag.BoolVarP(&exclude, "exclude-internal", "x", false, "exclude internal messages")
	pflag.BoolVarP(&useSameLog, "same-log", "l", false, "use the same log for both internal and external messages")
	pflag.StringVarP(&logFile, "logfile", "f", "", "log file to write to (default to stderr)")
	pflag.BoolVarP(&appendLogs, "append", "A", false, "append to log file instead of overwriting")
	pflag.BoolVarP(&useConsole, "console", "c", false, "use console output as well as defined logger ")
}
