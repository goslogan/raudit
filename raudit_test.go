package main_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/goslogan/raudit/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestServerCreation(t *testing.T) {

	var ibuf, mbuf bytes.Buffer
	internalLogger := zerolog.New(&ibuf).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	mainLogger := zerolog.New(&mbuf)

	srv := server.NewServer(internalLogger, mainLogger)
	assert.NotNil(t, srv)

	err := srv.Listen("127.0.0.1:54021")
	assert.Nil(t, err)

	err = srv.Start()
	assert.Nil(t, err)

	srv.Stop()
}

func TestServerListen(t *testing.T) {

	var ibuf, mbuf bytes.Buffer
	internalLogger := zerolog.New(&ibuf).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	mainLogger := zerolog.New(&mbuf)

	srv := server.NewServer(internalLogger, mainLogger)
	srv.Listen("127.0.0.1:54021")

	assert.NotNil(t, srv)

}

func TestSendMessage(t *testing.T) {

	var ibuf, mbuf bytes.Buffer
	internalLogger := zerolog.New(&ibuf).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	mainLogger := zerolog.New(&mbuf)

	srv := server.NewServer(internalLogger, mainLogger)
	assert.NotNil(t, srv)

	err := srv.Listen("127.0.0.1:54021")
	assert.Nil(t, err)

	err = srv.Start()
	assert.Nil(t, err)

	client, err := NewClient("127.0.0.1:54021")
	assert.NotNil(t, client)
	assert.Nil(t, err)

	ts := time.Now()
	sourceEvent := client.SendNewConn(ts)
	outputEvent := NewConnEvent{}

	err = json.Unmarshal(mbuf.Bytes(), &outputEvent)
	assert.NotNil(t, err)
	assert.Equal(t, sourceEvent, outputEvent)

	srv.Stop()
}
