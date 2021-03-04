package models

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/zqjzqj/instantCustomer/global"
)

type FieldsExtendsJsonLongType struct {
	ExtendsJson	string	`gorm:"type:longtext;comment:扩展字段"`
}

func (e *FieldsExtendsJsonLongType) GetExtendsJson(key string) gjson.Result {
	return gjson.Get(e.ExtendsJson, key)
}

func (e *FieldsExtendsJsonLongType) SetExtendsJson(key string, value interface{}) {
	r := global.Json2Map(e.ExtendsJson)
	r[key] = value
	eJson, _ := json.Marshal(r)
	e.ExtendsJson = string(eJson)
}

