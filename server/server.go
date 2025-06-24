package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Build and run the server

type Server struct {
	listeners      []net.Listener
	logger         zerolog.Logger
	filters        []Filter
	internalLogger zerolog.Logger
	wg             sync.WaitGroup
	done           chan bool
}

type Filter func(*Server, map[string]any) bool

// NewServer returns a server object ready to be configured.
func NewServer(internalLogger, mainLogger zerolog.Logger) *Server {

	server := &Server{
		logger:         mainLogger,
		internalLogger: internalLogger,
		done:           make(chan bool),
		filters:        []Filter{},
	}

	return server

}

// Listen adds an address and port for the server to listen on.
func (server *Server) Listen(addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	server.done = make(chan bool)
	server.listeners = append(server.listeners, listener)
	return nil
}

// AddFilters adds one or more filter functions  to the server, allowing some messages to be ignored.
func (server *Server) AddFilters(filters ...Filter) {
	server.filters = append(server.filters, filters...)
}

// Start begins the async listen/accept loop for the server, creating a
// a goroutine for each listener.
func (server *Server) Start() error {

	for _, listener := range server.listeners {
		server.goAcceptConnection(listener)
	}

	return nil
}

// Stop closes the server cleanly.
func (server *Server) Stop() {
	close(server.done)

	done := make(chan struct{})
	go func() {
		server.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		server.internalLogger.Error().Msg("Timed out waiting for connections to finish.")
		return
	}
}

// goAcceptConnection runs the connection loop for a single listener
func (server *Server) goAcceptConnection(listener net.Listener) {
	server.wg.Add(1)
	go func(listener net.Listener) {
	loop:
		for {
			select {
			case <-server.done:
				break loop
			default:
			}
			connection, err := listener.Accept()
			if err != nil {
				continue
			}

			server.readConnection(connection)
		}

		server.wg.Done()
	}(listener)
}

// readConnction reads a message on a connection and then writes it to the
// main log (after filtering).
func (server *Server) readConnection(connection net.Conn) {

	defer connection.Close()

	var remoteAddr string
	remote := connection.RemoteAddr()
	if remote == nil {
		remoteAddr = "<unknown>"
	} else {
		remoteAddr = remote.String()
	}

	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, connection)
	if err != nil {
		server.internalLogger.Error().Err(err).Str("remote", remoteAddr).Msg("Error reading from connection")
		return
	}

	server.writeLog(remoteAddr, *buf)
}

func (server *Server) writeLog(remote string, buf bytes.Buffer) {

	// if we have no filters, don't bother to parse the input.
	if len(server.filters) == 0 {
		server.logger.Info().Str("client", remote).RawJSON("message", buf.Bytes()).Send()
	} else {
		message := map[string]any{}
		err := json.Unmarshal(buf.Bytes(), &message)
		if err != nil {
			server.internalLogger.Error().Err(err).Str("client", remote).Bytes("message", buf.Bytes()).Msg("Error unmarshalling log message")
			return
		} else {
			for _, filter := range server.filters {
				if !filter(server, message) {
					server.internalLogger.Debug().Str("client", remote).Bytes("message", buf.Bytes()).Msg("Log message filtered out")
					return
				}
			}
			server.logger.Info().Str("client", remote).RawJSON("message", buf.Bytes()).Send()
		}
	}

}
