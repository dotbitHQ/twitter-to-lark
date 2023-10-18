package config

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/scorpiotzh/toolib"
)

var (
	Cfg CfgServer
	log = logger.NewLogger("config", logger.LevelDebug)
)

func InitCfg(configFilePath string) error {
	if configFilePath == "" {
		configFilePath = "../config/config.yaml"
	}
	log.Info("config file：", configFilePath)
	if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
		return fmt.Errorf("UnmarshalYamlFile err:%s", err.Error())
	}
	log.Info("config file：ok")
	return nil
}

func AddCfgFileWatcher(configFilePath string) (*fsnotify.Watcher, error) {
	if configFilePath == "" {
		configFilePath = "../config/config.yaml"
	}
	return toolib.AddFileWatcher(configFilePath, func() {
		log.Info("update config file：", configFilePath)
		if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
			log.Error("UnmarshalYamlFile err:", err.Error())
		}
		log.Info("update config file：ok")
	})
}

type CfgServer struct {
	Twitter struct {
		BearerToken string `json:"bearer_token" yaml:"bearer_token"`
		TwitterList []struct {
			LarkKey  string `json:"lark_key" yaml:"lark_key"`
			Username string `json:"username" yaml:"username"`
		} `json:"twitter_list" yaml:"twitter_list"`
	} `json:"twitter" yaml:"twitter"`

	Cache struct {
		Redis struct {
			Addr     string `json:"addr" yaml:"addr"`
			Password string `json:"password" yaml:"password"`
			DbNum    int    `json:"db_num" yaml:"db_num"`
		} `json:"redis" yaml:"redis"`
	} `json:"cache" yaml:"cache"`
}
