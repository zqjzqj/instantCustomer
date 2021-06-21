package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/zqjzqj/instantCustomer/appWeb/routes"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/logs"
	"github.com/zqjzqj/instantCustomer/migrates"
	"os"
	"regexp"
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
	routes.RegisterWsRoutes(app)

	err := ListenWeb(app)
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

func ListenWeb(appWeb *iris.Application) error {
	//注册api路由
	routes.RegisterApiRoutes(appWeb)

	logs.PrintlnInfo("Http API List:")
	port := config.GetWebCfg().GetPort()
	_regexp, _ := regexp.Compile("^/admin")
	for _, r := range appWeb.GetRoutes() {
		if r.Method != "OPTIONS" {
			if _regexp.MatchString(r.Path) {
				continue
			}
			logs.PrintlnInfo(fmt.Sprintf("[%s] http://127.0.0.1:%d%s", r.Method, port, r.Path))
		}
	}

	//监听http
	err := appWeb.Run(iris.Addr(fmt.Sprintf(":%d", port)), iris.WithConfiguration(iris.Configuration{
		TimeFormat: global.DateTimeFormatStr,
	}))
	if err != nil {
		return err
	}
	return nil
}
