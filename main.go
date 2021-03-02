package main

import (
	"flag"
	"github.com/kataras/iris/v12"
	"github.com/zqjzqj/instantCustomer/appWeb/routes"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/logs"
	"github.com/zqjzqj/instantCustomer/migrates"
	"os"
)

var configPath = flag.String("config", "./ic_config.yml", "配置文件路径")
var migrateCmd = flag.String("migrate", "", "迁移参数 run执行迁移 rollback回滚迁移")
var mRollbackId = flag.String("mRollbackId", "", "回滚迁移所需要的版本号")

func init() {
	flag.Parse()
	err := config.LoadConfigJson(*configPath)
	if err != nil {
		logs.Fatal("配置文件载入错误", err)
	}
}

func main() {
	migrateFunc()

	app := iris.New()
	//注册api路由
	routes.RegisterApiRoutes(app)
	err := app.Run(iris.Addr(":" + config.GetWebCfg().GetPort()), iris.WithConfiguration(iris.Configuration{
		TimeFormat: global.DateTimeFormatStr,
	}))
	if err != nil {
		logs.Fatal(err)
	}
}

func migrateFunc() {
	if *migrateCmd == "" {
		return
	}
	if *migrateCmd == "run" {
		migrates.Migrate()
	} else if *migrateCmd == "rollback" {
		if *mRollbackId == "" {
			logs.Fatal("无效的回退版本号【请填写参数 mRollbackId】")
		}
		migrates.Rollback(*mRollbackId)
	}
	os.Exit(0)
}
