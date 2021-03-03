package controller

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12/websocket"
	"github.com/zqjzqj/instantCustomer/logs"
	"github.com/zqjzqj/instantCustomer/models"
)

type WebsocketController struct {
	*websocket.NSConn `stateless:"true"`
	Namespace         string
	MAccount *models.MchAccount
}



func (c *WebsocketController) OnNamespaceDisconnect(msg websocket.Message) error {
	// This will call the "OnVisit" event on all clients, except the current one,
	// (it can't because it's left but for any case use this type of design)
	c.Conn.Server().Broadcast(nil, websocket.Message{
		Namespace: msg.Namespace,
		Event:     "OnVisit",
		Body:      []byte(fmt.Sprintf("%d", c.MAccount.ID)),
	})

	return nil
}

func (c *WebsocketController) OnNamespaceConnected(msg websocket.Message) error {
	// This will call the "OnVisit" event on all clients, including the current one,
	// with the 'newCount' variable.
	//
	// There are many ways that u can do it and faster, for example u can just send a new visitor
	// and client can increment itself, but here we are just "showcasing" the websocket controller.
	if c.MAccount.ID < 3 {
		_, err := c.JoinRoom(context.Background(), "cc")
		if err != nil {
			c.Conn.Close()
		}
	}

	c.Conn.Server().Broadcast(nil, websocket.Message{
		Namespace: msg.Namespace,
		Event:     "OnVisit",
		Body:      []byte(fmt.Sprintf("%d", c.MAccount.ID)),
	})

	return nil
}

func (c *WebsocketController) OnChat(msg websocket.Message) error {
	ctx := websocket.GetContext(c.Conn)
	ctx.Application().Logger().Infof("[IP: %s] [ID: %s]  broadcast to other clients the message [%s]",
		ctx.RemoteAddr(), c, string(msg.Body))
	logs.PrintlnSuccess(c.Conn.ID())
	c.Conn.Write(websocket.Message{
		Namespace: c.Namespace,
		Event: "OnChat",
		Body: []byte(c.MAccount.RealName),
	})
	room := c.Room("cc")
	if room != nil {
		logs.PrintlnSuccess(c.Rooms())
	} else {
		logs.PrintErr("未加入到CC房间")
	}
	c.Conn.Server().Broadcast(nil, websocket.Message{
		Namespace:c.Namespace,
		Event:"OnChat",
		Room:"cc",
		Body:[]byte("ROOM"),
	})
	return nil
}
