package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Build and run the server

type Server struct {
	port           int
	address        string
	logger         zerolog.Logger
	internalLogger zerolog.Logger
	listener       net.Listener
	wg             sync.WaitGroup
	shutdown       chan struct{}
	connection     chan net.Conn
}

func NewServer(address string, port int, internalLogger, mainLogger zerolog.Logger) (*Server, error) {

	server := &Server{
		port:           port,
		address:        address,
		logger:         mainLogger,
		internalLogger: internalLogger,
		shutdown:       make(chan struct{}),
		connection:     make(chan net.Conn),
	}

	err := server.buildListener()
	if err != nil {
		return nil, err
	}

	return server, nil

}

func (server *Server) buildListener() error {

	addressPort := fmt.Sprintf("%s:%d", server.address, server.port)
	listener, err := net.Listen("tcp", addressPort)
	if err != nil {
		server.internalLogger.Error().Err(err).Str("address", address).Msg("failed to create listener")
		return err
	}

	server.listener = listener

	return nil
}

func (server *Server) acceptConnections() {
	defer server.wg.Done()

	for {
		select {
		case <-server.shutdown:
			return
		default:
			conn, err := server.listener.Accept()
			if err != nil {
				continue
			}
			server.connection <- conn
		}
	}
}

func (server *Server) handleConnections() {
	defer server.wg.Done()

	for {
		select {
		case <-server.shutdown:
			return
		case conn := <-server.connection:
			go server.handleConnection(conn)
		}
	}
}

func (server *Server) handleConnection(conn net.Conn) {

	defer conn.Close()
	server.internalLogger.Debug().Str("address", conn.RemoteAddr().String()).Msg("connection accepted")

	var buffer bytes.Buffer
	n, err := io.Copy(&buffer, conn)

	if err != nil {
		server.internalLogger.Error().Err(err).Str("address", conn.RemoteAddr().String()).Msg("error reading from connection")
		return
	}

	fmt.Printf("bytes = %d, buffer = %s\n", n, buffer.String())

}

func (server *Server) Start() {
	server.wg.Add(2)
	go server.acceptConnections()
	go server.handleConnections()
}

func (server *Server) Stop() {
	close(server.shutdown)
	server.listener.Close()

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
