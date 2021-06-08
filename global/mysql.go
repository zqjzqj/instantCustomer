package global

import (
	"github.com/zqjzqj/instantCustomer/sErr"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var mysqlDefDb *MysqlDb
var mysqlDbs map[string]*MysqlDb

type MysqlDb struct {
	Host          string
	Port          string
	Database      string
	Charset       string
	UserName      string
	Password      string
	MaxIdleCounts int
	*gorm.DB
}

func init() {
	mysqlDbs = make(map[string]*MysqlDb)
}

func SetMysql(k string, db *MysqlDb, isDef bool) {
	mysqlDbs[k] = db
	if isDef {
		mysqlDefDb = db
	}
}

func GetMysql(k string) (*MysqlDb, bool) {
	ok, d := mysqlDbs[k]
	return ok, d
}

func GetMysqlDef() *MysqlDb {
	return mysqlDefDb
}

func NewDatabaseMysql(host, port, database, charset, username, password string, maxIdleCounts int, maxOpenCounts int) (*MysqlDb, error) {
	db, err := gorm.Open(mysql.Open(username+":"+password+"@("+host+":"+port+")/"+database+"?charset="+charset+"&parseTime=True&loc=Local"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
	})

	if err != nil {
		return nil, sErr.NewByError(err)
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
	sqlDB.SetConnMaxLifetime(1 * time.Minute)
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &MysqlDb{
		Host:          host,
		Port:          port,
		Database:      database,
		Charset:       charset,
		UserName:      username,
		Password:      password,
		MaxIdleCounts: maxIdleCounts,
		DB:            db,
	}, nil
}

func (db *MysqlDb) LockForUpdate() *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}
