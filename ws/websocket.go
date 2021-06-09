package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zqjzqj/instantCustomer/logs"
	"log"
	"net/http"
)

var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.PrintErr(err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s, type:%d", message, mt)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

type WsConn struct {
	id   string
	conn *websocket.Conn
}

func NewWsConn(w http.ResponseWriter, r *http.Request) (*WsConn, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	_uuid := uuid.New()
	return &WsConn{
		id:   _uuid.String(),
		conn: c,
	}, nil
}

func (wsc *WsConn) ID() string {
	return wsc.id
}

func (wsc *WsConn) Close() error {
	return wsc.conn.Close()
}
