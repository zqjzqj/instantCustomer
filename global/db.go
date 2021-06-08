package global

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func LockForUpdate(tx *gorm.DB) *gorm.DB {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"})
}

//这里是默认的db 后续如果跟换数据库驱动的话 可以在这里写代码
func GetDb() *gorm.DB {
	return GetMysqlDef().DB
}