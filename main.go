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

var listen []string
var useSyslog bool
var facility string
var useSameLog bool
var logFile string
var appendLogs bool
var useConsole bool
var excludeStatus []uint
var excludeInternalConn bool
var excludeNewConn bool
var excludeCloseConn bool
var excludeAuth bool

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

	srv := server.NewServer(internalLogger, logger)
	for _, l := range listen {
		err = srv.Listen(l)
		if err != nil {
			internalLogger.Fatal().Str("listen", l).Err(err).Msg("unable to listen")
		}
	}

	addFilters(srv, internalLogger)

	srv.Start()

	// Wait for a SIGINT or SIGTERM signal to gracefully shut down the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	internalLogger.Info().Msg("Shutting down server...")
	srv.Stop()
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

// Add any required filters
func addFilters(srv *server.Server, internalLogger zerolog.Logger) {

	filters := []server.Filter{}

	if excludeAuth {
		filters = append(filters, server.ExcludeAuthConnection())
	}

	if excludeNewConn {
		filters = append(filters, server.ExcludeNewConnection())
	}

	if excludeInternalConn {
		filters = append(filters, server.ExcludeInternalConnection())
	}

	if excludeCloseConn {
		filters = append(filters, server.ExcludeCloseConnection())
	}

	if len(excludeStatus) > 0 {
		s := []server.AuthStatus{}

		for _, n := range excludeStatus {
			if n < uint(server.AUTH_STATUS_MIN) || n > uint(server.AUTH_STATUS_MAX) {
				internalLogger.Fatal().Uint("status", n).Msg("invalid authentication status")
			}
			s = append(s, server.AuthStatus(n))
		}

		filters = append(filters, server.ExcludeAuthConnectionStatus(s...))
	}

	srv.AddFilters(filters...)

}

func init() {
	pflag.StringSliceVarP(&listen, "listen", "L", []string{"127.0.0.1:29001"}, "address:port (s) to listen on")
	pflag.BoolVarP(&useSyslog, "syslog", "s", false, "send logs to syslog")
	pflag.StringVarP(&facility, "facility", "F", "LOCAL0", "syslog facility to use")
	pflag.BoolVarP(&useSameLog, "same-log", "l", false, "use the same log for both internal and external messages")
	pflag.StringVarP(&logFile, "logfile", "f", "", "log file to write to (default to stderr)")
	pflag.BoolVarP(&appendLogs, "append", "A", false, "append to log file instead of overwriting")
	pflag.BoolVarP(&useConsole, "console", "c", false, "use console output as well as defined logger ")
	pflag.UintSliceVar(&excludeStatus, "exclude-auth-status", []uint{}, "exclude authentication status by code")
	pflag.BoolVar(&excludeNewConn, "exclude-new-connection", false, "exclude new connections")
	pflag.BoolVar(&excludeInternalConn, "exclude-internal-connection", false, "exclude internal connections")
	pflag.BoolVar(&excludeCloseConn, "exclude-close-connection", false, "exclude closing connections")
	pflag.BoolVar(&excludeAuth, "exclude-auth-connection", false, "exclude authentication connections")
}
