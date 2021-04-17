// https://github.com/gobwas/ws/blob/master/example/autobahn/autobahn.go

package ws

import (
	ej "encoding/json"
	"io"
	"net"
	"sync"
	"time"

	"authkey/pkg/lib/log"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	// WSManager define a ws server manager
	Manager = &WSManager{
		Broadcast:  make(chan []byte, 10),
		Register:   make(chan *Client, 10),
		Unregister: make(chan *Client, 10),
		Clients:    make(map[string]*Client, 20),
	}
)

// WSManager is a websocket manager
type WSManager struct {
	lock       sync.RWMutex
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	Type int64               `json:"type"`
	Raw  jsoniter.RawMessage `json:"raw"`
}

type Message2 struct {
	Type int64         `json:"type"`
	Raw  ej.RawMessage `json:"raw"`
}

func (manager *WSManager) Run() {
	log.Glog.Info("ws manager run")

	for {
		select {
		case client := <-manager.Register:
			log.Glog.Info("ws register", zap.String("id", client.Id))

			manager.lock.Lock()

			Manager.Clients[client.Id] = client

			manager.lock.Unlock()
		case client := <-manager.Unregister:
			log.Glog.Info("ws unregister", zap.String("id", client.Id))

			manager.lock.Lock()

			if _, ok := manager.Clients[client.Id]; ok {
				delete(manager.Clients, client.Id)

				if err := client.Close(); err != nil {
					log.Glog.Info("client close", zap.String("id", client.Id))
				}
			}

			manager.lock.Unlock()
		case m := <-Manager.Broadcast:
			manager.lock.Lock()

			for _, client := range Manager.Clients {
				select {
				case client.SendChan <- m:
				}
			}

			manager.lock.Unlock()
		}
	}
}

// Client 单个 websocket 信息
type Client struct {
	Id       string
	Conn     net.Conn
	SendChan chan []byte
	Reader   *wsutil.Reader
	Writer   *wsutil.Writer
}

func (c *Client) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	close(c.SendChan)

	return nil
}

// 读信息，从 websocket 连接直接读取数据
// 数据message和close message是单独发送的, close message允许携带reason
// 逻辑from wsutil.ReadClientData(conn)
func (c *Client) Read() {
	defer func() {
		log.Glog.Info("client reader exit", zap.String("id", c.Id))

		Manager.Unregister <- c
	}()

	for {
		hdr, err := c.Reader.NextFrame()
		if err != nil {
			log.Glog.Error("client read header error", zap.String("id", c.Id), zap.Error(err))

			return
		}

		payload := make([]byte, hdr.Length) // len=0 when hdr.OpCode == ws.OpClose
		if _, err = io.ReadFull(c.Reader, payload); err != nil {
			log.Glog.Error("client read msg error", zap.String("id", c.Id), zap.Error(err))

			return
		}

		log.Glog.Debug("playload", zap.String("id", c.Id), zap.Any("op", hdr.OpCode), zap.Int64("len", hdr.Length), zap.String("data", string(payload)))

		if hdr.OpCode == ws.OpClose {
			code, reason := ws.ParseCloseFrameData(payload)
			log.Glog.Info("client close reason", zap.String("id", c.Id), zap.Any("code", code), zap.String("reason", reason))

			c.Conn.Write(ws.CompiledClose)
			return
		}

		// todo
		c.SendChan <- payload
	}
}

// 写信息，从 SendChan 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		log.Glog.Info("client Write exit", zap.String("id", c.Id))

		Manager.Unregister <- c
	}()

	for {
		select {
		case msg, ok := <-c.SendChan:
			if !ok {
				c.Conn.Write(ws.CompiledClose)
				return
			}
			log.Glog.Debug("send", zap.String("id", c.Id), zap.String("data", string(msg)))

			if err := wsutil.WriteServerMessage(c.Conn, ws.OpText, msg); err != nil {
				log.Glog.Error("client write error", zap.String("id", c.Id), zap.Error(err))

				return
			}
		}
	}
}

func (c *Client) WriteTest() {
	// 测试单个 client 发送数据 by wsutil.WriteServerMessage
	c.SendChan <- []byte("Send message ----" + time.Now().Format("2006-01-02 15:04:05"))

	// write by wsutil.NewWriter
	c.Writer.Reset(c.Conn, ws.StateServerSide, ws.OpText)
	c.Writer.Write([]byte("test2"))
	if err := c.Writer.Flush(); err != nil {
		log.Glog.Error("client write flush error", zap.String("id", c.Id), zap.Error(err))
	}
}
