package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = &websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsConnectsMux sync.Mutex
var wsConnects = make(map[string]*WsConn, 50)

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
	wsc := &WsConn{
		id:   _uuid.String(),
		conn: c,
	}
	wsConnectsMux.Lock()
	wsConnects[wsc.id] = wsc
	wsConnectsMux.Unlock()
	return wsc, nil
}

func FindWsConn(id string) (*WsConn, bool) {
	wsc, ok := wsConnects[id]
	return wsc, ok
}

func CloseById(id string) error {
	wsc, ok := FindWsConn(id)
	if ok {
		return wsc.Close()
	}
	return nil
}

func (wsc *WsConn) ListenMsg(msgChan chan *Message) error {
	var rErr error
	for {
		mt, message, err := wsc.conn.ReadMessage()
		if err != nil {
			rErr = err
			break
		}
		msg, err := NewMessageByJsonBytes(message, mt, "", "")
		if err != nil {
			rErr = err
			break
		}
		msgChan <- msg
	}
	_ = wsc.Close()
	return rErr
}

func (wsc *WsConn) WriteMsg(msg *Message) error {
	return wsc.conn.WriteMessage(msg.Mt, msg.ToMsgBytes())
}

func (wsc *WsConn) ID() string {
	return wsc.id
}

func (wsc *WsConn) Close() error {
	err := wsc.conn.Close()
	if err != nil {
		return err
	}
	wsConnectsMux.Lock()
	delete(wsConnects, wsc.id)
	wsConnectsMux.Unlock()
	return nil
}
