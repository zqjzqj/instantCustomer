package config

import (
	"github.com/zqjzqj/instantCustomer/sErr"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	host string
	port string
	database string
	charset string
	userName string
	password string
	MaxIdleCounts int
	*gorm.DB
}

func NewDatabase(host, port, database, charset, username, password string, maxIdleCounts int, maxOpenCounts int) (*Database, error) {
	db, err := gorm.Open(mysql.Open(username+":"+password+"@("+host+":"+port+")/"+database+"?charset="+charset+"&parseTime=True&loc=Local"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating:true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return  nil, sErr.NewByError(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if maxIdleCounts > 0 {
		sqlDB.SetMaxIdleConns(maxIdleCounts)
	}
	if maxOpenCounts > 0 {
		sqlDB.SetMaxOpenConns(maxOpenCounts)
	}
	//sqlDB.SetConnMaxLifetime(2 * time.Minute)
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return &Database{
		host:          host,
		port:          port,
		database:      database,
		charset:       charset,
		userName:      username,
		password:      password,
		MaxIdleCounts: maxIdleCounts,
		DB:            db,
	}, nil
}

func (db *Database) LockForUpdate() *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}
func  LockForUpdate(db *gorm.DB) *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}