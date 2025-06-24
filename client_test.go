package main_test

import (
	"encoding/json"
	"math/rand/v2"
	"net"
	"time"
)

type Client struct {
	conn net.Conn
}

type NewConnection struct {
	ID      uint64 `json:"id"`
	Srcip   string `json:"srcip"`
	Srcp    string `json:"srcp"`
	Trgip   string `json:"trgip"`
	Trgp    string `json:"trgp"`
	Hname   string `json:"hname"`
	BdbName string `json:"bdb_name"`
	BdbUID  string `json:"bdb_uid"`
}
type AuthAudit struct {
	TS      int64         `json:"ts"`
	NewConn NewConnection `json:"new_conn"`
}

func NewClient(target string) (*Client, error) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil

}

func (client *Client) Send(input any) {

	go func() {
		buf, _ := json.Marshal(input)
		client.conn.Write(buf)
	}()

}

func (client *Client) SendAuth() {

	event := AuthAudit{
		TS: time.Now().Unix(),
		NewConn: NewConnection{
			ID:      rand.Uint64(),
			Srcip:   "127.0.0.1",
			Srcp:    "62301",
			Trgip:   "127.1.0.1",
			Trgp:    "120001",
			Hname:   "test.local",
			BdbName: "test",
			BdbUID:  "1",
		},
	}

	client.Send(event)
}
