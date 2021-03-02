package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/zqjzqj/instantCustomer/sErr"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var cfg *Cfg

func init() {
	cfg = &Cfg{
		dbs:make(map[string]*Database),
		defaultDb:&Database{},
	}
}

type Cfg struct {
	dbs map[string]*Database
	defaultDb *Database
	logs *LogsCfg
	web *Web
}

type LogsCfg struct {
	IsPrint bool
	LogFilePath string
	LogFile *os.File
}

func GetDbs() map[string]*Database {
	return cfg.dbs
}

func GetDb(key string) (*Database, error) {
	if db, ok := cfg.dbs[key]; ok {
		return db, nil
	}
	return nil, sErr.New("无效的key")
}

func GetDbDefault() *Database {
	return cfg.defaultDb
}

func GetLogCfg() *LogsCfg {
	return cfg.logs
}

func GetWebCfg() *Web {
	return cfg.web
}

func LoadConfigJson(p string) error {
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
		db, err := NewDatabase(dbConf.Get("host").String(), dbConf.Get("port").String(), dbConf.Get("database").String(), charset, dbConf.Get("username").String(), dbConf.Get("password").String(), maxIdleCounts, int(dbConf.Get("max_open_counts").Int()))
		if err != nil {
			return sErr.NewByError(err)
		}
		if cfg.defaultDb.host == "" && dbConf.Get("default").Bool() == true {
			cfg.defaultDb = db
		}
		cfg.dbs[key] = db
	}
	if cfg.defaultDb.host == "" {
		cfg.defaultDb, err = GetDb("instant_customer")
		if err != nil {
			return err
		}
	}

	//载入web配置
	w := &Web{}
	w.port = viper.GetString("web.port")
	cfg.web = w

	//载入日志配置
	logCfg := &LogsCfg{}
	logCfg.IsPrint = viper.GetBool("log.isPrint")
	logCfg.LogFilePath = viper.GetString("log.logFilePath")
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