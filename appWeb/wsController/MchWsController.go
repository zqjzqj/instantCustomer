package wsController

import (
	"github.com/kataras/iris/v12/websocket"
	"github.com/zqjzqj/instantCustomer/models"
)

type MchWsController struct {
	*websocket.NSConn `stateless:"true"`
	Namespace         string
	MAccount *models.MchAccount
}



func (c *MchWsController) OnNamespaceDisconnect(msg websocket.Message) error {
	return nil
}

func (c *MchWsController) OnNamespaceConnected(msg websocket.Message) error {
	return nil
}

func (c *MchWsController) OnChat(msg websocket.Message) error {
	return nil
}