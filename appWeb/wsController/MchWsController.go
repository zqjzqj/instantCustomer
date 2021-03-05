package wsController

import (
	"context"
	"github.com/kataras/iris/v12/websocket"
	"github.com/zqjzqj/instantCustomer/models"
)

type MchWsController struct {
	*websocket.NSConn `stateless:"true"`
	Namespace         string
	MAccount *models.MchAccount
}



func (c *MchWsController) OnNamespaceDisconnect(msg websocket.Message) error {
	if c.MAccount.Room != nil {
		_ = c.MAccount.Room.Leave(context.Background())
	}
	return nil
}

func (c *MchWsController) OnNamespaceConnected(msg websocket.Message) error {
	c.MAccount.Room, _ = c.JoinRoom(context.Background(), c.MAccount.GetWsRoomId())
	return nil
}

func (c *MchWsController) OnChat(msg websocket.Message) error {
	return nil
}