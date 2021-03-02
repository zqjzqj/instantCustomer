package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/zqjzqj/instantCustomer/models"
)

type MerchantController struct {
	MchAccount *models.MchAccount
}

func (m *MerchantController) Get(ctx iris.Context) interface{} {
	return m.MchAccount
}
