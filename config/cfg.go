package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/zqjzqj/instantCustomer/global"
	"github.com/zqjzqj/instantCustomer/logs"
	"github.com/zqjzqj/instantCustomer/sErr"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var cfg *Cfg

func init() {
	cfg = &Cfg{}
}

type Cfg struct {
	logs *LogsCfg
	web  *Web
}

type LogsCfg struct {
	IsPrint     bool
	LogFilePath string
	LogFile     *os.File
}

func GetLogCfg() *LogsCfg {
	return cfg.logs
}

func GetWebCfg() *Web {
	return cfg.web
}

func LoadConfigJson(p string) error {
	logs.PrintlnInfo("reload config.....")
	defer logs.PrintlnSuccess("reload config success!")
	paths, fileName := filepath.Split(p)
	viper.SetConfigName(fileName)
	viper.AddConfigPath(paths)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	//载入db配置
	dbConfigs := viper.GetStringMap("db")
	for key, v := range dbConfigs {
		b := make([]byte, 0)
		b, err := json.Marshal(v)
		if err != nil {
			return sErr.NewByError(err)
		}
		dbConf := gjson.ParseBytes(b)
		maxIdleCounts, _ := strconv.Atoi(dbConf.Get("max_idle_counts").String())
		charset := dbConf.Get("charset").String()
		if charset == "" {
			charset = "utf8"
		}
		db, err := global.NewDatabaseMysql(dbConf.Get("host").String(), dbConf.Get("port").String(), dbConf.Get("database").String(), charset, dbConf.Get("username").String(), dbConf.Get("password").String(), maxIdleCounts, int(dbConf.Get("max_open_counts").Int()))
		if err != nil {
			return sErr.NewByError(err)
		}
		if dbConf.Get("default").Bool() == true {
			global.SetMysql("default", db, true)
		} else {
			global.SetMysql(key, db, false)
		}
	}
	if global.GetMysqlDef() == nil {
		iDB, ok := global.GetMysql("instant_customer")
		if !ok {
			return sErr.New("not default database")
		}
		global.SetMysql("default", iDB, true)
	}

	//载入web配置
	w := &Web{}
	w.port = viper.GetUint64("web.port")
	if w.port == 0 {
		w.port = 80
	}
	cfg.web = w

	//载入日志配置
	logCfg := &LogsCfg{}
	logCfg.IsPrint = viper.GetBool("log.isPrint")
	logCfg.LogFilePath = viper.GetString("log.logFilePath")
	logs.IsPrintLog = logCfg.IsPrint
	if logCfg.LogFilePath != "" {
		logCfg.LogFile, err = os.OpenFile(logCfg.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		os.Stdout = logCfg.LogFile
		os.Stderr = logCfg.LogFile
		go func() {
			t := time.NewTicker(36 * time.Hour)
			defer t.Stop()
			log.Println("创建自动清除日志文件内容")
			for {
				select {
				case <-t.C:
					err = logCfg.LogFile.Truncate(0)
					if err != nil {
						log.Println("清空文件失败")
					}
					_, _ = logCfg.LogFile.Seek(0, 0)

				}
			}
		}()
	}
	cfg.logs = logCfg
	return nil
}
