package main_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/goslogan/raudit/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestServerCreation(t *testing.T) {

	var ibuf, mbuf bytes.Buffer
	internalLogger := zerolog.New(&ibuf).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	mainLogger := zerolog.New(&mbuf)

	srv, err := server.NewServer(internalLogger, mainLogger)
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	err = srv.Listen("127.0.0.1:54021")
	assert.Nil(t, err)

	err = srv.Start()
	assert.Nil(t, err)

	srv.Stop()
}

func TestServerListen(t *testing.T) {

	var ibuf, mbuf bytes.Buffer
	internalLogger := zerolog.New(&ibuf).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	mainLogger := zerolog.New(&mbuf)

	srv, err := server.NewServer(internalLogger, mainLogger)
	srv.Listen("127.0.0.1:54021")

	assert.Nil(t, err)
	assert.NotNil(t, srv)

}

func TestSendMessage(t *testing.T) {

	var ibuf, mbuf bytes.Buffer
	internalLogger := zerolog.New(&ibuf).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	mainLogger := zerolog.New(&mbuf)

	srv, err := server.NewServer(internalLogger, mainLogger)
	assert.Nil(t, err)
	assert.NotNil(t, srv)

	err = srv.Listen("127.0.0.1:54021")
	assert.Nil(t, err)

	err = srv.Start()
	assert.Nil(t, err)

	client, err := NewClient("127.0.0.1:54021")
	assert.NotNil(t, client)
	assert.Nil(t, err)

	go func() {
		client.SendAuth()
		assert.Nil(t, err)
	}()

	srv.Stop()

	fmt.Println(mbuf.String())
}
